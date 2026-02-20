package files

import (
	"github.com/onas/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

// Repository handles file-related database operations
type Repository struct {
}

// NewRepository creates a new instance of Repository
func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) FileExists(db *gorm.DB, fileID int64) (bool, error) {
	var file models.File
	err := db.Where("id = ?", fileID).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
