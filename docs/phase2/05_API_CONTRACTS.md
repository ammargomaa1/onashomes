# Phase 2 — API Contracts (Example Payloads)

## 1. Create Product V2

**`POST /api/admin/products/v2`**

### Request
```json
{
  "name_en": "Premium Cotton T-Shirt",
  "name_ar": "تيشيرت قطن فاخر",
  "description_en": "Soft premium cotton t-shirt with modern fit",
  "description_ar": "تيشيرت قطن ناعم بقصة عصرية",
  "brand_id": 5,
  "category_id": 12,
  "supplier_id": null,
  "is_internal_supplier": true,
  "attribute_type": "size",
  "store_front_ids": [1, 3],
  "is_featured": true,
  "is_new": true,
  "seo": {
    "meta_title_en": "Premium Cotton T-Shirt | Onas Homes",
    "meta_title_ar": "تيشيرت قطن فاخر | اوناس هومز",
    "meta_description_en": "Shop our premium cotton t-shirt collection. Soft, breathable, and stylish.",
    "meta_description_ar": "تسوق مجموعة تيشيرتات القطن الفاخرة. ناعمة ومريحة وأنيقة.",
    "meta_keywords": "cotton,t-shirt,premium,men",
    "og_title": "Premium Cotton T-Shirt",
    "og_description": "Soft premium cotton t-shirt with modern fit"
  }
}
```

### Response — `201 Created`
```json
{
  "status": "success",
  "message": "Product created successfully",
  "data": {
    "id": 42,
    "name_en": "Premium Cotton T-Shirt",
    "name_ar": "تيشيرت قطن فاخر",
    "slug": "premium-cotton-t-shirt",
    "status": "draft",
    "is_published": false,
    "attribute_type": "size",
    "brand": { "id": 5, "name_en": "BrandX" },
    "category": { "id": 12, "name_en": "T-Shirts" },
    "store_fronts": [
      { "id": 1, "name": "Store 1", "domain": "store1.com" },
      { "id": 3, "name": "Store 3", "domain": "store3.com" }
    ],
    "variants": [],
    "seo": {
      "meta_title_en": "Premium Cotton T-Shirt | Onas Homes",
      "meta_description_en": "Shop our premium cotton t-shirt collection..."
    }
  }
}
```

---

## 2. Create Variant V2

**`POST /api/admin/products/42/variants/v2`**

### Request
```json
{
  "sku": "PCOT-TS-XL-001",
  "attribute_value": "XL",
  "price": 149.99,
  "compare_at_price": 199.99,
  "cost_price": 45.00,
  "barcode": "6281234567890",
  "weight": 0.250,
  "is_active": true,
  "image_file_ids": [101, 102]
}
```

### Response — `201 Created`
```json
{
  "status": "success",
  "message": "Variant created successfully",
  "data": {
    "id": 88,
    "sku": "PCOT-TS-XL-001",
    "attribute_value": "XL",
    "price": 149.99,
    "compare_at_price": 199.99,
    "cost_price": 45.00,
    "barcode": "6281234567890",
    "weight": 0.250,
    "is_active": true,
    "images": [
      { "file_id": 101, "position": 0 },
      { "file_id": 102, "position": 1 }
    ]
  }
}
```

---

## 3. Adjust Inventory

**`POST /api/admin/inventory/adjust`**

### Request
```json
{
  "product_variant_id": 88,
  "store_front_id": 1,
  "adjustment": 50,
  "reason": "restock",
  "notes": "New shipment received - PO#12345"
}
```

### Response — `200 OK`
```json
{
  "status": "success",
  "message": "Inventory adjusted successfully",
  "data": {
    "id": 15,
    "product_variant_id": 88,
    "store_front_id": 1,
    "quantity": 50,
    "reserved_quantity": 0,
    "available_quantity": 50,
    "low_stock_threshold": 5,
    "is_low_stock": false,
    "adjustment": {
      "previous_quantity": 0,
      "new_quantity": 50,
      "adjustment_amount": 50,
      "reason": "restock"
    }
  }
}
```

---

## 4. Change Product Status

**`PATCH /api/admin/products/42/status`**

### Request
```json
{ "status": "active" }
```

### Error Response — `400 Bad Request` (if no active variant)
```json
{
  "status": "error",
  "message": "Cannot activate product: requires at least 1 active variant with inventory"
}
```

---

