# Postman Collection Reference — MVP

## Multi-Tenant Property Information System

| Property          | Value                                                                       |
| ----------------- | --------------------------------------------------------------------------- |
| **Document Type** | Postman Collection Reference                                                |
| **Version**       | 1.0.0 MVP                                                                   |
| **Date**          | 2026-06-26                                                                  |
| **Base URL**      | `http://localhost:8080/api/v1` (local) / `http://localhost/api/v1` (Docker) |
| **Reference**     | `08-API-Contract.md`                                                        |

---

## 1. Collection Setup

### 1.1 Variables

| Variable             | Local Value                    | Docker Value              |
| -------------------- | ------------------------------ | ------------------------- |
| `{{base_url}}`       | `http://localhost:8080/api/v1` | `http://localhost/api/v1` |
| `{{token_admin}}`    | _(auto-filled by login)_       | _(auto-filled by login)_  |
| `{{token_tenant}}`   | _(auto-filled by login)_       | _(auto-filled by login)_  |
| `{{token_salesman}}` | _(auto-filled by login)_       | _(auto-filled by login)_  |
| `{{token_buyer}}`    | _(auto-filled by login)_       | _(auto-filled by login)_  |
| `{{listing_id}}`     | _(auto-filled by create)_      | _(auto-filled by create)_ |
| `{{tenant_id}}`      | _(auto-filled by list)_        | _(auto-filled by list)_   |

### 1.2 Collection Pre-request Script

```javascript
// No default pre-request — auth handled per-folder
```

### 1.3 Collection Post-response Script

```javascript
// Global: log response status for debugging
console.log(`[${pm.response.code}] ${pm.request.method} ${pm.request.url}`);
```

---

## 2. Folder: Auth (No Auth)

### 🔹 Register Buyer

```http
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "name": "Test Buyer",
  "email": "testbuyer@email.com",
  "phone": "081300000099",
  "password": "Buyer@123"
}
```

**Post-response Script:**

```javascript
if (pm.response.code === 201) {
  console.log("✅ Buyer registered:", pm.response.json().data.email);
}
```

---

### 🔹 Login (Each Role)

```http
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "admin@propertyhub.id",
  "password": "Admin@123"
}
```

**Post-response Script:**

```javascript
if (pm.response.code === 200) {
  const data = pm.response.json().data;
  pm.collectionVariables.set("token_admin", data.token);
  console.log("✅ Logged in as:", data.user.role, "-", data.user.name);
}
```

**Test logins to save as separate requests:**

| Request Name           | Email                  | Password    | Variable Set         |
| ---------------------- | ---------------------- | ----------- | -------------------- |
| Login — Platform Admin | `admin@propertyhub.id` | `Admin@123` | `{{token_admin}}`    |
| Login — Tenant Admin   | `budi@propertijaya.id` | `Budi@123`  | `{{token_tenant}}`   |
| Login — Salesman       | `andi@propertijaya.id` | `Andi@123`  | `{{token_salesman}}` |
| Login — Buyer          | `rina@email.com`       | `Rina@123`  | `{{token_buyer}}`    |

---

## 3. Folder: Public (No Auth)

### 🔹 List Properties

```http
GET {{base_url}}/properties?page=1&per_page=12&city=Jakarta+Selatan&property_type=house&listing_type=sale
```

**Post-response Script:**

```javascript
if (pm.response.code === 200) {
  const listings = pm.response.json().data;
  if (listings.length > 0) {
    pm.collectionVariables.set("listing_id", listings[0].id);
  }
  pm.test("Has listings", () =>
    pm.expect(listings.length).to.be.greaterThan(0),
  );
}
```

---

### 🔹 Property Detail

```http
GET {{base_url}}/properties/{{listing_id}}
```

**Tests:**

```javascript
pm.test("Status 200", () => pm.response.to.have.status(200));
pm.test("Has title", () => pm.expect(pm.response.json().data.title).to.exist);
pm.test(
  "Has salesman",
  () => pm.expect(pm.response.json().data.salesman).to.exist,
);
```

---

### 🔹 Featured Properties

```http
GET {{base_url}}/properties/featured?city=Jakarta+Selatan&limit=6
```

