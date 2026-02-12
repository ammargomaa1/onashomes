package dtos

import (
	"time"

	"github.com/onas/ecommerce-api/internal/utils"
)

type CategoryDTO struct {
	ID        int64     `json:"id"`
	SectionID int64     `json:"section_id"`
	NameAr    string    `json:"name_ar"`
	NameEn    string    `json:"name_en"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *CategoryDTO) GetName() string {
	return utils.GetLocalizedStringFromContext(c.NameAr, c.NameEn)
}