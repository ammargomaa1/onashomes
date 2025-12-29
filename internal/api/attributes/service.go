package attributes

import (
	"context"
	"fmt"

	"github.com/onas/ecommerce-api/internal/utils"
)

// AttributeRequest represents the payload to create or update an attribute.
type AttributeRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// AttributeListItem is a lightweight representation used in list responses.
type AttributeListItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// AttributeValueRequest represents the payload to create or update an attribute value.
type AttributeValueRequest struct {
	Value     string `json:"value" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
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

// Service defines the business logic for managing attributes and their values.
type Service interface {
	CreateAttribute(ctx context.Context, req AttributeRequest) (*AttributeDetail, error)
	UpdateAttribute(ctx context.Context, id int64, req AttributeRequest) (*AttributeDetail, error)
	DeleteAttribute(ctx context.Context, id int64) error
	ListAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error)
	GetAttributeByID(ctx context.Context, id int64) (*AttributeDetail, error)
	ListDeletedAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error)
	RecoverAttribute(ctx context.Context, id int64) (*AttributeDetail, error)

	ListAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error)
	CreateAttributeValue(ctx context.Context, attributeID int64, req AttributeValueRequest) (*AttributeValue, error)
	UpdateAttributeValue(ctx context.Context, attributeID, valueID int64, req AttributeValueRequest) (*AttributeValue, error)
	DeleteAttributeValue(ctx context.Context, attributeID, valueID int64) error
	ListDeletedAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error)
	RecoverAttributeValue(ctx context.Context, attributeID, valueID int64) (*AttributeValue, error)
}

type service struct {
	repo Repository
}

// NewService creates a new attribute service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAttribute(ctx context.Context, req AttributeRequest) (*AttributeDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	id, err := s.repo.CreateAttribute(ctx, req)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAttributeByID(ctx, id)
}

func (s *service) UpdateAttribute(ctx context.Context, id int64, req AttributeRequest) (*AttributeDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if err := s.repo.UpdateAttribute(ctx, id, req); err != nil {
		return nil, err
	}
	return s.repo.GetAttributeByID(ctx, id)
}

func (s *service) DeleteAttribute(ctx context.Context, id int64) error {
	return s.repo.SoftDeleteAttribute(ctx, id)
}

func (s *service) ListAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	return s.repo.ListAttributes(ctx, pagination)
}

func (s *service) GetAttributeByID(ctx context.Context, id int64) (*AttributeDetail, error) {
	return s.repo.GetAttributeByID(ctx, id)
}

func (s *service) ListDeletedAttributes(ctx context.Context, pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	return s.repo.ListDeletedAttributes(ctx, pagination)
}

func (s *service) RecoverAttribute(ctx context.Context, id int64) (*AttributeDetail, error) {
	if err := s.repo.RestoreAttribute(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.GetAttributeByID(ctx, id)
}

func (s *service) ListAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	return s.repo.ListAttributeValues(ctx, attributeID, pagination)
}

func (s *service) CreateAttributeValue(ctx context.Context, attributeID int64, req AttributeValueRequest) (*AttributeValue, error) {
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}

	id, err := s.repo.CreateAttributeValue(ctx, attributeID, req)
	if err != nil {
		return nil, err
	}

	return &AttributeValue{
		ID:          id,
		AttributeID: attributeID,
		Value:       req.Value,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}, nil
}

func (s *service) UpdateAttributeValue(ctx context.Context, attributeID, valueID int64, req AttributeValueRequest) (*AttributeValue, error) {
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}

	if err := s.repo.UpdateAttributeValue(ctx, attributeID, valueID, req); err != nil {
		return nil, err
	}

	return &AttributeValue{
		ID:          valueID,
		AttributeID: attributeID,
		Value:       req.Value,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}, nil
}

func (s *service) DeleteAttributeValue(ctx context.Context, attributeID, valueID int64) error {
	return s.repo.SoftDeleteAttributeValue(ctx, attributeID, valueID)
}

func (s *service) ListDeletedAttributeValues(ctx context.Context, attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	return s.repo.ListDeletedAttributeValues(ctx, attributeID, pagination)
}

func (s *service) RecoverAttributeValue(ctx context.Context, attributeID, valueID int64) (*AttributeValue, error) {
	if err := s.repo.RestoreAttributeValue(ctx, attributeID, valueID); err != nil {
		return nil, err
	}
	// After restore, simply return a struct that reflects the latest state; for now we echo the input IDs.
	// Callers that need full detail can fetch via List or GetAttributeByID if needed.
	return &AttributeValue{
		ID:          valueID,
		AttributeID: attributeID,
	}, nil
}
