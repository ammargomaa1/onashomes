package sections

import (
	"errors"

	"github.com/onas/ecommerce-api/internal/api/sections/requests"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	db   *gorm.DB
}

func NewService(repo *Repository, db *gorm.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) ListSections(pagination *utils.Pagination) utils.IResource {
	sections, err := s.repo.ListSections(s.db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch sections", err)
	}

	return utils.NewPaginatedOKResource("Sections Retrieved Successfully", sections, pagination.GetMeta())
}

func (s *Service) ListDeletedSections(pagination *utils.Pagination) utils.IResource {
	sections, err := s.repo.ListDeletedSections(s.db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch deleted sections", err)
	}

	return utils.NewPaginatedOKResource("Deleted Sections Retrieved Successfully", sections, pagination.GetMeta())
}

func (s *Service) GetSectionByID(id int64) utils.IResource {
	section, err := s.repo.GetSectionByID(s.db, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewNotFoundResource("section not found", nil)
	}

	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch section", err)
	}


	return utils.NewOKResource(section.NameEn+" Section Retrieved Successfully", section)
}

func (s *Service) CreateSection(req requests.CreateSectionRequest) utils.IResource {
	section := &models.Section{
		NameAr: req.NameAr,
		NameEn: req.NameEn,
		IconID: req.IconID,
	}

	if err := s.repo.Create(s.db, section); err != nil {
		return utils.NewInternalErrorResource("failed to create section", err)
	}

	return utils.NewOKResource(section.NameEn+" Section Created Successfully", section)
}

func (s *Service) UpdateSection(id int64, req requests.CreateSectionRequest) utils.IResource {
	section, err := s.repo.GetSectionByID(s.db, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewNotFoundResource("section not found", nil)
	}

	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch section", err)
	}

	updated := &models.Section{
		ID:     id,
		NameAr: req.NameAr,
		NameEn: req.NameEn,
		IconID: req.IconID,
	}

	if err := s.repo.Update(s.db, updated); err != nil {
		return utils.NewInternalErrorResource("failed to update section", err)
	}

	return utils.NewOKResource("section updated successfully", section)
}

func (s *Service) DeleteSection(id int64) utils.IResource {
	_, err := s.repo.GetSectionByID(s.db, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewNotFoundResource("section not found", nil)
	}

	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch section", err)
	}

	if err := s.repo.Delete(s.db, id); err != nil {
		return utils.NewInternalErrorResource("failed to delete section", err)
	}

	return utils.NewOKResource("section deleted successfully", nil)
}

func (s *Service) RestoreSection(id int64) utils.IResource {
	if err := s.repo.Restore(s.db, id); err != nil {
		return utils.NewInternalErrorResource("failed to restore section", err)
	}

	return utils.NewOKResource("section restored successfully", nil)
}

func (s *Service) ListSectionsForDropdown() utils.IResource {
	sections, err := s.repo.ListSectionsForDropdown(s.db)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch sections", err)
	}

	return utils.NewOKResource("sections fetched successfully", sections)
}