---

### 🔹 Nearby Properties

```http
GET {{base_url}}/properties/nearby?latitude=-6.2431&longitude=106.7988&radius_km=5
```

---

### 🔹 List Cities

```http
GET {{base_url}}/locations/cities
```

---

## 4. Folder: Me (Auth Required)

### 🔹 Get Profile

```http
GET {{base_url}}/me/profile
Authorization: Bearer {{token_buyer}}
```

**Tests:**

```javascript
pm.test("Status 200", () => pm.response.to.have.status(200));
pm.test("No password field", () => {
  const data = pm.response.json().data;
  pm.expect(data.password_hash).to.be.undefined;
  pm.expect(data.password).to.be.undefined;
});
```

---

### 🔹 Update Profile

```http
PUT {{base_url}}/me/profile
Authorization: Bearer {{token_buyer}}
Content-Type: application/json

{
  "name": "Rina Wijaya Updated",
  "phone": "081300000099"
}
```

---

## 5. Folder: Buyer (Role: buyer)

### 🔹 List Saved Properties

```http
GET {{base_url}}/me/saved?page=1&per_page=20
Authorization: Bearer {{token_buyer}}
```

---

### 🔹 Save Property

```http
POST {{base_url}}/me/saved/{{listing_id}}
Authorization: Bearer {{token_buyer}}
```

---

### 🔹 Remove Saved Property

```http
DELETE {{base_url}}/me/saved/{{listing_id}}
Authorization: Bearer {{token_buyer}}
```

---

## 6. Folder: Salesman (Role: salesman, tenant_admin)

### 🔹 Dashboard

```http
GET {{base_url}}/salesman/dashboard
Authorization: Bearer {{token_salesman}}
```

**Tests:**

```javascript
pm.test("Status 200", () => pm.response.to.have.status(200));
pm.test(
  "Has quota info",
  () => pm.expect(pm.response.json().data.quota).to.exist,
);
```

---

### 🔹 List My Listings

```http
GET {{base_url}}/salesman/listings?status=draft&page=1&per_page=20
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Create Listing (Draft)

```http
POST {{base_url}}/salesman/listings
Authorization: Bearer {{token_salesman}}
Content-Type: application/json

{
  "title": "Rumah Test di Jakarta Selatan",
  "description": "Rumah testing untuk API validation.",
  "price": 1500000000,
  "listing_type": "sale",
  "property_type": "house",
  "source_type": "regular",
  "address": "Jl. Test No. 123",
  "city": "Jakarta Selatan",
  "province": "DKI Jakarta",
  "land_area": 150,
  "building_area": 100,
  "bedrooms": 3,
  "bathrooms": 2,
  "floors": 1,
  "certificate_type": "SHM",
  "facilities": {
    "carport": true,
    "garden": true
  }
}
```

**Post-response Script:**

```javascript
if (pm.response.code === 201) {
  pm.collectionVariables.set("listing_id", pm.response.json().data.id);
  console.log("✅ Created listing:", pm.response.json().data.title);
}
```

---

### 🔹 Get My Listing Detail

```http
GET {{base_url}}/salesman/listings/{{listing_id}}
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Update Listing

```http
PUT {{base_url}}/salesman/listings/{{listing_id}}
Authorization: Bearer {{token_salesman}}
Content-Type: application/json

{
  "title": "Rumah Test Updated",
  "price": 1600000000
}
```

---

### 🔹 Submit for Review

```http
POST {{base_url}}/salesman/listings/{{listing_id}}/submit
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Deactivate

```http
POST {{base_url}}/salesman/listings/{{listing_id}}/deactivate
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Mark Sold

```http
POST {{base_url}}/salesman/listings/{{listing_id}}/mark-sold
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Mark Rented

```http
POST {{base_url}}/salesman/listings/{{listing_id}}/mark-rented
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Upload Photos

```http
POST {{base_url}}/salesman/listings/{{listing_id}}/photos
Authorization: Bearer {{token_salesman}}
Content-Type: multipart/form-data

# Form fields:
# photos: [file1.jpg, file2.jpg]
```

---

### 🔹 Delete Photo

