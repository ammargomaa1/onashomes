package customers

import (
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

// Search finds customers by first name, last name, or phone
func (r *Repository) Search(query string, storeFrontID *int64, limit int) ([]models.Customer, error) {
	var customers []models.Customer
	searchPattern := "%" + query + "%"

	dbQuery := r.db.Where("first_name ILIKE ? OR last_name ILIKE ? OR phone ILIKE ?", searchPattern, searchPattern, searchPattern)

	if storeFrontID != nil {
		dbQuery = dbQuery.Where("store_front_id = ?", *storeFrontID)
	}

	err := dbQuery.Debug().Order("created_at DESC").
		Limit(limit).
		Find(&customers).Error

	return customers, err
}

// List retrieves customers with pagination and filtering
func (r *Repository) List(pagination *utils.Pagination, filter string, storeFrontID *int64) ([]models.Customer, error) {
	var customers []models.Customer
	query := r.db.Model(&models.Customer{})

	if filter != "" {
		searchPattern := "%" + filter + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if storeFrontID != nil {
		query = query.Where("store_front_id = ?", storeFrontID)
	}

	// Count total records
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.SetTotal(total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query = query.Offset(offset).Limit(pagination.Limit)

	// Apply sorting
	if pagination.Sort != "" {
		orderClause := pagination.Sort
		if pagination.Order != "" {
			orderClause += " " + pagination.Order
		}
		query = query.Order(orderClause)
	} else {
		query = query.Order("created_at desc")
	}

	err := query.Find(&customers).Error
	return customers, err
}

// Create adds a new customer
func (r *Repository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

// GetByPhone finds a customer by exact phone number
func (r *Repository) GetByPhone(phone string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByID finds a customer by ID
func (r *Repository) GetByID(id int64) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).
		Preload("Orders.OrderStatus").
		Preload("Orders.PaymentStatus").
		Preload("Orders.FulfillmentStatus").
		Preload("Orders.Currency").
		First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Update modifies an existing customer
func (r *Repository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}
