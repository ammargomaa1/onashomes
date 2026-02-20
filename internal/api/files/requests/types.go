package requests

// UploadRequest represents file upload request configuration
// (currently unused in handlers but kept for future JSON-based configs).
type UploadRequest struct {
	MaxFileSize int64 `json:"max_file_size,omitempty"`
}
