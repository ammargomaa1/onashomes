package suppliers

import (
	"time"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(companyName, contactPersonName, contactNumber string, createdBy int64) utils.IResource {
	supplier := &models.Supplier{
		CompanyName:       companyName,
		ContactPersonName: contactPersonName,
		ContactNumber:     contactNumber,
		IsActive:          true,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
	}

	if err := s.repo.Create(supplier); err != nil {
		return utils.NewInternalErrorResource("Failed to create supplier", err)
	}

	createdSupplier, err := s.repo.GetByID(supplier.ID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve supplier", err)
	}

	return utils.NewCreatedResource("Supplier created successfully", createdSupplier)
}

func (s *Service) Update(id int64, companyName, contactPersonName, contactNumber string, updatedBy int64) utils.IResource {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewNotFoundResource("supplier not found", nil)
	}

	supplier.CompanyName = companyName
	supplier.ContactPersonName = contactPersonName
	supplier.ContactNumber = contactNumber
	updatedAt := time.Now()
	supplier.UpdatedAt = updatedAt

	if err := s.repo.Update(supplier); err != nil {
		return utils.NewInternalErrorResource("Failed to update supplier", err)
	}

	updatedSupplier, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve supplier", err)
	}

	return utils.NewOKResource("Supplier updated successfully", updatedSupplier)
}

func (s *Service) GetByID(id int64) utils.IResource {
	supplier, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewNotFoundResource("supplier not found", nil)
	}

	return utils.NewOKResource("Supplier retrieved successfully", supplier)
}

func (s *Service) List(pagination *utils.Pagination) utils.IResource {
	suppliers, total, err := s.repo.List(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve suppliers", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Suppliers retrieved successfully", suppliers, pagination.GetMeta())
}

func (s *Service) Activate(id int64) utils.IResource {
	if err := s.repo.ToggleStatus(id, true); err != nil {
		return utils.NewInternalErrorResource("Failed to activate supplier", err)
	}

	return utils.NewNoContentResource()
}

func (s *Service) Deactivate(id int64) utils.IResource {
	if err := s.repo.ToggleStatus(id, false); err != nil {
		return utils.NewInternalErrorResource("Failed to deactivate supplier", err)
	}

	return utils.NewNoContentResource()
}
