package admins

import (
	"errors"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service interface {
	Create(email, password, firstName, lastName string, roleID int64) utils.IResource
	Login(email, password string) utils.IResource
	RefreshToken(refreshToken string) utils.IResource
	GetByID(id int64) utils.IResource
	Update(id int64, firstName, lastName string, roleID int64) utils.IResource
	Delete(id int64) utils.IResource
	List(pagination *utils.Pagination) utils.IResource
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(email, password, firstName, lastName string, roleID int64) utils.IResource {
	// Check if admin already exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return utils.NewBadRequestResource("admin with this email already exists", nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	admin := &models.Admin{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		RoleID:    roleID,
		IsActive:  true,
	}

	if err := s.repo.Create(admin); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	createdAdmin, err := s.repo.FindByID(admin.ID)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewCreatedResource("Admin created successfully", createdAdmin)
}

func (s *service) Login(email, password string) utils.IResource {
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewUnauthorizedResource("invalid credentials", nil)
		}
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	if !admin.IsActive {
		return utils.NewUnauthorizedResource("account is inactive", nil)
	}

	if !utils.CheckPasswordHash(password, admin.Password) {
		return utils.NewUnauthorizedResource("invalid credentials", nil)
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(admin.ID, utils.EntityAdmin, &admin.RoleID, utils.AccessToken)
	if err != nil {
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	refreshToken, err := utils.GenerateToken(admin.ID, utils.EntityAdmin, &admin.RoleID, utils.RefreshToken)
	if err != nil {
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	data := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return utils.NewOKResource("Login successful", data)
}

func (s *service) RefreshToken(refreshToken string) utils.IResource {
	claims, err := utils.ValidateToken(refreshToken, utils.RefreshToken)
	if err != nil {
		return utils.NewUnauthorizedResource("invalid refresh token", nil)
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(claims.EntityID, claims.EntityType, claims.RoleID, utils.AccessToken)
	if err != nil {
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	data := map[string]string{
		"access_token": accessToken,
	}

	return utils.NewOKResource("Token refreshed successfully", data)
}

func (s *service) GetByID(id int64) utils.IResource {
	admin, err := s.repo.FindByID(id)
	if err != nil {
		return utils.NewNotFoundResource("Admin not found", nil)
	}

	return utils.NewOKResource("Admin retrieved successfully", admin)
}

func (s *service) Update(id int64, firstName, lastName string, roleID int64) utils.IResource {
	admin, err := s.repo.FindByID(id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	admin.FirstName = firstName
	admin.LastName = lastName
	admin.RoleID = roleID

	if err := s.repo.Update(admin); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	updatedAdmin, err := s.repo.FindByID(id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Admin updated successfully", updatedAdmin)
}

func (s *service) Delete(id int64) utils.IResource {
	if err := s.repo.Delete(id); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Admin deleted successfully", nil)
}

func (s *service) List(pagination *utils.Pagination) utils.IResource {
	admins, total, err := s.repo.List(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve admins", err)
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Admins retrieved successfully", admins, pagination.GetMeta())
}
