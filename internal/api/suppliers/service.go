package suppliers

import (
	"errors"
	"time"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Service interface {
	Create(companyName, contactPersonName, contactNumber string, createdBy int64) (*models.Supplier, error)
	Update(id int64, companyName, contactPersonName, contactNumber string, updatedBy int64) (*models.Supplier, error)
	GetByID(id int64) (*models.Supplier, error)
	List(pagination *utils.Pagination) ([]models.Supplier, int64, error)
	Activate(id int64, updatedBy int64) error
	Deactivate(id int64, updatedBy int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(companyName, contactPersonName, contactNumber string, createdBy int64) (*models.Supplier, error) {
	supplier := &models.Supplier{
		CompanyName:       companyName,
		ContactPersonName: contactPersonName,
		ContactNumber:     contactNumber,
		IsActive:         true,
		CreatedBy:        createdBy,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.Create(supplier); err != nil {
		return nil, err
	}

	// Fetch the created supplier with related data
	return s.repo.GetByID(supplier.ID)
}

func (s *service) Update(id int64, companyName, contactPersonName, contactNumber string, updatedBy int64) (*models.Supplier, error) {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("supplier not found")
	}

	supplier.CompanyName = companyName
	supplier.ContactPersonName = contactPersonName
	supplier.ContactNumber = contactNumber
	updatedAt := time.Now()
	supplier.UpdatedAt = updatedAt

	if err := s.repo.Update(supplier); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *service) GetByID(id int64) (*models.Supplier, error) {
	return s.repo.GetByID(id)
}

func (s *service) List(pagination *utils.Pagination) ([]models.Supplier, int64, error) {
	return s.repo.List(pagination)
}

func (s *service) Activate(id int64, updatedBy int64) error {
	return s.repo.ToggleStatus(id, true, updatedBy)
}

func (s *service) Deactivate(id int64, updatedBy int64) error {
	return s.repo.ToggleStatus(id, false, updatedBy)
}
