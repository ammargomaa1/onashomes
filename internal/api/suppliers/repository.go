package suppliers

import (
	"gorm.io/gorm"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(supplier *models.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *Repository) Update(supplier *models.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *Repository) GetByID(id int64) (*models.Supplier, error) {
	var supplier models.Supplier
	err := r.db.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *Repository) List(pagination *utils.Pagination) ([]models.Supplier, int64, error) {
	var suppliers []models.Supplier

	err := pagination.Paginate(r.db, &models.Supplier{}).
		Scan(&suppliers).Error
	if err != nil {
		return nil, pagination.Total, err
	}

	return suppliers, pagination.Total, nil
}

func (r *Repository) ToggleStatus(id int64, isActive bool) error {
	return r.db.Model(&models.Supplier{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": isActive,
		}).Error
}

func (r *Repository) ListForDropdown() ([]utils.DropdownItem, error) {
	var items []utils.DropdownItem
	err := r.db.Table("suppliers").
		Select("id, company_name AS name_en, company_name AS name_ar").
		Where("deleted_at IS NULL").
		Order("company_name").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
