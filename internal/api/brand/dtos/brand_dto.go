package dtos

import (
	"time"

	"github.com/onas/ecommerce-api/internal/utils"
)

type BrandDTO struct {
	ID        int64     `json:"id"`
	NameAr    string    `json:"name_ar"`
	NameEn    string    `json:"name_en"`
	LogoPath  string    `json:"logo_path"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func (b *BrandDTO) GetName() string {
	return utils.GetLocalizedStringFromContext(b.NameEn, b.NameAr)
}
