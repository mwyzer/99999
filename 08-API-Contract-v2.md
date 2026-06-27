# API Contract — v2

## Multi-Tenant Property Information System

| Property            | Value                                           |
| ------------------- | ----------------------------------------------- |
| **Document Type**   | API Contract / OpenAPI-Like Specification       |
| **Version**         | 2.0.0                                           |
| **Date**            | 2026-06-27                                      |
| **Base URL**        | `http://localhost:8080/api/v1`                  |
| **Reference**       | `07-ERD-Database-Schema-v2.md`, `02-SRS-MVP.md` |
| **Total Endpoints** | 52 (dari 46 v1)                                 |

---

## 1. API Overview

| Property         | Value                              |
| ---------------- | ---------------------------------- |
| **Protocol**     | HTTP/1.1 + HTTPS                   |
| **Base Path**    | `/api/v1`                          |
| **Content-Type** | `application/json`                 |
| **Auth**         | JWT Bearer Token                   |
| **Date Format**  | ISO 8601                           |
| **Decimal**      | String (anti floating-point error) |

### 1.1 Endpoint Role Mapping

| Prefix                 | Role Required              | Scope  |
| ---------------------- | -------------------------- | ------ |
| `/api/v1/auth/*`       | None (public)              | –      |
| `/api/v1/properties/*` | None (public)              | –      |
| `/api/v1/locations/*`  | None (public)              | –      |
| `/api/v1/me/*`         | Any authenticated          | Self   |
| `/api/v1/salesman/*`   | `salesman`, `tenant_admin` | Tenant |
| `/api/v1/tenant/*`     | `tenant_admin`             | Tenant |
| `/api/v1/admin/*`      | `platform_admin`           | Global |
| `/api/v1/buyer/*`      | `buyer`                    | Self   |

### 1.2 Generic Response

**Success (with pagination):**

```json
{ "success": true, "data": [...], "meta": { "page": 1, "per_page": 20, "total": 150, "total_pages": 8 } }
```

**Success (no pagination):**

```json
{ "success": true, "data": {} }
```

**Error:**

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": []
  }
}
```

---

## 2. Auth Endpoints

### 2.1 Register

```http
POST /api/v1/auth/register
```

**Body:** `name` (2-100), `email` (valid, unique), `phone` (8-15 digits), `password` (min 8)

**Response `201`:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Rina Wijaya",
    "email": "rina@email.com",
    "phone": "081300000001",
    "role": "buyer",
    "status": "active",
    "created_at": "..."
  }
}
```

### 2.2 Login

```http
POST /api/v1/auth/login
```

**Body:** `email`, `password`

**Response `200`:**

```json
{
  "success": true,
  "data": {
    "token": "jwt...",
    "user": {
      "id": "uuid",
      "name": "...",
      "email": "...",
      "role": "buyer",
      "tenant_id": null,
      "tenant_name": null
    }
  }
}
```

---

## 3. Public Property Endpoints

### 3.1 List Properties

```http
GET /api/v1/properties?page=1&per_page=20&property_type=house&source_type=regular&listing_type=sale&city=Jakarta&price_min=500000000&price_max=5000000000&q=rumah&sort=created_at&order=desc
```

