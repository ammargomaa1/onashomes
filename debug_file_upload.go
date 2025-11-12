package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/onas/ecommerce-api/internal/services"
	"github.com/onas/ecommerce-api/internal/database"
)

func main() {
	// Initialize database
	db := database.GetDB()
	
	// Create file service
	fileService := services.NewFileService(db)
	
	// Create a test file in memory
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Create a form file field
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}
	
	// Write test content
	testContent := "This is a test file for debugging upload"
	_, err = fileWriter.Write([]byte(testContent))
	if err != nil {
		log.Fatalf("Failed to write test content: %v", err)
	}
	
	// Close the writer
	writer.Close()
	
	// Create a multipart reader
	reader := multipart.NewReader(&buf, writer.Boundary())
	
	// Parse the form
	form, err := reader.ReadForm(32 << 20) // 32MB max
	if err != nil {
		log.Fatalf("Failed to read form: %v", err)
	}
	defer form.RemoveAll()
	
	// Get the file header
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		log.Fatalf("No file found in form")
	}
	
	fileHeader := fileHeaders[0]
	
	fmt.Printf("📁 File name: %s\n", fileHeader.Filename)
	fmt.Printf("📊 File size: %d bytes\n", fileHeader.Size)
	
	// Test opening the file
	file, err := fileHeader.Open()
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	
	// Read content to verify
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file content: %v", err)
	}
	
	fmt.Printf("📄 File content: %s\n", string(content))
	
	// Now test the upload service
	fmt.Println("\n🔄 Testing file upload service...")
	
	uploadedFile, err := fileService.UploadFile(fileHeader, 1, nil)
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}
	
	fmt.Printf("✅ Upload successful!\n")
	fmt.Printf("📁 File ID: %d\n", uploadedFile.ID)
	fmt.Printf("📁 Original name: %s\n", uploadedFile.OriginalName)
	fmt.Printf("📁 File path: %s\n", uploadedFile.FilePath)
	fmt.Printf("📁 MIME type: %s\n", uploadedFile.MimeType)
	fmt.Printf("📁 File type: %s\n", uploadedFile.FileType)
	
	// Check if file exists on disk
	if _, err := os.Stat("storage/" + uploadedFile.FilePath); err != nil {
		fmt.Printf("❌ File not found on disk: %v\n", err)
	} else {
		fmt.Printf("✅ File exists on disk\n")
		
		// Read the saved file content
		savedContent, err := os.ReadFile("storage/" + uploadedFile.FilePath)
		if err != nil {
			fmt.Printf("❌ Failed to read saved file: %v\n", err)
		} else {
			fmt.Printf("📄 Saved content: %s\n", string(savedContent))
		}
	}
}