## 5. Storefront Product List (Public)

**`GET /api/storefront/products?category_id=12&sort=newest&page=1&limit=20`**

*Header:* `Host: store1.com` (resolved by StoreFrontResolver middleware)

### Response — `200 OK`
```json
{
  "status": "success",
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 42,
      "name_en": "Premium Cotton T-Shirt",
      "name_ar": "تيشيرت قطن فاخر",
      "slug": "premium-cotton-t-shirt",
      "brand": { "id": 5, "name_en": "BrandX" },
      "category": { "id": 12, "name_en": "T-Shirts" },
      "is_featured": true,
      "is_new": true,
      "min_price": 129.99,
      "max_price": 149.99,
      "compare_at_price": 199.99,
      "thumbnail": "/storage/files/101.webp",
      "in_stock": true,
      "variant_count": 3
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

## 6. Storefront Product Detail (Public)

**`GET /api/storefront/products/premium-cotton-t-shirt`**

### Response — `200 OK`
```json
{
  "status": "success",
  "data": {
    "id": 42,
    "name_en": "Premium Cotton T-Shirt",
    "name_ar": "تيشيرت قطن فاخر",
    "slug": "premium-cotton-t-shirt",
    "description_en": "Soft premium cotton t-shirt with modern fit",
    "description_ar": "تيشيرت قطن ناعم بقصة عصرية",
    "brand": { "id": 5, "name_en": "BrandX" },
    "category": { "id": 12, "name_en": "T-Shirts" },
    "attribute_type": "size",
    "is_featured": true,
    "seo": {
      "meta_title_en": "Premium Cotton T-Shirt | Onas Homes",
      "meta_description_en": "Shop our premium cotton t-shirt collection.",
      "canonical_url": "https://store1.com/products/premium-cotton-t-shirt",
      "og_title": "Premium Cotton T-Shirt",
      "og_image": "/storage/files/101.webp"
    },
    "variants": [
      {
        "id": 86,
        "sku": "PCOT-TS-M-001",
        "attribute_value": "M",
        "price": 129.99,
        "compare_at_price": 199.99,
        "in_stock": true,
        "images": [{ "file_id": 103, "url": "/storage/files/103.webp", "position": 0 }]
      },
      {
        "id": 88,
        "sku": "PCOT-TS-XL-001",
        "attribute_value": "XL",
        "price": 149.99,
        "compare_at_price": 199.99,
        "in_stock": true,
        "images": [{ "file_id": 101, "url": "/storage/files/101.webp", "position": 0 }]
      }
    ]
  }
}
```

---

## 7. Create StoreFront

**`POST /api/admin/storefronts`**

### Request
```json
{
  "name": "Onas Saudi Arabia",
  "slug": "onas-sa",
  "domain": "sa.onashomes.com",
  "currency": "SAR",
  "default_language": "ar"
}
```

### Response — `201 Created`
```json
{
  "status": "success",
  "message": "Store front created successfully",
  "data": {
    "id": 1,
    "name": "Onas Saudi Arabia",
    "slug": "onas-sa",
    "domain": "sa.onashomes.com",
    "currency": "SAR",
    "default_language": "ar",
    "is_active": true
  }
}
```

---

## 8. JSON-LD Structured Data

**`GET /api/storefront/products/premium-cotton-t-shirt/structured-data`**

### Response — `200 OK`
```json
{
  "@context": "https://schema.org",
  "@type": "Product",
  "name": "Premium Cotton T-Shirt",
  "description": "Soft premium cotton t-shirt with modern fit",
  "image": "https://store1.com/storage/files/101.webp",
  "brand": { "@type": "Brand", "name": "BrandX" },
  "sku": "PCOT-TS-XL-001",
  "offers": {
    "@type": "AggregateOffer",
    "lowPrice": 129.99,
    "highPrice": 149.99,
    "priceCurrency": "SAR",
    "availability": "https://schema.org/InStock",
    "offerCount": 3
  },
  "breadcrumb": {
    "@type": "BreadcrumbList",
    "itemListElement": [
      { "@type": "ListItem", "position": 1, "name": "Home", "item": "https://store1.com" },
      { "@type": "ListItem", "position": 2, "name": "T-Shirts", "item": "https://store1.com/categories/t-shirts" },
      { "@type": "ListItem", "position": 3, "name": "Premium Cotton T-Shirt" }
    ]
  }
}
```
