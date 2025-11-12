package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/onas/ecommerce-api/internal/utils"
)

func main() {
	// Create file utility instance
	fileUtil := utils.NewFileUtil()

	// Ensure storage directory exists
	if err := fileUtil.EnsureStorageDir(); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Example 1: Save a text file
	content := []byte("Hello, World!\nThis is a test file.")
	relativePath, err := fileUtil.SaveFile("documents/test.txt", content)
	if err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}
	fmt.Printf("✅ File saved at: %s\n", relativePath)

	// Example 2: Check if file exists
	exists := fileUtil.FileExists("documents/test.txt")
	fmt.Printf("📁 File exists: %t\n", exists)

	// Example 3: Save file from reader
	reader := strings.NewReader("This content comes from a reader")
	relativePath2, err := fileUtil.SaveFileFromReader("uploads/from_reader.txt", reader)
	if err != nil {
		log.Fatalf("Failed to save file from reader: %v", err)
	}
	fmt.Printf("✅ File from reader saved at: %s\n", relativePath2)

	// Example 4: Get file info
	info, err := fileUtil.GetFileInfo("documents/test.txt")
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	fmt.Printf("📊 File size: %d bytes, Modified: %s\n", info.Size(), info.ModTime())

	// Example 5: List files in directory
	files, err := fileUtil.ListFiles("documents")
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	fmt.Printf("📂 Files in documents/: %v\n", files)

	// Example 6: Get full path
	fullPath := fileUtil.GetFullPath("documents/test.txt")
	fmt.Printf("🔗 Full path: %s\n", fullPath)

	// Example 7: Delete file
	err = fileUtil.DeleteFile("uploads/from_reader.txt")
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}
	fmt.Println("🗑️  File deleted successfully")

	// Example 8: Try to access invalid paths (security test)
	_, err = fileUtil.SaveFile("../../../etc/passwd", []byte("hacker content"))
	if err != nil {
		fmt.Printf("🔒 Security check passed: %v\n", err)
	}

	fmt.Println("\n✅ All file operations completed successfully!")
}
