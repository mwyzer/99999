# API Contract — MVP

## Multi-Tenant Property Information System

| Property          | Value                                                  |
| ----------------- | ------------------------------------------------------ |
| **Document Type** | API Contract / OpenAPI-Like Specification              |
| **Version**       | 1.0.0 MVP                                              |
| **Date**          | 2026-06-26                                             |
| **Base URL**      | `http://localhost:8080/api/v1`                         |
| **Reference**     | `02-SRS-MVP.md` Sec 7, `05-Error-Handling-Standard.md` |

---

## 1. API Overview

### 1.1 General

| Property           | Value                                           |
| ------------------ | ----------------------------------------------- |
| **Protocol**       | HTTP/1.1 + HTTPS (production)                   |
| **Base Path**      | `/api/v1`                                       |
| **Content-Type**   | `application/json`                              |
| **Charset**        | UTF-8                                           |
| **Auth Method**    | JWT Bearer Token                                |
| **Date Format**    | ISO 8601 (`2026-06-26T14:30:00+07:00`)          |
| **Decimal Format** | String (untuk menghindari floating-point error) |

### 1.2 Common Headers

#### Request Headers

| Header          | Required        | Value                |
| --------------- | --------------- | -------------------- |
| `Content-Type`  | Yes (POST/PUT)  | `application/json`   |
| `Authorization` | Yes (protected) | `Bearer <jwt_token>` |
| `Accept`        | Optional        | `application/json`   |

#### Response Headers

| Header                   | Value                        |
| ------------------------ | ---------------------------- |
| `Content-Type`           | `application/json`           |
| `X-Request-Id`           | UUID v4 (unique per request) |
| `X-Content-Type-Options` | `nosniff`                    |
| `X-Frame-Options`        | `DENY`                       |
| `X-XSS-Protection`       | `1; mode=block`              |

### 1.3 Authentication

Semua endpoint kecuali **Public Endpoints** memerlukan JWT Bearer token.

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

JWT Payload:

```json
{
  "sub": "uuid-user-id",
  "role": "buyer|salesman|tenant_admin|platform_admin",
  "tenant_id": "uuid-tenant-id|null",
  "exp": 1719440000,
  "iat": 1719353600
}
```

### 1.4 Endpoint Role Mapping

| Prefix                 | Role Required                       | Tenant Scope |
| ---------------------- | ----------------------------------- | ------------ |
| `/api/v1/auth/*`       | None (public)                       | –            |
| `/api/v1/properties/*` | None (public)                       | –            |
| `/api/v1/locations/*`  | None (public)                       | –            |
| `/api/v1/me/*`         | `buyer`, `salesman`, `tenant_admin` | Self         |
| `/api/v1/salesman/*`   | `salesman`, `tenant_admin`          | Tenant       |
| `/api/v1/tenant/*`     | `tenant_admin`                      | Tenant       |
| `/api/v1/admin/*`      | `platform_admin`                    | Global       |

---

## 2. Generic Response Patterns

### 2.1 Success Response

```json
{
  "success": true,
  "data": {},
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### 2.2 Success Response (No Pagination)

```json
{
  "success": true,
  "data": {}
}
```

### 2.3 Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message in Bahasa Indonesia",
    "details": []
  }
}
```

### 2.4 Pagination

Semua endpoint `GET` list mendukung pagination:

| Query Param | Type    | Default | Max | Deskripsi        |
| ----------- | ------- | ------- | --- | ---------------- |
| `page`      | integer | 1       | –   | Halaman saat ini |
| `per_page`  | integer | 20      | 100 | Item per halaman |

Response selalu menyertakan `meta`:

```json
{
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

## 3. Public Endpoints

### 3.1 Register Buyer

```http
POST /api/v1/auth/register
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Rina Wijaya",
  "email": "rina@email.com",
  "phone": "081300000001",
  "password": "Rina@123"
}
```

| Field      | Type   | Required | Rules                                          |
| ---------- | ------ | -------- | ---------------------------------------------- |
| `name`     | string | ✅       | 2–100 chars, trimmed                           |
| `email`    | string | ✅       | Valid email, unique, max 255 chars, lowercased |
| `phone`    | string | ✅       | Regex `^\+?[0-9]{8,15}$`                       |
| `password` | string | ✅       | Min 8 chars                                    |

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "id": "f0000000-0000-0000-0000-000000000001",
    "name": "Rina Wijaya",
    "email": "rina@email.com",
    "phone": "081300000001",
    "role": "buyer",
    "status": "active",
    "created_at": "2026-06-26T10:00:00+07:00"
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `VAL_INPUT_INVALID` | 422 | Field validasi gagal |
| `AUTH_EMAIL_REGISTERED` | 409 | Email sudah terdaftar |

---

### 3.2 Login

```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request Body:**

```json
{
  "email": "rina@email.com",
  "password": "Rina@123"
}
```

| Field      | Type   | Required | Rules       |
| ---------- | ------ | -------- | ----------- |
| `email`    | string | ✅       | Valid email |
| `password` | string | ✅       | –           |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "f0000000-0000-0000-0000-000000000001",
      "name": "Rina Wijaya",
      "email": "rina@email.com",
      "phone": "081300000001",
      "photo_url": null,
      "role": "buyer",
      "tenant_id": null,
      "tenant_name": null
    }
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `AUTH_INVALID_CREDENTIALS` | 401 | Email atau password salah |
| `AUTH_ACCOUNT_SUSPENDED` | 403 | Tenant user di-suspend |
| `AUTH_ACCOUNT_INACTIVE` | 403 | User tidak aktif |
| `RATE_LOGIN_LIMIT` | 429 | 5+ attempt dalam 1 menit |

---

### 3.3 List Properties (Public)

```http
GET /api/v1/properties?page=1&per_page=20&property_type=house&city=Jakarta+Selatan&price_min=500000000&price_max=5000000000&listing_type=sale&source_type=regular&sort=created_at&order=desc
```

**Query Parameters:**

