package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds db for stats queries
type Handler struct {
	db *gorm.DB
}

// NewHandler creates a stats handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ItemOrderCountRow is a single row in the item order counts response
type ItemOrderCountRow struct {
	VariantID         int64  `json:"variant_id"`
	ProductName       string `json:"product_name"`
	AttributeValue    string `json:"attribute_value"`
	CurrentStock      int    `json:"current_stock"`
	DeliveredQuantity int    `json:"delivered_quantity"`
	RequiredQuantity  int    `json:"required_quantity"`
}

// GetItemOrderCounts returns aggregated variant-level order counts
func (h *Handler) GetItemOrderCounts(c *gin.Context) {
	var rows []ItemOrderCountRow

	sql := `
		SELECT
			pv.id AS variant_id,
			p.name AS product_name,
			pv.attribute_value,
			COALESCE(SUM(vi.quantity - vi.reserved_quantity), 0) AS current_stock,
			COALESCE(delivered.qty, 0) AS delivered_quantity,
			COALESCE(required.qty, 0) AS required_quantity
		FROM product_variants pv
		JOIN products p ON p.id = pv.product_id AND p.deleted_at IS NULL
		LEFT JOIN variant_inventory vi ON vi.product_variant_id = pv.id
		LEFT JOIN (
			SELECT oi.product_variant_id, SUM(oi.quantity) AS qty
			FROM order_items oi
			JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
			JOIN order_statuses os ON os.id = o.order_status_id AND os.slug = 'completed'
			GROUP BY oi.product_variant_id
		) delivered ON delivered.product_variant_id = pv.id
		LEFT JOIN (
			SELECT oi.product_variant_id, SUM(oi.quantity) AS qty
			FROM order_items oi
			JOIN orders o ON o.id = oi.order_id AND o.deleted_at IS NULL
			JOIN fulfillment_statuses fs ON fs.id = o.fulfillment_status_id AND fs.slug != 'fulfilled'
			JOIN order_statuses os ON os.id = o.order_status_id AND os.slug NOT IN ('completed','cancelled','returned','refunded','draft')
			GROUP BY oi.product_variant_id
		) required ON required.product_variant_id = pv.id
		WHERE pv.deleted_at IS NULL
		GROUP BY pv.id, p.name, pv.attribute_value, delivered.qty, required.qty
		ORDER BY pv.id
	`

	if err := h.db.Raw(sql).Scan(&rows).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"data": []ItemOrderCountRow{}, "message": "No data available"})
		return
	}

	if rows == nil {
		rows = []ItemOrderCountRow{}
	}

	c.JSON(http.StatusOK, gin.H{"data": rows})
}
