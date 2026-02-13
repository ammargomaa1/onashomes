package brand

import (
	"github.com/onas/ecommerce-api/internal/api/brand/requests"
	"github.com/onas/ecommerce-api/internal/api/files"
	"github.com/onas/ecommerce-api/internal/database"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Add service methods here
func (s *Service) ListBrands(pagination *utils.Pagination) utils.IResource {
	db := database.GetDB() // Assuming you have a way to get the DB instance
	brands, err := s.repo.ListBrands(db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve brands", err)
	}

	return utils.NewPaginatedOKResource("Brands retrieved successfully", brands, pagination.GetMeta())
}

func (s *Service) GetBrandByID(id int64) utils.IResource {
	db := database.GetDB()
	brand, err := s.repo.GetBrandByID(db, id)
	if err != nil || brand == nil {
		return utils.NewNotFoundResource("Brand not found", nil)
	}

	return utils.NewOKResource("Brand retrieved successfully", brand)
}

func (s *Service) CreateBrand(request requests.CreateBrandRequest) utils.IResource {
	// validate logo and icon existence
	db := database.GetDB()

	fileRepo := files.NewRepository()

	logoExists, err := fileRepo.FileExists(db, request.LogoID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate logo", err)
	}
	if !logoExists {
		return utils.NewBadRequestResource("Logo file does not exist", nil)
	}

	iconExists, err := fileRepo.FileExists(db, request.IconID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate icon", err)
	}

	if !iconExists {
		return utils.NewBadRequestResource("Icon file does not exist", nil)
	}
	brand := &models.Brand{
		NameAr: request.NameAr,
		NameEn: request.NameEn,
		LogoID: request.LogoID,
		IconID: request.IconID,
	}

	err = s.repo.Create(db, brand)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to create brand", err)
	}

	return utils.NewOKResource("Brand created successfully", nil)
}

func (s *Service) UpdateBrand(id int64, request requests.CreateBrandRequest) utils.IResource {
	// Implementation for updating a brand
	db := database.GetDB()

	// is brand exists
	exists, err := s.repo.IsBrandIdExists(db, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to check brand existence", err)
	}
	if !exists {
		return utils.NewNotFoundResource("Brand not found", nil)
	}

	brand := &models.Brand{
		ID:     id,
		NameAr: request.NameAr,
		NameEn: request.NameEn,
		LogoID: request.LogoID,
		IconID: request.IconID,
	}
	//validate logo and icon existence
	fileRepo := files.NewRepository()

	logoExists, err := fileRepo.FileExists(db, request.LogoID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate logo", err)
	}
	if !logoExists {
		return utils.NewBadRequestResource("Logo file does not exist", nil)
	}

	iconExists, err := fileRepo.FileExists(db, request.IconID)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to validate icon", err)
	}

	if !iconExists {
		return utils.NewBadRequestResource("Icon file does not exist", nil)
	}

	err = s.repo.Update(db, brand)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to update brand", err)
	}

	return utils.NewOKResource("Brand updated successfully", nil)
}

func (s *Service) DeleteBrand(id int64) utils.IResource {
	// Implementation for deleting a brand
	db := database.GetDB()

	// is brand exists
	exists, err := s.repo.IsBrandIdExists(db, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to check brand existence", err)
	}
	if !exists {
		return utils.NewNotFoundResource("Brand not found", nil)
	}

	err = s.repo.Delete(db, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to delete brand", err)
	}

	return utils.NewNoContentResource()
}

func (s *Service) RestoreBrand(id int64) utils.IResource {
	// Implementation for restoring a brand
	db := database.GetDB()

	// is brand exists
	exists, err := s.repo.IsBrandDeletedIdExists(db, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to check brand existence", err)
	}
	if !exists {
		return utils.NewNotFoundResource("Brand not found", nil)
	}

	err = s.repo.Restore(db, id)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to restore brand", err)
	}

	return utils.NewOKResource("Brand restored successfully", nil)
}

func (s *Service) ListDeletedBrands(pagination *utils.Pagination) utils.IResource {
	db := database.GetDB()
	brands, err := s.repo.ListDeletedBrands(db, pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve deleted brands", err)
	}

	return utils.NewPaginatedOKResource("Deleted brands retrieved successfully", brands, pagination.GetMeta())
}

func (s *Service) ListBrandsForDropdown() utils.IResource {
	db := database.GetDB()
	items, err := s.repo.ListBrandsForDropdown(db)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve brands", err)
	}
	return utils.NewOKResource("Brands retrieved successfully", items)
}