| Param           | Type    | Required | Default      | Deskripsi                                                                         |
| --------------- | ------- | -------- | ------------ | --------------------------------------------------------------------------------- |
| `page`          | integer | ❌       | 1            | Halaman                                                                           |
| `per_page`      | integer | ❌       | 20           | Item per halaman (max 100)                                                        |
| `property_type` | string  | ❌       | –            | Filter: `house`, `land`, `apartment`, `shophouse`, `warehouse`, `office`, `villa` |
| `source_type`   | string  | ❌       | –            | Filter: `regular`, `bank_auction`, `company_asset`                                |
| `listing_type`  | string  | ❌       | –            | Filter: `sale`, `rent`                                                            |
| `city`          | string  | ❌       | –            | Filter by city (case-insensitive, ILIKE)                                          |
| `price_min`     | integer | ❌       | –            | Harga minimum (IDR)                                                               |
| `price_max`     | integer | ❌       | –            | Harga maksimum (IDR)                                                              |
| `q`             | string  | ❌       | –            | Text search pada title + description                                              |
| `sort`          | string  | ❌       | `created_at` | Sort by: `created_at`, `price`                                                    |
| `order`         | string  | ❌       | `desc`       | Order: `asc`, `desc`                                                              |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000001",
      "title": "Rumah Minimalis 2 Lantai di Kebayoran Baru",
      "price": "2500000000.00",
      "listing_type": "sale",
      "property_type": "house",
      "source_type": "regular",
      "city": "Jakarta Selatan",
      "province": "DKI Jakarta",
      "land_area": "200.00",
      "building_area": "150.00",
      "bedrooms": 4,
      "bathrooms": 3,
      "status": "approved",
      "main_photo_url": "/uploads/listings/1000.../abc123_medium.jpg",
      "salesman": {
        "id": "e0000000-0000-0000-0000-000000000001",
        "name": "Andi Pratama",
        "photo_url": null,
        "phone": "081200000002"
      },
      "tenant": {
        "id": "b0000000-0000-0000-0000-000000000001",
        "name": "PropertiJaya Agency",
        "logo_url": null
      },
      "created_at": "2026-06-20T10:00:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

**Catatan:** Hanya mengembalikan listing dengan `status = 'approved'`.

---

### 3.4 Get Property Detail (Public)

```http
GET /api/v1/properties/:id
```

**Path Parameters:**

| Param | Type | Deskripsi           |
| ----- | ---- | ------------------- |
| `id`  | UUID | Property listing ID |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000001",
    "title": "Rumah Minimalis 2 Lantai di Kebayoran Baru",
    "description": "Rumah minimalis modern 2 lantai dengan desain elegan...",
    "price": "2500000000.00",
    "listing_type": "sale",
    "property_type": "house",
    "source_type": "regular",
    "address": "Jl. Melati No. 15, Kebayoran Baru",
    "city": "Jakarta Selatan",
    "province": "DKI Jakarta",
    "latitude": "-6.2431000",
    "longitude": "106.7988000",
    "land_area": "200.00",
    "building_area": "150.00",
    "bedrooms": 4,
    "bathrooms": 3,
    "floors": 2,
    "certificate_type": "SHM",
    "facilities": {
      "carport": true,
      "garden": true,
      "security_24h": true,
      "furnished": "semi"
    },
    "status": "approved",
    "photos": [
      {
        "id": "ff000000-0000-0000-0000-000000000001",
        "original_url": "/uploads/listings/1000.../abc123_orig.jpg",
        "thumbnail_url": "/uploads/listings/1000.../abc123_thumb.jpg",
        "medium_url": "/uploads/listings/1000.../abc123_medium.jpg",
        "watermarked_url": "/uploads/listings/1000.../abc123_wm.jpg",
        "sort_order": 0
      }
    ],
    "salesman": {
      "id": "e0000000-0000-0000-0000-000000000001",
      "name": "Andi Pratama",
      "photo_url": null,
      "phone": "081200000002"
    },
    "tenant": {
      "id": "b0000000-0000-0000-0000-000000000001",
      "name": "PropertiJaya Agency",
      "logo_url": null,
      "phone": "0215551234"
    },
    "created_at": "2026-06-20T10:00:00+07:00"
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `RES_LISTING_NOT_FOUND` | 404 | Listing tidak ditemukan atau bukan `approved` |

---

### 3.5 Featured Properties

```http
GET /api/v1/properties/featured?city=Jakarta+Selatan&limit=6
```

**Query Parameters:**

| Param   | Type    | Required | Default | Deskripsi            |
| ------- | ------- | -------- | ------- | -------------------- |
| `city`  | string  | ❌       | –       | Filter by city       |
| `limit` | integer | ❌       | 6       | Jumlah item (max 12) |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "location": "Jakarta Selatan",
    "properties": [
      /* ... same structure as list ... */
    ]
  }
}
```

---

### 3.6 Nearby Properties

```http
GET /api/v1/properties/nearby?latitude=-6.2431&longitude=106.7988&radius_km=5&limit=10
```

**Query Parameters:**

| Param       | Type    | Required | Default | Deskripsi                 |
| ----------- | ------- | -------- | ------- | ------------------------- |
| `latitude`  | number  | ✅       | –       | Latitude (-90 s.d. 90)    |
| `longitude` | number  | ✅       | –       | Longitude (-180 s.d. 180) |
| `radius_km` | number  | ❌       | 10      | Radius dalam km (1–50)    |
| `limit`     | integer | ❌       | 10      | Jumlah item (max 20)      |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000001",
      "title": "Rumah Minimalis 2 Lantai di Kebayoran Baru",
      "price": "2500000000.00",
      "listing_type": "sale",
      "property_type": "house",
      "city": "Jakarta Selatan",
      "distance_km": 2.3,
      "main_photo_url": "/uploads/listings/1000.../abc123_medium.jpg",
      "salesman": {
        "name": "Andi Pratama",
        "phone": "081200000002"
      },
      "tenant": {
        "name": "PropertiJaya Agency",
        "logo_url": null
      }
    }
  ]
}
```

---

### 3.7 List Cities

```http
GET /api/v1/locations/cities
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    { "city": "Jakarta Selatan", "province": "DKI Jakarta" },
    { "city": "Jakarta Pusat", "province": "DKI Jakarta" },
    { "city": "Bandung", "province": "Jawa Barat" },
    { "city": "Tangerang", "province": "Banten" },
    { "city": "Depok", "province": "Jawa Barat" }
  ]
}
```

> Data diambil dari `SELECT DISTINCT city, province FROM property_listings WHERE status = 'approved' ORDER BY city`.

---

## 4. Authenticated Common Endpoints

Berlaku untuk semua role: `buyer`, `salesman`, `tenant_admin`, `platform_admin`.

### 4.1 Get My Profile

