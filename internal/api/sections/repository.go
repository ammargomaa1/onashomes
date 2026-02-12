package sections

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/sections/dtos"
	"github.com/onas/ecommerce-api/internal/api/sections/responses"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ListSections(db *gorm.DB, pagination *utils.Pagination) ([]responses.SectionResponse, error) {

	var list []dtos.SectionDTO
	query := db.Table("sections").
		Joins("LEFT JOIN files AS icon ON icon.id = sections.icon_id").
		Select(`
		sections.*,
		icon.file_path AS icon_path
	`).
		Where("sections.deleted_at IS NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	sectionsList := []responses.SectionResponse{}
	for _, section := range list {
		sectionsList = append(sectionsList, responses.SectionResponseFromDTO(section))
	}

	return sectionsList, nil
}

func (r *Repository) ListDeletedSections(db *gorm.DB, pagination *utils.Pagination) ([]responses.SectionResponse, error) {

	var list []dtos.SectionDTO
	query := db.Table("sections").
		Joins("LEFT JOIN files AS icon ON icon.id = sections.icon_id").
		Select(`
		sections.*,
		icon.file_path AS icon_path
	`).
		Where("sections.deleted_at IS NOT NULL")

	err := pagination.Paginate(query, nil).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	sectionsList := []responses.SectionResponse{}
	for _, section := range list {
		sectionsList = append(sectionsList, responses.SectionResponseFromDTO(section))
	}

	return sectionsList, nil
}

func (r *Repository) GetSectionByID(db *gorm.DB, id int64) (*dtos.SectionDTO, error) {
	var section dtos.SectionDTO
	err := db.Table("sections").
		Joins("LEFT JOIN files AS icon ON icon.id = sections.icon_id").
		Select(`
		sections.*,
		icon.file_path AS icon_path
	`).
		First(&section, "sections.id = ? AND sections.deleted_at IS NULL", id).Error

	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *Repository) Create(db *gorm.DB, section *models.Section) error {
	return db.Create(section).Error
}

func (r *Repository) Update(db *gorm.DB, section *models.Section) error {
	return db.Save(section).Error
}

func (r *Repository) Delete(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE sections SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", time.Now(), id).Error
}

func (r *Repository) Restore(db *gorm.DB, id int64) error {
	return db.Exec("UPDATE sections SET deleted_at = NULL WHERE id = ? AND deleted_at IS NOT NULL", id).Error
}

func (r *Repository) IsSectionIDExists(db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.Model(&models.Section{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListSectionsForDropdown(db *gorm.DB) ([]responses.SectionDropdownResponse, error) {

	var list []dtos.SectionDTO
	err := db.Table("sections").
		Select("id, name_ar, name_en").
		Where("deleted_at IS NULL").
		Order("name_en").
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	dropdownList := []responses.SectionDropdownResponse{}
	for _, section := range list {
		dropdownList = append(dropdownList, responses.SectionDropdownResponse{
			ID:   section.ID,
			Name: section.GetName(),
		})
	}

	return dropdownList, nil
}
