package requests

type CreateStoreFrontRequest struct {
	Name            string `json:"name" binding:"required"`
	Slug            string `json:"slug" binding:"required"`
	Domain          string `json:"domain" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	DefaultLanguage string `json:"default_language" binding:"required"`
}

type UpdateStoreFrontRequest struct {
	Name            string `json:"name" binding:"required"`
	Domain          string `json:"domain" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	DefaultLanguage string `json:"default_language" binding:"required"`
	IsActive        bool   `json:"is_active"`
}
