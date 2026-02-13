# Phase 2 — Repository Interfaces, Services & Controllers

## 1. Repository Interfaces

### Product Repository (Enhanced)

**File:** `internal/api/products/repository.go` — extend existing interface

```go
type Repository interface {
    // === EXISTING (keep as-is) ===
    CreateProduct(ctx context.Context, tx *sql.Tx, req requests.AdminProductRequest) (int64, error)
    UpdateProduct(ctx context.Context, tx *sql.Tx, id int64, req requests.AdminProductRequest) error
    SoftDeleteProduct(ctx context.Context, tx *sql.Tx, id int64) error
    // ... existing methods ...

    // === NEW PHASE 2 ===
    // Product CRUD (enhanced)
    CreateProductV2(ctx context.Context, tx *sql.Tx, req requests.CreateProductV2Request) (int64, error)
    UpdateProductV2(ctx context.Context, tx *sql.Tx, id int64, req requests.UpdateProductV2Request) error
    GetProductBySlugAndStore(ctx context.Context, slug string, storeFrontID int64) (*AdminProductDetail, error)
    IsSlugUniqueForStore(ctx context.Context, slug string, storeFrontID int64, excludeProductID int64) (bool, error)
    ListAdminProductsV2(ctx context.Context, filter ProductFilter, pagination *utils.Pagination) ([]AdminProductListItem, int64, error)

    // Store assignment
    AssignProductToStores(ctx context.Context, tx *sql.Tx, productID int64, storeFrontIDs []int64) error
    RemoveProductFromStores(ctx context.Context, tx *sql.Tx, productID int64) error

    // SEO
    UpsertProductSEO(ctx context.Context, tx *sql.Tx, productID int64, seo requests.ProductSEORequest) error
    GetProductSEO(ctx context.Context, productID int64) (*models.ProductSEO, error)

    // Storefront queries (public)
    ListStorefrontProducts(ctx context.Context, storeFrontID int64, filter StorefrontProductFilter, pagination *utils.Pagination) ([]StorefrontProduct, int64, error)
    GetStorefrontProduct(ctx context.Context, storeFrontID int64, slug string) (*StorefrontProductDetail, error)
}
```

### Inventory Repository (NEW)

**File:** `internal/api/inventory/repository.go`

```go
type Repository interface {
    GetVariantInventory(ctx context.Context, variantID, storeFrontID int64) (*models.VariantInventory, error)
    ListInventoryByStore(ctx context.Context, storeFrontID int64, pagination *utils.Pagination) ([]VariantInventoryItem, int64, error)
    AdjustInventory(ctx context.Context, tx *sql.Tx, req AdjustInventoryParams) (*models.VariantInventory, error)
    GetLowStockItems(ctx context.Context, storeFrontID int64, pagination *utils.Pagination) ([]VariantInventoryItem, int64, error)
    ListAdjustmentHistory(ctx context.Context, variantInventoryID int64, pagination *utils.Pagination) ([]models.InventoryAdjustment, int64, error)
    EnsureInventoryRecord(ctx context.Context, tx *sql.Tx, variantID, storeFrontID int64) (*models.VariantInventory, error)
}

// AdjustInventoryParams holds all data for an inventory adjustment
type AdjustInventoryParams struct {
    VariantID    int64
    StoreFrontID int64
    Adjustment   int    // positive = add, negative = subtract
    Reason       string // "adjustment", "correction", "restock"
    Notes        string
    AdjustedBy   int64  // admin ID
}
```

### StoreFront Repository (NEW)

**File:** `internal/api/storefronts/repository.go`

```go
type Repository interface {
    Create(ctx context.Context, storeFront *models.StoreFront) error
    Update(ctx context.Context, storeFront *models.StoreFront) error
    Delete(ctx context.Context, id int64) error
    GetByID(ctx context.Context, id int64) (*models.StoreFront, error)
    GetByDomain(ctx context.Context, domain string) (*models.StoreFront, error)
    GetBySlug(ctx context.Context, slug string) (*models.StoreFront, error)
    List(ctx context.Context, pagination *utils.Pagination) ([]models.StoreFront, int64, error)
    IsSlugUnique(ctx context.Context, slug string, excludeID int64) (bool, error)
    IsDomainUnique(ctx context.Context, domain string, excludeID int64) (bool, error)
}
```

### Product Filter Structs

