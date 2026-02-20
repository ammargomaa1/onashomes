package dashboard

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds db for dashboard queries
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a dashboard handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// --- Response types ---

type KPIs struct {
	RevenueToday    float64 `json:"revenue_today"`
	RevenueMonth    float64 `json:"revenue_month"`
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	TotalCustomers  int64   `json:"total_customers"`
	LowStockCount   int64   `json:"low_stock_count"`
}

type ChartPoint struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount,omitempty"`
	Count  int     `json:"count,omitempty"`
}

type TopProduct struct {
	Name      string `json:"name"`
	TotalSold int    `json:"total_sold"`
}

type CategorySales struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

type RecentOrder struct {
	ID           int64   `json:"id"`
	OrderNumber  string  `json:"order_number"`
	CustomerName string  `json:"customer_name"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
}

type RecentCustomer struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}

type LowStockItem struct {
	VariantID      int64  `json:"variant_id"`
	ProductName    string `json:"product_name"`
	AttributeValue string `json:"attribute_value"`
	Quantity       int    `json:"quantity"`
	Threshold      int    `json:"threshold"`
}

type OverviewResponse struct {
	KPIs            KPIs             `json:"kpis"`
	RevenueChart    []ChartPoint     `json:"revenue_chart"`
	OrdersChart     []ChartPoint     `json:"orders_chart"`
	TopProducts     []TopProduct     `json:"top_products"`
	SalesByCategory []CategorySales  `json:"sales_by_category"`
	RecentOrders    []RecentOrder    `json:"recent_orders"`
	RecentCustomers []RecentCustomer `json:"recent_customers"`
	LowStockItems   []LowStockItem   `json:"low_stock_items"`
}

// GetOverview returns all dashboard data in one response
func (h *Handler) GetOverview(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	last30Days := now.AddDate(0, 0, -30)

	var resp OverviewResponse

	// --- KPIs ---
	// Revenue today
	h.db.Raw(`SELECT COALESCE(SUM(total_amount), 0) FROM orders
		WHERE deleted_at IS NULL AND created_at >= ?
		AND order_status_id = (SELECT id FROM order_statuses WHERE slug = 'completed')`,
		todayStart).Scan(&resp.KPIs.RevenueToday)

	// Revenue this month
	h.db.Raw(`SELECT COALESCE(SUM(total_amount), 0) FROM orders
		WHERE deleted_at IS NULL AND created_at >= ?
		AND order_status_id = (SELECT id FROM order_statuses WHERE slug = 'completed')`,
		monthStart).Scan(&resp.KPIs.RevenueMonth)

	// Total orders
	h.db.Raw(`SELECT COUNT(*) FROM orders WHERE deleted_at IS NULL`).Scan(&resp.KPIs.TotalOrders)

	// Pending orders (confirmed but not completed/cancelled)
	h.db.Raw(`SELECT COUNT(*) FROM orders o
		JOIN order_statuses os ON os.id = o.order_status_id
		WHERE o.deleted_at IS NULL AND os.slug IN ('confirmed', 'fulfilled')`).Scan(&resp.KPIs.PendingOrders)

	// Completed orders
	h.db.Raw(`SELECT COUNT(*) FROM orders o
		JOIN order_statuses os ON os.id = o.order_status_id
		WHERE o.deleted_at IS NULL AND os.slug = 'completed'`).Scan(&resp.KPIs.CompletedOrders)

	// Total customers
	h.db.Raw(`SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL`).Scan(&resp.KPIs.TotalCustomers)

	// Low stock count
	h.db.Raw(`SELECT COUNT(*) FROM variant_inventory
		WHERE (quantity - reserved_quantity) <= low_stock_threshold`).Scan(&resp.KPIs.LowStockCount)

	// --- Charts ---
	// Revenue last 30 days
	var revenueChart []ChartPoint
	h.db.Raw(`SELECT TO_CHAR(o.created_at, 'YYYY-MM-DD') AS date, COALESCE(SUM(o.total_amount), 0) AS amount
		FROM orders o
		JOIN order_statuses os ON os.id = o.order_status_id
		WHERE o.deleted_at IS NULL AND o.created_at >= ? AND os.slug = 'completed'
		GROUP BY TO_CHAR(o.created_at, 'YYYY-MM-DD')
		ORDER BY date`, last30Days).Scan(&revenueChart)
	if revenueChart == nil {
		revenueChart = []ChartPoint{}
	}
	resp.RevenueChart = revenueChart

	// Orders per day last 30 days
	var ordersChart []ChartPoint
	h.db.Raw(`SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS date, COUNT(*)::int AS count
		FROM orders
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		ORDER BY date`, last30Days).Scan(&ordersChart)
	if ordersChart == nil {
		ordersChart = []ChartPoint{}
	}
	resp.OrdersChart = ordersChart

	// Top 10 selling products
	var topProducts []TopProduct
	h.db.Raw(`SELECT
			CONCAT(p.name, CASE WHEN pv.attribute_value != '' THEN ' - ' || pv.attribute_value ELSE '' END) AS name,
			SUM(oi.quantity)::int AS total_sold
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
		JOIN order_statuses os ON os.id = o.order_status_id AND os.slug = 'completed'
		JOIN product_variants pv ON pv.id = oi.product_variant_id
		JOIN products p ON p.id = pv.product_id
		GROUP BY p.name, pv.attribute_value
		ORDER BY total_sold DESC
		LIMIT 10`).Scan(&topProducts)
	if topProducts == nil {
		topProducts = []TopProduct{}
	}
	resp.TopProducts = topProducts

	// Sales by category
	var salesByCat []CategorySales
	h.db.Raw(`SELECT
			COALESCE(c.name_en, 'Uncategorized') AS category,
			COALESCE(SUM(oi.total_price), 0) AS total
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
		JOIN order_statuses os ON os.id = o.order_status_id AND os.slug = 'completed'
		JOIN products p ON p.id = oi.product_id
		LEFT JOIN categories c ON c.id = p.category_id
		GROUP BY c.name_en
		ORDER BY total DESC
		LIMIT 8`).Scan(&salesByCat)
	if salesByCat == nil {
		salesByCat = []CategorySales{}
	}
	resp.SalesByCategory = salesByCat

	// --- Recent Activity ---
	// Latest 10 orders
	var recentOrders []RecentOrder
	h.db.Raw(`SELECT o.id, o.order_number, o.customer_name, o.total_amount,
			os.name_en AS status, TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM orders o
		JOIN order_statuses os ON os.id = o.order_status_id
		WHERE o.deleted_at IS NULL
		ORDER BY o.created_at DESC
		LIMIT 10`).Scan(&recentOrders)
	if recentOrders == nil {
		recentOrders = []RecentOrder{}
	}
	resp.RecentOrders = recentOrders

	// Latest 5 customers
	var recentCustomers []RecentCustomer
	h.db.Raw(`SELECT id, first_name, last_name, phone,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 5`).Scan(&recentCustomers)
	if recentCustomers == nil {
		recentCustomers = []RecentCustomer{}
	}
	resp.RecentCustomers = recentCustomers

	// Low stock alerts
	var lowStock []LowStockItem
	h.db.Raw(`SELECT vi.product_variant_id AS variant_id, p.name AS product_name,
			pv.attribute_value, (vi.quantity - vi.reserved_quantity) AS quantity, vi.low_stock_threshold AS threshold
		FROM variant_inventory vi
		JOIN product_variants pv ON pv.id = vi.product_variant_id AND pv.deleted_at IS NULL
		JOIN products p ON p.id = pv.product_id AND p.deleted_at IS NULL
		WHERE (vi.quantity - vi.reserved_quantity) <= vi.low_stock_threshold
		ORDER BY (vi.quantity - vi.reserved_quantity) ASC
		LIMIT 10`).Scan(&lowStock)
	if lowStock == nil {
		lowStock = []LowStockItem{}
	}
	resp.LowStockItems = lowStock

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
