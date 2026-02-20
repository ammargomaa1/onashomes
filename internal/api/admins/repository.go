package admins

import (
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	Create(admin *models.Admin) error
	FindByEmail(email string) (*models.Admin, error)
	FindByID(id int64) (*models.Admin, error)
	Update(admin *models.Admin) error
	Delete(id int64) error
	List(pagination *utils.Pagination) ([]models.Admin, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *repository) FindByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Preload("Role.Permissions").Where("email = ?", email).First(&admin).Error
	return &admin, err
}

func (r *repository) FindByID(id int64) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Preload("Role.Permissions").Where("id = ?", id).First(&admin).Error
	return &admin, err
}

func (r *repository) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Delete(&models.Admin{}, id).Error
}

func (r *repository) List(pagination *utils.Pagination) ([]models.Admin, int64, error) {
	var admins []models.Admin
	var total int64

	// Count total records
	r.db.Model(&models.Admin{}).Count(&total)

	// Apply pagination
	err := pagination.Paginate(r.db, admins).Preload("Role").Find(&admins).Error

	return admins, total, err
}
