package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// StorageDir is the base directory for all file storage
	StorageDir = "../../storage"
)

// FileUtil handles file operations within the storage directory
type FileUtil struct {
	baseDir string
}

// NewFileUtil creates a new file utility instance
func NewFileUtil() *FileUtil {
	return &FileUtil{
		baseDir: StorageDir,
	}
}

// SaveFile saves content to a file in the storage directory
// Returns the relative path from storage root
func (f *FileUtil) SaveFile(relativePath string, content []byte) (string, error) {
	// Clean and validate the path
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return "", fmt.Errorf("invalid file path: %s", relativePath)
	}

	// Create full path
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Write file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %v", fullPath, err)
	}

	return cleanPath, nil
}

// SaveFileFromReader saves content from an io.Reader to a file
// Returns the relative path from storage root
func (f *FileUtil) SaveFileFromReader(relativePath string, reader io.Reader) (string, error) {
	// Clean and validate the path
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return "", fmt.Errorf("invalid file path: %s", relativePath)
	}

	// Create full path
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", fullPath, err)
	}
	defer file.Close()

	// Copy content from reader to file
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file content: %v", err)
	}

	return cleanPath, nil
}

// DeleteFile deletes a file from the storage directory
func (f *FileUtil) DeleteFile(relativePath string) error {
	// Clean and validate the path
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return fmt.Errorf("invalid file path: %s", relativePath)
	}

	// Create full path
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Check if file exists
	if !f.FileExists(relativePath) {
		return fmt.Errorf("file does not exist: %s", relativePath)
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %v", fullPath, err)
	}

	return nil
}

// FileExists checks if a file exists in the storage directory
func (f *FileUtil) FileExists(relativePath string) bool {
	// Clean and validate the path
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return false
	}

	// Create full path
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Check if file exists and is not a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// GetFullPath returns the full filesystem path for a relative path
func (f *FileUtil) GetFullPath(relativePath string) string {
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return ""
	}
	return filepath.Join(f.baseDir, cleanPath)
}

// GetFileInfo returns file information for a file in storage
func (f *FileUtil) GetFileInfo(relativePath string) (os.FileInfo, error) {
	fullPath := f.GetFullPath(relativePath)
	if fullPath == "" {
		return nil, fmt.Errorf("invalid file path: %s", relativePath)
	}

	return os.Stat(fullPath)
}

// ListFiles lists all files in a directory within storage
func (f *FileUtil) ListFiles(relativeDirPath string) ([]string, error) {
	cleanPath := f.cleanPath(relativeDirPath)
	fullPath := filepath.Join(f.baseDir, cleanPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", fullPath, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(cleanPath, entry.Name()))
		}
	}

	return files, nil
}

// EnsureStorageDir creates the storage directory if it doesn't exist
func (f *FileUtil) EnsureStorageDir() error {
	return os.MkdirAll(f.baseDir, 0755)
}

// cleanPath cleans and validates a file path to prevent directory traversal
func (f *FileUtil) cleanPath(path string) string {
	// Clean the path
	cleaned := filepath.Clean(path)

	// Remove leading slash if present
	cleaned = strings.TrimPrefix(cleaned, "/")

	// Check for directory traversal attempts
	if strings.Contains(cleaned, "..") {
		return ""
	}

	// Check for absolute paths
	if filepath.IsAbs(cleaned) {
		return ""
	}

	// Ensure path is not empty
	if cleaned == "." || cleaned == "" {
		return ""
	}

	return cleaned
}

// GetRelativePath returns the relative path from storage root
// This is useful when you have a full path and want the relative part
func (f *FileUtil) GetRelativePath(fullPath string) (string, error) {
	absStorageDir, err := filepath.Abs(f.baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute storage path: %v", err)
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute file path: %v", err)
	}

	relativePath, err := filepath.Rel(absStorageDir, absFullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %v", err)
	}

	// Ensure the file is within storage directory
	if strings.HasPrefix(relativePath, "..") {
		return "", fmt.Errorf("file is outside storage directory")
	}

	return relativePath, nil
}

// ReadFile reads the content of a file from storage
func (f *FileUtil) ReadFile(relativePath string) ([]byte, error) {
	// Clean and validate the path
	cleanPath := f.cleanPath(relativePath)
	if cleanPath == "" {
		return nil, fmt.Errorf("invalid file path: %s", relativePath)
	}

	// Create full path
	fullPath := filepath.Join(f.baseDir, cleanPath)

	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", fullPath, err)
	}

	return data, nil
}