package attributes

import (
	"github.com/onas/ecommerce-api/internal/api/attributes/requests"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// Repository defines persistence operations for attributes and attribute values.


type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new attribute Repository using the shared GORM DB
// while executing explicit SQL via the underlying *sql.DB.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAttribute(req requests.AttributeRequest) (int64, error) {

	var id int64

	err := r.db.Raw(
		"INSERT INTO attributes (name, is_active) VALUES ($1, $2) RETURNING id",
		req.Name, req.IsActive,
	).Scan(&id).Error


	return id, err
}

func (r *Repository) UpdateAttribute(id int64, req requests.AttributeRequest) error {

	err := r.db.Exec(
		"UPDATE attributes SET name = $1, is_active = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL",
		req.Name, req.IsActive, id,
	).Error

	return err
}

// RestoreAttributeValue clears the soft-delete flag for a specific attribute value.
func (r *Repository) RestoreAttributeValue(attributeID, valueID int64) error {

	err := r.db.Exec(
		"UPDATE attribute_values SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE id = $1 AND attribute_id = $2 AND deleted_at IS NOT NULL",
		valueID, attributeID,
	).Error
	
	return err
}

func (r *Repository) SoftDeleteAttribute(id int64) error {
	// Soft delete attribute
	err := r.db.Exec(
		"UPDATE attributes SET deleted_at = NOW(), is_active = FALSE WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Error

	if err != nil {
		return err
	}

	// Also soft delete its values
	err = r.db.Exec(
		"UPDATE attribute_values SET deleted_at = NOW(), is_active = FALSE WHERE attribute_id = $1 AND deleted_at IS NULL",
		id,
	).Error

	return nil
}

// RestoreAttribute clears soft-delete flags for an attribute and its values.
func (r *Repository) RestoreAttribute(id int64) error {
	// Restore attribute
	err := r.db.Exec(
		"UPDATE attributes SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL",
		id,
	).Error
	if err != nil {
		return err
	}


	// Restore its values as active
	err = r.db.Exec(
		"UPDATE attribute_values SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE attribute_id = $1 AND deleted_at IS NOT NULL",
		id,
	).Error

	return err
}

func (r *Repository) ListAttributes(pagination *utils.Pagination) ([]AttributeListItem, int64, error) {

	var total int64
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM attributes WHERE deleted_at IS NULL",
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	
	var list []AttributeListItem
	err := r.db.Raw(
		"SELECT id, name, is_active FROM attributes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT $1 OFFSET $2",
		pagination.Limit, offset,
	).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// ListDeletedAttributeValues returns soft-deleted values for a given attribute.
func (r *Repository) ListDeletedAttributeValues(attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {

	var total int64
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NOT NULL",
		attributeID,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit

	var list []AttributeValue

	err := r.db.Raw(
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NOT NULL ORDER BY deleted_at DESC LIMIT $2 OFFSET $3",
		attributeID, pagination.Limit, offset,
	).Scan(&list).Error


	return list, total, err
}

// ListDeletedAttributes returns attributes that have been soft-deleted.
func (r *Repository) ListDeletedAttributes(pagination *utils.Pagination) ([]AttributeListItem, int64, error) {

	var total int64
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM attributes WHERE deleted_at IS NOT NULL",
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	var list []AttributeListItem
	err := r.db.Raw(
		"SELECT id, name, is_active FROM attributes WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC LIMIT $1 OFFSET $2",
		pagination.Limit, offset,
	).Scan(&list).Error

	return list, total, err
}

func (r *Repository) GetAttributeByID(id int64) (*AttributeDetail, error) {

	var d AttributeDetail
	if err := r.db.Raw(
		"SELECT id, name, is_active FROM attributes WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(&d).Error; err != nil {
		return nil, err
	}

	var v AttributeValue
	err := r.db.Raw(
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL ORDER BY sort_order, value",
		id,
	).Scan(&v).Error

	return &d, err
}

func (r *Repository) ListAttributeValues(attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {

	var total int64
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL",
		attributeID,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit

	var list []AttributeValue

	err := r.db.Raw(
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL ORDER BY sort_order, value LIMIT $2 OFFSET $3",
		attributeID, pagination.Limit, offset,
	).Scan(&list).Error
	
	return list, total, err
}

func (r *Repository) CreateAttributeValue(attributeID int64, req requests.AttributeValueRequest) (int64, error) {

	// Ensure attribute exists and is not soft-deleted
	var exists bool
	if err := r.db.Raw(
		"SELECT EXISTS (SELECT 1 FROM attributes WHERE id = $1 AND deleted_at IS NULL)",
		attributeID,
	).Scan(&exists).Error; err != nil {
		return 0, err
	}

	var id int64
	err := r.db.Raw(
		"INSERT INTO attribute_values (attribute_id, value, sort_order, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		attributeID, req.Value, req.SortOrder, req.IsActive,
	).Scan(&id).Error

	return id, err
}

func (r *Repository) UpdateAttributeValue(attributeID, valueID int64, req requests.AttributeValueRequest) error {
	err := r.db.Raw(
		"UPDATE attribute_values SET value = $1, sort_order = $2, is_active = $3, updated_at = NOW() WHERE id = $4 AND attribute_id = $5 AND deleted_at IS NULL",
		req.Value, req.SortOrder, req.IsActive, valueID, attributeID,
	).Error


	return err
}

func (r *Repository) SoftDeleteAttributeValue(attributeID, valueID int64) error {
	err := r.db.Raw(
		"UPDATE attribute_values SET deleted_at = NOW(), is_active = FALSE WHERE id = $1 AND attribute_id = $2 AND deleted_at IS NULL",
		valueID, attributeID,
	).Error

	return err
}
