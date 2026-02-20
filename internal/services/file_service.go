package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/onas/ecommerce-api/internal/models"
	"github.com/onas/ecommerce-api/internal/utils"
	"gorm.io/gorm"
)

// FileService handles file upload operations
type FileService struct {
	db       *gorm.DB
	fileUtil *utils.FileUtil
}

// FileUploadConfig holds configuration for file uploads
type FileUploadConfig struct {
	MaxFileSize  int64 // in bytes
	AllowedTypes map[string][]string
}

// NewFileService creates a new file service
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db:       db,
		fileUtil: utils.NewFileUtil(),
	}
}

// GetDefaultConfig returns default file upload configuration
func (s *FileService) GetDefaultConfig() *FileUploadConfig {
	return &FileUploadConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedTypes: map[string][]string{
			"document": {
				"application/pdf",
				"application/msword",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"text/plain",
				"application/rtf",
			},
			"spreadsheet": {
				"application/vnd.ms-excel",
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"text/csv",
			},
			"image": {
				"image/jpeg",
				"image/jpg",
				"image/png",
				"image/gif",
				"image/bmp",
				"image/webp",
				"image/svg+xml",
				"image/tiff",
				"image/ico",
			},
		},
	}
}

// UploadFile handles file upload with validation
func (s *FileService) UploadFile(fileHeader *multipart.FileHeader, uploadedBy int64, config *FileUploadConfig) (*models.File, error) {
	if config == nil {
		config = s.GetDefaultConfig()
	}

	// Validate file size
	if fileHeader.Size > config.MaxFileSize {
		return nil, fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", fileHeader.Size, config.MaxFileSize)
	}

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Detect MIME type
	mimeType, err := s.detectMimeType(src)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file type: %v", err)
	}

	// Validate file type
	fileType, allowed := s.validateFileType(mimeType, config.AllowedTypes)
	if !allowed {
		return nil, fmt.Errorf("file type %s is not allowed", mimeType)
	}

	// Generate unique filename
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if extension == "" {
		extension = s.getExtensionFromMimeType(mimeType)
	}

	uniqueFileName := s.generateUniqueFileName(fileHeader.Filename, extension)

	// Determine storage path based on file type
	storagePath := fmt.Sprintf("%s/%s", fileType, uniqueFileName)

	// Reset file reader to beginning (seek to start instead of reopening)
	if _, err := src.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file reader: %v", err)
	}

	// Save file to storage
	relativePath, err := s.fileUtil.SaveFileFromReader(storagePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Create file record
	fileRecord := &models.File{
		OriginalName: fileHeader.Filename,
		FileName:     uniqueFileName,
		FilePath:     relativePath,
		FileSize:     fileHeader.Size,
		MimeType:     mimeType,
		FileType:     fileType,
		Extension:    extension,
		UploadedBy:   uploadedBy,
		IsActive:     true,
	}

	// Save to database
	if err := s.db.Create(fileRecord).Error; err != nil {
		// Clean up file if database save fails
		s.fileUtil.DeleteFile(relativePath)
		return nil, fmt.Errorf("failed to save file record: %v", err)
	}

	return fileRecord, nil
}

// GetFile retrieves a file record by ID
func (s *FileService) GetFile(id int64) (*models.File, error) {
	var file models.File
	if err := s.db.Preload("UploadedByAdmin").First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

//GetFileByPath retrieves an actual file from storage by its path
func (s *FileService) GetFileByPath(filePath string) ([]byte, error) {
	data, err := s.fileUtil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ListFiles returns paginated list of files
func (s *FileService) ListFiles(page, limit int, fileType string) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := s.db.Model(&models.File{}).Where("is_active = ?", true)

	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Preload("UploadedByAdmin").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// DeleteFile soft deletes a file and removes it from storage
func (s *FileService) DeleteFile(id int64, deletedBy int64) error {
	var file models.File
	if err := s.db.First(&file, id).Error; err != nil {
		return err
	}

	// Soft delete from database
	if err := s.db.Delete(&file).Error; err != nil {
		return err
	}

	// Remove from storage
	if err := s.fileUtil.DeleteFile(file.FilePath); err != nil {
		// Log error but don't fail the operation
		// The file record is already deleted from DB
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	return nil
}

// detectMimeType detects the MIME type of a file
func (s *FileService) detectMimeType(file multipart.File) (string, error) {
	// Read first 512 bytes for MIME type detection
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset file position
	file.Seek(0, 0)

	// Use http.DetectContentType
	mimeType := detectContentType(buffer)
	return mimeType, nil
}

// validateFileType checks if the MIME type is allowed
func (s *FileService) validateFileType(mimeType string, allowedTypes map[string][]string) (string, bool) {
	for fileType, mimeTypes := range allowedTypes {
		for _, allowedMime := range mimeTypes {
			if mimeType == allowedMime {
				return fileType, true
			}
		}
	}
	return "", false
}

// generateUniqueFileName creates a unique filename
func (s *FileService) generateUniqueFileName(originalName, extension string) string {
	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(originalName, filepath.Ext(originalName))
	// Sanitize filename
	baseName = sanitizeFileName(baseName)
	return fmt.Sprintf("%s_%d%s", baseName, timestamp, extension)
}

// getExtensionFromMimeType returns file extension based on MIME type
func (s *FileService) getExtensionFromMimeType(mimeType string) string {
	extensions := map[string]string{
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"text/plain":               ".txt",
		"application/rtf":          ".rtf",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"text/csv":      ".csv",
		"image/jpeg":    ".jpg",
		"image/jpg":     ".jpg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/bmp":     ".bmp",
		"image/webp":    ".webp",
		"image/svg+xml": ".svg",
		"image/tiff":    ".tiff",
		"image/ico":     ".ico",
	}

	if ext, exists := extensions[mimeType]; exists {
		return ext
	}
	return ".bin"
}

// Helper functions
func detectContentType(data []byte) string {
	// Simple MIME type detection - in production, use a proper library
	if len(data) < 4 {
		return "application/octet-stream"
	}

	// PDF
	if string(data[:4]) == "%PDF" {
		return "application/pdf"
	}

	// PNG
	if len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}

	// JPEG
	if len(data) >= 3 && string(data[:3]) == "\xff\xd8\xff" {
		return "image/jpeg"
	}

	// GIF
	if len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a") {
		return "image/gif"
	}

	// Default fallback
	return "application/octet-stream"
}

func sanitizeFileName(filename string) string {
	// Remove or replace unsafe characters
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := filename
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}
