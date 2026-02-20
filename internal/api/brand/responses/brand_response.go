package responses

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/brand/dtos"
)

type BrandResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LogoPath  string    `json:"logo_path"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func BrandResponseFromDTO(dto dtos.BrandDTO) BrandResponse {
	return BrandResponse{
		ID:        dto.ID,
		Name:      dto.GetName(),
		LogoPath:  dto.LogoPath,
		IconPath:  dto.IconPath,
		CreatedAt: dto.CreatedAt,
	}
}
