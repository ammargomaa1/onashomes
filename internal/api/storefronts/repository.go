package storefronts

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

func (r *Repository) Create(storeFront *models.StoreFront) error {
	return r.db.Create(storeFront).Error
}

func (r *Repository) Update(storeFront *models.StoreFront) error {
	return r.db.Save(storeFront).Error
}

func (r *Repository) Delete(id int64) error {
	return r.db.Delete(&models.StoreFront{}, id).Error
}

func (r *Repository) GetByID(id int64) (*models.StoreFront, error) {
	var sf models.StoreFront
	if err := r.db.First(&sf, id).Error; err != nil {
		return nil, err
	}
	return &sf, nil
}

func (r *Repository) GetByDomain(domain string) (*models.StoreFront, error) {
	var sf models.StoreFront
	if err := r.db.Where("domain = ? AND is_active = true", domain).First(&sf).Error; err != nil {
		return nil, err
	}
	return &sf, nil
}

func (r *Repository) GetBySlug(slug string) (*models.StoreFront, error) {
	var sf models.StoreFront
	if err := r.db.Where("slug = ?", slug).First(&sf).Error; err != nil {
		return nil, err
	}
	return &sf, nil
}

func (r *Repository) List(pagination *utils.Pagination) ([]models.StoreFront, int64, error) {
	var list []models.StoreFront
	query := r.db.Model(&models.StoreFront{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(pagination.Limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *Repository) IsSlugUnique(slug string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&models.StoreFront{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *Repository) IsDomainUnique(domain string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.Model(&models.StoreFront{}).Where("domain = ?", domain)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
