package responses

import (
	"time"

	"github.com/onas/ecommerce-api/internal/api/categories/dtos"
)

type CategoryResponse struct {
	ID        int64     `json:"id"`
	SectionID int64     `json:"section_id"`
	Name      string    `json:"name"`
	IconPath  string    `json:"icon_path"`
	CreatedAt time.Time `json:"created_at"`
}

func CategoryResponseFromDTO(dto dtos.CategoryDTO) CategoryResponse {
	return CategoryResponse{
		ID:        dto.ID,
		SectionID: dto.SectionID,
		Name:      dto.GetName(),
		IconPath:  dto.IconPath,
		CreatedAt: dto.CreatedAt,
	}
}