```http
GET /api/v1/me/profile
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "e0000000-0000-0000-0000-000000000001",
    "name": "Andi Pratama",
    "email": "andi@propertijaya.id",
    "phone": "081200000002",
    "photo_url": null,
    "role": "salesman",
    "status": "active",
    "tenant": {
      "id": "b0000000-0000-0000-0000-000000000001",
      "name": "PropertiJaya Agency",
      "logo_url": null
    },
    "created_at": "2026-06-20T09:00:00+07:00"
  }
}
```

---

### 4.2 Update My Profile

```http
PUT /api/v1/me/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Andi Pratama Updated",
  "phone": "081200000099"
}
```

| Field       | Type   | Required | Rules                         |
| ----------- | ------ | -------- | ----------------------------- |
| `name`      | string | ❌       | 2–100 chars                   |
| `phone`     | string | ❌       | Regex `^\+?[0-9]{8,15}$`      |
| `photo_url` | string | ❌       | – (diisi setelah upload foto) |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "e0000000-0000-0000-0000-000000000001",
    "name": "Andi Pratama Updated",
    "email": "andi@propertijaya.id",
    "phone": "081200000099",
    "photo_url": null,
    "role": "salesman",
    "status": "active"
  }
}
```

---

## 5. Buyer Endpoints

**Role:** `buyer` | **Auth:** JWT

### 5.1 List Saved Properties

```http
GET /api/v1/me/saved?page=1&per_page=20
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000001",
      "title": "Rumah Minimalis 2 Lantai di Kebayoran Baru",
      "price": "2500000000.00",
      "listing_type": "sale",
      "property_type": "house",
      "city": "Jakarta Selatan",
      "main_photo_url": "/uploads/listings/1000.../abc123_medium.jpg",
      "salesman": {
        "name": "Andi Pratama",
        "phone": "081200000002"
      },
      "tenant": {
        "name": "PropertiJaya Agency",
        "logo_url": null
      },
      "saved_at": "2026-06-25T15:30:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### 5.2 Save / Bookmark Property

```http
POST /api/v1/me/saved/:propertyId
Authorization: Bearer <token>
```

**Path Parameters:**

| Param        | Type | Deskripsi           |
| ------------ | ---- | ------------------- |
| `propertyId` | UUID | Property listing ID |

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "message": "Properti berhasil disimpan ke favorit."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `RES_LISTING_NOT_FOUND` | 404 | Listing tidak ditemukan atau bukan `approved` |
| `RES_ALREADY_SAVED` | 409 | Sudah di-save sebelumnya |

---

### 5.3 Remove Saved Property

```http
DELETE /api/v1/me/saved/:propertyId
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Properti berhasil dihapus dari favorit."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `RES_NOT_FOUND` | 404 | Tidak ada saved record |

---

## 6. Salesman Endpoints

**Role:** `salesman`, `tenant_admin` | **Auth:** JWT | **Scope:** Salesman-scoped (sendiri) atau tenant-scoped (tenant_admin)

### 6.1 Salesman Dashboard

```http
GET /api/v1/salesman/dashboard
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "total_listings": 8,
    "status_breakdown": {
      "draft": 1,
      "pending": 1,
      "approved": 3,
      "rejected": 1,
      "sold": 1,
      "rented": 0,
      "inactive": 1
    },
    "active_count": 5,
    "quota": {
      "used": 5,
      "max": 5,
      "remaining": 0,
      "plan_type": "free"
    }
  }
}
```

**Catatan:** `active_count` = `draft` + `pending` + `approved`. Untuk `tenant_admin`, data aggregasi seluruh salesman dalam tenant.

---

### 6.2 List My Listings

```http
GET /api/v1/salesman/listings?status=draft&page=1&per_page=20
Authorization: Bearer <token>
```

**Query Parameters:**

| Param      | Type    | Required | Default | Deskripsi                                                                        |
| ---------- | ------- | -------- | ------- | -------------------------------------------------------------------------------- |
| `status`   | string  | ❌       | –       | Filter: `draft`, `pending`, `approved`, `rejected`, `sold`, `rented`, `inactive` |
| `page`     | integer | ❌       | 1       | Halaman                                                                          |
| `per_page` | integer | ❌       | 20      | Item per halaman                                                                 |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000005",
      "title": "Tanah Kavling Strategis Sentul City",
      "price": "450000000.00",
      "listing_type": "sale",
      "property_type": "land",
      "source_type": "regular",
      "city": "Bogor",
      "status": "draft",
      "main_photo_url": null,
      "created_at": "2026-06-26T08:00:00+07:00",
      "updated_at": null
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 6.3 Create Listing

```http
POST /api/v1/salesman/listings
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "title": "Rumah Minimalis Baru di BSD",
  "description": "Rumah baru siap huni di BSD City...",
  "price": 1500000000,
  "listing_type": "sale",
  "property_type": "house",
  "source_type": "regular",
  "address": "BSD City, Jl. Anggrek No. 20",
  "city": "Tangerang Selatan",
  "province": "Banten",
  "latitude": -6.3015,
  "longitude": 106.6708,
  "land_area": 120,
  "building_area": 80,
  "bedrooms": 3,
  "bathrooms": 2,
  "floors": 1,
  "certificate_type": "SHM",
  "facilities": {
    "carport": true,
    "garden": true,
    "security_24h": true
  }
}
```

| Field              | Type    | Required | Rules                                                                     |
| ------------------ | ------- | -------- | ------------------------------------------------------------------------- |
| `title`            | string  | ✅       | 5–300 chars                                                               |
| `description`      | string  | ❌       | Max 5000 chars                                                            |
| `price`            | number  | ✅       | > 0, max 999999999999                                                     |
| `listing_type`     | string  | ✅       | `sale` / `rent`                                                           |
| `property_type`    | string  | ✅       | `house`, `land`, `apartment`, `shophouse`, `warehouse`, `office`, `villa` |
| `source_type`      | string  | ❌       | `regular` (default), `bank_auction`, `company_asset`                      |
| `address`          | string  | ❌       | Max 500 chars                                                             |
| `city`             | string  | ❌       | 2–100 chars                                                               |
| `province`         | string  | ❌       | 2–100 chars                                                               |
| `latitude`         | number  | ❌       | -90 s.d. 90                                                               |
| `longitude`        | number  | ❌       | -180 s.d. 180                                                             |
| `land_area`        | number  | ❌       | > 0, max 999999                                                           |
| `building_area`    | number  | ❌       | > 0, max 999999                                                           |
| `bedrooms`         | integer | ❌       | 0–99                                                                      |
| `bathrooms`        | integer | ❌       | 0–99                                                                      |
| `floors`           | integer | ❌       | 0–200                                                                     |
| `certificate_type` | string  | ❌       | `SHM`, `SHGB`, `Girik`, `Lainnya`                                         |
| `facilities`       | object  | ❌       | Valid JSON, max 20 keys                                                   |

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000009",
    "title": "Rumah Minimalis Baru di BSD",
    "price": "1500000000.00",
    "listing_type": "sale",
    "property_type": "house",
    "source_type": "regular",
    "city": "Tangerang Selatan",
    "status": "draft",
    "created_at": "2026-06-26T09:00:00+07:00"
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `VAL_INPUT_INVALID` | 422 | Validasi field gagal |
| `BIZ_QUOTA_EXCEEDED` | 422 | Kuota penuh |

---

### 6.4 Get My Listing Detail

```http
GET /api/v1/salesman/listings/:id
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000005",
    "tenant_id": "b0000000-0000-0000-0000-000000000001",
    "salesman_id": "e0000000-0000-0000-0000-000000000001",
    "title": "Tanah Kavling Strategis Sentul City",
    "description": "Tanah kavling siap bangun...",
    "price": "450000000.00",
    "listing_type": "sale",
    "property_type": "land",
    "source_type": "regular",
    "address": "Sentul City, Babakan Madang",
    "city": "Bogor",
    "province": "Jawa Barat",
    "latitude": "-6.5717000",
    "longitude": "106.8618000",
    "land_area": "300.00",
    "building_area": null,
    "bedrooms": null,
    "bathrooms": null,
    "floors": null,
    "certificate_type": "SHM",
    "facilities": {},
    "status": "draft",
    "reject_reason": null,
    "photos": [],
    "created_at": "2026-06-26T08:00:00+07:00",
    "updated_at": null
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `RES_LISTING_NOT_FOUND` | 404 | Listing tidak ditemukan |
| `AUTHZ_NOT_OWNER` | 403 | Bukan listing miliknya (salesman role) |

