package users

import (
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
	List(pagination *utils.Pagination) ([]models.User, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) List(pagination *utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total records
	r.db.Model(&models.User{}).Count(&total)

	// Apply pagination
	err := pagination.Paginate(r.db).Find(&users).Error

	return users, total, err
}
