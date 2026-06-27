# Security Requirements — MVP

## Multi-Tenant Property Information System

| Property          | Value                                      |
| ----------------- | ------------------------------------------ |
| **Document Type** | Security Requirements Specification        |
| **Version**       | 1.0.0 MVP                                  |
| **Date**          | 2026-06-26                                 |
| **Reference**     | `02-SRS-MVP.md`, `03-Permission-Matrix.md` |

---

## 1. Security Principles

| Principle              | Description                                                               |
| ---------------------- | ------------------------------------------------------------------------- |
| **Defense in Depth**   | Multiple layers of security: transport, application, data, infrastructure |
| **Least Privilege**    | Every user and service gets only the minimum permissions needed           |
| **Fail Securely**      | Errors and exceptions default to denying access, not granting it          |
| **Never Trust Client** | All authorization, validation, and quota checks happen server-side        |
| **Secure by Default**  | Default configurations are locked down; explicit action needed to open    |
| **Data Isolation**     | Strict tenant-level data separation enforced at query layer               |

---

## 2. Authentication Requirements

### 2.1 Password Policy

| Requirement | Detail                                                                      |
| ----------- | --------------------------------------------------------------------------- |
| **SEC-A01** | Minimum password length: 8 characters                                       |
| **SEC-A02** | Must contain at least: 1 uppercase letter OR 1 digit OR 1 special character |
| **SEC-A03** | Password hashed with bcrypt, cost factor 12                                 |
| **SEC-A04** | Password never logged, never returned in API response                       |
| **SEC-A05** | Password reset not in MVP scope (manual admin reset)                        |

### 2.2 JWT Configuration

| Requirement | Detail                                                                                                         |
| ----------- | -------------------------------------------------------------------------------------------------------------- |
| **SEC-A06** | JWT signing algorithm: HS256                                                                                   |
| **SEC-A07** | JWT secret: 256-bit random string, stored in environment variable `JWT_SECRET`                                 |
| **SEC-A08** | JWT expiry: 24 hours                                                                                           |
| **SEC-A09** | JWT payload: `{ sub: user_id, role: role_code, tenant_id: tenant_id \| null, exp: timestamp, iat: timestamp }` |
| **SEC-A10** | JWT sent via `Authorization: Bearer <token>` header                                                            |
| **SEC-A11** | No JWT refresh token in MVP (re-login after expiry)                                                            |

### 2.3 Login Security

| Requirement | Detail                                                                                                                         |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------ |
| **SEC-A12** | Rate limiting on `/api/v1/auth/login`: max 5 attempts per IP per 1 minute                                                      |
| **SEC-A13** | Rate limit exceeded → 429 Too Many Requests with message: "Terlalu banyak percobaan login. Silakan coba lagi dalam 1 menit."   |
| **SEC-A14** | Generic error message on failed login: "Email atau password salah." (no distinction between invalid email vs invalid password) |
| **SEC-A15** | Account lockout not in MVP scope                                                                                               |

---

## 3. Authorization Requirements (RBAC)

### 3.1 Role-Based Access Control

| Requirement | Detail                                                                                                                   |
| ----------- | ------------------------------------------------------------------------------------------------------------------------ |
| **SEC-Z01** | All protected endpoints must pass through JWT middleware                                                                 |
| **SEC-Z02** | JWT middleware must verify signature, expiry, and extract claims                                                         |
| **SEC-Z03** | RBAC middleware must check `role` claim against endpoint's required role(s)                                              |
| **SEC-Z04** | Invalid JWT → 401 Unauthorized                                                                                           |
| **SEC-Z05** | Valid JWT but insufficient role → 403 Forbidden with message: "Anda tidak memiliki izin untuk mengakses resource ini."   |
| **SEC-Z06** | Suspended tenant users → 403 Forbidden with message: "Akun organisasi Anda sedang dinonaktifkan. Hubungi administrator." |

### 3.2 Tenant Data Isolation

| Requirement | Detail                                                                        |
| ----------- | ----------------------------------------------------------------------------- |
| **SEC-Z07** | All tenant-scoped queries MUST include `WHERE tenant_id = :current_tenant_id` |
| **SEC-Z08** | `tenant_id` extracted from JWT claims, never from client request body/params  |
| **SEC-Z09** | Salesman queries additionally scoped to `salesman_id = :current_user_id`      |
| **SEC-Z10** | Tenant Admin bypasses salesman scope but stays within tenant scope            |
| **SEC-Z11** | Platform Admin bypasses all tenant scoping                                    |
| **SEC-Z12** | Cross-tenant data access must be impossible through normal API usage          |

### 3.3 Verification Matrix

