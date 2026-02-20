# Phase 2 — Multi-Store E-Commerce Platform Design

> **Architecture Team Design Document**
> Date: 2026-02-12 | Version: 1.0
> Based on existing: Go 1.23, GIN, GORM, PostgreSQL, RBAC

---

## Quick Navigation

| Document | Contents |
|----------|----------|
| [01_DATABASE_SCHEMA.md](./01_DATABASE_SCHEMA.md) | ER diagram, full SQL migration, indexing strategy, normalization decisions |
| [02_MODELS.md](./02_MODELS.md) | All GORM models (StoreFront, Product, Variant, Inventory, SEO, Pivots) |
| [03_REPOSITORIES_SERVICES.md](./03_REPOSITORIES_SERVICES.md) | Repository interfaces, service contracts, request DTOs, filter structs |
| [04_CONTROLLERS_ROUTES_PERMISSIONS.md](./04_CONTROLLERS_ROUTES_PERMISSIONS.md) | Controller endpoints, route registration, permission matrix |
| [05_API_CONTRACTS.md](./05_API_CONTRACTS.md) | Full example request/response payloads for every endpoint |
| [06_VALIDATION_SEO_SSR_CACHING_RISKS.md](./06_VALIDATION_SEO_SSR_CACHING_RISKS.md) | Validation rules, SEO, SSR, caching, admin dashboard, risks |

---

## Architecture Summary

### New Modules

```
internal/
├── api/
│   ├── products/         ← ENHANCED (V2 endpoints + storefront queries)
│   ├── inventory/        ← NEW
│   │   ├── controller.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── routes.go
│   │   └── requests/
│   ├── storefronts/      ← NEW
│   │   ├── controller.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── routes.go
│   │   └── requests/
│   └── ... (existing modules unchanged)
├── middleware/
│   ├── storefront.go     ← NEW (domain-based store resolver)
│   └── ... (existing)
├── models/
│   ├── store_front.go    ← NEW
│   ├── product.go        ← NEW (replaces inline structs)
│   ├── product_variant.go ← NEW
│   ├── variant_inventory.go ← NEW
│   ├── inventory_adjustment.go ← NEW
│   ├── product_seo.go    ← NEW
│   ├── pivots.go         ← NEW (M2M pivot models)
│   └── ... (existing unchanged)
```

### New Database Tables (7)

| Table | Type | Purpose |
|-------|------|---------|
| `store_fronts` | Core | Multi-store configuration |
| `product_storefront` | Pivot | Product ↔ Store M2M |
| `category_storefront` | Pivot | Category ↔ Store M2M |
| `brand_storefront` | Pivot | Brand ↔ Store M2M |
| `supplier_storefront` | Pivot | Supplier ↔ Store M2M |
| `variant_inventory` | Core | Per-store, per-variant inventory |
| `inventory_adjustments` | Audit | Inventory change history |
| `product_seo` | Core | SEO metadata (1:1 with product) |

### Altered Tables (2)

| Table | Changes |
|-------|---------|
| `products` | +14 columns (name_en/ar, slug, brand_id, category_id, supplier_id, status, attribute_type, marketing flags) |
| `product_variants` | +5 columns (price, compare_at_price, cost_price, barcode, weight, attribute_value) |

### New Permissions (4)

| Permission | Description |
|-----------|-------------|
| `inventory.adjust` | Adjust product inventory quantities |
| `inventory.view` | View inventory levels |
| `storefronts.manage` | Full CRUD on store fronts |
| `seo.manage` | Manage product SEO metadata |

### Key Business Rules

1. **Attribute exclusivity**: A product can have `size` OR `color` OR no attribute — **never both**
2. **Slug uniqueness**: Per-store, not globally
3. **Inventory is variant-based AND store-scoped**: Not product-level
4. **Negative stock prevention**: DB CHECK constraint + service validation
5. **Status transitions are guarded**: Cannot activate without variant + inventory

### Implementation Order

1. Database migration → 2. Models → 3. StoreFront CRUD → 4. Product V2 → 5. Variant V2 → 6. Inventory → 7. SEO → 8. Storefront public API → 9. Permissions → 10. Admin dashboard
