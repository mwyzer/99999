# Product Requirements Document (PRD) - MVP
## Multi-Tenant Property Information System

| Property | Value |
|---|---|
| **Project Name** | Multi-Tenant Property Information System |
| **Version** | 1.0.0 MVP |
| **Author** | System |
| **Date** | 2026-06-26 |
| **Status** | Draft |

---

## 1. Executive Summary

### 1.1 Product Vision
A multi-tenant web platform that enables property agencies, banks, and companies to list, manage, and showcase their property portfolios (sale & rent). Each tenant operates independently with their own branding, sales team, and property listings. End users (buyers/renters) can browse approved listings, contact sales agents via WhatsApp, and discover nearby properties.

### 1.2 Problem Statement
- Property agencies need a simple, branded platform to showcase listings without building their own website.
- Banks need a channel to publish auction assets.
- Companies need to list corporate property assets.
- Buyers/renters want a unified experience to browse properties from multiple sources with direct WhatsApp contact.

### 1.3 Target Users
| # | Role | Description |
|---|------|-------------|
| 1 | **Guest** | Unauthenticated visitor browsing public listings |
| 2 | **Buyer / Renter** | Registered user who can save/bookmark properties |
| 3 | **Salesman** | Agent under a tenant who creates and manages listings |
| 4 | **Tenant Admin / Agency Owner** | Manages one tenant account, sales team, and subscription |
| 5 | **Platform Admin** | Super admin managing all tenants, approvals, and platform health |

---

## 2. User Stories (MVP)

### 2.1 Guest
| ID | Story | Priority |
|----|-------|----------|
| US-G01 | As a guest, I can browse approved property listings without login | P0 |
| US-G02 | As a guest, I can filter properties by type, location, price range, source_type | P0 |
| US-G03 | As a guest, I can view property detail with photos, description, specs | P0 |
| US-G04 | As a guest, I can see sales agent photo and tenant logo on property card | P1 |
| US-G05 | As a guest, I can click WhatsApp button on property card to contact sales | P0 |
| US-G06 | As a guest, I can view a carousel of featured properties by location | P1 |
| US-G07 | As a guest, I can allow browser location access to see nearest properties | P2 |
| US-G08 | As a guest, I can register as Buyer/Renter | P0 |

### 2.2 Buyer / Renter
| ID | Story | Priority |
|----|-------|----------|
| US-B01 | As a buyer, I can login with email & password | P0 |
| US-B02 | As a buyer, I can save/bookmark favorite properties | P1 |
| US-B03 | As a buyer, I can view my saved properties list | P1 |
| US-B04 | As a buyer, I can view my profile and edit basic info | P2 |

### 2.3 Salesman
| ID | Story | Priority |
|----|-------|----------|
| US-S01 | As a salesman, I can login and see my dashboard | P0 |
| US-S02 | As a salesman, I can create property listing (draft/pending) | P0 |
| US-S03 | As a salesman, I can upload property photos with automatic watermark | P0 |
| US-S04 | As a salesman, I can edit my own listings | P0 |
| US-S05 | As a salesman, I can delete/disable my own listings | P0 |
| US-S06 | As a salesman, I can see my active listing count vs quota | P0 |
| US-S07 | As a salesman, I cannot exceed quota limit (5 for Free plan) | P0 |
| US-S08 | As a salesman, I can view status of my listings (draft, pending, approved, rejected) | P0 |
| US-S09 | As a salesman, I can see my profile and edit basic info | P2 |

### 2.4 Tenant Admin / Agency Owner
| ID | Story | Priority |
|----|-------|----------|
| US-T01 | As a tenant admin, I can login and see agency dashboard | P0 |
| US-T02 | As a tenant admin, I can manage my sales team (add/remove salesmen) | P0 |
| US-T03 | As a tenant admin, I can enforce max 5 salesmen on Free plan | P0 |
| US-T04 | As a tenant admin, I can view all listings under my tenant | P0 |
| US-T05 | As a tenant admin, I can manage tenant profile (logo, name, contact info) | P0 |
| US-T06 | As a tenant admin, I can view subscription plan and quota usage | P0 |
| US-T07 | As a tenant admin, I can upgrade from Free to Premium plan | P1 |

### 2.5 Platform Admin
| ID | Story | Priority |
|----|-------|----------|
| US-P01 | As a platform admin, I can login to admin panel | P0 |
| US-P02 | As a platform admin, I can view all tenants and their status | P0 |
| US-P03 | As a platform admin, I can create new tenant accounts | P0 |
| US-P04 | As a platform admin, I can approve or reject property listings | P0 |
| US-P05 | As a platform admin, I can suspend/activate tenant accounts | P0 |
| US-P06 | As a platform admin, I can manage subscription plans | P1 |
| US-P07 | As a platform admin, I can view audit logs | P1 |