---

### 6.5 Update Listing

```http
PUT /api/v1/salesman/listings/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Rules:** Hanya bisa edit listing dengan status `draft` atau `rejected`.

**Request Body:** (semua field opsional — partial update)

```json
{
  "title": "Tanah Kavling Premium Sentul City",
  "price": 500000000,
  "description": "Updated description..."
}
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000005",
    "title": "Tanah Kavling Premium Sentul City",
    "price": "500000000.00",
    "status": "draft",
    "updated_at": "2026-06-26T09:30:00+07:00"
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `BIZ_LISTING_NOT_EDITABLE` | 422 | Status bukan `draft`/`rejected` |

---

### 6.6 Delete Listing (Soft)

```http
DELETE /api/v1/salesman/listings/:id
Authorization: Bearer <token>
```

**Rules:** Hanya bisa hapus listing dengan status `draft` atau `rejected`.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Listing berhasil dihapus."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `BIZ_LISTING_NOT_DELETABLE` | 422 | Status bukan `draft`/`rejected` |

---

### 6.7 Submit Listing for Review

```http
POST /api/v1/salesman/listings/:id/submit
Authorization: Bearer <token>
```

**Rules:**

- Hanya listing dengan status `draft` atau `rejected`.
- Quota check: jika kuota penuh, tidak bisa submit.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000005",
    "status": "pending",
    "message": "Listing berhasil diajukan untuk review."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `BIZ_LISTING_NOT_SUBMITTABLE` | 422 | Status bukan `draft`/`rejected` |
| `BIZ_QUOTA_EXCEEDED` | 422 | Kuota penuh |

---

### 6.8 Deactivate Listing

```http
POST /api/v1/salesman/listings/:id/deactivate
Authorization: Bearer <token>
```

**Rules:** Hanya listing dengan status `approved`.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000001",
    "status": "inactive",
    "message": "Listing berhasil dinonaktifkan."
  }
}
```

---

### 6.9 Mark Listing as Sold

```http
POST /api/v1/salesman/listings/:id/mark-sold
Authorization: Bearer <token>
```

**Rules:** Hanya listing dengan status `approved`.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000001",
    "status": "sold",
    "message": "Listing berhasil ditandai sebagai terjual."
  }
}
```

---

### 6.10 Mark Listing as Rented

```http
POST /api/v1/salesman/listings/:id/mark-rented
Authorization: Bearer <token>
```

**Rules:** Hanya listing dengan status `approved`.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000001",
    "status": "rented",
    "message": "Listing berhasil ditandai sebagai tersewa."
  }
}
```

---

### 6.11 Upload Photos

```http
POST /api/v1/salesman/listings/:id/photos
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Request:** Multipart form

| Field    | Type   | Required | Rules                                               |
| -------- | ------ | -------- | --------------------------------------------------- |
| `photos` | file[] | ✅       | Max 10 files, JPEG/PNG/WebP, masing-masing max 5 MB |

**Rules:**

- Max 10 foto per listing (termasuk yang sudah ada).
- Foto diproses: generate thumbnail (400×300), medium (800×600), watermark.
- File disimpan dengan nama UUID random.
- EXIF metadata di-strip.

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "uploaded": 3,
    "photos": [
      {
        "id": "ff000000-0000-0000-0000-000000000011",
        "original_url": "/uploads/listings/1000.../def456_orig.jpg",
        "thumbnail_url": "/uploads/listings/1000.../def456_thumb.jpg",
        "medium_url": "/uploads/listings/1000.../def456_medium.jpg",
        "watermarked_url": "/uploads/listings/1000.../def456_wm.jpg",
        "sort_order": 1
      }
    ]
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `VAL_FILE_TOO_LARGE` | 413 | File > 5 MB |
| `VAL_FILE_INVALID_TYPE` | 415 | Bukan JPEG/PNG/WebP |
| `VAL_FILE_COUNT_EXCEEDED` | 422 | > 10 foto |

---

### 6.12 Delete Photo

```http
DELETE /api/v1/salesman/listings/:id/photos/:photoId
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Foto berhasil dihapus."
  }
}
```

---

### 6.13 Reorder Photos

```http
PUT /api/v1/salesman/listings/:id/photos/reorder
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "photo_ids": [
    "ff000000-0000-0000-0000-000000000013",
    "ff000000-0000-0000-0000-000000000011",
    "ff000000-0000-0000-0000-000000000012"
  ]
}
```

