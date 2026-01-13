package users

import (
	"errors"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service interface {
	Register(email, password, firstName, lastName string) utils.IResource
	Login(email, password string) utils.IResource
	RefreshToken(refreshToken string) utils.IResource
	GetProfile(id int64) utils.IResource
	UpdateProfile(id int64, firstName, lastName string) utils.IResource
	List(pagination *utils.Pagination) utils.IResource
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(email, password, firstName, lastName string) utils.IResource {
	// Check if user already exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return utils.NewBadRequestResource("user with this email already exists", nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	user := &models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	if err := s.repo.Create(user); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewCreatedResource("User registered successfully", user)
}

func (s *service) Login(email, password string) utils.IResource {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewUnauthorizedResource("invalid credentials", nil)
		}
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	if !user.IsActive {
		return utils.NewUnauthorizedResource("account is inactive", nil)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return utils.NewUnauthorizedResource("invalid credentials", nil)
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, utils.EntityUser, nil, utils.AccessToken)
	if err != nil {
		return utils.NewUnauthorizedResource(err.Error(), nil)
	}

	refreshToken, err := utils.GenerateToken(user.ID, utils.EntityUser, nil, utils.RefreshToken)
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

func (s *service) GetProfile(id int64) utils.IResource {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return utils.NewNotFoundResource("User not found", nil)
	}

	return utils.NewOKResource("Profile retrieved successfully", user)
}

func (s *service) UpdateProfile(id int64, firstName, lastName string) utils.IResource {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	user.FirstName = firstName
	user.LastName = lastName

	if err := s.repo.Update(user); err != nil {
		return utils.NewBadRequestResource(err.Error(), nil)
	}

	return utils.NewOKResource("Profile updated successfully", user)
}

func (s *service) List(pagination *utils.Pagination) utils.IResource {
	users, total, err := s.repo.List(pagination)
	if err != nil {
		return utils.NewInternalErrorResource("Failed to retrieve users", err.Error())
	}

	pagination.SetTotal(total)
	return utils.NewPaginatedOKResource("Users retrieved successfully", users, pagination.GetMeta())
}
