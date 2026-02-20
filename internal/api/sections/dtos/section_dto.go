package dtos

import (
	"time"

	"github.com/onas/ecommerce-api/internal/utils"
)

type SectionDTO struct {
	ID        int64     `json:"id"`
	NameAr    string    `json:"name_ar"`
	NameEn    string    `json:"name_en"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *SectionDTO) GetName() string {
	return utils.GetLocalizedStringFromContext(s.NameAr, s.NameEn)
}
