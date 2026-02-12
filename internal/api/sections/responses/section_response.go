package responses

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/sections/dtos"
)

type SectionResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func SectionResponseFromDTO(dto dtos.SectionDTO) SectionResponse {
	return SectionResponse{
		ID:        dto.ID,
		Name:      dto.GetName(),
		IconPath:  dto.IconPath,
		CreatedAt: dto.CreatedAt,
	}
}
