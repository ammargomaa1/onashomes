package customers

import (
	"errors"

	"regexp"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListCustomers retrieves a paginated list of customers
func (s *Service) ListCustomers(pagination *utils.Pagination, filter string, storeFrontID *int64) ([]models.Customer, map[string]interface{}, error) {
	customers, err := s.repo.List(pagination, filter, storeFrontID)
	if err != nil {
		return nil, nil, err
	}
	return customers, pagination.GetMeta(), nil
}

// SearchCustomers searches for customers by query
func (s *Service) SearchCustomers(query string, storeFrontID *int64) ([]models.Customer, error) {
	return s.repo.Search(query, storeFrontID, 10) // Default limit 10 for dropdown search
}

// CreateCustomer creates a new customer after validating uniqueness
// Checks if phone already exists (unless deleted)
func (s *Service) CreateCustomer(customer *models.Customer) error {
	// 1. Basic Validation
	if customer.FirstName == "" {
		return errors.New("first name is required")
	}
	if customer.Phone == "" {
		return errors.New("phone is required")
	}

	// Validate Phone Format (010, 011, 015 and 11 digits)
	matched, _ := regexp.MatchString(`^(010|011|015)\d{8}$`, customer.Phone)
	if !matched {
		return errors.New("phone must start with 010, 011, or 015 and be exactly 11 digits")
	}

	if customer.StoreFrontID <= 0 {
		return errors.New("store_front_id is required")
	}

	// 2. Check for duplicates by phone
	existing, err := s.repo.GetByPhone(customer.Phone)
	if err == nil && existing != nil {
		return errors.New("customer with this phone number already exists")
	}

	// 3. Create
	return s.repo.Create(customer)
}

// GetCustomer gets a customer by ID
func (s *Service) GetCustomer(id int64) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

// UpdateCustomer updates an existing customer
func (s *Service) UpdateCustomer(id int64, input *models.Customer) error {
	// 1. Get existing customer
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 2. Basic Validation
	if input.FirstName == "" {
		return errors.New("first name is required")
	}
	if input.Phone == "" {
		return errors.New("phone is required")
	}
	if input.StoreFrontID <= 0 {
		return errors.New("store_front_id is required")
	}

	// Validate Phone Format
	matched, _ := regexp.MatchString(`^(010|011|015)\d{8}$`, input.Phone)
	if !matched {
		return errors.New("phone must start with 010, 011, or 015 and be exactly 11 digits")
	}

	// 3. Check for duplicates if phone changed
	if input.Phone != existing.Phone {
		duplicate, err := s.repo.GetByPhone(input.Phone)
		if err == nil && duplicate != nil {
			return errors.New("customer with this phone number already exists")
		}
	}

	// 4. Update fields
	existing.FirstName = input.FirstName
	existing.LastName = input.LastName
	existing.Email = input.Email
	existing.Phone = input.Phone
	existing.StoreFrontID = input.StoreFrontID

	// 5. Save
	return s.repo.Update(existing)
}
