package storefronts

import (
	"github.com/onas/ecommerce-api/internal/api/storefronts/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req requests.CreateStoreFrontRequest) utils.IResource {
	// Validate slug uniqueness
	unique, err := s.repo.IsSlugUnique(req.Slug, 0)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate slug", err)
	}
	if !unique {
		return utils.NewBadRequestResource("Slug already in use", nil)
	}

	// Validate domain uniqueness
	domainUnique, err := s.repo.IsDomainUnique(req.Domain, 0)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate domain", err)
	}
	if !domainUnique {
		return utils.NewBadRequestResource("Domain already in use", nil)
	}

	sf := &models.StoreFront{
		Name:            req.Name,
		Slug:            req.Slug,
		Domain:          req.Domain,
		Currency:        req.Currency,
		DefaultLanguage: req.DefaultLanguage,
		IsActive:        true,
	}

	if err := s.repo.Create(sf); err != nil {
		return utils.NewInternalErrorResource("Failed to create store front", err)
	}

	return utils.NewCreatedResource("Store front created successfully", sf)
}

func (s *Service) Update(id int64, req requests.UpdateStoreFrontRequest) utils.IResource {
	sf, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Store front not found", nil)
	}

	// Validate domain uniqueness if changed
	domainUnique, err := s.repo.IsDomainUnique(req.Domain, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate domain", err)
	}
	if !domainUnique {
		return utils.NewBadRequestResource("Domain already in use", nil)
	}

	sf.Name = req.Name
	sf.Domain = req.Domain
	sf.Currency = req.Currency
	sf.DefaultLanguage = req.DefaultLanguage
	sf.IsActive = req.IsActive

	if err := s.repo.Update(sf); err != nil {
		return utils.NewInternalErrorResource("Failed to update store front", err)
	}

	return utils.NewOKResource("Store front updated successfully", sf)
}

func (s *Service) Delete(id int64) utils.IResource {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Store front not found", nil)
	}

	if err := s.repo.Delete(id); err != nil {
		return utils.NewInternalErrorResource("Failed to delete store front", err)
	}

	return utils.NewNoContentResource()
}

func (s *Service) GetByID(id int64) utils.IResource {
	sf, err := s.repo.GetByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Store front not found", nil)
	}

	return utils.NewOKResource("Store front retrieved successfully", sf)
}

func (s *Service) List(pagination *utils.Pagination) utils.IResource {
	list, total, err := s.repo.List(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve store fronts", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Store fronts retrieved successfully", list, pagination.GetMeta())
}

func (s *Service) ResolveByDomain(domain string) utils.IResource {
	sf, err := s.repo.GetByDomain(domain)
	if err != nil {
		return utils.NewNotFoundResource("Store not found for domain", nil)
	}

	return utils.NewOKResource("Store resolved successfully", sf)
}