**Response `200`:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Rumah Minimalis Bintaro",
      "price": "Rp 1,85M",
      "listing_type": "sale",
      "property_type": "house",
      "source_type": "regular",
      "city": "Tangerang Selatan",
      "province": "Banten",
      "land_area": "150.00",
      "building_area": "200.00",
      "bedrooms": 4,
      "bathrooms": 3,
      "main_photo_url": "/uploads/.../wm_abc.jpg",
      "salesman": {
        "id": "uuid",
        "name": "Andi Pratama",
        "photo_url": null,
        "phone": "081300000001"
      },
      "tenant": {
        "id": "uuid",
        "name": "PropertiJaya Agency",
        "logo_url": null
      },
      "status": "approved",
      "created_at": "..."
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 150, "total_pages": 8 }
}
```

### 3.2 Property Detail

```http
GET /api/v1/properties/:id
```

**Response `200`:** (includes ALL fields + photos + facilities + auction/company details if applicable)

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Rumah Minimalis Bintaro",
    "description": "Rumah modern minimalis...",
    "price": "Rp 1,85M",
    "listing_type": "sale",
    "property_type": "house",
    "source_type": "regular",
    "rent_period": null,
    "address": "Jl. Bintaro Utama 7 No. 15",
    "city": "Tangerang Selatan",
    "province": "Banten",
    "latitude": "-6.2765000",
    "longitude": "106.7183000",
    "land_area": "150.00",
    "building_area": "200.00",
    "bedrooms": 4,
    "bathrooms": 3,
    "floors": 2,
    "certificate_type": "SHM",
    "status": "approved",
    "photos": [
      {
        "id": "uuid",
        "image_url": "...",
        "watermarked_image_url": "...",
        "watermark_status": "skipped",
        "thumbnail_url": "...",
        "medium_url": "...",
        "is_primary": true,
        "sort_order": 0
      }
    ],
    "facilities": [
      {
        "id": "uuid",
        "facility_id": "uuid",
        "name": "Carport",
        "icon": "car",
        "value": "2 mobil"
      }
    ],
    "bank_auction_details": null,
    "company_asset_details": null,
    "salesman": {
      "id": "uuid",
      "name": "Andi Pratama",
      "photo_url": null,
      "phone": "081300000001",
      "whatsapp_number": null,
      "show_whatsapp": true
    },
    "tenant": {
      "id": "uuid",
      "name": "PropertiJaya Agency",
      "logo_url": null,
      "whatsapp_number": null,
      "show_whatsapp": true
    }
  }
}
```

### 3.3 Featured Properties

```http
GET /api/v1/properties/featured?city=Jakarta%20Selatan&limit=6
```

### 3.4 Nearby Properties

```http
GET /api/v1/properties/nearby?lat=-6.2765&lng=106.7183&radius=10&limit=10
```

**Response:** includes `distance_km` per item.

---

## 4. Master Data Endpoints (Public)

### 4.1 List Property Types

```http
GET /api/v1/property-types
```

**Response `200`:**

```json
{ "success": true, "data": [{ "id": "uuid", "name": "Rumah", "slug": "house" }, ...] }
```

### 4.2 List Facilities

```http
GET /api/v1/facilities
```

**Response `200`:**

```json
{ "success": true, "data": [{ "id": "uuid", "name": "Carport", "icon": "car" }, ...] }
```

### 4.3 List Locations

```http
GET /api/v1/locations?province=DKI+Jakarta&q=Banda
```

**Response `200`:**

```json
{ "success": true, "data": [{ "id": "uuid", "city": "Jakarta Selatan", "province": "DKI Jakarta", "latitude": "-6.2436000", "longitude": "106.8000000" }, ...] }
```

---

## 5. Profile Endpoints (All Roles)

### 5.1 Get Profile

```http
GET /api/v1/me/profile
Authorization: Bearer <token>
```

**Response `200`:** (buyer tidak punya tenant info)

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "...",
    "email": "...",
    "phone": "...",
    "photo_url": null,
    "whatsapp_number": null,
    "show_whatsapp": true,
    "role": "salesman",
    "tenant": { "id": "uuid", "name": "PropertiJaya Agency" }
  }
}
```

### 5.2 Update Profile

```http
PUT /api/v1/me/profile
```

**Body:** `name`, `phone`, `photo_url`, `whatsapp_number`, `show_whatsapp` (all optional)

---

## 6. Buyer Endpoints

### 6.1 List Favorites

```http
GET /api/v1/buyer/favorites?page=1&per_page=12
Authorization: Bearer <token>
```

**Response:** array of property cards + pagination.

### 6.2 Save Favorite

```http
POST /api/v1/buyer/favorites/:propertyId
```

**Response `201`:** `{ "success": true, "data": { "id": "uuid", "property_id": "uuid" } }`

### 6.3 Remove Favorite

```http
DELETE /api/v1/buyer/favorites/:propertyId
```

**Response `200`:** `{ "success": true, "message": "Dihapus dari favorit." }`

### 6.4 Send Inquiry

```http
POST /api/v1/buyer/inquiries
```

**Body:** `property_id` (UUID, required), `message` (text, required, max 2000)

**Response `201`:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "property_id": "uuid",
    "message": "Apakah harga bisa nego?",
    "status": "unread",
    "created_at": "..."
  }
}
```

