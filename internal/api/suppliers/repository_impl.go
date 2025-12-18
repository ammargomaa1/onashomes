package suppliers

import (
	"gorm.io/gorm"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(supplier *models.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *repository) Update(supplier *models.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *repository) GetByID(id int64) (*models.Supplier, error) {
	var supplier models.Supplier
	err := r.db.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *repository) List(pagination *utils.Pagination) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier
	var total int64

	err := pagination.Paginate(r.db).Preload("CreatedByAdmin").Model(suppliers).Count(&total).Error

	if err != nil {
		return suppliers, total, err
	}
	

	return suppliers, total, nil
}

func (r *repository) ToggleStatus(id int64, isActive bool, updatedBy int64) error {
	return r.db.Model(&models.Supplier{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  isActive,
			"updated_by": updatedBy,
		}).Error
}
