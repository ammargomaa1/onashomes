# Automatic Permission System

This document explains the automatic permission scanning and management system implemented in the OnasHomes API.

## Overview

The system automatically scans your API routes and discovers permissions, then syncs them to the database and assigns them to the Super Admin role on every application startup.

## How It Works

### 1. **Route Scanning**
- The system scans all registered Gin routes when the application starts
- It looks for routes that use the `middleware.RequirePermission()` middleware
- It also infers permissions from route patterns and HTTP methods

### 2. **Permission Discovery**
The scanner discovers permissions in two ways:

#### A. **Predefined Permissions**
Permissions that are explicitly defined in the code:
```go
// From routes.go files
middleware.RequirePermission("admins.view")
middleware.RequirePermission("admins.create")
middleware.RequirePermission("users.view")
```

#### B. **Inferred Permissions**
Permissions inferred from route patterns:
- `GET /api/users` → `users.view`
- `POST /api/users` → `users.create`
- `PUT /api/users/:id` → `users.update`
- `DELETE /api/users/:id` → `users.delete`

### 3. **Database Synchronization**
- New permissions are automatically added to the database
- Existing permissions are left unchanged
- The Super Admin role automatically gets all new permissions

### 4. **Automatic Assignment**
- All discovered permissions are automatically assigned to the "Super Admin" role
- This ensures the Super Admin always has access to all functionality

## Permission Naming Convention

Permissions follow the pattern: `{module}.{action}`

### Modules
- `users` - User management
- `admins` - Admin management  
- `roles` - Role management
- `permissions` - Permission management
- `products` - Product management (future)

### Actions
- `view` - Read/view operations (GET requests)
- `create` - Create operations (POST requests)
- `update` - Update operations (PUT/PATCH requests)
- `delete` - Delete operations (DELETE requests)

## Examples

### Current Permissions
```
admins.view       - View and list admins
admins.create     - Create new admins
admins.update     - Update existing admins
admins.delete     - Delete existing admins

users.view        - View and list users
users.create      - Create new users
users.update      - Update existing users
users.delete      - Delete existing users
```

### Adding New Permissions

#### Method 1: Use RequirePermission Middleware
```go
// In your routes file
authenticated.GET("/reports", middleware.RequirePermission("reports.view"), controller.GetReports)
authenticated.POST("/reports", middleware.RequirePermission("reports.create"), controller.CreateReport)
```

#### Method 2: Add to Predefined List
```go
// In internal/permissions/scanner.go
{Name: "reports.view", Description: "View reports", Module: "reports", Action: "view"},
{Name: "reports.create", Description: "Create reports", Module: "reports", Action: "create"},
```

## System Startup Flow

1. **Database Migration** - Tables are created/updated
2. **Seed Default Data** - Roles and default admin user are created
3. **Route Registration** - All API routes are registered
4. **Permission Scanning** - Routes are scanned for permissions
5. **Database Sync** - New permissions are added to database
6. **Super Admin Assignment** - New permissions are assigned to Super Admin
7. **Server Start** - API server starts listening

## Logs

The system provides detailed logging during startup:

```
🔍 Scanning routes for permissions...
🔍 Found 16 permissions to sync
➕ New permission: products.view - View and list products
➕ New permission: products.create - Create new products
✅ Added 4 new permissions
🔐 Assigned 4 permissions to Super Admin
```

## Configuration

### Default Roles Created
- **Super Admin**: Gets all permissions automatically
- **Admin**: Limited permissions (manually assigned)

### Default Admin User
- **Email**: `admin@onas.com`
- **Password**: `password`
- **Role**: Admin (limited permissions)

## Adding Custom Permissions

You can add custom permissions programmatically:

```go
scanner := permissions.NewScanner(db)
scanner.AddCustomPermission("reports.export", "Export reports to PDF/Excel")
scanner.ScanAndSync(router)
```

## Benefits

1. **Automatic Discovery** - No need to manually maintain permission lists
2. **Always Up-to-Date** - Permissions are synced on every startup
3. **Super Admin Safety** - Super Admin always has access to everything
4. **Developer Friendly** - Just add `RequirePermission()` middleware and it's discovered
5. **Future Proof** - New modules automatically get their permissions

## Security Considerations

- Super Admin role gets ALL permissions automatically
- Regular admin roles need manual permission assignment
- Permissions are only added, never removed automatically
- Use the Admin panel to manage role permissions for non-Super Admin roles

## Troubleshooting

### Permission Not Discovered
1. Check if route uses `RequirePermission()` middleware
2. Verify route path starts with `/api/`
3. Add permission to predefined list if needed

### Super Admin Missing Permissions
- Super Admin automatically gets all permissions on startup
- If missing, restart the application
- Check logs for permission scanning errors

### Custom Module Permissions
Add your module permissions to the predefined list in `internal/permissions/scanner.go`
