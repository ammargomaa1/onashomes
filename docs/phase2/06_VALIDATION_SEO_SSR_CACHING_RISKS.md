# Phase 2 — Validation, SEO, SSR, Caching & Risks

---

## 1. Validation Rules

### Product Validation

| Rule | Layer | Implementation |
|------|-------|----------------|
| `name_en` required | Controller (binding) + Service | `binding:"required"` |
| `name_ar` required | Controller (binding) + Service | `binding:"required"` |
| Only one attribute type per product (`size` OR `color` OR `nil`) | **Service** | Validate `attribute_type ∈ {nil, "size", "color"}` |
| All variants follow same attribute type | **Service** | When creating variant, check product's `attribute_type`; if nil, reject variant with attribute_value |
| SKU unique globally | **DB** (unique index) + Service pre-check | `ux_product_variants_sku` |
| Slug unique per store | **Service** | Query `product_storefront` JOIN `products` for duplicates per `store_front_id` |
| Cannot publish without ≥1 active variant | **Service** | `UpdateProductStatus` checks variant count |
| Cannot activate if no inventory | **Service** | Check `variant_inventory` has records with `quantity > 0` |
| Cannot delete if has orders | **Service** | Check `order_items` table (future-proofed with error message) |
| Barcode unique (when not null) | **DB** partial unique index | `ux_product_variants_barcode` |
| `meta_title` 50-60 chars | **Service** | Warn if outside range, allow override |
| `meta_description` 140-160 chars | **Service** | Warn if outside range, allow override |
| Inventory cannot go negative | **DB** CHECK constraint + **Service** | `CHECK (quantity >= 0)` |
| Adjustment reason required | Controller binding | `binding:"required,oneof=adjustment correction restock"` |

### Attribute Rule Enforcement (Critical)

```go
// In product service:
func validateAttributeRule(product *Product, variant *CreateVariantV2Request) error {
    if product.AttributeType == nil {
        // Simple product — variant MUST NOT have attribute_value
        if variant.AttributeValue != "" {
            return errors.New("simple product cannot have attribute values")
        }
        return nil
    }

    // Product has attribute_type — variant MUST have attribute_value
    if variant.AttributeValue == "" {
        return fmt.Errorf("variant must specify attribute_value for %s product", *product.AttributeType)
    }

    // Verify no cross-type mixing:
    // All existing variants must have same attribute_type (enforced at product level)
    return nil
}
```

### Status Transition Matrix

```
┌─────────┐     requires ≥1 active    ┌──────────┐
│  draft   │ ──── variant + inventory ─→ │  active  │
└─────────┘                             └──────────┘
     ↑                                    │       │
     │  re-open                           │       │
     │                                    ▼       ▼
┌──────────┐                         ┌──────────┐
│ archived │ ←───────────────────────│ inactive │
└──────────┘                         └──────────┘
```

---

## 2. SEO Implementation

### Auto-Suggest Logic

```go
func (s *service) SuggestSEO(ctx context.Context, productID int64) utils.IResource {
    product, _ := s.repo.GetProductByID(ctx, productID)

    suggested := ProductSEORequest{
        MetaTitleEn:      truncate(product.NameEn + " | " + product.Brand.NameEn, 60),
        MetaTitleAr:      truncate(product.NameAr + " | " + product.Brand.NameAr, 60),
        MetaDescriptionEn: truncate("Shop " + product.NameEn + " from " +
            product.Brand.NameEn + " in " + product.Category.NameEn + ". " +
            product.DescriptionEn, 160),
        MetaDescriptionAr: truncate("تسوق " + product.NameAr + " من " +
            product.Brand.NameAr + " في " + product.Category.NameAr + ". " +
            product.DescriptionAr, 160),
        MetaKeywords: strings.Join([]string{
            product.NameEn, product.Brand.NameEn,
            product.Category.NameEn,
        }, ","),
        CanonicalURL: fmt.Sprintf("/products/%s", product.Slug),
        OgTitle:      product.NameEn,
        OgDescription: truncate(product.DescriptionEn, 200),
    }

    return utils.NewOKResource("SEO suggestions generated", suggested)
}
```

### Slug Generation

```go
func Slugify(input string) string {
    s := strings.ToLower(input)
    s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
    s = regexp.MustCompile(`[\s]+`).ReplaceAllString(s, "-")
    s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}
```

### Structured Data Generation

The `/structured-data` endpoint returns JSON-LD for React SSR injection into `<script type="application/ld+json">`. Includes:
- **Product schema** with name, description, brand, SKU, images
- **AggregateOffer** with low/highPrice from variants, currency from store
- **BreadcrumbList** from category hierarchy

---

## 3. SSR Store Resolution

### StoreFront Resolver Middleware

```go
// internal/middleware/storefront.go

func StoreFrontResolver() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract domain from Host header
        host := c.Request.Host
        // Strip port if present
        if idx := strings.Index(host, ":"); idx != -1 {
            host = host[:idx]
        }

        // 2. Resolve store from DB (cache this in Redis later)
        db := database.GetDB()
        var storeFront models.StoreFront
        if err := db.Where("domain = ? AND is_active = true", host).
            First(&storeFront).Error; err != nil {
            utils.ErrorResponse(c, 404, "Store not found", nil)
            c.Abort()
            return
        }

        // 3. Set store context
        c.Set("store_front_id", storeFront.ID)
        c.Set("store_front", storeFront)
        c.Set("store_currency", storeFront.Currency)
        c.Set("store_language", storeFront.DefaultLanguage)

        c.Next()
    }
}
```