```go
// ProductFilter for admin list filtering
type ProductFilter struct {
    StoreFrontID  *int64
    CategoryID    *int64
    BrandID       *int64
    SupplierID    *int64
    Status        *string  // draft, active, inactive, archived
    AttributeType *string  // size, color, or empty for simple
    StockStatus   *string  // in_stock, low_stock, out_of_stock
    DateFrom      *time.Time
    DateTo        *time.Time
    Search        string   // search name_en, name_ar
}

// StorefrontProductFilter for public storefront
type StorefrontProductFilter struct {
    CategoryID *int64
    BrandID    *int64
    MinPrice   *float64
    MaxPrice   *float64
    Search     string
    SortBy     string // price_asc, price_desc, newest, best_seller
}
```

---

## 2. Service Layer

### Product Service V2

**File:** `internal/api/products/service.go` — extend existing

```go
// New methods added to Service interface
type Service interface {
    // ... existing methods ...

    // Phase 2
    CreateProductV2(ctx context.Context, req requests.CreateProductV2Request) utils.IResource
    UpdateProductV2(ctx context.Context, id int64, req requests.UpdateProductV2Request) utils.IResource
    UpdateProductStatus(ctx context.Context, id int64, status string) utils.IResource
    UpdateProductSEO(ctx context.Context, id int64, seo requests.ProductSEORequest) utils.IResource
    SuggestSEO(ctx context.Context, id int64) utils.IResource

    // Storefront public
    StorefrontListProducts(ctx context.Context, storeFrontID int64, filter StorefrontProductFilter, pagination *utils.Pagination) utils.IResource
    StorefrontGetProduct(ctx context.Context, storeFrontID int64, slug string) utils.IResource
    GetProductStructuredData(ctx context.Context, storeFrontID int64, slug string) utils.IResource
}
```

**Key Service Business Logic:**

```go
func (s *service) CreateProductV2(ctx context.Context, req requests.CreateProductV2Request) utils.IResource {
    // 1. Validate required fields (name_en, name_ar)
    // 2. Validate attribute_type rule: only "size", "color", or nil
    // 3. Auto-generate slug from name_en (slugify)
    // 4. Validate slug uniqueness per each assigned store
    // 5. Validate brand_id, category_id, supplier_id existence
    // 6. Begin transaction
    //    a. Insert product
    //    b. Assign to store_fronts (pivot inserts)
    //    c. Upsert SEO if provided
    // 7. Commit
    // 8. Return created product detail
}

func (s *service) UpdateProductStatus(ctx context.Context, id int64, status string) utils.IResource {
    // Validate status transitions:
    //   draft    → active (requires ≥1 active variant with inventory)
    //   active   → inactive, archived
    //   inactive → active (re-validate), archived
    //   archived → draft (re-open)
    // Cannot publish without at least 1 active variant
    // Cannot activate if no inventory exists
}
```

### Inventory Service (NEW)

**File:** `internal/api/inventory/service.go`

```go
type Service interface {
    AdjustInventory(ctx context.Context, req requests.AdjustInventoryRequest, adminID int64) utils.IResource
    GetVariantInventory(ctx context.Context, variantID, storeFrontID int64) utils.IResource
    ListInventory(ctx context.Context, storeFrontID int64, pagination *utils.Pagination) utils.IResource
    GetLowStockAlerts(ctx context.Context, storeFrontID int64, pagination *utils.Pagination) utils.IResource
    GetAdjustmentHistory(ctx context.Context, variantInventoryID int64, pagination *utils.Pagination) utils.IResource
}
```

**Key Inventory Business Logic:**

```go
func (s *service) AdjustInventory(ctx context.Context, req requests.AdjustInventoryRequest, adminID int64) utils.IResource {
    // 1. Validate reason is one of allowed values
    // 2. Validate adjustment won't result in negative quantity
    // 3. Begin transaction
    //    a. EnsureInventoryRecord (upsert)
    //    b. Lock row (SELECT ... FOR UPDATE)
    //    c. Check: current_qty + adjustment >= 0
    //    d. Update quantity
    //    e. Insert inventory_adjustment audit record
    // 4. Commit
    // 5. Return updated inventory
}
```

### StoreFront Service (NEW)

**File:** `internal/api/storefronts/service.go`

```go
type Service interface {
    Create(ctx context.Context, req requests.CreateStoreFrontRequest) utils.IResource
    Update(ctx context.Context, id int64, req requests.UpdateStoreFrontRequest) utils.IResource
    Delete(ctx context.Context, id int64) utils.IResource
    GetByID(ctx context.Context, id int64) utils.IResource
    List(ctx context.Context, pagination *utils.Pagination) utils.IResource
    ResolveByDomain(ctx context.Context, domain string) utils.IResource
}
```

