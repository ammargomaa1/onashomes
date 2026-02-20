package locations

import (
	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCountries() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Find(&countries).Error
	return countries, err
}

func (r *Repository) GetGovernorates(countryID int64) ([]models.Governorate, error) {
	var governorates []models.Governorate
	err := r.db.Where("country_id = ?", countryID).Find(&governorates).Error
	return governorates, err
}

func (r *Repository) GetCities(governorateID int64) ([]models.City, error) {
	var cities []models.City
	err := r.db.Where("governorate_id = ?", governorateID).Find(&cities).Error
	return cities, err
}
