# E-Commerce API - Go Starter Pack

A production-ready Go API starter pack for e-commerce applications with OAuth2 JWT authentication, role-based access control (RBAC), and a clean architecture following the Controller-Service-Repository pattern.

## 🚀 Features

- **Gin Framework** - Fast HTTP web framework
- **PostgreSQL** with GORM - Robust database with ORM
- **Auto Migration** - Automatic database schema migration
- **OAuth2 JWT Authentication** - Secure token-based authentication
- **Dual Entity System** - Separate authentication for Users and Admins
- **RBAC** - Role-based access control with permissions
- **Permission Middleware** - Fine-grained access control
- **Unified API Response** - Consistent response structure
- **Pagination** - Built-in pagination support
- **Clean Architecture** - Controller-Service-Repository pattern
- **Modular Structure** - Separate modules for Users and Admins
- **Migration Scripts** - Easy database migration management
- **Build Scripts** - Cross-platform build support

## 📁 Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── config/
│   └── config.go                   # Configuration management
├── internal/
│   ├── api/
│   │   ├── users/                  # User module
│   │   │   ├── controller.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── routes.go
│   │   └── admins/                 # Admin module
│   │       ├── controller.go
│   │       ├── service.go
│   │       ├── repository.go
│   │       └── routes.go
│   ├── database/
│   │   ├── connection.go           # Singleton DB connection
│   │   └── migrations.go           # Auto-migration logic
│   ├── middleware/
│   │   ├── auth.go                 # JWT authentication
│   │   ├── permission.go           # Permission checking
│   │   └── cors.go                 # CORS handling
│   ├── models/
│   │   ├── user.go
│   │   ├── admin.go
│   │   ├── role.go
│   │   └── permission.go
│   └── utils/
│       ├── jwt.go                  # JWT utilities
│       ├── hash.go                 # Password hashing
│       ├── pagination.go           # Pagination helper
│       └── response.go             # Unified responses
├── scripts/
│   ├── create_migration.sh         # Create migration files
│   ├── migrate_down.sh             # Rollback migrations
│   └── build.sh                    # Build executables
├── migrations/                     # Migration files
├── .env.example                    # Environment template
└── README.md
```

## 🛠️ Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   cd /home/ammar/work/OnasHomes
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Create database**
   ```bash
   createdb ecommerce_db
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

The API will start on `http://localhost:8080`

## 🔧 Configuration

Edit `.env` file with your settings:

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret
JWT_EXPIRY=3600
JWT_REFRESH_EXPIRY=604800

# Pagination
DEFAULT_PAGE_SIZE=10
MAX_PAGE_SIZE=100
```

## 📚 API Endpoints

### Health Check
- `GET /health` - API health status

### User Endpoints
- `POST /api/users/register` - Register new user
- `POST /api/users/login` - User login
- `POST /api/users/refresh` - Refresh access token
- `GET /api/users/profile` - Get user profile (authenticated)
- `PUT /api/users/profile` - Update user profile (authenticated)
- `GET /api/users` - List users (admin only)

### Admin Endpoints
- `POST /api/admins/login` - Admin login
- `POST /api/admins/refresh` - Refresh access token
- `GET /api/admins` - List admins (requires `admins.view`)
- `GET /api/admins/:id` - Get admin by ID (requires `admins.view`)
- `POST /api/admins` - Create admin (requires `admins.create`)
- `PUT /api/admins/:id` - Update admin (requires `admins.update`)
- `DELETE /api/admins/:id` - Delete admin (requires `admins.delete`)

## 🔐 Authentication

### User Registration
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### User Login
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Admin Login
```bash
curl -X POST http://localhost:8080/api/admins/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Using JWT Token
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📄 Pagination

All list endpoints support pagination:

```bash
GET /api/users?page=1&limit=10&sort=created_at&order=desc
```

Parameters:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `sort` - Sort field (default: created_at)
- `order` - Sort order: asc/desc (default: desc)

Response includes metadata:
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

## 🔑 Permissions System

Default permissions:
- `users.view`, `users.create`, `users.update`, `users.delete`
- `admins.view`, `admins.create`, `admins.update`, `admins.delete`
- `roles.view`, `roles.create`, `roles.update`, `roles.delete`

Default roles:
- **Super Admin** - All permissions
- **Admin** - User management permissions

## 🗄️ Database Migrations

### Auto Migration
Runs automatically on application start.

### Manual Migrations

**Create new migration:**
```bash
chmod +x scripts/create_migration.sh
./scripts/create_migration.sh add_products_table
```

**Rollback migrations:**
```bash
chmod +x scripts/migrate_down.sh
./scripts/migrate_down.sh 1  # Rollback last migration
```

## 🏗️ Building

### Build for all platforms:
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

Executables will be in `bin/` directory:
- `ecommerce-api-linux-amd64` - Linux
- `ecommerce-api-windows-amd64.exe` - Windows
- `ecommerce-api-darwin-amd64` - macOS Intel
- `ecommerce-api-darwin-arm64` - macOS Apple Silicon

### Run executable:
```bash
./bin/ecommerce-api-linux-amd64
```

## 🧪 Testing

```bash
go test ./...
```

## 📦 Adding New Modules

1. Create module directory: `internal/api/products/`
2. Add files: `repository.go`, `service.go`, `controller.go`, `routes.go`
3. Register routes in `cmd/api/main.go`
4. Add models in `internal/models/`
5. Update migrations in `internal/database/migrations.go`

## 🔒 Security Best Practices

- ✅ Passwords are hashed with bcrypt
- ✅ JWT tokens with expiration
- ✅ Separate access and refresh tokens
- ✅ Role-based access control
- ✅ Permission-based authorization
- ✅ CORS middleware included
- ✅ SQL injection protection via GORM
- ⚠️ Change JWT secrets in production
- ⚠️ Use HTTPS in production
- ⚠️ Enable rate limiting (not included)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📝 License

MIT License

## 👨‍💻 Author

Onas Homes Development Team

---

**Happy Coding! 🚀**
