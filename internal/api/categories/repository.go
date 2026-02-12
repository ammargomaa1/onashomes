package categories

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/categories/dtos"
	"github.com/onas/ecommerce-api/internal/api/categories/responses"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ListCategories(db *gorm.DB, pagination *utils.Pagination) ([]responses.CategoryResponse, error) {

	var list []dtos.CategoryDTO
	query := db.Table("categories").
		Joins("LEFT JOIN files AS icon ON icon.id = categories.icon_id").
		Select(`
		categories.*,
		icon.file_path AS icon_path
	`).
		Where("categories.deleted_at IS NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	categoriesList := []responses.CategoryResponse{}
	for _, category := range list {
		categoriesList = append(categoriesList, responses.CategoryResponseFromDTO(category))
	}

	return categoriesList, nil
}

func (r *Repository) ListCategoriesBySection(db *gorm.DB, sectionID int64, pagination *utils.Pagination) ([]responses.CategoryResponse, error) {

	var list []dtos.CategoryDTO
	query := db.Table("categories").
		Joins("LEFT JOIN files AS icon ON icon.id = categories.icon_id").
		Select(`
		categories.*,
		icon.file_path AS icon_path
	`).
		Where("categories.section_id = ? AND categories.deleted_at IS NULL", sectionID)

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	categoriesList := []responses.CategoryResponse{}
	for _, category := range list {
		categoriesList = append(categoriesList, responses.CategoryResponseFromDTO(category))
	}

	return categoriesList, nil
}

func (r *Repository) ListDeletedCategories(db *gorm.DB, pagination *utils.Pagination) ([]responses.CategoryResponse, error) {

	var list []dtos.CategoryDTO
	query := db.Table("categories").
		Joins("LEFT JOIN files AS icon ON icon.id = categories.icon_id").
		Select(`
		categories.*,
		icon.file_path AS icon_path
	`).
		Where("categories.deleted_at IS NOT NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	categoriesList := []responses.CategoryResponse{}
	for _, category := range list {
		categoriesList = append(categoriesList, responses.CategoryResponseFromDTO(category))
	}

	return categoriesList, nil
}

func (r *Repository) GetCategoryByID(db *gorm.DB, id int64) (*dtos.CategoryDTO, error) {
	var category dtos.CategoryDTO
	err := db.Table("categories").
		Joins("LEFT JOIN files AS icon ON icon.id = categories.icon_id").
		Select(`
		categories.*,
		icon.file_path AS icon_path
	`).
		First(&category, "categories.id = ? AND categories.deleted_at IS NULL", id).Error

	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) Create(db *gorm.DB, category *models.Category) error {
	return db.Create(category).Error
}

func (r *Repository) Update(db *gorm.DB, category *models.Category) error {
	return db.Save(category).Error
}

func (r *Repository) Delete(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE categories SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", time.Now(), id).Error
}

func (r *Repository) Restore(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE categories SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL", id).Error
}

func (r *Repository) IsCategoryIDExists(db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.Model(&models.Category{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