### 6.5 List My Inquiries

```http
GET /api/v1/buyer/inquiries?page=1&per_page=20
```

**Response:** array of inquiries with property title + status.

---

## 7. Salesman Endpoints

### 7.1 Dashboard

```http
GET /api/v1/salesman/dashboard
```

**Response:**

```json
{
  "success": true,
  "data": {
    "total_listings": 12,
    "active_count": 8,
    "status_breakdown": {
      "draft": 2,
      "pending": 1,
      "approved": 8,
      "rejected": 0,
      "sold": 1,
      "rented": 0,
      "inactive": 0
    },
    "quota": { "used": 8, "max": 5 }
  }
}
```

### 7.2 List My Listings

```http
GET /api/v1/salesman/listings?page=1&per_page=20&status=approved
```

### 7.3 Create Listing

```http
POST /api/v1/salesman/listings
```

**Body:**

```json
{
  "title": "Rumah Minimalis Bintaro",
  "description": "...",
  "price": 1850000000,
  "listing_type": "sale",
  "property_type_id": "uuid-of-house",
  "source_type": "regular",
  "rent_period": null,
  "location_id": "uuid-of-tangsel",
  "address": "Jl. Bintaro Utama 7 No. 15",
  "latitude": -6.2765,
  "longitude": 106.7183,
  "land_area": 150,
  "building_area": 200,
  "bedrooms": 4,
  "bathrooms": 3,
  "floors": 2,
  "certificate_type": "SHM",
  "facility_ids": ["uuid1", "uuid2"],
  "bank_auction": null,
  "company_asset": null
}
```

**When `source_type = "bank_auction"`:**

```json
{
  "bank_auction": {
    "bank_name": "BankMaju",
    "auction_number": "AUC-2026-001",
    "auction_limit_price": 500000000,
    "auction_deposit": 50000000,
    "auction_date": "2026-07-15T10:00:00+07:00",
    "auction_location": "Kantor Cabang BankMaju, Jakarta"
  }
}
```

**When `source_type = "company_asset"`:**

```json
{
  "company_asset": {
    "company_name": "CiptaGraha Corporindo",
    "company_asset_code": "CG-ASSET-001",
    "disposal_type": "sale",
    "pic_name": "Ratna Sari Dewi",
    "pic_phone": "081200000005",
    "pic_whatsapp_number": "081200000005"
  }
}
```

### 7.4 Get Listing

```http
GET /api/v1/salesman/listings/:id
```

**Response:** full listing detail (own listing only).

### 7.5 Update Listing

```http
PUT /api/v1/salesman/listings/:id
```

Same body as create. Only `draft` & `rejected` status allowed.

### 7.6 Delete Listing

```http
DELETE /api/v1/salesman/listings/:id
```

### 7.7 Submit for Review

```http
POST /api/v1/salesman/listings/:id/submit
```

Changes status from `draft`/`rejected` → `pending`.

### 7.8 Deactivate Listing

```http
POST /api/v1/salesman/listings/:id/deactivate
```

Changes status → `inactive`.

### 7.9 Mark as Sold

```http
POST /api/v1/salesman/listings/:id/mark-sold
```

### 7.10 Mark as Rented

```http
POST /api/v1/salesman/listings/:id/mark-rented
```

### 7.11 Upload Photos

```http
POST /api/v1/salesman/listings/:id/photos
Content-Type: multipart/form-data
```

