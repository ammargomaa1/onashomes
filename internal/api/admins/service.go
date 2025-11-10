package admins

import (
	"errors"

	"github.com/google/uuid"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Service interface {
	Create(email, password, firstName, lastName string, roleID uuid.UUID) (*models.Admin, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	GetByID(id uuid.UUID) (*models.Admin, error)
	Update(id uuid.UUID, firstName, lastName string, roleID uuid.UUID) (*models.Admin, error)
	Delete(id uuid.UUID) error
	List(pagination *utils.Pagination) ([]models.Admin, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(email, password, firstName, lastName string, roleID uuid.UUID) (*models.Admin, error) {
	// Check if admin already exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("admin with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	admin := &models.Admin{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		RoleID:    roleID,
		IsActive:  true,
	}

	err = s.repo.Create(admin)
	if err != nil {
		return nil, err
	}

	// Fetch the admin with role preloaded
	return s.repo.FindByID(admin.ID)
}

func (s *service) Login(email, password string) (string, string, error) {
	admin, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid credentials")
		}
		return "", "", err
	}

	if !admin.IsActive {
		return "", "", errors.New("account is inactive")
	}

	if !utils.CheckPasswordHash(password, admin.Password) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(admin.ID, utils.EntityAdmin, &admin.RoleID, utils.AccessToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(admin.ID, utils.EntityAdmin, &admin.RoleID, utils.RefreshToken)
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

func (s *service) GetByID(id uuid.UUID) (*models.Admin, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(id uuid.UUID, firstName, lastName string, roleID uuid.UUID) (*models.Admin, error) {
	admin, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	admin.FirstName = firstName
	admin.LastName = lastName
	admin.RoleID = roleID

	err = s.repo.Update(admin)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *service) List(pagination *utils.Pagination) ([]models.Admin, int64, error) {
	return s.repo.List(pagination)
}