| Field       | Type   | Required | Deskripsi                                         |
| ----------- | ------ | -------- | ------------------------------------------------- |
| `photo_ids` | UUID[] | ✅       | Array ID foto dalam urutan baru. Index 0 = cover. |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Urutan foto berhasil diubah."
  }
}
```

---

### 6.14 Get Quota Usage

```http
GET /api/v1/salesman/quota
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "used": 5,
    "max": 5,
    "remaining": 0,
    "plan_type": "free",
    "status_breakdown": {
      "draft": 1,
      "pending": 1,
      "approved": 3
    }
  }
}
```

---

## 7. Tenant Admin Endpoints

**Role:** `tenant_admin` | **Auth:** JWT | **Scope:** Tenant-scoped

### 7.1 Tenant Dashboard

```http
GET /api/v1/tenant/dashboard
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "total_listings": 25,
    "active_listings": 18,
    "total_salesmen": 3,
    "max_salesmen": 5,
    "status_breakdown": {
      "draft": 2,
      "pending": 3,
      "approved": 15,
      "rejected": 2,
      "sold": 1,
      "rented": 0,
      "inactive": 2
    },
    "plan": {
      "type": "free",
      "max_salesmen": 5,
      "max_listings_per_salesman": 5
    }
  }
}
```

---

### 7.2 Get Tenant Profile

```http
GET /api/v1/tenant/profile
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "b0000000-0000-0000-0000-000000000001",
    "organization_name": "PropertiJaya Agency",
    "subdomain_slug": "propertijaya",
    "logo_url": null,
    "description": "Agensi properti terpercaya sejak 2010...",
    "phone": "0215551234",
    "address": "Jl. Sudirman No. 123, Jakarta Pusat",
    "status": "active",
    "created_at": "2026-06-20T09:00:00+07:00"
  }
}
```

---

### 7.3 Update Tenant Profile

```http
PUT /api/v1/tenant/profile
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "organization_name": "PropertiJaya Agency Updated",
  "description": "Agensi properti #1 di Jabodetabek...",
  "phone": "0215559999",
  "address": "Jl. Sudirman No. 500, Jakarta Pusat"
}
```

| Field               | Type   | Required | Rules                    |
| ------------------- | ------ | -------- | ------------------------ |
| `organization_name` | string | ❌       | 2–200 chars              |
| `description`       | text   | ❌       | Max 2000 chars           |
| `phone`             | string | ❌       | Regex `^\+?[0-9]{8,15}$` |
| `address`           | string | ❌       | Max 500 chars            |

> Upload logo menggunakan endpoint terpisah: `POST /api/v1/tenant/profile/logo` (multipart).

---

### 7.4 List Salesmen

```http
GET /api/v1/tenant/salesmen?page=1&per_page=20
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "e0000000-0000-0000-0000-000000000001",
      "name": "Andi Pratama",
      "email": "andi@propertijaya.id",
      "phone": "081200000002",
      "photo_url": null,
      "status": "active",
      "listing_count": {
        "total": 8,
        "active": 5
      },
      "created_at": "2026-06-20T09:30:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

---

### 7.5 Add Salesman

```http
POST /api/v1/tenant/salesmen
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Dewi Lestari",
  "email": "dewi@propertijaya.id",
  "phone": "081200000005",
  "password": "Dewi@123"
}
```

| Field      | Type   | Required | Rules                              |
| ---------- | ------ | -------- | ---------------------------------- |
| `name`     | string | ✅       | 2–100 chars                        |
| `email`    | string | ✅       | Valid email, unique, max 255 chars |
| `phone`    | string | ✅       | Regex `^\+?[0-9]{8,15}$`           |
| `password` | string | ✅       | Min 8 chars                        |

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "id": "e0000000-0000-0000-0000-000000000004",
    "name": "Dewi Lestari",
    "email": "dewi@propertijaya.id",
    "phone": "081200000005",
    "role": "salesman",
    "status": "active",
    "created_at": "2026-06-26T10:00:00+07:00"
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `AUTH_EMAIL_REGISTERED` | 409 | Email sudah terdaftar |
| `BIZ_SALESMAN_LIMIT` | 422 | Salesman limit tercapai |

---

### 7.6 Remove Salesman

```http
DELETE /api/v1/tenant/salesmen/:id
Authorization: Bearer <token>
```

**Catatan:** Salesman akan di-set `status = 'inactive'` (soft deactivate), bukan dihapus.

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Salesman berhasil dinonaktifkan."
  }
}
```

---

### 7.7 List All Tenant Listings

```http
GET /api/v1/tenant/listings?status=approved&salesman_id=e0000000-...&page=1&per_page=20
Authorization: Bearer <token>
```

**Query Parameters:**

| Param         | Type    | Required | Default | Deskripsi                |
| ------------- | ------- | -------- | ------- | ------------------------ |
| `status`      | string  | ❌       | –       | Filter by listing status |
| `salesman_id` | UUID    | ❌       | –       | Filter by salesman       |
| `page`        | integer | ❌       | 1       | Halaman                  |
| `per_page`    | integer | ❌       | 20      | Item per halaman         |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000001",
      "title": "Rumah Minimalis 2 Lantai di Kebayoran Baru",
      "price": "2500000000.00",
      "listing_type": "sale",
      "property_type": "house",
      "status": "approved",
      "salesman": {
        "id": "e0000000-...",
        "name": "Andi Pratama"
      },
      "created_at": "2026-06-20T10:00:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

---

### 7.8 View Subscription

```http
GET /api/v1/tenant/subscription
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "plan_type": "free",
    "max_salesmen": 5,
    "max_listings_per_salesman": 5,
    "start_date": "2026-06-20T09:00:00+07:00",
    "end_date": null,
    "usage": {
      "salesmen_used": 3,
      "salesmen_max": 5,
      "total_active_listings": 18
    }
  }
}
```

---

### 7.9 Request Plan Upgrade

```http
POST /api/v1/tenant/subscription/upgrade
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "plan_type": "premium"
}
```

| Field       | Type   | Required | Rules     |
| ----------- | ------ | -------- | --------- |
| `plan_type` | string | ✅       | `premium` |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "message": "Permintaan upgrade ke Premium telah dikirim. Tim kami akan menghubungi Anda dalam 1×24 jam.",
    "requested_plan": "premium"
  }
}
```

> MVP: Upgrade request disimpan sebagai log. Proses approval manual oleh platform admin.

---

## 8. Platform Admin Endpoints

**Role:** `platform_admin` | **Auth:** JWT | **Scope:** Global