**Body:** `photos` (file[], max 5MB each, JPG/PNG/WebP)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "image_url": "/uploads/.../orig.jpg",
      "watermarked_image_url": "/uploads/.../wm.jpg",
      "watermark_status": "pending",
      "thumbnail_url": "/uploads/.../thumb.jpg",
      "medium_url": "/uploads/.../med.jpg",
      "is_primary": false,
      "sort_order": 1
    }
  ]
}
```

### 7.12 Delete Photo

```http
DELETE /api/v1/salesman/listings/:id/photos/:photoId
```

### 7.13 Reorder Photos

```http
PUT /api/v1/salesman/listings/:id/photos/reorder
```

**Body:** `{ "photo_ids": ["uuid1", "uuid2", "uuid3"] }`

### 7.14 Get Quota

```http
GET /api/v1/salesman/quota
```

**Response:** `{ "used": 4, "max": 5, "plan": "free" }`

### 7.15 List Inquiries (for my listings)

```http
GET /api/v1/salesman/inquiries?page=1&per_page=20&status=unread
```

### 7.16 Reply/Update Inquiry

```http
PUT /api/v1/salesman/inquiries/:id
```

**Body:** `{ "status": "replied" }`

---

## 8. Tenant Admin Endpoints

### 8.1 Dashboard

```http
GET /api/v1/tenant/dashboard
```

**Response:**

```json
{
  "success": true,
  "data": {
    "total_listings": 45,
    "active_listings": 32,
    "total_salesmen": 3,
    "max_salesmen": 5,
    "plan": {
      "type": "free",
      "name": "Free",
      "max_salesmen": 5,
      "max_listings_per_salesman": 5
    }
  }
}
```

### 8.2 List Salesmen

```http
GET /api/v1/tenant/salesmen?page=1&per_page=50
```

### 8.3 Add Salesman

```http
POST /api/v1/tenant/salesmen
```

**Body:** `name`, `email`, `phone`, `password` (min 8)

### 8.4 Remove/Deactivate Salesman

```http
DELETE /api/v1/tenant/salesmen/:id
```

### 8.5 List Tenant Listings

```http
GET /api/v1/tenant/listings?page=1&per_page=20&status=approved&salesman_id=uuid
```

### 8.6 Get Tenant Profile

```http
GET /api/v1/tenant/profile
```

**Response:** `{ "organization_name": "...", "subdomain_slug": "...", "logo_url": null, "whatsapp_number": null, "show_whatsapp": true, "description": "...", "phone": "...", "address": "..." }`

### 8.7 Update Tenant Profile

```http
PUT /api/v1/tenant/profile
```

**Body:** `organization_name`, `description`, `phone`, `address`, `logo_url`, `whatsapp_number`, `show_whatsapp`

### 8.8 Get Subscription

```http
GET /api/v1/tenant/subscription
```

**Response:**

```json
{
  "success": true,
  "data": {
    "plan": {
      "id": "uuid",
      "name": "Free",
      "slug": "free",
      "max_salesmen": 5,
      "max_listings_per_salesman": 5
    },
    "start_date": "...",
    "end_date": null,
    "status": "active"
  }
}
```

### 8.9 Request Upgrade

```http
POST /api/v1/tenant/subscription/upgrade
```

### 8.10 List All Inquiries (tenant scope)

```http
GET /api/v1/tenant/inquiries?page=1&per_page=20&status=unread&property_id=uuid
```

---

## 9. Platform Admin Endpoints

### 9.1 Dashboard

```http
GET /api/v1/admin/dashboard
```

**Response:** `{ "total_tenants": 12, "active_tenants": 10, "suspended_tenants": 2, "total_users": 345, "total_listings": 1200, "pending_reviews": 15, "listings_by_status": { "draft": 50, "pending": 15, "approved": 900, ... } }`

### 9.2 List Tenants

```http
GET /api/v1/admin/tenants?page=1&per_page=20&status=active&q=Jaya
```

**Response:** array of tenant + subscription info.

### 9.3 Create Tenant

```http
POST /api/v1/admin/tenants
```

**Body:**

```json
{
  "organization_name": "PT Maju Bersama",
  "subdomain_slug": "majubersama",
  "admin_name": "Budi Santoso",
  "admin_email": "budi@majubersama.id",
  "admin_phone": "081200000001",
  "admin_password": "Budi@123",
  "plan_type": "free"
}
```

### 9.4 Get Tenant Detail

```http
GET /api/v1/admin/tenants/:id
```

**Response:** tenant + subscription + salesmen count + listing count.

### 9.5 Suspend Tenant

```http
POST /api/v1/admin/tenants/:id/suspend
```

### 9.6 Activate Tenant

```http
POST /api/v1/admin/tenants/:id/activate
```

### 9.7 Change Plan

```http
PUT /api/v1/admin/tenants/:id/plan
```

**Body:** `{ "plan_type": "premium" }`

### 9.8 Change Plan by Plan ID

```http
PUT /api/v1/admin/tenants/:id/subscription
```

**Body:** `{ "plan_id": "uuid-of-premium-plan" }`

### 9.9 List Pending Reviews

```http
GET /api/v1/admin/listings/pending?page=1&per_page=20
```

### 9.10 Approve Listing

```http
POST /api/v1/admin/listings/:id/approve
```

### 9.11 Reject Listing

```http
POST /api/v1/admin/listings/:id/reject
```

**Body:** `{ "reason": "Foto kurang jelas, silakan upload ulang." }` (min 10, max 500)

### 9.12 List All Listings (admin view)

```http
GET /api/v1/admin/listings?page=1&per_page=20&tenant_id=uuid&status=pending
```

### 9.13 Audit Logs

```http
GET /api/v1/admin/audit-logs?page=1&per_page=20&action=APPROVE_PROPERTY&module=property&user_id=uuid&tenant_id=uuid
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "user_name": "Super Admin",
      "tenant_id": "uuid",
      "action": "APPROVE_PROPERTY",
      "module": "property",
      "reference_id": "uuid",
      "description": "Admin menyetujui listing",
      "old_data": { "status": "pending" },
      "new_data": { "status": "approved" },
      "ip_address": "127.0.0.1",
      "created_at": "..."
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 500, "total_pages": 25 }
}
```

### 9.14 Subscription Plans (CRUD)

```http
GET /api/v1/admin/subscription-plans
POST /api/v1/admin/subscription-plans
PUT /api/v1/admin/subscription-plans/:id
```

### 9.15 Property Types (CRUD)

```http
GET /api/v1/admin/property-types
POST /api/v1/admin/property-types
PUT /api/v1/admin/property-types/:id
```

### 9.16 Facilities (CRUD)

```http
GET /api/v1/admin/facilities
POST /api/v1/admin/facilities
PUT /api/v1/admin/facilities/:id
```

### 9.17 Locations (CRUD)

```http
GET /api/v1/admin/locations
POST /api/v1/admin/locations
PUT /api/v1/admin/locations/:id
```

---

## 10. Property Views (Auto)

### 10.1 Log View

```http
POST /api/v1/properties/:id/view
```

Auto-called by frontend when detail page loads. No auth required. Logs `user_id` (if authenticated), `ip_address`, `user_agent`.

---

## 11. Error Codes Reference

| Code                       | HTTP | Condition                     |
| -------------------------- | ---- | ----------------------------- |
| `VAL_INPUT_INVALID`        | 422  | Validasi field gagal          |
| `AUTH_EMAIL_REGISTERED`    | 409  | Email sudah dipakai           |
| `AUTH_INVALID_CREDENTIALS` | 401  | Email/password salah          |
| `AUTH_TOKEN_MISSING`       | 401  | Tidak ada token               |
| `AUTH_TOKEN_EXPIRED`       | 401  | Token expired                 |
| `AUTH_FORBIDDEN`           | 403  | Role tidak diizinkan          |
| `RES_NOT_FOUND`            | 404  | Resource tidak ditemukan      |
| `RES_ALREADY_EXISTS`       | 409  | Resource sudah ada            |
| `BIZ_QUOTA_EXCEEDED`       | 422  | Melebihi kuota                |
| `BIZ_MAX_PHOTOS_EXCEEDED`  | 422  | Foto > 10                     |
| `BIZ_LISTING_NOT_EDITABLE` | 422  | Status tidak mengizinkan edit |
| `BIZ_SAME_PLAN`            | 422  | Sudah pakai plan yang sama    |
| `BIZ_TENANT_SUSPENDED`     | 403  | Tenant dinonaktifkan          |
| `SYS_INTERNAL_ERROR`       | 500  | Server error                  |

---

## 12. Endpoint Summary

| #   | Method | Endpoint                                 | Auth | Role                   |
| --- | ------ | ---------------------------------------- | ---- | ---------------------- |
| 1   | POST   | `/auth/register`                         | No   | Public                 |
| 2   | POST   | `/auth/login`                            | No   | Public                 |
| 3   | GET    | `/properties`                            | No   | Public                 |
| 4   | GET    | `/properties/featured`                   | No   | Public                 |
| 5   | GET    | `/properties/nearby`                     | No   | Public                 |
| 6   | GET    | `/properties/:id`                        | No   | Public                 |
| 7   | POST   | `/properties/:id/view`                   | No   | Public                 |
| 8   | GET    | `/property-types`                        | No   | Public                 |
| 9   | GET    | `/facilities`                            | No   | Public                 |
| 10  | GET    | `/locations`                             | No   | Public                 |
| 11  | GET    | `/me/profile`                            | Yes  | All                    |
| 12  | PUT    | `/me/profile`                            | Yes  | All                    |
| 13  | GET    | `/buyer/favorites`                       | Yes  | buyer                  |
| 14  | POST   | `/buyer/favorites/:propertyId`           | Yes  | buyer                  |
| 15  | DELETE | `/buyer/favorites/:propertyId`           | Yes  | buyer                  |
| 16  | POST   | `/buyer/inquiries`                       | Yes  | buyer                  |
| 17  | GET    | `/buyer/inquiries`                       | Yes  | buyer                  |
| 18  | GET    | `/salesman/dashboard`                    | Yes  | salesman, tenant_admin |
| 19  | GET    | `/salesman/listings`                     | Yes  | salesman, tenant_admin |
| 20  | POST   | `/salesman/listings`                     | Yes  | salesman, tenant_admin |
| 21  | GET    | `/salesman/listings/:id`                 | Yes  | salesman, tenant_admin |
| 22  | PUT    | `/salesman/listings/:id`                 | Yes  | salesman, tenant_admin |
| 23  | DELETE | `/salesman/listings/:id`                 | Yes  | salesman, tenant_admin |
| 24  | POST   | `/salesman/listings/:id/submit`          | Yes  | salesman, tenant_admin |
| 25  | POST   | `/salesman/listings/:id/deactivate`      | Yes  | salesman, tenant_admin |
| 26  | POST   | `/salesman/listings/:id/mark-sold`       | Yes  | salesman, tenant_admin |
| 27  | POST   | `/salesman/listings/:id/mark-rented`     | Yes  | salesman, tenant_admin |
| 28  | POST   | `/salesman/listings/:id/photos`          | Yes  | salesman, tenant_admin |
| 29  | DELETE | `/salesman/listings/:id/photos/:photoId` | Yes  | salesman, tenant_admin |
| 30  | PUT    | `/salesman/listings/:id/photos/reorder`  | Yes  | salesman, tenant_admin |
| 31  | GET    | `/salesman/quota`                        | Yes  | salesman, tenant_admin |
| 32  | GET    | `/salesman/inquiries`                    | Yes  | salesman, tenant_admin |
| 33  | PUT    | `/salesman/inquiries/:id`                | Yes  | salesman, tenant_admin |
| 34  | GET    | `/tenant/dashboard`                      | Yes  | tenant_admin           |
| 35  | GET    | `/tenant/salesmen`                       | Yes  | tenant_admin           |
| 36  | POST   | `/tenant/salesmen`                       | Yes  | tenant_admin           |
| 37  | DELETE | `/tenant/salesmen/:id`                   | Yes  | tenant_admin           |
| 38  | GET    | `/tenant/listings`                       | Yes  | tenant_admin           |
| 39  | GET    | `/tenant/profile`                        | Yes  | tenant_admin           |
| 40  | PUT    | `/tenant/profile`                        | Yes  | tenant_admin           |
| 41  | GET    | `/tenant/subscription`                   | Yes  | tenant_admin           |
| 42  | POST   | `/tenant/subscription/upgrade`           | Yes  | tenant_admin           |
| 43  | GET    | `/tenant/inquiries`                      | Yes  | tenant_admin           |
| 44  | GET    | `/admin/dashboard`                       | Yes  | platform_admin         |
| 45  | GET    | `/admin/tenants`                         | Yes  | platform_admin         |
| 46  | POST   | `/admin/tenants`                         | Yes  | platform_admin         |
| 47  | GET    | `/admin/tenants/:id`                     | Yes  | platform_admin         |
| 48  | POST   | `/admin/tenants/:id/suspend`             | Yes  | platform_admin         |
| 49  | POST   | `/admin/tenants/:id/activate`            | Yes  | platform_admin         |
| 50  | PUT    | `/admin/tenants/:id/plan`                | Yes  | platform_admin         |
| 51  | PUT    | `/admin/tenants/:id/subscription`        | Yes  | platform_admin         |
| 52  | GET    | `/admin/listings/pending`                | Yes  | platform_admin         |
| 53  | POST   | `/admin/listings/:id/approve`            | Yes  | platform_admin         |
| 54  | POST   | `/admin/listings/:id/reject`             | Yes  | platform_admin         |
| 55  | GET    | `/admin/listings`                        | Yes  | platform_admin         |
| 56  | GET    | `/admin/audit-logs`                      | Yes  | platform_admin         |
| 57  | GET    | `/admin/subscription-plans`              | Yes  | platform_admin         |
| 58  | POST   | `/admin/subscription-plans`              | Yes  | platform_admin         |
| 59  | PUT    | `/admin/subscription-plans/:id`          | Yes  | platform_admin         |
| 60  | GET    | `/admin/property-types`                  | Yes  | platform_admin         |
| 61  | POST   | `/admin/property-types`                  | Yes  | platform_admin         |
| 62  | PUT    | `/admin/property-types/:id`              | Yes  | platform_admin         |
| 63  | GET    | `/admin/facilities`                      | Yes  | platform_admin         |
| 64  | POST   | `/admin/facilities`                      | Yes  | platform_admin         |
| 65  | PUT    | `/admin/facilities/:id`                  | Yes  | platform_admin         |
| 66  | GET    | `/admin/locations`                       | Yes  | platform_admin         |
| 67  | POST   | `/admin/locations`                       | Yes  | platform_admin         |
| 68  | PUT    | `/admin/locations/:id`                   | Yes  | platform_admin         |

**Total: 68 endpoints**

---

## 13. Change Log from v1 (46 endpoints)

| Change                  | Detail                                                                                                                                                              |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **+22 new endpoints**   | inquiries (buyer+salesman+tenant), property-types CRUD, facilities CRUD, locations CRUD, subscription-plans CRUD, property view log, tenant subscription by plan_id |
| **Renamed fields**      | `main_photo_url` → `watermarked_image_url`, `saved_properties` → `favorites`                                                                                        |
| **New response fields** | `watermark_status`, `is_primary`, `rent_period`, `whatsapp_number`, `show_whatsapp`, `property_type_id`, `location_id`, `facility_ids`                              |
| **New request fields**  | `bank_auction` nested object, `company_asset` nested object, `facility_ids` array                                                                                   |
| **Deprecated**          | `property_listings` → `properties`, `property_photos` → `property_images`, `saved_properties` → `favorites`                                                         |

---

**📄 API Contract v2 complete.** 68 endpoints, 10 public + 5 buyer + 16 salesman + 10 tenant + 25 admin.