```http
DELETE {{base_url}}/salesman/listings/{{listing_id}}/photos/:photoId
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Reorder Photos

```http
PUT {{base_url}}/salesman/listings/{{listing_id}}/photos/reorder
Authorization: Bearer {{token_salesman}}
Content-Type: application/json

{
  "photo_ids": [
    "photo-uuid-3",
    "photo-uuid-1",
    "photo-uuid-2"
  ]
}
```

---

### 🔹 Get Quota

```http
GET {{base_url}}/salesman/quota
Authorization: Bearer {{token_salesman}}
```

---

### 🔹 Delete Listing

```http
DELETE {{base_url}}/salesman/listings/{{listing_id}}
Authorization: Bearer {{token_salesman}}
```

---

## 7. Folder: Tenant Admin (Role: tenant_admin)

### 🔹 Dashboard

```http
GET {{base_url}}/tenant/dashboard
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 Get Profile

```http
GET {{base_url}}/tenant/profile
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 Update Profile

```http
PUT {{base_url}}/tenant/profile
Authorization: Bearer {{token_tenant}}
Content-Type: application/json

{
  "phone": "0215558888",
  "description": "Updated description"
}
```

---

### 🔹 List Salesmen

```http
GET {{base_url}}/tenant/salesmen?page=1&per_page=20
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 Add Salesman

```http
POST {{base_url}}/tenant/salesmen
Authorization: Bearer {{token_tenant}}
Content-Type: application/json

{
  "name": "Salesman Baru",
  "email": "salesman.baru@propertijaya.id",
  "phone": "081200000099",
  "password": "Sales@123"
}
```

---

### 🔹 Remove Salesman

```http
DELETE {{base_url}}/tenant/salesmen/:salesmanId
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 List Tenant Listings

```http
GET {{base_url}}/tenant/listings?status=approved&page=1&per_page=20
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 View Subscription

```http
GET {{base_url}}/tenant/subscription
Authorization: Bearer {{token_tenant}}
```

---

### 🔹 Request Upgrade

```http
POST {{base_url}}/tenant/subscription/upgrade
Authorization: Bearer {{token_tenant}}
Content-Type: application/json

{
  "plan_type": "premium"
}
```

---

## 8. Folder: Platform Admin (Role: platform_admin)

### 🔹 Dashboard

```http
GET {{base_url}}/admin/dashboard
Authorization: Bearer {{token_admin}}
```

---

### 🔹 List Tenants

```http
GET {{base_url}}/admin/tenants?page=1&per_page=20
Authorization: Bearer {{token_admin}}
```

**Post-response Script:**

```javascript
if (pm.response.code === 200) {
  const tenants = pm.response.json().data;
  if (tenants.length > 0) {
    pm.collectionVariables.set("tenant_id", tenants[0].id);
  }
}
```

---

### 🔹 Create Tenant

```http
POST {{base_url}}/admin/tenants
Authorization: Bearer {{token_admin}}
Content-Type: application/json

{
  "organization_name": "Agensi Test Baru",
  "subdomain_slug": "agensitest",
  "admin_name": "Admin Test",
  "admin_email": "admin@agensitest.id",
  "admin_phone": "081100000099",
  "admin_password": "Admin@123",
  "plan_type": "free"
}
```

---

### 🔹 Get Tenant Detail

```http
GET {{base_url}}/admin/tenants/{{tenant_id}}
Authorization: Bearer {{token_admin}}
```

---

### 🔹 Suspend Tenant

```http
POST {{base_url}}/admin/tenants/{{tenant_id}}/suspend
Authorization: Bearer {{token_admin}}
```

---

### 🔹 Activate Tenant

```http
POST {{base_url}}/admin/tenants/{{tenant_id}}/activate
Authorization: Bearer {{token_admin}}
```

---

### 🔹 Change Plan

```http
PUT {{base_url}}/admin/tenants/{{tenant_id}}/plan
Authorization: Bearer {{token_admin}}
Content-Type: application/json

{
  "plan_type": "premium"
}
```

---

### 🔹 List Pending

```http
GET {{base_url}}/admin/listings/pending?page=1&per_page=20
Authorization: Bearer {{token_admin}}
```

**Post-response Script:**