---

## 3. Functional Requirements (MVP)

### 3.1 Property Types
System supports the following property types:
- `house` — Rumah
- `land` — Tanah
- `apartment` — Apartemen
- `shophouse` — Ruko
- `warehouse` — Gudang
- `office` — Office Space
- `villa` — Villa

### 3.2 Source Types
- `regular` — Regular agency listing
- `bank_auction` — Bank auction asset (Lelang Bank)
- `company_asset` — Company property asset

### 3.3 Listing Statuses
| Status | Description | Counts toward quota |
|--------|-------------|---------------------|
| `draft` | Being created, not yet submitted | ✅ Yes |
| `pending` | Submitted for review | ✅ Yes |
| `approved` | Approved & publicly visible | ✅ Yes |
| `rejected` | Rejected by platform admin | ❌ No |
| `sold` | Property sold | ❌ No |
| `rented` | Property rented out | ❌ No |
| `inactive` | Manually deactivated | ❌ No |
| `deleted` | Soft deleted | ❌ No |

### 3.4 Quota System
| Plan | Max Salesmen | Max Active Listings per Salesman |
|------|-------------|----------------------------------|
| **Free** | 5 | 5 |
| **Premium** | Unlimited | Unlimited |

### 3.5 Property Card Display
Each property card must include:
- Property main photo (with watermark)
- Property title
- Price
- Location
- Property type badge
- Source type badge
- Sales agent photo (thumbnail)
- Tenant logo/photo
- WhatsApp contact button

### 3.6 Watermark
- All property photos on listing cards must have a semi-transparent watermark overlay.
- Watermark text: Tenant name or platform name.

### 3.7 WhatsApp Integration
- `https://wa.me/{phone_number}?text={encoded_message}`
- Pre-filled message format: "Halo, saya tertarik dengan properti {property_title} yang saya lihat di {platform_name}"

### 3.8 Location Features
- Carousel: "Properti Pilihan di [Lokasi]" based on user's selected location.
- Nearest properties: Browser Geolocation API to get user coordinates, then query properties near that location.

---

## 4. Non-Functional Requirements

### 4.1 Performance
- API response time < 500ms for listing queries
- Image optimization (thumbnail generation)
- Pagination on list endpoints (default 20 items/page)

### 4.2 Security
- bcrypt password hashing (cost factor 12)
- JWT authentication with environment-variable secret
- Role-based access control (RBAC)
- CORS whitelist configuration
- No password in any API response
- GORM prepared statements (parameterized queries)
- Input validation on all endpoints
- Audit logging for CUD + approve/reject operations

### 4.3 Scalability
- Multi-tenant data isolation (tenant_id on all tenant-scoped tables)
- Stateless API design for horizontal scaling

### 4.4 Reliability
- Database transactions for critical operations
- Graceful error handling

---

## 5. MVP Scope Boundaries

### IN SCOPE (MVP)
- All 5 user roles
- All 7 property types + 3 source types
- Free & Premium plans
- CRUD listings
- Public property browsing
- WhatsApp contact
- Photo watermark
- JWT auth + RBAC
- Admin approval workflow
- Audit logging
- Responsive web frontend (React)

### OUT OF SCOPE (Future)
- Payment gateway integration
- Real-time chat
- Property comparison
- Advanced analytics dashboard
- Mobile native apps (React Native/Flutter)
- Email/SMS notifications
- Property visit scheduling
- Multi-language (i18n)
- Social media sharing

---

## 6. Technology Stack

| Layer | Technology |
|-------|-----------|
| **Backend Language** | Go 1.22+ |
| **Web Framework** | Gin |
| **ORM** | GORM |
| **Database** | PostgreSQL 16 |
| **Auth** | JWT (golang-jwt) |
| **Password** | bcrypt |
| **Frontend** | React 18 + Vite |
| **Styling** | Tailwind CSS 3 |
| **HTTP Client** | Axios / Fetch |
| **Deployment** | Docker + Docker Compose |

---

## 7. Success Metrics (MVP)

| Metric | Target |
|--------|--------|
| Tenant onboarding | < 5 minutes |
| Listing creation | < 3 minutes |
| Property search load | < 2 seconds |
| System uptime | 99.5% |
| Zero critical security vulnerabilities | 100% |

---

## 8. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Photo storage scalability | High | Use local storage with path reference; plan S3 migration later |
| Tenant data leak | Critical | Strict tenant_id filtering on all queries |
| Quota bypass | Medium | Server-side quota enforcement, never trust client |
| JWT secret exposure | Critical | Environment variable only; never hardcode |

---

*Dokumen ini adalah Tahap 1 dari 11 tahap. Lanjut ke Tahap 2 setelah review.*