---

## 3. Request DTOs

### Product Requests V2

**File:** `internal/api/products/requests/v2.go`

```go
type CreateProductV2Request struct {
    NameEn             string   `json:"name_en" binding:"required"`
    NameAr             string   `json:"name_ar" binding:"required"`
    DescriptionEn      string   `json:"description_en"`
    DescriptionAr      string   `json:"description_ar"`
    BrandID            *int64   `json:"brand_id"`
    CategoryID         *int64   `json:"category_id"`
    SupplierID         *int64   `json:"supplier_id"`
    IsInternalSupplier bool     `json:"is_internal_supplier"`
    AttributeType      *string  `json:"attribute_type"` // "size" | "color" | null
    StoreFrontIDs      []int64  `json:"store_front_ids" binding:"required,min=1"`
    IsFeatured         bool     `json:"is_featured"`
    IsNew              bool     `json:"is_new"`
    SEO                *ProductSEORequest `json:"seo"`
}

type UpdateProductV2Request struct {
    NameEn             string   `json:"name_en" binding:"required"`
    NameAr             string   `json:"name_ar" binding:"required"`
    DescriptionEn      string   `json:"description_en"`
    DescriptionAr      string   `json:"description_ar"`
    BrandID            *int64   `json:"brand_id"`
    CategoryID         *int64   `json:"category_id"`
    SupplierID         *int64   `json:"supplier_id"`
    IsInternalSupplier bool     `json:"is_internal_supplier"`
    AttributeType      *string  `json:"attribute_type"`
    StoreFrontIDs      []int64  `json:"store_front_ids" binding:"required,min=1"`
    IsFeatured         bool     `json:"is_featured"`
    IsNew              bool     `json:"is_new"`
    IsBestSeller       bool     `json:"is_best_seller"`
    SEO                *ProductSEORequest `json:"seo"`
}

type ProductSEORequest struct {
    MetaTitleEn      string `json:"meta_title_en" binding:"max=60"`
    MetaTitleAr      string `json:"meta_title_ar" binding:"max=60"`
    MetaDescriptionEn string `json:"meta_description_en" binding:"max=160"`
    MetaDescriptionAr string `json:"meta_description_ar" binding:"max=160"`
    MetaKeywords     string `json:"meta_keywords"`
    CanonicalURL     string `json:"canonical_url"`
    OgTitle          string `json:"og_title"`
    OgDescription    string `json:"og_description"`
    OgImage          string `json:"og_image"`
}

type CreateVariantV2Request struct {
    SKU            string   `json:"sku" binding:"required"`
    AttributeValue string   `json:"attribute_value"` // "XL", "Red"
    Price          *float64 `json:"price" binding:"required"`
    CompareAtPrice *float64 `json:"compare_at_price"`
    CostPrice      *float64 `json:"cost_price"`
    Barcode        string   `json:"barcode"`
    Weight         *float64 `json:"weight"`
    IsActive       bool     `json:"is_active"`
    ImageFileIDs   []int64  `json:"image_file_ids"`
}

type UpdateProductStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=draft active inactive archived"`
}
```

### Inventory Requests

**File:** `internal/api/inventory/requests/types.go`

```go
type AdjustInventoryRequest struct {
    ProductVariantID int64  `json:"product_variant_id" binding:"required"`
    StoreFrontID     int64  `json:"store_front_id" binding:"required"`
    Adjustment       int    `json:"adjustment" binding:"required"` // +/- value
    Reason           string `json:"reason" binding:"required,oneof=adjustment correction restock"`
    Notes            string `json:"notes"`
}
```

### StoreFront Requests

**File:** `internal/api/storefronts/requests/types.go`

```go
type CreateStoreFrontRequest struct {
    Name            string `json:"name" binding:"required"`
    Slug            string `json:"slug" binding:"required"`
    Domain          string `json:"domain" binding:"required"`
    Currency        string `json:"currency" binding:"required"`
    DefaultLanguage string `json:"default_language" binding:"required"`
}

type UpdateStoreFrontRequest struct {
    Name            string `json:"name" binding:"required"`
    Domain          string `json:"domain" binding:"required"`
    Currency        string `json:"currency" binding:"required"`
    DefaultLanguage string `json:"default_language" binding:"required"`
    IsActive        bool   `json:"is_active"`
}
```
