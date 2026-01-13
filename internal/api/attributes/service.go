package attributes

import (
	"database/sql"
	"errors"
	"time"

	"github.com/onas/ecommerce-api/internal/api/attributes/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

// AttributeListItem is a lightweight representation used in list responses.
type AttributeListItem struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// AttributeValue represents a single attribute value.
type AttributeValue struct {
	ID          int64  `json:"id"`
	AttributeID int64  `json:"attribute_id"`
	Value       string `json:"value"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// AttributeDetail represents an attribute with its values.
type AttributeDetail struct {
	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	IsActive bool             `json:"is_active"`
	Values   []AttributeValue `json:"values"`
}

type Service struct {
	repo Repository
}

// NewService creates a new attribute service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAttribute(req requests.AttributeRequest) utils.IResource {
	if req.Name == "" {
		return utils.NewBadRequestResource("name is required", nil)
	}

	tx := s.repo.db.Begin()

	id, err := s.repo.CreateAttribute(tx, req)
	if err != nil {
		tx.Rollback()
		return utils.NewBadRequestResource("name has to be unique", nil)
	}

	detailIem, err := s.repo.GetAttributeByID(tx, id)
	if err != nil {
		tx.Rollback()
		return utils.NewInternalErrorResource("Failed to retrieve attribute", err.Error())
	}

	_, err = s.repo.CreateAttributeValuesBulk(tx, id, req.AttributeValues)
	if err != nil {
		tx.Rollback()
		return utils.NewInternalErrorResource("Failed to create attribute values", err.Error())
	}

	attributeValues, err := s.repo.GetAttributeValuesResponse(id)
	if err != nil {
		tx.Rollback()
		return utils.NewInternalErrorResource("Failed to retrieve attribute values", err.Error())
	}


	if err := tx.Commit().Error; err != nil {
		return utils.NewInternalErrorResource("Failed to proceed", err.Error())
	}

	return utils.NewCreatedResource("Attribute created successfully", requests.AttributeResponse{
		ID:       detailIem.ID,
		Name:     detailIem.Name,
		IsActive: detailIem.IsActive,
		Values:   attributeValues,
	})
}

func (s *Service) UpdateAttribute(id int64, req requests.AttributeRequest) utils.IResource {
	if req.Name == "" {
		return utils.NewBadRequestResource("name is required", nil)
	}

	if err := s.repo.UpdateAttribute(id, req); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	detail, err := s.repo.GetAttributeByID(s.repo.db, id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Attribute updated successfully", detail)
}

func (s *Service) DeleteAttribute(id int64) utils.IResource {
	if err := s.repo.SoftDeleteAttribute(id); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewNoContentResource()
}

func (s *Service) ListAttributes(pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListAttributes(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve attributes", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Attributes retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) GetAttributeByID(id int64) utils.IResource {
	detail, err := s.repo.GetAttributeByID(s.repo.db,id)
	if err != nil {
		return utils.NewNotFoundResource("Attribute not found", nil)
	}

	return utils.NewOKResource("Attribute retrieved successfully", detail)
}

func (s *Service) ListDeletedAttributes(pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListDeletedAttributes(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve deleted attributes", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Deleted attributes retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) RecoverAttribute(id int64) utils.IResource {
	if err := s.repo.RestoreAttribute(id); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	detail, err := s.repo.GetAttributeByID(s.repo.db, id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Attribute recovered successfully", detail)
}

func (s *Service) ListAttributeValues(attributeID int64, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListAttributeValues(attributeID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve attribute values", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Attribute values retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) CreateAttributeValue(attributeID int64, req requests.AttributeValueRequest) utils.IResource {
	if req.Value == "" {
		return utils.NewBadRequestResource("value is required", nil)
	}

	id, err := s.repo.CreateAttributeValue(attributeID, req)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	value := &AttributeValue{
		ID:          id,
		AttributeID: attributeID,
		Value:       req.Value,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	return utils.NewCreatedResource("Attribute value created successfully", value)
}

func (s *Service) UpdateAttributeValue(attributeID, valueID int64, req requests.AttributeValueRequest) utils.IResource {
	if req.Value == "" {
		return utils.NewBadRequestResource("value is required", nil)
	}

	if err := s.repo.UpdateAttributeValue(attributeID, valueID, req); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("some entities can not be found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	value := &AttributeValue{
		ID:          valueID,
		AttributeID: attributeID,
		Value:       req.Value,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}

	return utils.NewOKResource("Attribute value updated successfully", value)
}

func (s *Service) DeleteAttributeValue(attributeID, valueID int64) utils.IResource {
	if err := s.repo.SoftDeleteAttributeValue(attributeID, valueID); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewNoContentResource()
}

func (s *Service) ListDeletedAttributeValues(attributeID int64, pagination *utils.Pagination) utils.IResource {
	items, total, err := s.repo.ListDeletedAttributeValues(attributeID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve deleted attribute values", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Deleted attribute values retrieved successfully", items, pagination.GetMeta())
}

func (s *Service) RecoverAttributeValue(attributeID, valueID int64) utils.IResource {
	if err := s.repo.RestoreAttributeValue(attributeID, valueID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewNotFoundResource("some entities can not be found", nil)
		}
		return utils.NewBadRequestResource(err.Error(), nil)
	}
	// After restore, simply return a struct that reflects the latest state; for now we echo the input IDs.
	// Callers that need full detail can fetch via List or GetAttributeByID if needed.
	value := &AttributeValue{
		ID:          valueID,
		AttributeID: attributeID,
	}

	return utils.NewOKResource("Attribute value recovered successfully", value)
}