### 8.1 Platform Dashboard

```http
GET /api/v1/admin/dashboard
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "total_tenants": 2,
    "active_tenants": 2,
    "suspended_tenants": 0,
    "total_users": 8,
    "total_listings": 8,
    "pending_reviews": 1,
    "listings_by_status": {
      "draft": 1,
      "pending": 1,
      "approved": 4,
      "rejected": 1,
      "sold": 1,
      "rented": 0,
      "inactive": 0
    }
  }
}
```

---

### 8.2 List All Tenants

```http
GET /api/v1/admin/tenants?status=active&page=1&per_page=20
Authorization: Bearer <token>
```

**Query Parameters:**

| Param      | Type    | Required | Default | Deskripsi              |
| ---------- | ------- | -------- | ------- | ---------------------- |
| `status`   | string  | ❌       | –       | `active` / `suspended` |
| `page`     | integer | ❌       | 1       | Halaman                |
| `per_page` | integer | ❌       | 20      | Item per halaman       |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "b0000000-0000-0000-0000-000000000001",
      "organization_name": "PropertiJaya Agency",
      "subdomain_slug": "propertijaya",
      "phone": "0215551234",
      "status": "active",
      "plan_type": "free",
      "total_listings": 7,
      "total_users": 4,
      "created_at": "2026-06-20T09:00:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2,
    "total_pages": 1
  }
}
```

---

### 8.3 Create Tenant

```http
POST /api/v1/admin/tenants
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "organization_name": "PropertiBaru Agency",
  "subdomain_slug": "propertibaru",
  "admin_name": "Ahmad Fauzi",
  "admin_email": "ahmad@propertibaru.id",
  "admin_phone": "081100000003",
  "admin_password": "Ahmad@123",
  "plan_type": "free"
}
```

| Field               | Type   | Required | Rules                                               |
| ------------------- | ------ | -------- | --------------------------------------------------- |
| `organization_name` | string | ✅       | 2–200 chars                                         |
| `subdomain_slug`    | string | ✅       | Lowercase, alphanumeric + dash, 3–100 chars, unique |
| `admin_name`        | string | ✅       | 2–100 chars                                         |
| `admin_email`       | string | ✅       | Valid email, unique                                 |
| `admin_phone`       | string | ✅       | Regex `^\+?[0-9]{8,15}$`                            |
| `admin_password`    | string | ✅       | Min 8 chars                                         |
| `plan_type`         | string | ❌       | `free` (default), `premium`                         |

**Success Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "tenant": {
      "id": "b0000000-0000-0000-0000-000000000003",
      "organization_name": "PropertiBaru Agency",
      "subdomain_slug": "propertibaru",
      "status": "active"
    },
    "admin": {
      "id": "d0000000-0000-0000-0000-000000000003",
      "name": "Ahmad Fauzi",
      "email": "ahmad@propertibaru.id",
      "role": "tenant_admin"
    },
    "subscription": {
      "plan_type": "free",
      "max_salesmen": 5,
      "max_listings_per_salesman": 5
    }
  }
}
```

> Endpoint ini melakukan transaksi: insert `tenants` + `users` (tenant_admin) + `subscriptions` dalam satu transaksi database.

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `AUTH_EMAIL_REGISTERED` | 409 | Admin email sudah terdaftar |
| `RES_ALREADY_EXISTS` | 409 | subdomain_slug sudah digunakan |

---

### 8.4 Get Tenant Detail

```http
GET /api/v1/admin/tenants/:id
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "b0000000-0000-0000-0000-000000000001",
    "organization_name": "PropertiJaya Agency",
    "subdomain_slug": "propertijaya",
    "logo_url": null,
    "description": "Agensi properti terpercaya sejak 2010...",
    "phone": "0215551234",
    "address": "Jl. Sudirman No. 123, Jakarta Pusat",
    "status": "active",
    "subscription": {
      "plan_type": "free",
      "max_salesmen": 5,
      "max_listings_per_salesman": 5,
      "start_date": "2026-06-20T09:00:00+07:00"
    },
    "stats": {
      "total_users": 4,
      "total_salesmen": 3,
      "total_listings": 7,
      "pending_listings": 1
    },
    "created_at": "2026-06-20T09:00:00+07:00"
  }
}
```

---

### 8.5 Suspend Tenant

```http
POST /api/v1/admin/tenants/:id/suspend
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "b0000000-0000-0000-0000-000000000001",
    "status": "suspended",
    "message": "Tenant berhasil dinonaktifkan. Semua pengguna tidak dapat login."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `BIZ_INVALID_STATUS_TRANSITION` | 422 | Tenant sudah suspended |

---

### 8.6 Activate Tenant

```http
POST /api/v1/admin/tenants/:id/activate
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "b0000000-0000-0000-0000-000000000001",
    "status": "active",
    "message": "Tenant berhasil diaktifkan kembali."
  }
}
```

---

### 8.7 Change Tenant Plan

```http
PUT /api/v1/admin/tenants/:id/plan
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "plan_type": "premium"
}
```

| Field       | Type   | Required | Rules              |
| ----------- | ------ | -------- | ------------------ |
| `plan_type` | string | ✅       | `free` / `premium` |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "tenant_id": "b0000000-0000-0000-0000-000000000001",
    "plan_type": "premium",
    "max_salesmen": 999999,
    "max_listings_per_salesman": 999999,
    "message": "Paket tenant berhasil diubah ke Premium."
  }
}
```

---

### 8.8 List Pending Listings

