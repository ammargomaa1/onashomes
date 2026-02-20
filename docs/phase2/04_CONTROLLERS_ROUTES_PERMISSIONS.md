# Phase 2 — Controllers, Routes & Permissions

## 1. Controller Endpoints

### Product Controller V2 (extend existing)

**File:** `internal/api/products/controller.go`

| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `POST` | `/api/admin/products/v2` | `products.create` | Create product with full V2 fields |
| `PUT` | `/api/admin/products/v2/:id` | `products.update` | Update product V2 |
| `PATCH` | `/api/admin/products/:id/status` | `products.update` | Change product status |
| `PUT` | `/api/admin/products/:id/seo` | `seo.manage` | Update product SEO |
| `GET` | `/api/admin/products/:id/seo/suggest` | `seo.manage` | Auto-suggest SEO fields |
| `POST` | `/api/admin/products/:id/variants/v2` | `products.update` | Create variant with pricing |
| `PUT` | `/api/admin/products/:id/variants/:variantId/v2` | `products.update` | Update variant with pricing |

### Inventory Controller (NEW)

**File:** `internal/api/inventory/controller.go`

| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `POST` | `/api/admin/inventory/adjust` | `inventory.adjust` | Adjust variant inventory |
| `GET` | `/api/admin/inventory/store/:storeFrontId` | `inventory.adjust` | List inventory by store |
| `GET` | `/api/admin/inventory/variant/:variantId/store/:storeFrontId` | `inventory.adjust` | Get specific variant inventory |
| `GET` | `/api/admin/inventory/low-stock/:storeFrontId` | `inventory.adjust` | Get low stock alerts |
| `GET` | `/api/admin/inventory/:inventoryId/history` | `inventory.adjust` | Adjustment history |

### StoreFront Controller (NEW)

**File:** `internal/api/storefronts/controller.go`

| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `POST` | `/api/admin/storefronts` | `storefronts.manage` | Create store front |
| `PUT` | `/api/admin/storefronts/:id` | `storefronts.manage` | Update store front |
| `DELETE` | `/api/admin/storefronts/:id` | `storefronts.manage` | Delete store front |
| `GET` | `/api/admin/storefronts` | `storefronts.manage` | List store fronts |
| `GET` | `/api/admin/storefronts/:id` | `storefronts.manage` | Get store front |

### Public Storefront API (Domain-resolved)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/storefront/products` | None (domain-resolved) | List published products for store |
| `GET` | `/api/storefront/products/:slug` | None | Get product by slug |
| `GET` | `/api/storefront/products/:slug/structured-data` | None | JSON-LD structured data |
| `GET` | `/api/storefront/categories` | None | List store categories |
| `GET` | `/api/storefront/brands` | None | List store brands |

---

## 2. Route Definitions

### Updated Product Routes

```go
// internal/api/products/routes.go — Phase 2 additions

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
    // ... existing routes unchanged ...

    // Phase 2: Enhanced admin routes
    adminRoutes := router.Group("/admin/products")
    adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
    {
        // V2 Product CRUD
        adminRoutes.POST("/v2", middleware.RequirePermission("products.create"), controller.AdminCreateProductV2)
        adminRoutes.PUT("/v2/:id", middleware.RequirePermission("products.update"), controller.AdminUpdateProductV2)
        adminRoutes.PATCH("/:id/status", middleware.RequirePermission("products.update"), controller.AdminUpdateProductStatus)

        // SEO
        adminRoutes.PUT("/:id/seo", middleware.RequirePermission("seo.manage"), controller.AdminUpdateProductSEO)
        adminRoutes.GET("/:id/seo/suggest", middleware.RequirePermission("seo.manage"), controller.AdminSuggestSEO)

        // V2 Variants (with pricing)
        adminRoutes.POST("/:id/variants/v2", middleware.RequirePermission("products.update"), controller.AdminCreateVariantV2)
        adminRoutes.PUT("/:id/variants/:variantId/v2", middleware.RequirePermission("products.update"), controller.AdminUpdateVariantV2)
    }

    // Public storefront routes (domain-resolved)
    sfRoutes := router.Group("/storefront")
    sfRoutes.Use(middleware.StoreFrontResolver()) // NEW middleware
    {
        sfRoutes.GET("/products", controller.StorefrontListProducts)
        sfRoutes.GET("/products/:slug", controller.StorefrontGetProduct)
        sfRoutes.GET("/products/:slug/structured-data", controller.StorefrontGetStructuredData)
    }
}
```

### Inventory Routes (NEW)

```go
// internal/api/inventory/routes.go

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
    adminRoutes := router.Group("/admin/inventory")
    adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
    {
        adminRoutes.POST("/adjust", middleware.RequirePermission("inventory.adjust"), controller.AdjustInventory)
        adminRoutes.GET("/store/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.ListInventoryByStore)
        adminRoutes.GET("/variant/:variantId/store/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.GetVariantInventory)
        adminRoutes.GET("/low-stock/:storeFrontId", middleware.RequirePermission("inventory.adjust"), controller.GetLowStockAlerts)
        adminRoutes.GET("/:inventoryId/history", middleware.RequirePermission("inventory.adjust"), controller.GetAdjustmentHistory)
    }
}
```

### StoreFront Routes (NEW)

```go
// internal/api/storefronts/routes.go

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
    adminRoutes := router.Group("/admin/storefronts")
    adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
    {
        adminRoutes.POST("", middleware.RequirePermission("storefronts.manage"), controller.Create)
        adminRoutes.PUT("/:id", middleware.RequirePermission("storefronts.manage"), controller.Update)
        adminRoutes.DELETE("/:id", middleware.RequirePermission("storefronts.manage"), controller.Delete)
        adminRoutes.GET("", middleware.RequirePermission("storefronts.manage"), controller.List)
        adminRoutes.GET("/:id", middleware.RequirePermission("storefronts.manage"), controller.GetByID)
    }
}
```

---

## 3. New Permissions

Add to `internal/permissions/scanner.go` `addPredefinedPermissions()`:

```go
// Phase 2 permissions
{Name: "inventory.adjust", Description: "Adjust product inventory", Module: "inventory", Action: "adjust"},
{Name: "inventory.view", Description: "View inventory levels", Module: "inventory", Action: "view"},
{Name: "storefronts.manage", Description: "Manage store fronts", Module: "storefronts", Action: "manage"},
{Name: "seo.manage", Description: "Manage product SEO metadata", Module: "seo", Action: "manage"},
```

### Full Permission Matrix

| Permission | Description | Used By |
|-----------|-------------|---------|
| `products.create` | Create products | Admin |
| `products.update` | Update products & variants | Admin |
| `products.delete` | Delete products | Admin |
| `products.view` | View product list & details | Admin |
| `inventory.adjust` | Adjust inventory quantities | Warehouse/Admin |
| `inventory.view` | View inventory levels | Warehouse/Admin |
| `storefronts.manage` | Full CRUD on store fronts | Super Admin |
| `seo.manage` | Manage SEO metadata | Marketing/Admin |

---

## 4. main.go Registration

```go
// In cmd/api/main.go — add inside the api group:

// StoreFront module
sfRepo := storefronts.NewRepository(db)
sfService := storefronts.NewService(sfRepo)
sfController := storefronts.NewController(sfService)
storefronts.RegisterRoutes(api, sfController)

// Inventory module
invRepo := inventory.NewRepository(db)
invService := inventory.NewService(sqlDB, invRepo)
invController := inventory.NewController(invService)
inventory.RegisterRoutes(api, invController)
```
