package categories

import (
	"github.com/onas/ecommerce-api/internal/api/categories/requests"
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

func (s *Service) ListCategories(pagination *utils.Pagination) utils.IResource {
	categories, err := s.repo.ListCategories(s.db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch categories", err)
	}

	return utils.NewPaginatedOKResource("Categories Retrived Successfully", categories, pagination.GetMeta())
}

func (s *Service) ListCategoriesBySection(sectionID int64, pagination *utils.Pagination) utils.IResource {
	categories, err := s.repo.ListCategoriesBySection(s.db, sectionID, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch categories", err)
	}

	return utils.NewPaginatedOKResource("Categories Retrieved Successfully", categories, pagination.GetMeta())
}

func (s *Service) ListDeletedCategories(pagination *utils.Pagination) utils.IResource {
	categories, err := s.repo.ListDeletedCategories(s.db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch deleted categories", err)
	}

	return utils.NewPaginatedOKResource("Categories Retrieved Successfully", categories, pagination.GetMeta())
}

func (s *Service) GetCategoryByID(id int64) utils.IResource {
	category, err := s.repo.GetCategoryByID(s.db, id)
	if err != nil {
		return utils.NewInternalErrorResource("category not found", err)
	}

	return utils.NewOKResource("Category retrieved successfully", category)
}

func (s *Service) CreateCategory(req requests.CreateCategoryRequest) utils.IResource {
	category := &models.Category{
		SectionID: req.SectionID,
		NameAr:    req.NameAr,
		NameEn:    req.NameEn,
		IconID:    req.IconID,
	}

	if err := s.repo.Create(s.db, category); err != nil {
		return utils.NewInternalErrorResource("failed to create category", err)
	}

	return utils.NewOKResource("Category retrieved successfully", category)
}

func (s *Service) UpdateCategory(id int64, req requests.CreateCategoryRequest) utils.IResource {
	category, err := s.repo.GetCategoryByID(s.db, id)
	if err != nil {
		return utils.NewInternalErrorResource("category not found", err)
	}

	updated := &models.Category{
		ID:        id,
		SectionID: req.SectionID,
		NameAr:    req.NameAr,
		NameEn:    req.NameEn,
		IconID:    req.IconID,
	}

	if err := s.repo.Update(s.db, updated); err != nil {
		return utils.NewInternalErrorResource("failed to update category", err)
	}

	return utils.NewOKResource("category updated successfully", category)
}

func (s *Service) DeleteCategory(id int64) utils.IResource {
	_, err := s.repo.GetCategoryByID(s.db, id)
	if err != nil {
		return utils.NewInternalErrorResource("category not found", err)
	}

	if err := s.repo.Delete(s.db, id); err != nil {
		return utils.NewInternalErrorResource("failed to delete category", err)
	}

	return utils.NewOKResource("category deleted successfully", nil)
}

func (s *Service) RestoreCategory(id int64) utils.IResource {
	if err := s.repo.Restore(s.db, id); err != nil {
		return utils.NewInternalErrorResource("failed to restore category", err)
	}

	return utils.NewOKResource("category restored successfully", nil)
}

func (s *Service) ListCategoriesForDropdown() utils.IResource {
	items, err := s.repo.ListCategoriesForDropdown(s.db)
	if err != nil {
		return utils.NewInternalErrorResource("failed to fetch categories", err)
	}
	return utils.NewOKResource("categories fetched successfully", items)
}
