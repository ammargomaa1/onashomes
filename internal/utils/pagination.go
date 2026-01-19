package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/config"
	"gorm.io/gorm"
)

type Pagination struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	Sort       string `json:"sort,omitempty"`
	Order      string `json:"order,omitempty"`
}

// ParsePaginationParams extracts pagination parameters from request
func ParsePaginationParams(c *gin.Context) *Pagination {
	cfg := config.AppConfig

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(cfg.Pagination.DefaultPageSize)))
	if limit < 1 {
		limit = cfg.Pagination.DefaultPageSize
	}
	if limit > cfg.Pagination.MaxPageSize {
		limit = cfg.Pagination.MaxPageSize
	}

	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	// Validate order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return &Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

func (p *Pagination) Paginate(db *gorm.DB, model interface{}) *gorm.DB {
	var total int64

	// Count the total records BEFORE applying offset/limit
	// This works for both simple models and complex queries with JOINs
	countQuery := db
	if model != nil {
		countQuery = db.Model(model)
	}

	if err := countQuery.Count(&total).Error; err == nil {
		p.SetTotal(total)
	}

	offset := (p.Page - 1) * p.Limit

	// Apply offset and limit to the original query
	query := db.Offset(offset).Limit(p.Limit)

	// Apply sorting
	if p.Sort != "" {
		orderClause := p.Sort
		if p.Order != "" {
			orderClause += " " + p.Order
		}
		query = query.Order(orderClause)
	}

	return query
}

// SetTotal sets the total count and calculates total pages
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
}

// GetMeta returns pagination metadata
func (p *Pagination) GetMeta() map[string]interface{} {
	return map[string]interface{}{
		"page":        p.Page,
		"limit":       p.Limit,
		"total":       p.Total,
		"total_pages": p.TotalPages,
	}
}
