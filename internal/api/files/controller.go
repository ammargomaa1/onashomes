package files

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/services"
)

// Controller handles file-related HTTP requests
type Controller struct {
	fileService *services.FileService
}

// NewController creates a new file controller
func NewController(fileService *services.FileService) *Controller {
	return &Controller{
		fileService: fileService,
	}
}

// UploadResponse represents file upload response
type UploadResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	File    *FileResponse `json:"file,omitempty"`
}

// FileResponse represents file data in API responses
type FileResponse struct {
	ID              int64      `json:"id"`
	OriginalName    string     `json:"original_name"`
	FileName        string     `json:"file_name"`
	FilePath        string     `json:"file_path"`
	FileSize        int64      `json:"file_size"`
	MimeType        string     `json:"mime_type"`
	FileType        string     `json:"file_type"`
	Extension       string     `json:"extension"`
	UploadedBy      int64      `json:"uploaded_by"`
	UploadedByAdmin *AdminInfo `json:"uploaded_by_admin,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

// AdminInfo represents basic admin info in responses
type AdminInfo struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ListResponse represents file list response
type ListResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Files   []FileResponse `json:"files"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
}

// Upload handles file upload
func (c *Controller) Upload(ctx *gin.Context) {
	// Get admin ID from context (set by auth middleware)
	adminID, exists := ctx.Get("entity_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, UploadResponse{
			Success: false,
			Message: "Admin authentication required",
		})
		return
	}

	adminIDInt, ok := adminID.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Invalid admin ID format",
		})
		return
	}

	// Get uploaded file
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "No file uploaded or invalid file",
		})
		return
	}

	// Get custom config if provided
	config := c.fileService.GetDefaultConfig()

	// Check for custom max file size in form data
	if maxSizeStr := ctx.PostForm("max_file_size"); maxSizeStr != "" {
		if maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil && maxSize > 0 {
			config.MaxFileSize = maxSize
		}
	}

	// Upload file
	file, err := c.fileService.UploadFile(fileHeader, adminIDInt, config)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	fileResponse := c.convertToFileResponse(file)

	ctx.JSON(http.StatusCreated, UploadResponse{
		Success: true,
		Message: "File uploaded successfully",
		File:    fileResponse,
	})
}

// GetFile retrieves a file by ID
func (c *Controller) GetFile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file ID"})
		return
	}

	file, err := c.fileService.GetFile(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		return
	}

	fileResponse := c.convertToFileResponse(file)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "File retrieved successfully", "data": fileResponse})
}

// ListFiles returns paginated list of files
func (c *Controller) ListFiles(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	fileType := ctx.Query("type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	files, total, err := c.fileService.ListFiles(page, limit, fileType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ListResponse{
			Success: false,
			Message: "Failed to retrieve files",
		})
		return
	}

	// Convert to response format
	fileResponses := make([]FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = *c.convertToFileResponse(&file)
	}

	ctx.JSON(http.StatusOK, ListResponse{
		Success: true,
		Message: "Files retrieved successfully",
		Files:   fileResponses,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// DeleteFile deletes a file
func (c *Controller) DeleteFile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid file ID"})
		return
	}

	// Get admin ID from context
	adminID, exists := ctx.Get("admin_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Admin authentication required"})
		return
	}

	adminIDInt, ok := adminID.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid admin ID format"})
		return
	}

	if err := c.fileService.DeleteFile(id, adminIDInt); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found or already deleted"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "File deleted successfully"})
}

// GetUploadConfig returns the current upload configuration
func (c *Controller) GetUploadConfig(ctx *gin.Context) {
	config := c.fileService.GetDefaultConfig()

	response := map[string]interface{}{
		"max_file_size":    config.MaxFileSize,
		"max_file_size_mb": float64(config.MaxFileSize) / (1024 * 1024),
		"allowed_types":    config.AllowedTypes,
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Upload configuration retrieved", "data": response})
}

// convertToFileResponse converts model to response format
func (c *Controller) convertToFileResponse(file *models.File) *FileResponse {
	response := &FileResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		MimeType:     file.MimeType,
		FileType:     file.FileType,
		Extension:    file.Extension,
		UploadedBy:   file.UploadedBy,
		IsActive:     file.IsActive,
		CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    file.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if file.UploadedByAdmin != nil {
		response.UploadedByAdmin = &AdminInfo{
			ID:        file.UploadedByAdmin.ID,
			Email:     file.UploadedByAdmin.Email,
			FirstName: file.UploadedByAdmin.FirstName,
			LastName:  file.UploadedByAdmin.LastName,
		}
	}

	return response
}
