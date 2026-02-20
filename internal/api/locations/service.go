package locations

import (
	"github.com/onas/ecommerce-api/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCountries() ([]models.Country, error) {
	return s.repo.GetCountries()
}

func (s *Service) GetGovernorates(countryID int64) ([]models.Governorate, error) {
	return s.repo.GetGovernorates(countryID)
}

func (s *Service) GetCities(governorateID int64) ([]models.City, error) {
	return s.repo.GetCities(governorateID)
}