| Scenario                                                           | Expected Result                          |
| ------------------------------------------------------------------ | ---------------------------------------- |
| Salesman A (Tenant X) requests Salesman B's listing (Tenant X)     | 403 Forbidden                            |
| Salesman A (Tenant X) requests Salesman C's listing (Tenant Y)     | 403 Forbidden                            |
| Tenant Admin X requests listing from Tenant Y                      | 403 Forbidden                            |
| Salesman requests listing with `tenant_id` changed in request body | 403 (value from JWT, not body)           |
| Platform Admin requests any listing                                | 200 OK                                   |
| Guest requests listing detail with `status=pending`                | 404 Not Found (only `approved` returned) |

---

## 4. Input Validation & Sanitization

### 4.1 General Rules

| Requirement | Detail                                                               |
| ----------- | -------------------------------------------------------------------- |
| **SEC-V01** | All input validated server-side regardless of client-side validation |
| **SEC-V02** | Whitelist approach: define what is allowed, reject everything else   |
| **SEC-V03** | Validation happens at API handler level before business logic        |
| **SEC-V04** | Validation errors return 422 with structured error details           |

### 4.2 Field-Level Validation

| Field              | Rule                                                                  |
| ------------------ | --------------------------------------------------------------------- |
| `email`            | RFC 5322 format, max 255 chars, unique in system, trimmed, lowercased |
| `password`         | Min 8 chars, string                                                   |
| `name`             | Min 2 chars, max 100 chars, trimmed, no HTML tags                     |
| `phone`            | Regex `^\+?[0-9]{8,15}$`, trimmed                                     |
| `title` (listing)  | Min 5 chars, max 200 chars, trimmed                                   |
| `description`      | Max 5000 chars, no HTML tags (plain text)                             |
| `price`            | Positive number, max 999,999,999,999                                  |
| `property_type`    | Enum: house, land, apartment, shophouse, warehouse, office, villa     |
| `source_type`      | Enum: regular, bank_auction, company_asset                            |
| `listing_type`     | Enum: sale, rent                                                      |
| `city`             | Min 2 chars, max 100 chars                                            |
| `latitude`         | Range -90 to 90                                                       |
| `longitude`        | Range -180 to 180                                                     |
| `land_area`        | Positive number, max 999,999                                          |
| `building_area`    | Positive number, max 999,999                                          |
| `bedrooms`         | Integer 0–99                                                          |
| `bathrooms`        | Integer 0–99                                                          |
| `floors`           | Integer 0–200                                                         |
| `certificate_type` | Enum: SHM, SHGB, Girik, Lainnya                                       |
| `facilities`       | Valid JSON object, max 20 keys                                        |
| `sort_order`       | Integer ≥ 0                                                           |
| `reject_reason`    | Min 10 chars, max 500 chars                                           |

### 4.3 File Upload Validation

| Requirement | Detail                                                                        |
| ----------- | ----------------------------------------------------------------------------- |
| **SEC-V05** | Allowed MIME types: `image/jpeg`, `image/png`, `image/webp`                   |
| **SEC-V06** | Max file size: 5 MB per image                                                 |
| **SEC-V07** | Validate MIME type by magic bytes (file signature), not by extension          |
| **SEC-V08** | Reject files with double extensions (e.g., `photo.jpg.php`)                   |
| **SEC-V09** | Strip EXIF metadata from uploaded images                                      |
| **SEC-V10** | Max 10 photos per listing                                                     |
| **SEC-V11** | Generate random filename on storage (UUID-based), never use original filename |

### 4.4 SQL Injection Prevention

| Requirement | Detail                                                                        |
| ----------- | ----------------------------------------------------------------------------- |
| **SEC-V12** | All database queries use GORM with parameterized queries                      |
| **SEC-V13** | Never concatenate user input into raw SQL strings                             |
| **SEC-V14** | Dynamic `ORDER BY` or `GROUP BY` must use whitelist mapping (never raw input) |

---

## 5. Data Protection

### 5.1 Sensitive Data Handling

| Requirement | Detail                                                                          |
| ----------- | ------------------------------------------------------------------------------- |
| **SEC-D01** | Password hash only in `users.password_hash` column; never in responses          |
| **SEC-D02** | JWT secret only in environment variable; never in code, config files, or logs   |
| **SEC-D03** | Database credentials only in environment variables                              |
| **SEC-D04** | No PII (Personally Identifiable Information) in log files                       |
| **SEC-D05** | Error responses in production never expose stack traces, SQL, or internal paths |

### 5.2 Data at Rest

| Requirement | Detail                                                                   |
| ----------- | ------------------------------------------------------------------------ |
| **SEC-D06** | Passwords stored as bcrypt hash (one-way, not decryptable)               |
| **SEC-D07** | Database files on encrypted disk (infrastructure level)                  |
| **SEC-D08** | Soft delete for major entities — data retained, not physically destroyed |

### 5.3 Data in Transit