```http
GET /api/v1/admin/listings/pending?page=1&per_page=20
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "10000000-0000-0000-0000-000000000004",
      "title": "Villa Mewah Pemandangan Gunung Lembang",
      "price": "1800000000.00",
      "listing_type": "sale",
      "property_type": "villa",
      "source_type": "regular",
      "city": "Bandung",
      "status": "pending",
      "tenant": {
        "id": "b0000000-...",
        "name": "PropertiJaya Agency"
      },
      "salesman": {
        "id": "e0000000-...",
        "name": "Siti Nurhaliza"
      },
      "photos": [
        {
          "thumbnail_url": "/uploads/listings/.../thumb.jpg"
        }
      ],
      "created_at": "2026-06-25T14:00:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 8.9 Approve Listing

```http
POST /api/v1/admin/listings/:id/approve
Authorization: Bearer <token>
```

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000004",
    "status": "approved",
    "approved_by": "a0000000-0000-0000-0000-000000000001",
    "approved_at": "2026-06-26T11:00:00+07:00",
    "message": "Listing berhasil disetujui. Sekarang tampil di halaman publik."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `RES_LISTING_NOT_FOUND` | 404 | Listing tidak ditemukan |
| `BIZ_INVALID_STATUS_TRANSITION` | 422 | Status bukan `pending` |

---

### 8.10 Reject Listing

```http
POST /api/v1/admin/listings/:id/reject
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "reason": "Foto yang diunggah tidak sesuai dengan deskripsi properti. Mohon unggah foto asli."
}
```

| Field    | Type   | Required | Rules        |
| -------- | ------ | -------- | ------------ |
| `reason` | string | ✅       | 10–500 chars |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "id": "10000000-0000-0000-0000-000000000004",
    "status": "rejected",
    "reject_reason": "Foto yang diunggah tidak sesuai dengan deskripsi properti. Mohon unggah foto asli.",
    "message": "Listing berhasil ditolak. Salesman dapat mengedit dan mengajukan ulang."
  }
}
```

**Error Responses:**
| Code | HTTP | Condition |
|------|------|-----------|
| `BIZ_INVALID_STATUS_TRANSITION` | 422 | Status bukan `pending` |
| `BIZ_REJECT_REASON_REQUIRED` | 422 | Reason < 10 karakter |

---

### 8.11 View Audit Logs

```http
GET /api/v1/admin/audit-logs?user_id=e00000...&action=approve&entity_type=listing&from=2026-06-01&to=2026-06-30&page=1&per_page=50
Authorization: Bearer <token>
```

**Query Parameters:**

| Param         | Type    | Required | Default | Deskripsi                                                                |
| ------------- | ------- | -------- | ------- | ------------------------------------------------------------------------ |
| `user_id`     | UUID    | ❌       | –       | Filter by user                                                           |
| `action`      | string  | ❌       | –       | `create`, `update`, `delete`, `approve`, `reject`, `suspend`, `activate` |
| `entity_type` | string  | ❌       | –       | `listing`, `user`, `tenant`, `subscription`                              |
| `from`        | date    | ❌       | –       | Dari tanggal (ISO date)                                                  |
| `to`          | date    | ❌       | –       | Sampai tanggal (ISO date)                                                |
| `page`        | integer | ❌       | 1       | Halaman                                                                  |
| `per_page`    | integer | ❌       | 50      | Item per halaman (max 100)                                               |

**Success Response:** `200 OK`

```json
{
  "success": true,
  "data": [
    {
      "id": "aa000000-0000-0000-0000-000000000001",
      "user": {
        "id": "a0000000-0000-0000-0000-000000000001",
        "name": "Super Admin",
        "role": "platform_admin"
      },
      "action": "approve",
      "entity_type": "listing",
      "entity_id": "10000000-0000-0000-0000-000000000001",
      "old_values": { "status": "pending" },
      "new_values": { "status": "approved" },
      "ip_address": "192.168.1.1",
      "created_at": "2026-06-20T10:00:00+07:00"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 50,
    "total": 25,
    "total_pages": 1
  }
}
```

---

## 9. Error Response Catalog

| Code                            | HTTP | Message                                                                                             |
| ------------------------------- | ---- | --------------------------------------------------------------------------------------------------- |
| `AUTH_INVALID_CREDENTIALS`      | 401  | Email atau password salah.                                                                          |
| `AUTH_TOKEN_MISSING`            | 401  | Token autentikasi diperlukan. Silakan login terlebih dahulu.                                        |
| `AUTH_TOKEN_EXPIRED`            | 401  | Sesi Anda telah berakhir. Silakan login kembali.                                                    |
| `AUTH_TOKEN_INVALID`            | 401  | Token autentikasi tidak valid.                                                                      |
| `AUTH_ACCOUNT_SUSPENDED`        | 403  | Akun organisasi Anda sedang dinonaktifkan. Hubungi administrator.                                   |
| `AUTH_ACCOUNT_INACTIVE`         | 403  | Akun Anda tidak aktif. Hubungi administrator.                                                       |
| `AUTH_EMAIL_REGISTERED`         | 409  | Email sudah terdaftar. Silakan gunakan email lain atau login.                                       |
| `AUTHZ_FORBIDDEN`               | 403  | Anda tidak memiliki izin untuk mengakses resource ini.                                              |
| `AUTHZ_NOT_OWNER`               | 403  | Anda tidak dapat mengubah listing milik sales lain.                                                 |
| `AUTHZ_CROSS_TENANT`            | 403  | Data tidak dapat diakses lintas organisasi.                                                         |
| `VAL_INPUT_INVALID`             | 422  | Beberapa field tidak valid. Silakan periksa kembali.                                                |
| `VAL_FILE_TOO_LARGE`            | 413  | Ukuran file maksimal 5 MB per foto.                                                                 |
| `VAL_FILE_INVALID_TYPE`         | 415  | Format file tidak didukung. Gunakan JPEG, PNG, atau WebP.                                           |
| `VAL_FILE_COUNT_EXCEEDED`       | 422  | Maksimal 10 foto per listing.                                                                       |
| `RES_NOT_FOUND`                 | 404  | Data tidak ditemukan.                                                                               |
| `RES_LISTING_NOT_FOUND`         | 404  | Listing properti tidak ditemukan.                                                                   |
| `RES_ALREADY_SAVED`             | 409  | Properti sudah ada di daftar favorit Anda.                                                          |
| `RES_ALREADY_EXISTS`            | 409  | Data dengan field tersebut sudah ada.                                                               |
| `BIZ_QUOTA_EXCEEDED`            | 422  | Kuota listing Anda sudah penuh ({current}/{max}). Upgrade ke Premium untuk listing unlimited.       |
| `BIZ_SALESMAN_LIMIT`            | 422  | Jumlah salesman sudah mencapai batas ({current}/{max}). Upgrade ke Premium untuk menambah salesman. |
| `BIZ_LISTING_NOT_EDITABLE`      | 422  | Listing dengan status {status} tidak dapat diedit.                                                  |
| `BIZ_LISTING_NOT_DELETABLE`     | 422  | Listing dengan status {status} tidak dapat dihapus.                                                 |
| `BIZ_LISTING_NOT_SUBMITTABLE`   | 422  | Hanya listing dengan status draft atau rejected yang dapat diajukan.                                |
| `BIZ_INVALID_STATUS_TRANSITION` | 422  | Tidak dapat mengubah status dari {from} ke {to}.                                                    |
| `BIZ_REJECT_REASON_REQUIRED`    | 422  | Alasan penolakan wajib diisi (minimal 10 karakter).                                                 |
| `RATE_LOGIN_LIMIT`              | 429  | Terlalu banyak percobaan login. Silakan coba lagi dalam 1 menit.                                    |
| `RATE_GLOBAL_LIMIT`             | 429  | Terlalu banyak permintaan. Silakan coba lagi nanti.                                                 |
| `SRV_INTERNAL_ERROR`            | 500  | Terjadi kesalahan pada server. Silakan coba lagi nanti.                                             |

