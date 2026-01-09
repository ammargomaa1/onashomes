package attributes

import (
	"fmt"

	"github.com/onas/ecommerce-api/internal/api/attributes/requests"
	"github.com/onas/ecommerce-api/internal/utils"
)

// AttributeListItem is a lightweight representation used in list responses.
type AttributeListItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
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

func (s *Service) CreateAttribute(req requests.AttributeRequest) (*AttributeDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	id, err := s.repo.CreateAttribute(req)
	if err != nil {
		return nil, err
	}

	//create attribute values
	s.repo.CreateAttribute(req)
	return s.repo.GetAttributeByID(id)
}

func (s *Service) UpdateAttribute(id int64, req requests.AttributeRequest) (*AttributeDetail, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if err := s.repo.UpdateAttribute(id, req); err != nil {
		return nil, err
	}
	return s.repo.GetAttributeByID(id)
}

func (s *Service) DeleteAttribute(id int64) error {
	return s.repo.SoftDeleteAttribute(id)
}

func (s *Service) ListAttributes(pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	return s.repo.ListAttributes(pagination)
}

func (s *Service) GetAttributeByID(id int64) (*AttributeDetail, error) {
	return s.repo.GetAttributeByID(id)
}

func (s *Service) ListDeletedAttributes(pagination *utils.Pagination) ([]AttributeListItem, int64, error) {
	return s.repo.ListDeletedAttributes(pagination)
}

func (s *Service) RecoverAttribute(id int64) (*AttributeDetail, error) {
	if err := s.repo.RestoreAttribute(id); err != nil {
		return nil, err
	}
	return s.repo.GetAttributeByID(id)
}

func (s *Service) ListAttributeValues(attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	return s.repo.ListAttributeValues(attributeID, pagination)
}

func (s *Service) CreateAttributeValue(attributeID int64, req requests.AttributeValueRequest) (*AttributeValue, error) {
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}

	id, err := s.repo.CreateAttributeValue(attributeID, req)
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

func (s *Service) UpdateAttributeValue(attributeID, valueID int64, req requests.AttributeValueRequest) (*AttributeValue, error) {
	if req.Value == "" {
		return nil, fmt.Errorf("value is required")
	}

	if err := s.repo.UpdateAttributeValue(attributeID, valueID, req); err != nil {
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

func (s *Service) DeleteAttributeValue(attributeID, valueID int64) error {
	return s.repo.SoftDeleteAttributeValue(attributeID, valueID)
}

func (s *Service) ListDeletedAttributeValues(attributeID int64, pagination *utils.Pagination) ([]AttributeValue, int64, error) {
	return s.repo.ListDeletedAttributeValues(attributeID, pagination)
}

func (s *Service) RecoverAttributeValue(attributeID, valueID int64) (*AttributeValue, error) {
	if err := s.repo.RestoreAttributeValue(attributeID, valueID); err != nil {
		return nil, err
	}
	// After restore, simply return a struct that reflects the latest state; for now we echo the input IDs.
	// Callers that need full detail can fetch via List or GetAttributeByID if needed.
	return &AttributeValue{
		ID:          valueID,
		AttributeID: attributeID,
	}, nil
}