| Requirement | Detail                                                                                      |
| ----------- | ------------------------------------------------------------------------------------------- |
| **SEC-D09** | HTTPS enforced in production (TLS 1.2+)                                                     |
| **SEC-D10** | JWT transmitted only over HTTPS in production                                               |
| **SEC-D11** | CORS `Access-Control-Allow-Origin` set to explicit frontend domain, never `*` in production |

---

## 6. CORS Configuration

| Environment | Allowed Origin                 |
| ----------- | ------------------------------ |
| Development | `http://localhost:5173`        |
| Staging     | `https://staging.{domain}.com` |
| Production  | `https://{domain}.com`         |

| Requirement | Detail                                                          |
| ----------- | --------------------------------------------------------------- |
| **SEC-C01** | `Access-Control-Allow-Methods`: GET, POST, PUT, DELETE, OPTIONS |
| **SEC-C02** | `Access-Control-Allow-Headers`: Content-Type, Authorization     |
| **SEC-C03** | `Access-Control-Allow-Credentials`: true                        |
| **SEC-C04** | `Access-Control-Max-Age`: 86400 (24 hours)                      |

---

## 7. API Security

### 7.1 Headers

| Requirement   | Detail                                                   |
| ------------- | -------------------------------------------------------- |
| **SEC-API01** | `X-Content-Type-Options: nosniff`                        |
| **SEC-API02** | `X-Frame-Options: DENY`                                  |
| **SEC-API03** | `X-XSS-Protection: 1; mode=block`                        |
| **SEC-API04** | `Referrer-Policy: strict-origin-when-cross-origin`       |
| **SEC-API05** | `Content-Security-Policy` header configured for frontend |

### 7.2 Rate Limiting

| Endpoint                     | Limit | Window            |
| ---------------------------- | ----- | ----------------- |
| `POST /api/v1/auth/login`    | 5     | 1 minute per IP   |
| `POST /api/v1/auth/register` | 3     | 10 minutes per IP |
| All other endpoints          | 100   | 1 minute per IP   |

### 7.3 Request Size Limits

| Requirement   | Detail                                                       |
| ------------- | ------------------------------------------------------------ |
| **SEC-API06** | Max request body size: 10 MB                                 |
| **SEC-API07** | Max multipart form size: 55 MB (10 photos × 5 MB + overhead) |

---

## 8. Audit & Monitoring

| Requirement  | Detail                                                                                                                       |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| **SEC-AM01** | Log all authentication events: login success, login failure, logout                                                          |
| **SEC-AM02** | Log all CUD operations with actor identity                                                                                   |
| **SEC-AM03** | Log all approve/reject actions                                                                                               |
| **SEC-AM04** | Log all tenant suspend/activate actions                                                                                      |
| **SEC-AM05** | Log all plan changes                                                                                                         |
| **SEC-AM06** | Audit logs are append-only; no update or delete on audit records                                                             |
| **SEC-AM07** | Audit log entries include: timestamp, user_id, user_role, action, entity_type, entity_id, old_values, new_values, ip_address |

---

## 9. Infrastructure Security (MVP)

| Requirement | Detail                                                                                |
| ----------- | ------------------------------------------------------------------------------------- |
| **SEC-I01** | PostgreSQL port not exposed to public internet                                        |
| **SEC-I02** | Database credentials in environment variables, not in docker-compose.yml (use `.env`) |
| **SEC-I03** | Go backend runs as non-root user in Docker container                                  |
| **SEC-I04** | `.env` and `.env.*` files in `.gitignore`                                             |
| **SEC-I05** | `.env.example` provided with placeholder values only                                  |
| **SEC-I06** | No secrets committed to git repository                                                |

---

## 10. Security Checklist — Go-Live Gate

| #   | Check                                           | Status |
| --- | ----------------------------------------------- | ------ |
| 1   | bcrypt cost factor = 12                         | ☐      |
| 2   | JWT secret is 256-bit random string             | ☐      |
| 3   | JWT secret NOT hardcoded                        | ☐      |
| 4   | CORS whitelist set (no `*`)                     | ☐      |
| 5   | All queries use parameterized statements        | ☐      |
| 6   | Password never in API response                  | ☐      |
| 7   | Error responses hide internal details           | ☐      |
| 8   | Rate limiting on login endpoint                 | ☐      |
| 9   | File upload: MIME type validated by magic bytes | ☐      |
| 10  | EXIF metadata stripped from uploads             | ☐      |
| 11  | HTTPS enforced in production                    | ☐      |
| 12  | `.env` in `.gitignore`                          | ☐      |
| 13  | Security headers configured                     | ☐      |
| 14  | No hardcoded secrets in codebase                | ☐      |
| 15  | Tenant isolation verified (cross-tenant test)   | ☐      |

---

_Dokumen ini adalah bagian dari Tahap 1. Lanjut ke `05-Error-Handling-Standard.md`._
