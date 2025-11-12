package users

import (
	"errors"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service interface {
	Register(email, password, firstName, lastName string) (*models.User, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	GetProfile(id int64) (*models.User, error)
	UpdateProfile(id int64, firstName, lastName string) (*models.User, error)
	List(pagination *utils.Pagination) ([]models.User, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(email, password, firstName, lastName string) (*models.User, error) {
	// Check if user already exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid credentials")
		}
		return "", "", err
	}

	if !user.IsActive {
		return "", "", errors.New("account is inactive")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, utils.EntityUser, nil, utils.AccessToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(user.ID, utils.EntityUser, nil, utils.RefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken, utils.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(claims.EntityID, claims.EntityType, claims.RoleID, utils.AccessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *service) GetProfile(id int64) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateProfile(id int64, firstName, lastName string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) List(pagination *utils.Pagination) ([]models.User, int64, error) {
	return s.repo.List(pagination)
}