---

## 10. Endpoint Summary Matrix

| #   | Method   | Path                                            | Auth | Role                                          |
| --- | -------- | ----------------------------------------------- | :--: | --------------------------------------------- |
| 1   | `POST`   | `/api/v1/auth/register`                         |  ❌  | Public                                        |
| 2   | `POST`   | `/api/v1/auth/login`                            |  ❌  | Public                                        |
| 3   | `GET`    | `/api/v1/properties`                            |  ❌  | Public                                        |
| 4   | `GET`    | `/api/v1/properties/featured`                   |  ❌  | Public                                        |
| 5   | `GET`    | `/api/v1/properties/nearby`                     |  ❌  | Public                                        |
| 6   | `GET`    | `/api/v1/properties/:id`                        |  ❌  | Public                                        |
| 7   | `GET`    | `/api/v1/locations/cities`                      |  ❌  | Public                                        |
| 8   | `GET`    | `/api/v1/me/profile`                            |  ✅  | buyer, salesman, tenant_admin, platform_admin |
| 9   | `PUT`    | `/api/v1/me/profile`                            |  ✅  | buyer, salesman, tenant_admin, platform_admin |
| 10  | `GET`    | `/api/v1/me/saved`                              |  ✅  | buyer                                         |
| 11  | `POST`   | `/api/v1/me/saved/:propertyId`                  |  ✅  | buyer                                         |
| 12  | `DELETE` | `/api/v1/me/saved/:propertyId`                  |  ✅  | buyer                                         |
| 13  | `GET`    | `/api/v1/salesman/dashboard`                    |  ✅  | salesman, tenant_admin                        |
| 14  | `GET`    | `/api/v1/salesman/listings`                     |  ✅  | salesman, tenant_admin                        |
| 15  | `POST`   | `/api/v1/salesman/listings`                     |  ✅  | salesman, tenant_admin                        |
| 16  | `GET`    | `/api/v1/salesman/listings/:id`                 |  ✅  | salesman, tenant_admin                        |
| 17  | `PUT`    | `/api/v1/salesman/listings/:id`                 |  ✅  | salesman, tenant_admin                        |
| 18  | `DELETE` | `/api/v1/salesman/listings/:id`                 |  ✅  | salesman, tenant_admin                        |
| 19  | `POST`   | `/api/v1/salesman/listings/:id/submit`          |  ✅  | salesman, tenant_admin                        |
| 20  | `POST`   | `/api/v1/salesman/listings/:id/deactivate`      |  ✅  | salesman, tenant_admin                        |
| 21  | `POST`   | `/api/v1/salesman/listings/:id/mark-sold`       |  ✅  | salesman, tenant_admin                        |
| 22  | `POST`   | `/api/v1/salesman/listings/:id/mark-rented`     |  ✅  | salesman, tenant_admin                        |
| 23  | `POST`   | `/api/v1/salesman/listings/:id/photos`          |  ✅  | salesman, tenant_admin                        |
| 24  | `DELETE` | `/api/v1/salesman/listings/:id/photos/:photoId` |  ✅  | salesman, tenant_admin                        |
| 25  | `PUT`    | `/api/v1/salesman/listings/:id/photos/reorder`  |  ✅  | salesman, tenant_admin                        |
| 26  | `GET`    | `/api/v1/salesman/quota`                        |  ✅  | salesman, tenant_admin                        |
| 27  | `GET`    | `/api/v1/tenant/dashboard`                      |  ✅  | tenant_admin                                  |
| 28  | `GET`    | `/api/v1/tenant/profile`                        |  ✅  | tenant_admin                                  |
| 29  | `PUT`    | `/api/v1/tenant/profile`                        |  ✅  | tenant_admin                                  |
| 30  | `GET`    | `/api/v1/tenant/salesmen`                       |  ✅  | tenant_admin                                  |
| 31  | `POST`   | `/api/v1/tenant/salesmen`                       |  ✅  | tenant_admin                                  |
| 32  | `DELETE` | `/api/v1/tenant/salesmen/:id`                   |  ✅  | tenant_admin                                  |
| 33  | `GET`    | `/api/v1/tenant/listings`                       |  ✅  | tenant_admin                                  |
| 34  | `GET`    | `/api/v1/tenant/subscription`                   |  ✅  | tenant_admin                                  |
| 35  | `POST`   | `/api/v1/tenant/subscription/upgrade`           |  ✅  | tenant_admin                                  |
| 36  | `GET`    | `/api/v1/admin/dashboard`                       |  ✅  | platform_admin                                |
| 37  | `GET`    | `/api/v1/admin/tenants`                         |  ✅  | platform_admin                                |
| 38  | `POST`   | `/api/v1/admin/tenants`                         |  ✅  | platform_admin                                |
| 39  | `GET`    | `/api/v1/admin/tenants/:id`                     |  ✅  | platform_admin                                |
| 40  | `POST`   | `/api/v1/admin/tenants/:id/suspend`             |  ✅  | platform_admin                                |
| 41  | `POST`   | `/api/v1/admin/tenants/:id/activate`            |  ✅  | platform_admin                                |
| 42  | `PUT`    | `/api/v1/admin/tenants/:id/plan`                |  ✅  | platform_admin                                |
| 43  | `GET`    | `/api/v1/admin/listings/pending`                |  ✅  | platform_admin                                |
| 44  | `POST`   | `/api/v1/admin/listings/:id/approve`            |  ✅  | platform_admin                                |
| 45  | `POST`   | `/api/v1/admin/listings/:id/reject`             |  ✅  | platform_admin                                |
| 46  | `GET`    | `/api/v1/admin/audit-logs`                      |  ✅  | platform_admin                                |

**Total: 46 endpoints** (7 public + 39 protected)

---

_Dokumen ini adalah Tahap 3. Lanjut ke Tahap 4 setelah mengetik: **LANJUT**._
