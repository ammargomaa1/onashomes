# Quick Start Guide

## 🚀 Get Started in 5 Minutes

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Setup Database
Create a PostgreSQL database:
```bash
createdb ecommerce_db
```

Or using psql:
```sql
CREATE DATABASE ecommerce_db;
```

### 3. Configure Environment
The `.env` file has been created. Update it with your database credentials:
```bash
nano .env
```

Minimum required changes:
```env
DB_PASSWORD=your_postgres_password
JWT_SECRET=change-this-to-a-random-secret-key
JWT_REFRESH_SECRET=change-this-to-another-random-secret-key
```

### 4. Run the Application
```bash
go run cmd/api/main.go
```

The API will:
- ✅ Connect to the database
- ✅ Run auto-migrations
- ✅ Seed default roles and permissions
- ✅ Start on http://localhost:8080

### 5. Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Register a User:**
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

**Get Profile:**
```bash
curl http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📝 Creating Your First Admin

Since admins require a role_id, you need to:

1. Get the role ID from the database:
```sql
SELECT id, name FROM roles;
```

2. Create an admin via database or create an endpoint for super admin creation.

For quick testing, you can manually insert an admin:
```sql
INSERT INTO admins (id, email, password, first_name, last_name, role_id, is_active, created_at, updated_at)
VALUES (
  gen_random_uuid(),
  'admin@example.com',
  '$2a$10$...',  -- Use bcrypt hash of your password
  'Admin',
  'User',
  (SELECT id FROM roles WHERE name = 'Super Admin'),
  true,
  NOW(),
  NOW()
);
```

Or use this Go snippet to generate a bcrypt hash:
```go
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "admin123"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

## 🔧 Common Commands

**Run in development:**
```bash
go run cmd/api/main.go
```

**Build for production:**
```bash
./scripts/build.sh
```

**Run built binary:**
```bash
./bin/ecommerce-api-linux-amd64
```

**Create migration:**
```bash
./scripts/create_migration.sh add_products_table
```

**Rollback migration:**
```bash
./scripts/migrate_down.sh 1
```

## 🎯 Next Steps

1. **Add your business models** in `internal/models/`
2. **Create new modules** following the users/admins pattern
3. **Add custom permissions** in `internal/database/migrations.go`
4. **Implement business logic** in service layers
5. **Add validation** in controllers
6. **Write tests** for your endpoints

## 🐛 Troubleshooting

**Database connection failed:**
- Check PostgreSQL is running: `sudo systemctl status postgresql`
- Verify credentials in `.env`
- Ensure database exists: `psql -l`

**Port already in use:**
- Change `SERVER_PORT` in `.env`
- Or kill the process: `lsof -ti:8080 | xargs kill -9`

**Import errors:**
- Run: `go mod tidy`
- Clear cache: `go clean -modcache`

## 📚 Learn More

- Read the full [README.md](README.md)
- Check the [API Documentation](#) (add your API docs link)
- Review the code structure in `internal/`

---

**Happy Building! 🎉**
