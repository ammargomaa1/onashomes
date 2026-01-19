package brand

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/brand/dtos"
	"github.com/onas/ecommerce-api/internal/api/brand/responses"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
	// Add repository fields and methods here
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ListBrands(db *gorm.DB, pagination *utils.Pagination) ([]responses.BrandResponse, error) {

	var list []dtos.BrandDTO
	query := db.Table("brands").
		Joins("LEFT JOIN files AS logo ON logo.id = brands.logo_id").
		Joins("LEFT JOIN files AS icon ON icon.id = brands.icon_id").
		Select(`
		brands.*,
		logo.file_path AS logo_path,
		icon.file_path AS icon_path
	`).
	Where(" brands.deleted_at IS NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	brandsList := []responses.BrandResponse{}
	for _, brand := range list {
		brandsList = append(brandsList, responses.BrandResponseFromDTO(brand))
	}

	return brandsList, nil
}

func (r *Repository) ListDeletedBrands(db *gorm.DB, pagination *utils.Pagination) ([]responses.BrandResponse, error) {

	var list []dtos.BrandDTO
	query := db.Table("brands").
		Joins("LEFT JOIN files AS logo ON logo.id = brands.logo_id").
		Joins("LEFT JOIN files AS icon ON icon.id = brands.icon_id").
		Select(`
		brands.*,
		logo.file_path AS logo_path,
		icon.file_path AS icon_path
	`).
	Where(" brands.deleted_at IS NOT NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	brandsList := []responses.BrandResponse{}
	for _, brand := range list {
		brandsList = append(brandsList, responses.BrandResponseFromDTO(brand))
	}

	return brandsList, nil
}

func (r *Repository) GetBrandByID(db *gorm.DB, id int64) (*dtos.BrandDTO, error) {
	var brand dtos.BrandDTO
	err := db.Table("brands").
		Joins("LEFT JOIN files AS logo ON logo.id = brands.logo_id").
		Joins("LEFT JOIN files AS icon ON icon.id = brands.icon_id").
		Select(`
		brands.*,
		logo.file_path AS logo_path,
		icon.file_path AS icon_path
	`).
		First(&brand, "brands.id = ? AND brands.deleted_at IS NULL", id).Error

	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *Repository) Create(db *gorm.DB, brand *models.Brand) error {
	return db.Create(brand).Error
}

func (r *Repository) Update(db *gorm.DB, brand *models.Brand) error {
	return db.Save(brand).Error
}

func (r *Repository) Delete(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE brands SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", time.Now(), id).Error
}

func (r *Repository) Restore(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE brands SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL", id).Error
}

func (r *Repository) IsBrandIdExists(db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.Model(&models.Brand{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


func (r *Repository) IsBrandDeletedIdExists(db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.Model(&models.Brand{}).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
