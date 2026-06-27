# 🏠 PropertyHub — Multi-Tenant Property Information System

<p align="center">
  <strong>Platform multi-tenant untuk listing properti dari agensi, bank, dan perusahaan di Indonesia.</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=white" alt="React">
  <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/License-Proprietary-red" alt="License">
</p>

---

## 📋 Daftar Isi

- [🚀 Quick Start (Docker)](#-quick-start-docker)
- [📋 Service Ports](#-service-ports)
- [🔑 Akun Demo](#-akun-demo)
- [📁 Struktur Proyek](#-struktur-proyek)
- [🛠️ Tech Stack](#️-tech-stack)
- [🐳 Docker Development](#-docker-development)
- [🔧 Local Development (Tanpa Docker)](#-local-development-tanpa-docker)
- [⚙️ Environment Variables](#️-environment-variables)
- [📖 Dokumentasi Lengkap](#-dokumentasi-lengkap)
- [🗄️ Database](#️-database)
- [🔐 Authentication & Authorization](#-authentication--authorization)
- [📡 API Overview](#-api-overview)
- [🧪 Testing](#-testing)
- [🚀 Deployment](#-deployment)
- [🔍 Troubleshooting](#-troubleshooting)
- [📄 License](#-license)

---

## 🚀 Quick Start (Docker)

Cara termudah menjalankan seluruh stack:

```bash
# 1. Clone & masuk ke folder proyek
cd 99999

# 2. Salin environment variables
copy .env.docker .env

# 3. Jalankan semua service (PostgreSQL + Backend + Frontend)
docker compose up -d

# 4. Cek status service
docker compose ps

# 5. Akses aplikasi
# Frontend : http://localhost
# Backend  : http://localhost:8080
# Health   : http://localhost:8080/health
```

Untuk menghentikan semua service:

```bash
docker compose down
```

Untuk menghentikan dan menghapus volume (reset database):

```bash
docker compose down -v
```

---

## 📋 Service Ports

| Service                  | Port | URL                     | Keterangan                        |
| ------------------------ | ---- | ----------------------- | --------------------------------- |
| Frontend (React + Nginx) | 80   | `http://localhost`      | SPA static files via Nginx        |
| Backend (Go + Gin)       | 8080 | `http://localhost:8080` | REST API server                   |
| PostgreSQL               | 5432 | `localhost:5432`        | Database (bisa diakses dari host) |

---

## 🔑 Akun Demo

Saat pertama kali dijalankan di environment `development`, backend akan otomatis melakukan seed data berikut:

| Role               | Email                  | Password    | Tenant            |
| ------------------ | ---------------------- | ----------- | ----------------- |
| **Platform Admin** | `admin@propertyhub.id` | `Admin@123` | N/A (Super Admin) |
| **Tenant Admin**   | `budi@propertijaya.id` | `Budi@123`  | PropertiJaya      |
| **Salesman**       | `andi@propertijaya.id` | `Andi@123`  | PropertiJaya      |
| **Buyer**          | `rina@email.com`       | `Rina@123`  | N/A (Public)      |

---

## 📁 Struktur Proyek

```
99999/
├── 01-PRD-MVP.md                 # Product Requirements Document
├── 02-SRS-MVP.md                 # Software Requirements Specification
├── 03-Permission-Matrix.md       # RBAC permission matrix
├── 04-Security-Requirements.md   # Security specification
├── 05-Error-Handling-Standard.md # Error handling standard
├── 06-MVP-Scope-Acceptance.md    # Scope & acceptance criteria
├── 07-ERD-Database-Schema.md     # Database design (ERD, schema, indexes)
├── 08-API-Contract.md            # API specification (46 endpoints)
├── 09-Testing-QA-Plan.md         # Testing & QA plan
├── 10-Postman-Collection.md      # Postman collection documentation
├── 11-Business-Workflows.md      # Business workflow documentation
├── README.md                     # This file
├── postman-collection.json       # Postman collection (importable)
├── .env.docker                   # Environment template for Docker
├── docker-compose.yml            # Docker Compose (production)
├── docker-compose.dev.yml        # Docker Compose (development override)
├── docker-compose.override.yml   # Docker Compose (auto-applied safe defaults)
├── db/
│   ├── init.sql                  # PostgreSQL extensions
│   └── migration.sql             # Full DDL migration script
├── backend/                      # Go 1.22+ | Gin | GORM
│   ├── main.go                   # Entry point
│   ├── Dockerfile                # Multi-stage production build
│   ├── Dockerfile.dev            # Development hot-reload
│   ├── .dockerignore
│   ├── .env.example              # Local environment template
│   ├── config/                   # Configuration loader (env vars)
│   ├── database/                 # DB connection, auto-migration, seed
│   ├── models/                   # 7 GORM models
│   ├── handlers/                 # 5 handler groups (public, buyer, salesman, tenant, platform)
│   ├── middleware/               # Auth, RBAC, CORS, Scope, Rate Limit, Security Headers
│   ├── routes/                   # 46 route definitions
│   ├── services/                 # Business logic layer
│   ├── repositories/             # Data access layer
│   ├── dto/                      # Request/Response structs
│   └── utils/                    # JWT, Password, Response, Pagination, Image
└── frontend/                     # React 18 | Vite | Tailwind CSS
    ├── Dockerfile                # Multi-stage build (Node → Nginx)
    ├── Dockerfile.dev            # Development hot-reload (Vite)
    ├── .dockerignore
    ├── nginx.conf                # Nginx config (SPA + API proxy)
    ├── index.html
    ├── package.json
    ├── vite.config.js            # Vite config with API proxy
    └── src/
        ├── main.jsx              # React entry point
        ├── App.jsx               # Root component + routes
        ├── api/                  # Axios instance + 7 API modules
        ├── context/              # AuthContext (React Context API)
        ├── hooks/                # Custom React hooks
        ├── components/           # Navbar, Footer, PropertyCard, PropertyFilter, Loading
        └── pages/                # 10 pages (Home, Detail, Login, Register, Profile, Dashboards, etc.)
```

---

## 🛠️ Tech Stack

| Layer           | Technology                                       |
| --------------- | ------------------------------------------------ |
| **Backend**     | Go 1.22+, Gin v1.10, GORM v1.25                  |
| **Frontend**    | React 18, Vite 5, Tailwind CSS 3, React Router 6 |
| **Database**    | PostgreSQL 16                                    |
| **Auth**        | JWT (HS256), bcrypt (cost 12)                    |
| **Validation**  | go-playground/validator                          |
| **HTTP Client** | Axios                                            |
| **Icons**       | Lucide React                                     |
| **Toast**       | React Hot Toast                                  |
| **Infra**       | Docker, Docker Compose, Nginx 1.27               |

---

## 🐳 Docker Development

### Mode Development (Hot-Reload)

Untuk development dengan hot-reload (perubahan kode langsung terlihat):

```bash
# Jalankan dengan development override
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Lihat logs
docker compose logs -f backend
docker compose logs -f frontend
```

Saat mode development:

- Backend menggunakan `go run .` — restart otomatis? Gunakan [Air](https://github.com/air-verse/air) untuk live-reload.
- Frontend menggunakan Vite dev server di port 5173 — hot-reload bawaan.
- Source code di-mount sebagai volume — edit di host, terlihat di container.

### Mode Production

```bash
# Build & jalankan production
docker compose up -d --build

# Atau build dulu, lalu jalankan
docker compose build
docker compose up -d
```

### Mendapatkan Shell di Container

```bash
# Backend
docker compose exec backend sh

# Frontend
docker compose exec frontend sh

# PostgreSQL
docker compose exec postgres psql -U propertyhub -d propertyhub
```

### Membersihkan Docker Resources

```bash
# Hentikan container
docker compose down

# Hentikan + hapus volumes (reset DB & uploads)
docker compose down -v

# Hentikan + hapus volumes + hapus images
docker compose down -v --rmi all
```

---

## 🔧 Local Development (Tanpa Docker)

### Prerequisites

- **Go** 1.22+
- **Node.js** 20+
- **PostgreSQL** 16
- **npm** 9+

### Backend

```bash
cd backend

# Setup environment
cp .env.example .env
# Edit .env — sesuaikan DB_HOST, DB_USER, DB_PASSWORD

# Install dependencies
go mod tidy

# Jalankan
go run main.go        # → http://localhost:8080
```

### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Jalankan dev server
npm run dev           # → http://localhost:5173 (proxy API ke :8080)
```

### Database

```bash
# Pastikan PostgreSQL 16 berjalan
# Buat database
psql -U postgres -c "CREATE DATABASE propertyhub;"
psql -U postgres -c "CREATE USER propertyhub WITH PASSWORD 'propertyhub_secret';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE propertyhub TO propertyhub;"

# Jalankan migrasi (opsional — backend auto-migrate)
psql -U propertyhub -d propertyhub -f db/migration.sql

# Backend akan auto-migrate & seed saat startup (APP_ENV=development)
```

---

## ⚙️ Environment Variables

### Application

| Variable   | Default                 | Deskripsi                                |
| ---------- | ----------------------- | ---------------------------------------- |
| `APP_NAME` | `PropertyHub`           | Nama aplikasi                            |
| `APP_ENV`  | `development`           | Environment (`development`/`production`) |
| `APP_PORT` | `8080`                  | Port backend                             |
| `APP_URL`  | `http://localhost:8080` | Base URL backend                         |

### Database

| Variable               | Default              | Deskripsi                          |
| ---------------------- | -------------------- | ---------------------------------- |
| `DB_HOST`              | `localhost`          | Host database (Docker: `postgres`) |
| `DB_PORT`              | `5432`               | Port database                      |
| `DB_USER`              | `propertyhub`        | Username database                  |
| `DB_PASSWORD`          | `propertyhub_secret` | Password database                  |
| `DB_NAME`              | `propertyhub`        | Nama database                      |
| `DB_SSLMODE`           | `disable`            | SSL mode                           |
| `DB_TIMEZONE`          | `Asia/Jakarta`       | Timezone database                  |
| `DB_MAX_OPEN_CONNS`    | `50`                 | Maksimum koneksi terbuka           |
| `DB_MAX_IDLE_CONNS`    | `10`                 | Maksimum koneksi idle              |
| `DB_CONN_MAX_LIFETIME` | `30`                 | Lifetime koneksi (menit)           |

### JWT & Security

| Variable           | Default                   | Deskripsi              |
| ------------------ | ------------------------- | ---------------------- |
| `JWT_SECRET`       | `change-me-in-production` | Secret key JWT (HS256) |
| `JWT_EXPIRY_HOURS` | `24`                      | JWT expiry dalam jam   |
| `BCRYPT_COST`      | `12`                      | Bcrypt hashing cost    |

### CORS

| Variable               | Default                 | Deskripsi                                |
| ---------------------- | ----------------------- | ---------------------------------------- |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Origins yang diizinkan (comma-separated) |

### Upload

| Variable                 | Default     | Deskripsi                   |
| ------------------------ | ----------- | --------------------------- |
| `UPLOAD_DIR`             | `./uploads` | Direktori upload foto       |
| `MAX_UPLOAD_SIZE_MB`     | `5`         | Maksimum ukuran upload (MB) |
| `MAX_PHOTOS_PER_LISTING` | `10`        | Maksimum foto per listing   |

### Rate Limiting

| Variable                       | Default | Deskripsi                         |
| ------------------------------ | ------- | --------------------------------- |
| `RATE_LIMIT_LOGIN_PER_MINUTE`  | `5`     | Maksimum login attempt per menit  |
| `RATE_LIMIT_GLOBAL_PER_MINUTE` | `100`   | Maksimum request global per menit |

### PostgreSQL Container

| Variable            | Default              | Deskripsi          |
| ------------------- | -------------------- | ------------------ |
| `POSTGRES_USER`     | `propertyhub`        | Postgres superuser |
| `POSTGRES_PASSWORD` | `propertyhub_secret` | Postgres password  |
| `POSTGRES_DB`       | `propertyhub`        | Default database   |

---

## 📖 Dokumentasi Lengkap

Lihat file `.md` di root proyek untuk dokumentasi lengkap:

| File                              | Konten                                                                         |
| --------------------------------- | ------------------------------------------------------------------------------ |
| **01-PRD-MVP.md**                 | Product Requirements Document — visi, user stories, prioritas                  |
| **02-SRS-MVP.md**                 | Software Requirements Specification — functional & non-functional requirements |
| **03-Permission-Matrix.md**       | Matriks RBAC — 4 role × 46 endpoint                                            |
| **04-Security-Requirements.md**   | Spesifikasi keamanan — JWT, CORS, rate limiting, input validation              |
| **05-Error-Handling-Standard.md** | Standar format response error                                                  |
| **06-MVP-Scope-Acceptance.md**    | Scope MVP & acceptance criteria                                                |
| **07-ERD-Database-Schema.md**     | Desain database — ERD, skema tabel, indexes, relasi                            |
| **08-API-Contract.md**            | Kontrak API — 46 endpoint dengan request/response lengkap                      |
| **09-Testing-QA-Plan.md**         | Rencana testing & quality assurance                                            |
| **10-Postman-Collection.md**      | Dokumentasi Postman collection                                                 |
| **11-Business-Workflows.md**      | Dokumentasi business workflow                                                  |

---

## 🗄️ Database

### Skema

Proyek menggunakan PostgreSQL 16 dengan 7 tabel utama:

| Tabel               | Deskripsi                           |
| ------------------- | ----------------------------------- |
| `tenants`           | Organisasi/agency terdaftar         |
| `users`             | Semua user (buyer, salesman, admin) |
| `subscriptions`     | Paket subscription tenant           |
| `property_listings` | Listing properti (sale & rent)      |
| `property_photos`   | Foto-foto properti                  |
| `saved_properties`  | Bookmark properti oleh buyer        |
| `audit_logs`        | Log aktivitas untuk audit trail     |

### Migrasi

- **Auto-migrate**: Backend menjalankan `AutoMigrate()` via GORM saat startup
- **Manual migration**: `psql -U propertyhub -d propertyhub -f db/migration.sql`
- **Seed data**: Otomatis saat `APP_ENV=development` dan database kosong

---

## 🔐 Authentication & Authorization

### Flow Autentikasi

1. User register (`POST /api/v1/auth/register`)
2. User login (`POST /api/v1/auth/login`) → dapat JWT token
3. Sertakan token di header: `Authorization: Bearer <token>`
4. Middleware `AuthRequired` memvalidasi token

### Role-Based Access Control (RBAC)

| Role               | Deskripsi                        | Akses                                 |
| ------------------ | -------------------------------- | ------------------------------------- |
| **Buyer**          | User terdaftar, pencari properti | Browse, save/bookmark                 |
| **Salesman**       | Agent di bawah tenant            | CRUD listing, upload foto             |
| **Tenant Admin**   | Pemilik/manajer agency           | Manage tenant, salesmen, subscription |
| **Platform Admin** | Super admin platform             | Manage all tenants, approvals         |

Detail permission matrix: lihat `03-Permission-Matrix.md`

---

## 📡 API Overview

Base URL: `http://localhost:8080/api/v1`

### Public Endpoints (No Auth)

| Method | Endpoint               | Deskripsi                  |
| ------ | ---------------------- | -------------------------- |
| POST   | `/auth/register`       | Register user baru         |
| POST   | `/auth/login`          | Login, dapatkan JWT        |
| GET    | `/properties`          | List properti (filterable) |
| GET    | `/properties/featured` | Properti unggulan          |
| GET    | `/properties/nearby`   | Properti terdekat          |
| GET    | `/properties/:id`      | Detail properti            |
| GET    | `/locations/cities`    | List kota                  |

### Authenticated Endpoints

| Method | Endpoint                             | Role           | Deskripsi                   |
| ------ | ------------------------------------ | -------------- | --------------------------- |
| GET    | `/me/profile`                        | All            | Profil user sendiri         |
| PUT    | `/me/profile`                        | All            | Update profil sendiri       |
| GET    | `/me/saved`                          | Buyer          | List properti tersimpan     |
| POST   | `/me/saved/:propertyId`              | Buyer          | Simpan properti             |
| DELETE | `/me/saved/:propertyId`              | Buyer          | Hapus properti tersimpan    |
| GET    | `/salesman/dashboard`                | Salesman       | Dashboard salesman          |
| GET    | `/salesman/listings`                 | Salesman       | List listing milik sendiri  |
| POST   | `/salesman/listings`                 | Salesman       | Buat listing baru           |
| GET    | `/salesman/listings/:id`             | Salesman       | Detail listing              |
| PUT    | `/salesman/listings/:id`             | Salesman       | Update listing              |
| DELETE | `/salesman/listings/:id`             | Salesman       | Hapus listing               |
| POST   | `/salesman/listings/:id/submit`      | Salesman       | Submit listing untuk review |
| POST   | `/salesman/listings/:id/photos`      | Salesman       | Upload foto listing         |
| GET    | `/tenant/dashboard`                  | Tenant Admin   | Dashboard tenant            |
| GET    | `/tenant/salesmen`                   | Tenant Admin   | List salesman               |
| POST   | `/tenant/salesmen`                   | Tenant Admin   | Tambah salesman             |
| DELETE | `/tenant/salesmen/:id`               | Tenant Admin   | Hapus salesman              |
| GET    | `/tenant/subscription`               | Tenant Admin   | Info subscription           |
| GET    | `/platform/dashboard`                | Platform Admin | Dashboard platform          |
| GET    | `/platform/tenants`                  | Platform Admin | List semua tenant           |
| GET    | `/platform/tenants/:id`              | Platform Admin | Detail tenant               |
| PUT    | `/platform/tenants/:id`              | Platform Admin | Update tenant               |
| GET    | `/platform/tenants/:id/subscription` | Platform Admin | Subscription tenant         |
| GET    | `/platform/listings/pending`         | Platform Admin | List listing pending review |
| POST   | `/platform/listings/:id/approve`     | Platform Admin | Approve listing             |
| POST   | `/platform/listings/:id/reject`      | Platform Admin | Reject listing              |
| GET    | `/platform/audit-logs`               | Platform Admin | Audit logs                  |

> **Detail lengkap**: Lihat `08-API-Contract.md` untuk request/response schema setiap endpoint.

---

## 🧪 Testing

### Backend

```bash
cd backend
go test ./...                    # Semua test
go test ./handlers/ -v           # Handler tests
go test ./utils/ -v              # Utility tests
go test -cover ./...             # Dengan coverage
```

### Frontend

```bash
cd frontend
npm run build                    # Build check (pastikan tidak error)
```

### Postman Collection

Import `postman-collection.json` ke Postman untuk testing API manual.
Lihat `10-Postman-Collection.md` untuk dokumentasi lengkap.

---

## 🚀 Deployment

### Production Checklist

- [ ] Ganti semua password default di `.env`
- [ ] Generate `JWT_SECRET` yang kuat (`openssl rand -base64 32`)
- [ ] Set `APP_ENV=production`
- [ ] Set `CORS_ALLOWED_ORIGINS` ke domain frontend production
- [ ] Set `FRONTEND_URL` ke domain frontend production
- [ ] Gunakan PostgreSQL managed service (jangan container di production)
- [ ] Setup SSL/TLS (via reverse proxy: Nginx, Caddy, Traefik)
- [ ] Setup monitoring & alerting
- [ ] Backup database secara berkala

### Menggunakan Reverse Proxy

Contoh Nginx di depan container:

```nginx
server {
    listen 443 ssl;
    server_name propertyhub.id;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:80;    # Frontend
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/ {
        proxy_pass http://localhost:8080;  # Backend
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /uploads/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }
}
```

---

## 🔍 Troubleshooting

### Container gagal start

```bash
# Cek status semua service
docker compose ps

# Lihat logs
docker compose logs postgres
docker compose logs backend
docker compose logs frontend

# Cek health postgres
docker compose exec postgres pg_isready -U propertyhub -d propertyhub
```

### Backend tidak bisa connect database

```bash
# Pastikan postgres sehat
docker compose ps postgres

# Cek network
docker compose exec backend ping postgres

# Cek environment variables
docker compose exec backend env | grep DB_
```

### Frontend tidak bisa akses API

```bash
# Pastikan nginx config benar
docker compose exec frontend cat /etc/nginx/conf.d/default.conf

# Cek nginx bisa resolve backend
docker compose exec frontend nslookup backend
```

### Reset database ke data awal

```bash
docker compose down -v
docker compose up -d
```

### Port sudah digunakan

```bash
# Windows: cek port 5432, 8080, 80
netstat -ano | findstr :5432
netstat -ano | findstr :8080
netstat -ano | findstr :80

# Ubah port mapping di docker-compose.override.yml atau .env
```

---

## 📄 License

Proprietary — Internal use only.
