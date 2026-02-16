package orders

import (
	"fmt"
	"strings"

	"github.com/onas/ecommerce-api/internal/api/orders/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateOrder creates an order and its items transactionally
// Note: This expects to be called within a transaction (tx)
func (r *Repository) CreateOrder(tx *gorm.DB, order *models.Order) error {
	if err := tx.Create(order).Error; err != nil {
		return err
	}
	return nil
}

// CreateOrderItems batch inserts order items
func (r *Repository) CreateOrderItems(tx *gorm.DB, items []models.OrderItem) error {
	return tx.CreateInBatches(items, 100).Error
}

// GetOrderByID retrieves an order with preloaded items and storefront
func (r *Repository) GetOrderByID(id int64) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").Preload("StoreFront").
		Preload("OrderStatus").Preload("PaymentStatus").Preload("FulfillmentStatus").Preload("Currency").
		Preload("CreatedBy").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus updates the order status ID
func (r *Repository) UpdateStatus(tx *gorm.DB, id int64, statusID int64) error {
	return tx.Model(&models.Order{}).Where("id = ?", id).Update("order_status_id", statusID).Error
}

// Helper methods to get status IDs
func (r *Repository) GetOrderStatusBySlug(slug string) (*models.OrderStatus, error) {
	var status models.OrderStatus
	err := r.db.Where("slug = ?", slug).First(&status).Error
	return &status, err
}

func (r *Repository) GetPaymentStatusBySlug(slug string) (*models.PaymentStatus, error) {
	var status models.PaymentStatus
	err := r.db.Where("slug = ?", slug).First(&status).Error
	return &status, err
}

func (r *Repository) GetFulfillmentStatusBySlug(slug string) (*models.FulfillmentStatus, error) {
	var status models.FulfillmentStatus
	err := r.db.Where("slug = ?", slug).First(&status).Error
	return &status, err
}

func (r *Repository) GetCurrencyByCode(code string) (*models.Currency, error) {
	var currency models.Currency
	err := r.db.Where("code = ?", code).First(&currency).Error
	return &currency, err
}

// ListOrders retrieves orders with filtering and pagination
func (r *Repository) ListOrders(filter requests.OrderFilterRequest, pagination *utils.Pagination) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{}).
		Preload("StoreFront").
		Preload("OrderStatus").Preload("PaymentStatus").Preload("FulfillmentStatus").Preload("Currency").
		Preload("CreatedBy")

	// Filtering
	if filter.StoreFrontID > 0 {
		query = query.Where("orders.store_front_id = ?", filter.StoreFrontID)
	}

	if filter.Status != "" {
		// Join order_statuses to filter by slug.
		query = query.Joins("JOIN order_statuses ON orders.order_status_id = order_statuses.id").
			Where("order_statuses.slug = ?", filter.Status)
	}

	if filter.PaymentStatus != "" {
		query = query.Joins("JOIN payment_statuses ON orders.payment_status_id = payment_statuses.id").
			Where("payment_statuses.slug = ?", filter.PaymentStatus)
	}

	if filter.ProductIDs != "" {
		ids := strings.Split(filter.ProductIDs, ",")
		if len(ids) > 0 {
			query = query.Joins("JOIN order_items ON orders.id = order_items.order_id").
				Where("order_items.product_id IN ?", ids).
				Group("orders.id")
		}
	}

	if filter.DateFrom != "" {
		query = query.Where("orders.created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("orders.created_at <= ?", filter.DateTo)
	}
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		// Search in Order Number, Customer Name, Customer Phone
		query = query.Where(
			r.db.Where("LOWER(orders.order_number) LIKE ?", search).
				Or("LOWER(orders.customer_name) LIKE ?", search).
				Or("orders.customer_phone LIKE ?", search),
		)
	}

	// Count total
	// Note: When using Group, Count behaves differently (counts groups).
	// To get total count of orders matching criteria, we might need a subquery or careful counting.
	// For simplicity, if Group is used, we count the results.
	// GORM's Count method with Group might be tricky.
	// A common workaround is counting distinct IDs.
	if filter.ProductIDs != "" {
		if err := query.Distinct("orders.id").Count(&total).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	// Pagination & Sorting
	offset := (pagination.Page - 1) * pagination.Limit
	query = query.Offset(offset).Limit(pagination.Limit)

	if pagination.Sort != "" {
		query = query.Order(fmt.Sprintf("orders.%s %s", pagination.Sort, pagination.Order))
	} else {
		query = query.Order("orders.created_at DESC")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ListOrderStatuses retrieves all order statuses
func (r *Repository) ListOrderStatuses() ([]models.OrderStatus, error) {
	var statuses []models.OrderStatus
	err := r.db.Find(&statuses).Error
	return statuses, err
}

// ListPaymentStatuses retrieves all payment statuses
func (r *Repository) ListPaymentStatuses() ([]models.PaymentStatus, error) {
	var statuses []models.PaymentStatus
	err := r.db.Find(&statuses).Error
	return statuses, err
}

// ListFulfillmentStatuses retrieves all fulfillment statuses
func (r *Repository) ListFulfillmentStatuses() ([]models.FulfillmentStatus, error) {
	var statuses []models.FulfillmentStatus
	err := r.db.Find(&statuses).Error
	return statuses, err
}

// ListCurrencies retrieves all currencies
func (r *Repository) ListCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.Find(&currencies).Error
	return currencies, err
}
