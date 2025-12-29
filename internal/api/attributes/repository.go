package attributes

import (
	"context"
	"database/sql"

	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// Repository defines persistence operations for attributes and attribute values.
type Repository interface {
	CreateAttribute(ctx context.Context, req AttributeRequest) (int64, error)
	UpdateAttribute(ctx context.Context, id int64, req AttributeRequest) error
	SoftDeleteAttribute(ctx context.Context, id int64) error
	ListAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error)
	GetAttributeByID(ctx context.Context, id int64) (*AttributeDetail, error)
	// Soft-deleted attributes
	ListDeletedAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error)
	RestoreAttribute(ctx context.Context, id int64) error

	ListAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error)
	CreateAttributeValue(ctx context.Context, attributeID int64, req AttributeValueRequest) (int64, error)
	UpdateAttributeValue(ctx context.Context, attributeID, valueID int64, req AttributeValueRequest) error
	SoftDeleteAttributeValue(ctx context.Context, attributeID, valueID int64) error
	// Soft-deleted attribute values
	ListDeletedAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error)
	RestoreAttributeValue(ctx context.Context, attributeID, valueID int64) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new attribute repository using the shared GORM DB
// while executing explicit SQL via the underlying *sql.DB.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) sqlDB() (*sql.DB, error) {
	return r.db.DB()
}

func (r *repository) CreateAttribute(ctx context.Context, req AttributeRequest) (int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return 0, err
	}
	row := sqlDB.QueryRowContext(ctx,
		"INSERT INTO attributes (name, is_active) VALUES ($1, $2) RETURNING id",
		req.Name, req.IsActive,
	)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) UpdateAttribute(ctx context.Context, id int64, req AttributeRequest) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attributes SET name = $1, is_active = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL",
		req.Name, req.IsActive, id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// RestoreAttributeValue clears the soft-delete flag for a specific attribute value.
func (r *repository) RestoreAttributeValue(ctx context.Context, attributeID, valueID int64) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attribute_values SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE id = $1 AND attribute_id = $2 AND deleted_at IS NOT NULL",
		valueID, attributeID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) SoftDeleteAttribute(ctx context.Context, id int64) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	// Soft delete attribute
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attributes SET deleted_at = NOW(), is_active = FALSE WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	// Also soft delete its values
	_, err = sqlDB.ExecContext(ctx,
		"UPDATE attribute_values SET deleted_at = NOW(), is_active = FALSE WHERE attribute_id = $1 AND deleted_at IS NULL",
		id,
	)
	return err
}

// RestoreAttribute clears soft-delete flags for an attribute and its values.
func (r *repository) RestoreAttribute(ctx context.Context, id int64) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	// Restore attribute
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attributes SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL",
		id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	// Restore its values as active
	_, err = sqlDB.ExecContext(ctx,
		"UPDATE attribute_values SET deleted_at = NULL, is_active = TRUE, updated_at = NOW() WHERE attribute_id = $1 AND deleted_at IS NOT NULL",
		id,
	)
	return err
}

func (r *repository) ListAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attributes WHERE deleted_at IS NULL",
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, name, is_active FROM attributes WHERE deleted_at IS NULL ORDER BY name LIMIT $1 OFFSET $2",
		pagination.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []AttributeListItem
	for rows.Next() {
		var item AttributeListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.IsActive); err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListDeletedAttributeValues returns soft-deleted values for a given attribute.
func (r *repository) ListDeletedAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NOT NULL",
		attributeID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NOT NULL ORDER BY deleted_at DESC LIMIT $2 OFFSET $3",
		attributeID, pagination.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []AttributeValue
	for rows.Next() {
		var v AttributeValue
		if err := rows.Scan(&v.ID, &v.Value, &v.SortOrder, &v.IsActive); err != nil {
			return nil, 0, err
		}
		v.AttributeID = attributeID
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListDeletedAttributes returns attributes that have been soft-deleted.
func (r *repository) ListDeletedAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attributes WHERE deleted_at IS NOT NULL",
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, name, is_active FROM attributes WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC LIMIT $1 OFFSET $2",
		pagination.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []AttributeListItem
	for rows.Next() {
		var item AttributeListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.IsActive); err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) GetAttributeByID(ctx context.Context, id int64) (*AttributeDetail, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return nil, err
	}

	var d AttributeDetail
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT id, name, is_active FROM attributes WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(&d.ID, &d.Name, &d.IsActive); err != nil {
		return nil, err
	}

	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL ORDER BY sort_order, value",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v AttributeValue
		if err := rows.Scan(&v.ID, &v.Value, &v.SortOrder, &v.IsActive); err != nil {
			return nil, err
		}
		v.AttributeID = id
		d.Values = append(d.Values, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) ListAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL",
		attributeID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, value, sort_order, is_active FROM attribute_values WHERE attribute_id = $1 AND deleted_at IS NULL ORDER BY sort_order, value LIMIT $2 OFFSET $3",
		attributeID, pagination.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []AttributeValue
	for rows.Next() {
		var v AttributeValue
		if err := rows.Scan(&v.ID, &v.Value, &v.SortOrder, &v.IsActive); err != nil {
			return nil, 0, err
		}
		v.AttributeID = attributeID
		list = append(list, v)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) CreateAttributeValue(ctx context.Context, attributeID int64, req AttributeValueRequest) (int64, error) {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return 0, err
	}

	// Ensure attribute exists and is not soft-deleted
	var exists bool
	if err := sqlDB.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM attributes WHERE id = $1 AND deleted_at IS NULL)",
		attributeID,
	).Scan(&exists); err != nil {
		return 0, err
	}
	if !exists {
		return 0, sql.ErrNoRows
	}

	row := sqlDB.QueryRowContext(ctx,
		"INSERT INTO attribute_values (attribute_id, value, sort_order, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		attributeID, req.Value, req.SortOrder, req.IsActive,
	)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) UpdateAttributeValue(ctx context.Context, attributeID, valueID int64, req AttributeValueRequest) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attribute_values SET value = $1, sort_order = $2, is_active = $3, updated_at = NOW() WHERE id = $4 AND attribute_id = $5 AND deleted_at IS NULL",
		req.Value, req.SortOrder, req.IsActive, valueID, attributeID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) SoftDeleteAttributeValue(ctx context.Context, attributeID, valueID int64) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	res, err := sqlDB.ExecContext(ctx,
		"UPDATE attribute_values SET deleted_at = NOW(), is_active = FALSE WHERE id = $1 AND attribute_id = $2 AND deleted_at IS NULL",
		valueID, attributeID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
