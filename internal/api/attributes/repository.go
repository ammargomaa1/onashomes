package attributes

import (
	"encoding/json"
	"time"

	"github.com/onas/ecommerce-api/internal/api/attributes/requests"
	"github.com/onas/ecommerce-api/internal/models"
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

func (r *Repository) CreateAttribute(db *gorm.DB, req requests.AttributeRequest) (int64, error) {

	var id int64

	err := db.Raw(
		"INSERT INTO attributes (name_ar, name_en) VALUES ($1, $2) RETURNING id",
		req.NameAr, req.NameEn,
	).Scan(&id).Error


	return id, err
}

func (r *Repository) UpdateAttribute(id int64, req requests.AttributeRequest) error {

	err := r.db.Exec(
		"UPDATE attributes SET name_ar = $1, name_en = $2, updated_at = $3 WHERE id = $4 AND deleted_at IS NULL",
		req.NameAr, req.NameEn, time.Now(), id,
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
		"UPDATE attributes SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
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
	
	var list []models.Attribute
	err := pagination.Paginate(r.db, models.Attribute{}).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	var listItems []AttributeListItem
	for _, attr := range list {
		listItems = append(listItems, AttributeListItem{
			ID:       attr.ID,
			Name:     utils.GetLocalizedStringFromContext(attr.NameAr, attr.NameEn),
			CreatedAt: attr.CreatedAt,
		})
	}

	return listItems, pagination.Total, nil
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

func (r *Repository) GetAttributeByID(db *gorm.DB, id int64) (*models.Attribute, error) {

	var d models.Attribute
	if err := db.Raw(
		"SELECT id, name_en, name_ar FROM attributes WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(&d).Error; err != nil {
		return nil, err
	}

	return &d, nil
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
		"INSERT INTO attribute_values (attribute_id, value_ar, value_en, sort_order, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		attributeID, req.ValueAr, req.ValueEn, req.SortOrder, req.IsActive,
	).Scan(&id).Error

	return id, err
}

func (r *Repository) CreateAttributeValuesBulk(
	db *gorm.DB,
	attributeID int64,
	reqs []requests.AttributeValueRequest,
) ([]int64, error) {

	if len(reqs) == 0 {
		return []int64{}, nil
	}

	// 1. Ensure attribute exists and is not soft-deleted
	var exists bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM attributes
			WHERE id = $1
			  AND deleted_at IS NULL
		)
	`, attributeID).Scan(&exists).Error; err != nil {
		return nil, err
	}

	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	// 2. Marshal requests into JSON
	payload, err := json.Marshal(reqs)
	if err != nil {
		return nil, err
	}

	// 3. Bulk insert using jsonb_to_recordset
	var ids []int64
	err = db.Raw(`
		INSERT INTO attribute_values (
			attribute_id,
			value_ar,
			value_en,
			sort_order,
			is_active
		)
		SELECT
			$1 AS attribute_id,
			v.value_ar,
			v.value_en,
			v.sort_order,
			v.is_active
		FROM jsonb_to_recordset($2::jsonb) AS v(
			value_ar TEXT,
			value_en TEXT,
			sort_order INT,
			is_active BOOLEAN
		)
		RETURNING id
	`, attributeID, payload).Scan(&ids).Error

	return ids, err
}

func (r *Repository) GetAttributeValuesResponse(
	attributeID int64,
) ([]requests.AttributeValuesResponse, error) {

	var values []requests.AttributeValuesResponse

	err := r.db.Raw(`
		SELECT
			id,
			value_ar,
			value_en,
			sort_order,
			is_active
		FROM attribute_values
		WHERE attribute_id = $1
		  AND deleted_at IS NULL
		ORDER BY sort_order
	`, attributeID).Scan(&values).Error

	return values, err
}


func (r *Repository) UpdateAttributeValue(attributeID, valueID int64, req requests.AttributeValueRequest) error {
	err := r.db.Raw(
		"UPDATE attribute_values SET value_ar = $1, value_en = $2, sort_order = $3, is_active = $4, updated_at = $5 WHERE id = $6 AND attribute_id = $7 AND deleted_at IS NULL",
		req.ValueAr, req.ValueEn, req.SortOrder, req.IsActive, time.Now(), valueID, attributeID,
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
