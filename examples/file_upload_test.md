# File Upload API Testing Guide

## Overview
The file upload system allows admins to upload documents, spreadsheets, and images with automatic validation and storage management.

## API Endpoints

### 1. Upload File
**POST** `/api/files/upload`

**Headers:**
```
Authorization: Bearer <admin_jwt_token>
Content-Type: multipart/form-data
```

**Form Data:**
- `file`: The file to upload (required)
- `max_file_size`: Custom max file size in bytes (optional)

**Example using curl:**
```bash
curl -X POST \
  http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "file=@/path/to/document.pdf"
```

**Response:**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "file": {
    "id": 1,
    "original_name": "document.pdf",
    "file_name": "document_1699824000.pdf",
    "file_path": "document/document_1699824000.pdf",
    "file_size": 1024000,
    "mime_type": "application/pdf",
    "file_type": "document",
    "extension": ".pdf",
    "uploaded_by": 1,
    "uploaded_by_admin": {
      "id": 1,
      "email": "admin@onas.com",
      "first_name": "Admin",
      "last_name": "Admin"
    },
    "is_active": true,
    "created_at": "2024-11-12 20:30:00",
    "updated_at": "2024-11-12 20:30:00"
  }
}
```

### 2. Get Upload Configuration
**GET** `/api/files/config`

**Headers:**
```
Authorization: Bearer <admin_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Upload configuration retrieved",
  "data": {
    "max_file_size": 5242880,
    "max_file_size_mb": 5.0,
    "allowed_types": {
      "document": [
        "application/pdf",
        "application/msword",
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "text/plain",
        "application/rtf"
      ],
      "spreadsheet": [
        "application/vnd.ms-excel",
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        "text/csv"
      ],
      "image": [
        "image/jpeg",
        "image/jpg",
        "image/png",
        "image/gif",
        "image/bmp",
        "image/webp",
        "image/svg+xml",
        "image/tiff",
        "image/ico"
      ]
    }
  }
}
```

### 3. List Files
**GET** `/api/files?page=1&limit=10&type=document`

**Headers:**
```
Authorization: Bearer <admin_jwt_token>
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `type`: Filter by file type (document, spreadsheet, image)

**Response:**
```json
{
  "success": true,
  "message": "Files retrieved successfully",
  "files": [
    {
      "id": 1,
      "original_name": "document.pdf",
      "file_name": "document_1699824000.pdf",
      "file_path": "document/document_1699824000.pdf",
      "file_size": 1024000,
      "mime_type": "application/pdf",
      "file_type": "document",
      "extension": ".pdf",
      "uploaded_by": 1,
      "uploaded_by_admin": {
        "id": 1,
        "email": "admin@onas.com",
        "first_name": "Admin",
        "last_name": "Admin"
      },
      "is_active": true,
      "created_at": "2024-11-12 20:30:00",
      "updated_at": "2024-11-12 20:30:00"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

### 4. Get File by ID
**GET** `/api/files/{id}`

**Headers:**
```
Authorization: Bearer <admin_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "File retrieved successfully",
  "data": {
    "id": 1,
    "original_name": "document.pdf",
    "file_name": "document_1699824000.pdf",
    "file_path": "document/document_1699824000.pdf",
    "file_size": 1024000,
    "mime_type": "application/pdf",
    "file_type": "document",
    "extension": ".pdf",
    "uploaded_by": 1,
    "uploaded_by_admin": {
      "id": 1,
      "email": "admin@onas.com",
      "first_name": "Admin",
      "last_name": "Admin"
    },
    "is_active": true,
    "created_at": "2024-11-12 20:30:00",
    "updated_at": "2024-11-12 20:30:00"
  }
}
```

### 5. Delete File
**DELETE** `/api/files/{id}`

**Headers:**
```
Authorization: Bearer <admin_jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

## File Validation Rules

### File Size
- **Default Maximum**: 5MB (5,242,880 bytes)
- **Configurable**: Can be adjusted per upload request
- **Validation**: Server validates actual file size

### Allowed File Types

#### Documents
- PDF (`.pdf`)
- Microsoft Word (`.doc`, `.docx`)
- Plain Text (`.txt`)
- Rich Text Format (`.rtf`)

#### Spreadsheets
- Microsoft Excel (`.xls`, `.xlsx`)
- CSV (`.csv`)

#### Images
- JPEG (`.jpg`, `.jpeg`)
- PNG (`.png`)
- GIF (`.gif`)
- BMP (`.bmp`)
- WebP (`.webp`)
- SVG (`.svg`)
- TIFF (`.tiff`)
- ICO (`.ico`)

### Security Features
- **MIME Type Detection**: Server validates actual file content
- **Path Sanitization**: Prevents directory traversal attacks
- **Unique Filenames**: Prevents filename conflicts
- **Storage Isolation**: All files stored in `storage/` directory

## Storage Structure

Files are organized by type in the storage directory:
```
storage/
├── document/
│   ├── report_1699824000.pdf
│   └── manual_1699824001.docx
├── spreadsheet/
│   ├── data_1699824002.xlsx
│   └── export_1699824003.csv
└── image/
    ├── logo_1699824004.png
    └── banner_1699824005.jpg
```

## Permissions Required

All file operations require admin authentication and specific permissions:

- **files.view**: View and list files, get upload config
- **files.create**: Upload new files
- **files.delete**: Delete existing files

## Error Responses

### File Too Large
```json
{
  "success": false,
  "message": "file size 6291456 bytes exceeds maximum allowed size of 5242880 bytes"
}
```

### Invalid File Type
```json
{
  "success": false,
  "message": "file type application/x-executable is not allowed"
}
```

### Authentication Required
```json
{
  "success": false,
  "message": "Admin authentication required"
}
```

### Permission Denied
```json
{
  "success": false,
  "message": "Permission denied: files.create required"
}
```

## Testing Steps

1. **Get Admin Token**: Login as admin to get JWT token
2. **Check Config**: GET `/api/files/config` to see upload limits
3. **Upload File**: POST `/api/files/upload` with a test file
4. **List Files**: GET `/api/files` to see uploaded files
5. **Get File**: GET `/api/files/{id}` to get specific file details
6. **Delete File**: DELETE `/api/files/{id}` to remove file

## Sample Test Files

Create test files for validation:
- `test.pdf` (small PDF document)
- `test.xlsx` (Excel spreadsheet)
- `test.png` (small image)
- `large.pdf` (>5MB file to test size limit)
- `test.exe` (executable to test type restriction)