### React SSR App Architecture (per store)

```
storefront-app/
├── src/
│   ├── pages/
│   │   ├── Home.tsx           ← featured + new products
│   │   ├── ProductList.tsx    ← category browsing
│   │   ├── ProductDetail.tsx  ← SEO meta + JSON-LD injection
│   │   └── CategoryList.tsx
│   ├── components/
│   │   ├── SEOHead.tsx        ← <Helmet> with meta tags
│   │   ├── StructuredData.tsx ← JSON-LD script injection
│   │   └── ProductCard.tsx
│   ├── services/
│   │   └── api.ts             ← calls /api/storefront/*
│   └── utils/
│       └── seo.ts             ← canonical URL helpers
```

Each store resolves via:
1. DNS → same server (reverse proxy by domain)
2. Middleware resolves `store_front_id` from `Host` header
3. All queries scoped to that store

---

## 4. Caching Strategy (Redis-Ready)

| Cache Key Pattern | TTL | Invalidation |
|-------------------|-----|--------------|
| `sf:{domain}` | 1h | On store update |
| `products:sf:{id}:list:{hash}` | 5m | On product create/update in store |
| `product:sf:{id}:slug:{slug}` | 10m | On product update |
| `product:{id}:seo` | 30m | On SEO update |
| `product:{id}:jsonld` | 30m | On product/SEO update |
| `inventory:{variant}:{store}` | 30s | On inventory adjust |
| `lowstock:{store}` | 2m | On inventory adjust |

**Current phase:** No Redis. Queries hit DB directly. Schema and service contracts are structured so adding a cache layer later requires zero API changes — only repository method wrappers.

---

## 5. Admin Dashboard Requirements

### Product List Page
- Filterable table with columns: Name, Status, Brand, Category, Store(s), Stock Status, Price Range, Created
- Filters sidebar: store_front, category, brand, supplier, status, attribute_type, stock_status, date_range
- Bulk actions: change status, assign to stores
- Sort by: name, created_at, price, inventory_quantity, sales (future)

### Product Detail Page
- Tabbed interface: Basic Info | Variants | Inventory | SEO | Store Assignment
- Inline variant management with pricing fields
- SEO tab with auto-suggest button and live preview (title length indicator)
- Status badge with transition actions

### Inventory Dashboard
- Store selector dropdown
- List of variant inventories per store with: SKU, Product Name, Qty, Reserved, Available, Threshold, Status
- Low stock alert badge
- Inline adjust button (opens modal: quantity +/-, reason, notes)
- Adjustment history expandable rows

### StoreFront Manager
- CRUD for store fronts
- Domain + slug configuration
- Active/inactive toggle

---

## 6. Risks & Trade-offs

| Risk | Mitigation | Trade-off |
|------|-----------|-----------|
| **Slug uniqueness per store** enforced at app level, not DB | Service-layer check inside transaction with `SELECT ... FOR UPDATE` | Slight race condition window vs clean DB constraint |
| **No Redis yet** — storefront queries hit DB | Proper indexing, query optimization, future Redis layer designed in | Slightly higher latency for storefront pages |
| **Monolith** — all modules in single binary | Modular package structure, interface-driven. Can extract to microservices later | Simpler deployment vs eventual scale limits |
| **Legacy `name`/`price` columns** on products | Keep during migration, backfill `name_en`, deprecate over 2 releases | Duplication of data temporarily |
| **Single attribute type per product** (not multi-attribute) | Simplified variant model matches business rule. Future: add multi-attribute matrix | Doesn't support "Size + Color" combos. By design per requirements |
| **Inventory per store** — no multi-warehouse | `variant_inventory` has `store_front_id`. Future: add `warehouse_id` FK | Extra column later; current schema supports it |
| **SEO in separate table** vs JSONB | Clean schema, indexable, type-safe. Extra JOIN is trivial on 1:1 | Cannot add arbitrary SEO fields without migration |
| **Permission granularity** — `inventory.adjust` covers both add/subtract | Fine for current team size. Future: split into `inventory.restock` and `inventory.deduct` if needed | Single permission for simplicity |

---

## 7. Future-Proofing Checklist

| Future Feature | How Current Design Supports It |
|---------------|-------------------------------|
| **Multi-warehouse** | Add `warehouse_id` to `variant_inventory`. Current index handles it |
| **Promotions engine** | `compare_at_price` already on variants. Add `promotions` table with rules referencing products/variants |
| **Product reviews** | Add `product_reviews` table with `product_id` FK. No schema changes needed |
| **Multi-currency** | `store_fronts.currency` already exists. Add currency conversion table + middleware |
| **Event-driven stock deduction** | `inventory_adjustments` audit trail ready. Add message queue producer in service layer |
| **Order system** | Products/variants have stable IDs. Add `orders` + `order_items` with `product_variant_id` FK |
| **Search/Elasticsearch** | Products have `name_en`, `name_ar`, `slug`, `status` — all indexable. Add async sync worker |

---

## 8. Implementation Order (Recommended)

1. **Database migration** — Run the SQL migration first
2. **Models** — Create new GORM model files
3. **StoreFront module** — CRUD + domain resolution middleware
4. **Product V2** — Enhanced create/update with store assignment
5. **Variant V2** — Pricing fields, attribute_value
6. **Inventory module** — Adjust + audit trail
7. **SEO module** — Upsert + auto-suggest + structured data
8. **Storefront public API** — Domain-resolved product list/detail
9. **Permissions** — Add new permissions to scanner
10. **Admin dashboard** — Build Angular components for each module