```javascript
if (pm.response.code === 200) {
  const listings = pm.response.json().data;
  if (listings.length > 0) {
    pm.collectionVariables.set("pending_listing_id", listings[0].id);
  }
}
```

---

### 🔹 Approve Listing

```http
POST {{base_url}}/admin/listings/:id/approve
Authorization: Bearer {{token_admin}}
```

---

### 🔹 Reject Listing

```http
POST {{base_url}}/admin/listings/:id/reject
Authorization: Bearer {{token_admin}}
Content-Type: application/json

{
  "reason": "Foto tidak sesuai dengan deskripsi properti. Mohon unggah foto asli."
}
```

---

### 🔹 Audit Logs

```http
GET {{base_url}}/admin/audit-logs?action=approve&from=2026-06-01&to=2026-06-30&page=1&per_page=50
Authorization: Bearer {{token_admin}}
```

---

## 9. Test Flow (Postman Runner / Newman)

### 9.1 Full E2E Test Sequence

Run requests in this order using Postman Collection Runner:

```
1.  Auth / Login — Platform Admin    → save {{token_admin}}
2.  Auth / Login — Tenant Admin      → save {{token_tenant}}
3.  Auth / Login — Salesman          → save {{token_salesman}}
4.  Auth / Login — Buyer             → save {{token_buyer}}

5.  Public / List Properties         → save {{listing_id}}
6.  Public / Property Detail         → verify
7.  Public / Featured                → verify
8.  Public / Cities                  → verify

9.  Salesman / Dashboard            → verify
10. Salesman / Create Listing       → save {{listing_id}}
11. Salesman / Get Listing Detail   → verify
12. Salesman / Update Listing       → verify
13. Salesman / Submit for Review    → verify status=pending

14. Admin / List Pending            → save listing
15. Admin / Approve Listing         → verify status=approved

16. Public / List Properties        → verify new listing appears

17. Buyer / Save Property           → verify 201
18. Buyer / List Saved              → verify listing in list
19. Buyer / Remove Saved            → verify removed

20. Admin / Audit Logs              → verify approve action logged
```

### 9.2 Newman CLI

```bash
# Export Postman collection as JSON, then:
npx newman run postman-collection.json \
  --env-var base_url=http://localhost:8080/api/v1 \
  --reporters cli,htmlextra
```

---

## 10. Collection Export Notes

### Folder Structure for Export

```
📁 PropertyHub MVP
├── 📁 1. Auth
│   ├── Register Buyer
│   ├── Login — Platform Admin
│   ├── Login — Tenant Admin
│   ├── Login — Salesman
│   └── Login — Buyer
├── 📁 2. Public
│   ├── List Properties
│   ├── Property Detail
│   ├── Featured Properties
│   ├── Nearby Properties
│   └── List Cities
├── 📁 3. Me (Profile)
│   ├── Get Profile
│   └── Update Profile
├── 📁 4. Buyer
│   ├── List Saved
│   ├── Save Property
│   └── Remove Saved
├── 📁 5. Salesman
│   ├── Dashboard
│   ├── List My Listings
│   ├── Create Listing
│   ├── Get Listing Detail
│   ├── Update Listing
│   ├── Submit for Review
│   ├── Deactivate
│   ├── Mark Sold
│   ├── Mark Rented
│   ├── Upload Photos
│   ├── Delete Photo
│   ├── Reorder Photos
│   ├── Get Quota
│   └── Delete Listing
├── 📁 6. Tenant Admin
│   ├── Dashboard
│   ├── Get/Update Profile
│   ├── List/Add/Remove Salesmen
│   ├── List Tenant Listings
│   ├── View Subscription
│   └── Request Upgrade
├── 📁 7. Platform Admin
│   ├── Dashboard
│   ├── List/Create/Get Tenants
│   ├── Suspend/Activate Tenant
│   ├── Change Plan
│   ├── List Pending
│   ├── Approve/Reject Listing
│   └── Audit Logs
└── 📁 8. E2E Flow
    └── Full Test Sequence (Runner)
```

---

_Dokumen ini adalah bagian dari Tahap 7. Impor struktur ini ke Postman untuk testing manual._
