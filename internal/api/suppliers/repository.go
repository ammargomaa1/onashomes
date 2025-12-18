package suppliers

import (
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Repository interface {
	Create(supplier *models.Supplier) error
	Update(supplier *models.Supplier) error
	GetByID(id int64) (*models.Supplier, error)
	List(pagination *utils.Pagination) ([]models.Supplier, int64, error)
	ToggleStatus(id int64, isActive bool, updatedBy int64) error
}
