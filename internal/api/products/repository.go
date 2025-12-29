package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, req AdminProductRequest) (int64, error)
	UpdateProduct(ctx context.Context, tx *sql.Tx, id int64, req AdminProductRequest) error
	SoftDeleteProduct(ctx context.Context, tx *sql.Tx, id int64) error
	ReplaceProductAttributes(ctx context.Context, tx *sql.Tx, productID int64, attrs []AdminProductAttributeRequest) error
	ListAdminProducts(ctx context.Context, pagination *utils.Pagination) ([]AdminProductListItem, int64, error)
	ListPublicProducts(ctx context.Context, pagination *utils.Pagination, q string) ([]AdminProductListItem, int64, error)
	GetAdminProductByID(ctx context.Context, id int64) (*AdminProductDetail, error)
	CreateVariant(ctx context.Context, tx *sql.Tx, productID int64, req AdminVariantRequest) (int64, error)
	UpdateVariant(ctx context.Context, tx *sql.Tx, productID, variantID int64, req AdminVariantRequest) error
	GetAdminVariantByID(ctx context.Context, productID, variantID int64) (*AdminVariant, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProduct(ctx context.Context, tx *sql.Tx, req AdminProductRequest) (int64, error) {
	row := tx.QueryRowContext(ctx,
		"INSERT INTO products (name, description, price, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Name, req.Description, req.Price, req.IsActive,
	)
	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repository) UpdateProduct(ctx context.Context, tx *sql.Tx, id int64, req AdminProductRequest) error {
	res, err := tx.ExecContext(ctx,
		"UPDATE products SET name = $1, description = $2, price = $3, is_active = $4, updated_at = NOW() WHERE id = $5 AND deleted_at IS NULL",
		req.Name, req.Description, req.Price, req.IsActive, id,
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

func (r *repository) SoftDeleteProduct(ctx context.Context, tx *sql.Tx, id int64) error {
	res, err := tx.ExecContext(ctx,
		"UPDATE products SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
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
	return nil
}

func (r *repository) ReplaceProductAttributes(ctx context.Context, tx *sql.Tx, productID int64, attrs []AdminProductAttributeRequest) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_attribute_values WHERE product_id = $1", productID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_attributes WHERE product_id = $1", productID); err != nil {
		return err
	}
	if len(attrs) == 0 {
		return nil
	}
	for _, a := range attrs {
		// basic existence check for attribute
		var exists bool
		if err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM attributes WHERE id = $1 AND deleted_at IS NULL)", a.AttributeID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("attribute %d does not exist", a.AttributeID)
		}
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO product_attributes (product_id, attribute_id, sort_order) VALUES ($1, $2, $3)",
			productID, a.AttributeID, a.SortOrder,
		); err != nil {
			return err
		}
		if len(a.AllowedValueIDs) == 0 {
			continue
		}
		// validate values belong to attribute and are active
		rows, err := tx.QueryContext(ctx,
			"SELECT id FROM attribute_values WHERE id = ANY($1::bigint[]) AND attribute_id = $2 AND deleted_at IS NULL",
			pqInt64Array(a.AllowedValueIDs), a.AttributeID,
		)
		if err != nil {
			return err
		}
		var okValues []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return err
			}
			okValues = append(okValues, id)
		}
		rows.Close()
		if len(okValues) != len(a.AllowedValueIDs) {
			return fmt.Errorf("one or more attribute values are invalid for attribute %d", a.AttributeID)
		}
		for _, v := range a.AllowedValueIDs {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO product_attribute_values (product_id, attribute_value_id) VALUES ($1, $2)",
				productID, v,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *repository) ListAdminProducts(ctx context.Context, pagination *utils.Pagination) ([]AdminProductListItem, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL").Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (pagination.Page - 1) * pagination.Limit
	rows, err := sqlDB.QueryContext(ctx,
		"SELECT id, name, price, is_active FROM products WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		pagination.Limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []AdminProductListItem
	for rows.Next() {
		var item AdminProductListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.IsActive); err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ListPublicProducts returns only active, non-deleted products for public APIs.
func (r *repository) ListPublicProducts(ctx context.Context, pagination *utils.Pagination, q string) ([]AdminProductListItem, int64, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, 0, err
	}

	// Count total active products matching optional search
	var total int64
	if q != "" {
		if err := sqlDB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND is_active = TRUE AND name ILIKE '%' || $1 || '%'",
			q,
		).Scan(&total); err != nil {
			return nil, 0, err
		}
	} else {
		if err := sqlDB.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM products WHERE deleted_at IS NULL AND is_active = TRUE",
		).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	offset := (pagination.Page - 1) * pagination.Limit
	var rows *sql.Rows
	if q != "" {
		rows, err = sqlDB.QueryContext(ctx,
			"SELECT id, name, price, is_active FROM products WHERE deleted_at IS NULL AND is_active = TRUE AND name ILIKE '%' || $1 || '%' ORDER BY created_at DESC LIMIT $2 OFFSET $3",
			q, pagination.Limit, offset,
		)
	} else {
		rows, err = sqlDB.QueryContext(ctx,
			"SELECT id, name, price, is_active FROM products WHERE deleted_at IS NULL AND is_active = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2",
			pagination.Limit, offset,
		)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []AdminProductListItem
	for rows.Next() {
		var item AdminProductListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.IsActive); err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) GetAdminProductByID(ctx context.Context, id int64) (*AdminProductDetail, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	var d AdminProductDetail
	var desc sql.NullString
	err = sqlDB.QueryRowContext(ctx,
		"SELECT id, name, description, price, is_active FROM products WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(&d.ID, &d.Name, &desc, &d.Price, &d.IsActive)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		d.Description = desc.String
	}
	// attributes and allowed values
	attrRows, err := sqlDB.QueryContext(ctx,
		`SELECT pa.attribute_id, a.name, pa.sort_order
		 FROM product_attributes pa
		 JOIN attributes a ON a.id = pa.attribute_id
		 WHERE pa.product_id = $1
		 ORDER BY pa.sort_order, a.name`, id,
	)
	if err != nil {
		return nil, err
	}
	defer attrRows.Close()
	attrMap := make(map[int64]*AdminProductAttribute)
	for attrRows.Next() {
		var attrID int64
		var name string
		var sortOrder int
		if err := attrRows.Scan(&attrID, &name, &sortOrder); err != nil {
			return nil, err
		}
		attrMap[attrID] = &AdminProductAttribute{
			AttributeID:   attrID,
			Name:          name,
			SortOrder:     sortOrder,
			AllowedValues: []AdminAttributeValue{},
		}
	}
	if err := attrRows.Err(); err != nil {
		return nil, err
	}
	valRows, err := sqlDB.QueryContext(ctx,
		`SELECT pav.attribute_value_id, av.attribute_id, av.value, av.sort_order
		 FROM product_attribute_values pav
		 JOIN attribute_values av ON av.id = pav.attribute_value_id
		 WHERE pav.product_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer valRows.Close()
	for valRows.Next() {
		var valueID, attrID int64
		var value string
		var sortOrder int
		if err := valRows.Scan(&valueID, &attrID, &value, &sortOrder); err != nil {
			return nil, err
		}
		attr := attrMap[attrID]
		if attr == nil {
			continue
		}
		attr.AllowedValues = append(attr.AllowedValues, AdminAttributeValue{
			ID:            valueID,
			AttributeID:   attrID,
			AttributeName: attr.Name,
			Value:         value,
			SortOrder:     sortOrder,
		})
	}
	if err := valRows.Err(); err != nil {
		return nil, err
	}
	// sort attributes
	for _, a := range attrMap {
		for i := range a.AllowedValues {
			a.AllowedValues[i].AttributeName = a.Name
		}
		d.Attributes = append(d.Attributes, *a)
	}
	sort.Slice(d.Attributes, func(i, j int) bool { return d.Attributes[i].SortOrder < d.Attributes[j].SortOrder })
	// variants
	vRows, err := sqlDB.QueryContext(ctx,
		"SELECT id, sku, is_active FROM product_variants WHERE product_id = $1 AND deleted_at IS NULL ORDER BY id",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer vRows.Close()
	for vRows.Next() {
		var v AdminVariant
		if err := vRows.Scan(&v.ID, &v.SKU, &v.IsActive); err != nil {
			return nil, err
		}
		// attribute values
		avRows, err := sqlDB.QueryContext(ctx,
			`SELECT av.id, av.attribute_id, a.name, av.value
			 FROM product_variant_attribute_values pvav
			 JOIN attribute_values av ON av.id = pvav.attribute_value_id
			 JOIN attributes a ON a.id = av.attribute_id
			 WHERE pvav.product_variant_id = $1`, v.ID,
		)
		if err != nil {
			return nil, err
		}
		for avRows.Next() {
			var valID, attrID int64
			var attrName, value string
			if err := avRows.Scan(&valID, &attrID, &attrName, &value); err != nil {
				avRows.Close()
				return nil, err
			}
			v.AttributeValues = append(v.AttributeValues, AdminVariantAttributeValue{
				AttributeID:   attrID,
				AttributeName: attrName,
				ValueID:       valID,
				Value:         value,
			})
		}
		avRows.Close()
		// images
		imgRows, err := sqlDB.QueryContext(ctx,
			"SELECT file_id, position FROM product_variant_images WHERE product_variant_id = $1 ORDER BY position",
			v.ID,
		)
		if err != nil {
			return nil, err
		}
		for imgRows.Next() {
			var img AdminVariantImage
			if err := imgRows.Scan(&img.FileID, &img.Position); err != nil {
				imgRows.Close()
				return nil, err
			}
			v.Images = append(v.Images, img)
		}
		imgRows.Close()
		// add-ons
		addRows, err := sqlDB.QueryContext(ctx,
			"SELECT pva.add_on_product_id, p.name FROM product_variant_add_ons pva JOIN products p ON p.id = pva.add_on_product_id WHERE pva.product_variant_id = $1 AND p.deleted_at IS NULL",
			v.ID,
		)
		if err != nil {
			return nil, err
		}
		for addRows.Next() {
			var ao AdminVariantAddOn
			if err := addRows.Scan(&ao.ProductID, &ao.Name); err != nil {
				addRows.Close()
				return nil, err
			}
			v.AddOns = append(v.AddOns, ao)
		}
		addRows.Close()
		d.Variants = append(d.Variants, v)
	}
	if err := vRows.Err(); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) CreateVariant(ctx context.Context, tx *sql.Tx, productID int64, req AdminVariantRequest) (int64, error) {
	if err := r.validateVariant(ctx, tx, productID, 0, req); err != nil {
		return 0, err
	}
	row := tx.QueryRowContext(ctx,
		"INSERT INTO product_variants (product_id, sku, is_active) VALUES ($1, $2, $3) RETURNING id",
		productID, req.SKU, req.IsActive,
	)
	var variantID int64
	if err := row.Scan(&variantID); err != nil {
		return 0, err
	}
	if err := r.replaceVariantDetails(ctx, tx, productID, variantID, req); err != nil {
		return 0, err
	}
	return variantID, nil
}

func (r *repository) UpdateVariant(ctx context.Context, tx *sql.Tx, productID, variantID int64, req AdminVariantRequest) error {
	if err := r.validateVariant(ctx, tx, productID, variantID, req); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx,
		"UPDATE product_variants SET sku = $1, is_active = $2, updated_at = NOW() WHERE id = $3 AND product_id = $4 AND deleted_at IS NULL",
		req.SKU, req.IsActive, variantID, productID,
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
	return r.replaceVariantDetails(ctx, tx, productID, variantID, req)
}

func (r *repository) replaceVariantDetails(ctx context.Context, tx *sql.Tx, productID, variantID int64, req AdminVariantRequest) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_variant_attribute_values WHERE product_variant_id = $1", variantID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_variant_images WHERE product_variant_id = $1", variantID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM product_variant_add_ons WHERE product_variant_id = $1", variantID); err != nil {
		return err
	}
	// attribute values
	for _, vID := range req.AttributeValueIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO product_variant_attribute_values (product_variant_id, attribute_value_id) VALUES ($1, $2)",
			variantID, vID,
		); err != nil {
			return err
		}
	}
	// images with order
	for idx, fileID := range req.ImageFileIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO product_variant_images (product_variant_id, file_id, position) VALUES ($1, $2, $3)",
			variantID, fileID, idx,
		); err != nil {
			return err
		}
	}
	// add-ons
	for _, addOnID := range req.AddOnProductIDs {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO product_variant_add_ons (product_variant_id, add_on_product_id) VALUES ($1, $2)",
			variantID, addOnID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) GetAdminVariantByID(ctx context.Context, productID, variantID int64) (*AdminVariant, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return nil, err
	}

	var v AdminVariant
	err = sqlDB.QueryRowContext(ctx,
		"SELECT id, sku, is_active FROM product_variants WHERE id = $1 AND product_id = $2 AND deleted_at IS NULL",
		variantID, productID,
	).Scan(&v.ID, &v.SKU, &v.IsActive)
	if err != nil {
		return nil, err
	}
	// attribute values
	avRows, err := sqlDB.QueryContext(ctx,
		`SELECT av.id, av.attribute_id, a.name, av.value
		 FROM product_variant_attribute_values pvav
		 JOIN attribute_values av ON av.id = pvav.attribute_value_id
		 JOIN attributes a ON a.id = av.attribute_id
		 WHERE pvav.product_variant_id = $1`, v.ID,
	)
	if err != nil {
		return nil, err
	}
	for avRows.Next() {
		var valID, attrID int64
		var attrName, value string
		if err := avRows.Scan(&valID, &attrID, &attrName, &value); err != nil {
			avRows.Close()
			return nil, err
		}
		v.AttributeValues = append(v.AttributeValues, AdminVariantAttributeValue{
			AttributeID:   attrID,
			AttributeName: attrName,
			ValueID:       valID,
			Value:         value,
		})
	}
	avRows.Close()
	// images
	imgRows, err := sqlDB.QueryContext(ctx,
		"SELECT file_id, position FROM product_variant_images WHERE product_variant_id = $1 ORDER BY position",
		v.ID,
	)
	if err != nil {
		return nil, err
	}
	for imgRows.Next() {
		var img AdminVariantImage
		if err := imgRows.Scan(&img.FileID, &img.Position); err != nil {
			imgRows.Close()
			return nil, err
		}
		v.Images = append(v.Images, img)
	}
	imgRows.Close()
	// add-ons
	addRows, err := sqlDB.QueryContext(ctx,
		"SELECT pva.add_on_product_id, p.name FROM product_variant_add_ons pva JOIN products p ON p.id = pva.add_on_product_id WHERE pva.product_variant_id = $1 AND p.deleted_at IS NULL",
		v.ID,
	)
	if err != nil {
		return nil, err
	}
	for addRows.Next() {
		var ao AdminVariantAddOn
		if err := addRows.Scan(&ao.ProductID, &ao.Name); err != nil {
			addRows.Close()
			return nil, err
		}
		v.AddOns = append(v.AddOns, ao)
	}
	addRows.Close()
	return &v, nil
}

func (r *repository) validateVariant(ctx context.Context, tx *sql.Tx, productID, excludeVariantID int64, req AdminVariantRequest) error {
	// ensure product exists
	var exists bool
	if err := tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM products WHERE id = $1 AND deleted_at IS NULL)", productID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	// load attribute values with attribute ids
	if len(req.AttributeValueIDs) == 0 {
		return errors.New("attribute_value_ids required")
	}
	rows, err := tx.QueryContext(ctx,
		"SELECT id, attribute_id FROM attribute_values WHERE id = ANY($1::bigint[]) AND deleted_at IS NULL",
		pqInt64Array(req.AttributeValueIDs),
	)
	if err != nil {
		return err
	}
	attrByValue := make(map[int64]int64)
	for rows.Next() {
		var id, attrID int64
		if err := rows.Scan(&id, &attrID); err != nil {
			rows.Close()
			return err
		}
		attrByValue[id] = attrID
	}
	rows.Close()
	if len(attrByValue) != len(req.AttributeValueIDs) {
		return fmt.Errorf("one or more attribute values are invalid")
	}
	// ensure values belong to product via product_attribute_values
	rows, err = tx.QueryContext(ctx,
		"SELECT attribute_value_id FROM product_attribute_values WHERE product_id = $1 AND attribute_value_id = ANY($2::bigint[])",
		productID, pqInt64Array(req.AttributeValueIDs),
	)
	if err != nil {
		return err
	}
	allowed := make(map[int64]struct{})
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		allowed[id] = struct{}{}
	}
	rows.Close()
	if len(allowed) != len(req.AttributeValueIDs) {
		return fmt.Errorf("one or more attribute values are not allowed for this product")
	}
	// ensure unique attribute per variant
	attrSeen := make(map[int64]struct{})
	for _, id := range req.AttributeValueIDs {
		attrID := attrByValue[id]
		if _, ok := attrSeen[attrID]; ok {
			return fmt.Errorf("duplicate attribute in variant definition")
		}
		attrSeen[attrID] = struct{}{}
	}
	// ensure unique combination across variants
	sortedIDs := append([]int64(nil), req.AttributeValueIDs...)
	sort.Slice(sortedIDs, func(i, j int) bool { return sortedIDs[i] < sortedIDs[j] })
	query := `SELECT pv.id, ARRAY_AGG(pvav.attribute_value_id ORDER BY pvav.attribute_value_id)
		FROM product_variants pv
		JOIN product_variant_attribute_values pvav ON pvav.product_variant_id = pv.id
		WHERE pv.product_id = $1 AND pv.deleted_at IS NULL`
	args := []interface{}{productID}
	if excludeVariantID > 0 {
		query += " AND pv.id <> $2"
		args = append(args, excludeVariantID)
	}
	query += " GROUP BY pv.id"
	rows, err = tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		var vid int64
		var vals []int64
		if err := rows.Scan(&vid, pqInt64ArrayScanner(&vals)); err != nil {
			rows.Close()
			return err
		}
		if len(vals) != len(sortedIDs) {
			continue
		}
		equal := true
		for i := range vals {
			if vals[i] != sortedIDs[i] {
				equal = false
				break
			}
		}
		if equal {
			rows.Close()
			return fmt.Errorf("variant with the same attribute combination already exists")
		}
	}
	rows.Close()
	// validate image files
	if len(req.ImageFileIDs) > 0 {
		rows, err = tx.QueryContext(ctx,
			"SELECT id FROM files WHERE id = ANY($1::bigint[]) AND is_active = TRUE",
			pqInt64Array(req.ImageFileIDs),
		)
		if err != nil {
			return err
		}
		count := 0
		for rows.Next() {
			count++
		}
		rows.Close()
		if count != len(req.ImageFileIDs) {
			return fmt.Errorf("one or more image_file_ids are invalid")
		}
	}
	// validate add-on products
	if len(req.AddOnProductIDs) > 0 {
		rows, err = tx.QueryContext(ctx,
			"SELECT id FROM products WHERE id = ANY($1::bigint[]) AND deleted_at IS NULL",
			pqInt64Array(req.AddOnProductIDs),
		)
		if err != nil {
			return err
		}
		count := 0
		for rows.Next() {
			count++
		}
		rows.Close()
		if count != len(req.AddOnProductIDs) {
			return fmt.Errorf("one or more add_on_product_ids are invalid")
		}
	}
	return nil
}

// pqInt64Array is a lightweight helper to pass []int64 to ANY($1)
func pqInt64Array(ids []int64) interface{} {
	if len(ids) == 0 {
		return "{}"
	}
	// Format as Postgres array literal: {1,2,3}
	b := make([]byte, 0, len(ids)*4)
	b = append(b, '{')
	for i, id := range ids {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, []byte(fmt.Sprintf("%d", id))...)
	}
	b = append(b, '}')
	return string(b)
}

// pqInt64ArrayScanner scans a Postgres int8[] into a []int64
func pqInt64ArrayScanner(dst *[]int64) interface{} {
	return &pqInt64ArrayScannerType{dst: dst}
}

type pqInt64ArrayScannerType struct {
	dst *[]int64
}

func (s *pqInt64ArrayScannerType) Scan(src interface{}) error {
	// src is expected to be []byte representing Postgres array literal, e.g. {1,2,3}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("unexpected type for int8[]: %T", src)
	}
	str := string(b)
	if len(str) < 2 {
		*s.dst = nil
		return nil
	}
	str = str[1 : len(str)-1]
	if str == "" {
		*s.dst = nil
		return nil
	}
	parts := strings.Split(str, ",")
	res := make([]int64, 0, len(parts))
	for _, p := range parts {
		var v int64
		_, err := fmt.Sscanf(p, "%d", &v)
		if err != nil {
			return err
		}
		res = append(res, v)
	}
	*s.dst = res
	return nil
}
