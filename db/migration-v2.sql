-- ============================================================
-- migration-v2.sql — Full Database Migration (17 Tables)
-- Multi-Tenant Property Information System — v2.0.0
-- PostgreSQL 16
-- Run: psql -U propertyhub -d propertyhub -f migration-v2.sql
-- ============================================================

BEGIN;

-- ============================================================
-- EXTENSIONS
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "cube";
CREATE EXTENSION IF NOT EXISTS "earthdistance";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";       -- for text search

-- ============================================================
-- ENUM TYPES (17 groups)
-- ============================================================

DO $$ BEGIN CREATE TYPE tenant_status AS ENUM ('active','suspended'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE user_role AS ENUM ('buyer','salesman','tenant_admin','platform_admin'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE user_status AS ENUM ('active','inactive','suspended'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE tenant_user_role AS ENUM ('tenant_admin','salesman'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE sub_status AS ENUM ('active','expired','cancelled','pending_upgrade'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE listing_type AS ENUM ('sale','rent'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE source_type AS ENUM ('regular','bank_auction','company_asset'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE rent_period AS ENUM ('daily','monthly','yearly'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE certificate_type AS ENUM ('SHM','SHGB','Girik','Lainnya'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE property_status AS ENUM ('draft','pending','approved','rejected','sold','rented','inactive','deleted'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE watermark_status AS ENUM ('pending','processed','failed','skipped'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE auction_status AS ENUM ('upcoming','open','closed','cancelled','sold'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE disposal_type AS ENUM ('sale','rent','lease'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE asset_status AS ENUM ('available','under_review','sold','rented','inactive'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE inquiry_status AS ENUM ('unread','read','replied','closed'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE audit_action AS ENUM ('CREATE_PROPERTY','UPDATE_PROPERTY','DELETE_PROPERTY','SUBMIT_PROPERTY','APPROVE_PROPERTY','REJECT_PROPERTY','CREATE_TENANT','UPDATE_TENANT','SUSPEND_TENANT','ACTIVATE_TENANT','INVITE_SALESMAN','DEACTIVATE_SALESMAN','UPDATE_SUBSCRIPTION','LOGIN_FAILED','LOGIN_SUCCESS','REGISTER_USER'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;
DO $$ BEGIN CREATE TYPE audit_module AS ENUM ('property','tenant','user','subscription','auth'); EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- ============================================================
-- PHASE 1: MASTER TABLES (no FK dependencies)
-- ============================================================

-- 1. subscription_plans
CREATE TABLE IF NOT EXISTS subscription_plans (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                      VARCHAR(100) NOT NULL,
    slug                      VARCHAR(50)  NOT NULL,
    max_salesmen              INTEGER      NOT NULL DEFAULT 5,
    max_listings_per_salesman INTEGER      NOT NULL DEFAULT 5,
    description               TEXT,
    is_active                 BOOLEAN      NOT NULL DEFAULT true,
    created_at                TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CONSTRAINT uq_plans_name UNIQUE (name),
    CONSTRAINT uq_plans_slug UNIQUE (slug)
);

-- 2. property_types
CREATE TABLE IF NOT EXISTS property_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(50) NOT NULL,
    slug        VARCHAR(50) NOT NULL,
    description TEXT,
    is_active   BOOLEAN     NOT NULL DEFAULT true,
    CONSTRAINT uq_property_types_name UNIQUE (name),
    CONSTRAINT uq_property_types_slug UNIQUE (slug)
);

-- 3. facilities
CREATE TABLE IF NOT EXISTS facilities (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name      VARCHAR(100) NOT NULL,
    icon      VARCHAR(50),
    is_active BOOLEAN      NOT NULL DEFAULT true,
    CONSTRAINT uq_facilities_name UNIQUE (name)
);

-- 4. locations
CREATE TABLE IF NOT EXISTS locations (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city      VARCHAR(100)  NOT NULL,
    province  VARCHAR(100)  NOT NULL,
    country   VARCHAR(50)   NOT NULL DEFAULT 'Indonesia',
    latitude  DECIMAL(10,7),
    longitude DECIMAL(10,7),
    is_active BOOLEAN       NOT NULL DEFAULT true,
    CONSTRAINT uq_locations_city_province UNIQUE (city, province)
);

-- 5. tenants
CREATE TABLE IF NOT EXISTS tenants (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_name VARCHAR(200) NOT NULL,
    subdomain_slug    VARCHAR(100) NOT NULL,
    logo_url          VARCHAR(500),
    whatsapp_number   VARCHAR(20),
    show_whatsapp     BOOLEAN      NOT NULL DEFAULT true,
    description       TEXT,
    phone             VARCHAR(20),
    address           TEXT,
    status            tenant_status NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT uq_tenants_subdomain UNIQUE (subdomain_slug)
);

-- 6. users
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    phone           VARCHAR(20),
    photo_url       VARCHAR(500),
    whatsapp_number VARCHAR(20),
    show_whatsapp   BOOLEAN      NOT NULL DEFAULT true,
    role            user_role    NOT NULL,
    status          user_status  NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ,
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT uq_users_email UNIQUE (email)
);

-- ============================================================
-- PHASE 2: JUNCTION + SUBSCRIPTION
-- ============================================================

-- 7. tenant_users
CREATE TABLE IF NOT EXISTS tenant_users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID                NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id     UUID                NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    tenant_role tenant_user_role    NOT NULL,
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT now(),
    CONSTRAINT uq_tenant_users UNIQUE (tenant_id, user_id)
);

-- 8. tenant_subscriptions
CREATE TABLE IF NOT EXISTS tenant_subscriptions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID              NOT NULL REFERENCES tenants(id)            ON DELETE CASCADE,
    plan_id    UUID              NOT NULL REFERENCES subscription_plans(id) ON DELETE RESTRICT,
    start_date TIMESTAMPTZ       NOT NULL DEFAULT now(),
    end_date   TIMESTAMPTZ,
    status     sub_status        NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ       NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT uq_tenant_subscriptions UNIQUE (tenant_id)
);

-- ============================================================
-- PHASE 3: PROPERTIES
-- ============================================================

-- 9. properties
CREATE TABLE IF NOT EXISTS properties (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID              NOT NULL REFERENCES tenants(id)        ON DELETE CASCADE,
    salesman_id       UUID              NOT NULL REFERENCES users(id)          ON DELETE RESTRICT,
    property_type_id  UUID              REFERENCES property_types(id)          ON DELETE SET NULL,
    location_id       UUID              REFERENCES locations(id)               ON DELETE SET NULL,
    title             VARCHAR(300)      NOT NULL,
    description       TEXT,
    price             DECIMAL(16,2)     NOT NULL,
    listing_type      listing_type      NOT NULL,
    source_type       source_type       NOT NULL DEFAULT 'regular',
    rent_period       rent_period,
    address           TEXT,
    latitude          DECIMAL(10,7),
    longitude         DECIMAL(10,7),
    land_area         DECIMAL(12,2),
    building_area     DECIMAL(12,2),
    bedrooms          INTEGER,
    bathrooms         INTEGER,
    floors            INTEGER,
    certificate_type  certificate_type,
    status            property_status   NOT NULL DEFAULT 'draft',
    reject_reason     TEXT,
    approved_by       UUID              REFERENCES users(id) ON DELETE SET NULL,
    approved_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ       NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ
);

-- ============================================================
-- PHASE 4: PROPERTY DETAILS
-- ============================================================

-- 10. property_photos
CREATE TABLE IF NOT EXISTS property_photos (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id            UUID            NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    file_name             VARCHAR(255),
    original_url          VARCHAR(500)    NOT NULL,
    watermarked_url       VARCHAR(500)    NOT NULL,
    watermark_status      watermark_status NOT NULL DEFAULT 'pending',
    thumbnail_url         VARCHAR(500),
    medium_url            VARCHAR(500),
    is_primary            BOOLEAN         NOT NULL DEFAULT false,
    sort_order            INTEGER         NOT NULL DEFAULT 0,
    created_at            TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- 11. property_facilities
CREATE TABLE IF NOT EXISTS property_facilities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    facility_id UUID NOT NULL REFERENCES facilities(id)  ON DELETE RESTRICT,
    value       TEXT,
    CONSTRAINT uq_property_facilities UNIQUE (property_id, facility_id)
);

-- 12. bank_auction_details
CREATE TABLE IF NOT EXISTS bank_auction_details (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id          UUID             NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    bank_name            VARCHAR(200)     NOT NULL,
    auction_number       VARCHAR(100),
    auction_limit_price  DECIMAL(16,2),
    auction_deposit      DECIMAL(16,2),
    auction_date         TIMESTAMPTZ,
    auction_location     TEXT,
    auction_document_url VARCHAR(500),
    auction_status       auction_status   NOT NULL DEFAULT 'upcoming',
    notes                TEXT,
    created_at           TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ,
    CONSTRAINT uq_bank_auction UNIQUE (property_id)
);

-- 13. company_asset_details
CREATE TABLE IF NOT EXISTS company_asset_details (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id          UUID             NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    company_name         VARCHAR(200)     NOT NULL,
    company_asset_code   VARCHAR(100),
    disposal_type        disposal_type    NOT NULL,
    asset_status         asset_status     NOT NULL DEFAULT 'available',
    pic_name             VARCHAR(200),
    pic_phone            VARCHAR(20),
    pic_whatsapp_number  VARCHAR(20),
    document_url         VARCHAR(500),
    internal_note        TEXT,
    created_at           TIMESTAMPTZ      NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ,
    CONSTRAINT uq_company_asset UNIQUE (property_id)
);

-- ============================================================
-- PHASE 5: INTERACTIONS
-- ============================================================

-- 14. inquiries
CREATE TABLE IF NOT EXISTS inquiries (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID            NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    buyer_id    UUID            NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
    message     TEXT,
    status      inquiry_status  NOT NULL DEFAULT 'unread',
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ
);

-- 15. saved_properties
CREATE TABLE IF NOT EXISTS saved_properties (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id    UUID         NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
    listing_id  UUID         NOT NULL REFERENCES properties(id)  ON DELETE CASCADE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CONSTRAINT uq_saved_properties UNIQUE (buyer_id, listing_id)
);

-- 16. property_views
CREATE TABLE IF NOT EXISTS property_views (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID          NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    user_id     UUID          REFERENCES users(id)               ON DELETE SET NULL,
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- ============================================================
-- PHASE 6: AUDIT
-- ============================================================

-- 17. audit_logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID          REFERENCES users(id)   ON DELETE SET NULL,
    tenant_id    UUID          REFERENCES tenants(id) ON DELETE SET NULL,
    action       audit_action  NOT NULL,
    module       audit_module  NOT NULL,
    reference_id VARCHAR(36)   NOT NULL,
    description  TEXT,
    old_data     JSONB,
    new_data     JSONB,
    ip_address   VARCHAR(45),
    user_agent   TEXT,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- ============================================================
-- INDEXES
-- ============================================================

-- FK indexes (GORM creates automatically, but explicit for safety)
CREATE INDEX IF NOT EXISTS idx_tenant_users_tenant ON tenant_users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_tenant_users_user   ON tenant_users(user_id);
CREATE INDEX IF NOT EXISTS idx_tenant_users_role   ON tenant_users(tenant_id, tenant_role);

CREATE INDEX IF NOT EXISTS idx_ts_tenant    ON tenant_subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ts_plan      ON tenant_subscriptions(plan_id);

CREATE INDEX IF NOT EXISTS idx_prop_tenant      ON properties(tenant_id);
CREATE INDEX IF NOT EXISTS idx_prop_salesman    ON properties(salesman_id);
CREATE INDEX IF NOT EXISTS idx_prop_type        ON properties(property_type_id);
CREATE INDEX IF NOT EXISTS idx_prop_location    ON properties(location_id);
CREATE INDEX IF NOT EXISTS idx_prop_approved_by ON properties(approved_by);

-- Query performance indexes
CREATE INDEX IF NOT EXISTS idx_prop_tenant_status       ON properties(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_prop_status_loc          ON properties(status, location_id);
CREATE INDEX IF NOT EXISTS idx_prop_status_type         ON properties(status, property_type_id);
CREATE INDEX IF NOT EXISTS idx_prop_status_source       ON properties(status, source_type);
CREATE INDEX IF NOT EXISTS idx_prop_status_price        ON properties(status, price);
CREATE INDEX IF NOT EXISTS idx_prop_status_created      ON properties(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_prop_tenant_salesman_st  ON properties(tenant_id, salesman_id, status);
CREATE INDEX IF NOT EXISTS idx_prop_city_status         ON properties(status);  -- used with location join

-- Geospatial index (nearby search)
CREATE INDEX IF NOT EXISTS idx_prop_latlng ON properties USING gist (
    ll_to_earth(COALESCE(latitude, 0), COALESCE(longitude, 0))
);

-- Text search index
CREATE INDEX IF NOT EXISTS idx_prop_text ON properties USING gin (
    (COALESCE(title,'') || ' ' || COALESCE(description,'')) gin_trgm_ops
);

-- Property photos
CREATE INDEX IF NOT EXISTS idx_pphot_listing_primary ON property_photos(listing_id, is_primary);
CREATE INDEX IF NOT EXISTS idx_pphot_listing_sort    ON property_photos(listing_id, sort_order);

-- Saved properties
CREATE INDEX IF NOT EXISTS idx_saved_buyer ON saved_properties(buyer_id, created_at DESC);

-- Property views
CREATE INDEX IF NOT EXISTS idx_pview_prop ON property_views(property_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_pview_user ON property_views(user_id, created_at DESC);

-- Audit logs
CREATE INDEX IF NOT EXISTS idx_audit_created  ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_entity   ON audit_logs(module, reference_id);
CREATE INDEX IF NOT EXISTS idx_audit_user     ON audit_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_tenant   ON audit_logs(tenant_id, created_at DESC);

-- Locations
CREATE INDEX IF NOT EXISTS idx_locations_prov_city ON locations(province, city);

-- ============================================================
-- SEED DATA
-- ============================================================

-- -- Uncomment to seed demo data:
-- \echo 'Seeding master data...'

-- 1. subscription_plans
INSERT INTO subscription_plans (id, name, slug, max_salesmen, max_listings_per_salesman, description) VALUES
  ('p0000000-0000-0000-0000-000000000001', 'Free',    'free',    5,      5,      'Paket gratis dengan batasan'),
  ('p0000000-0000-0000-0000-000000000002', 'Premium', 'premium', 999999, 999999, 'Paket premium tanpa batasan')
ON CONFLICT (slug) DO NOTHING;

-- 2. property_types
INSERT INTO property_types (id, name, slug) VALUES
  ('pt000000-0000-0000-0000-000000000001', 'Rumah',      'house'),
  ('pt000000-0000-0000-0000-000000000002', 'Tanah',      'land'),
  ('pt000000-0000-0000-0000-000000000003', 'Apartemen',  'apartment'),
  ('pt000000-0000-0000-0000-000000000004', 'Ruko',       'shophouse'),
  ('pt000000-0000-0000-0000-000000000005', 'Gudang',     'warehouse'),
  ('pt000000-0000-0000-0000-000000000006', 'Kantor',     'office'),
  ('pt000000-0000-0000-0000-000000000007', 'Villa',      'villa')
ON CONFLICT (slug) DO NOTHING;

-- 3. facilities
INSERT INTO facilities (id, name, icon) VALUES
  ('fc000000-0000-0000-0000-000000000001', 'Carport',          'car'),
  ('fc000000-0000-0000-0000-000000000002', 'Garasi',           'garage'),
  ('fc000000-0000-0000-0000-000000000003', 'Taman',            'tree-pine'),
  ('fc000000-0000-0000-0000-000000000004', 'Kolam Renang',     'waves'),
  ('fc000000-0000-0000-0000-000000000005', 'Keamanan 24 Jam',  'shield-check'),
  ('fc000000-0000-0000-0000-000000000006', 'CCTV',             'camera'),
  ('fc000000-0000-0000-0000-000000000007', 'Gym',              'dumbbell'),
  ('fc000000-0000-0000-0000-000000000008', 'AC',               'wind'),
  ('fc000000-0000-0000-0000-000000000009', 'Dapur',            'cooking-pot'),
  ('fc000000-0000-0000-0000-000000000010', 'Furnished',        'sofa'),
  ('fc000000-0000-0000-0000-000000000011', 'Listrik',          'zap'),
  ('fc000000-0000-0000-0000-000000000012', 'Air',              'droplet'),
  ('fc000000-0000-0000-0000-000000000013', 'Internet',         'wifi'),
  ('fc000000-0000-0000-0000-000000000014', 'Akses Jalan Lebar','truck'),
  ('fc000000-0000-0000-0000-000000000015', 'Lift',             'arrow-up-down')
ON CONFLICT (name) DO NOTHING;

-- 4. locations
INSERT INTO locations (id, city, province, latitude, longitude) VALUES
  ('lc000000-0000-0000-0000-000000000001', 'Jakarta Selatan',  'DKI Jakarta', -6.2436000, 106.8000000),
  ('lc000000-0000-0000-0000-000000000002', 'Jakarta Pusat',    'DKI Jakarta', -6.1865000, 106.8306000),
  ('lc000000-0000-0000-0000-000000000003', 'Tangerang Selatan','Banten',      -6.2889000, 106.7183000),
  ('lc000000-0000-0000-0000-000000000004', 'Bandung',          'Jawa Barat',  -6.9148000, 107.6098000),
  ('lc000000-0000-0000-0000-000000000005', 'Surabaya',         'Jawa Timur',  -7.2575000, 112.7521000),
  ('lc000000-0000-0000-0000-000000000006', 'Depok',            'Jawa Barat',  -6.4001000, 106.8186000),
  ('lc000000-0000-0000-0000-000000000007', 'Bogor',            'Jawa Barat',  -6.5950000, 106.8166000),
  ('lc000000-0000-0000-0000-000000000008', 'Jakarta Utara',    'DKI Jakarta', -6.1376000, 106.8759000)
ON CONFLICT (city, province) DO NOTHING;

-- 5. tenants (demo)
INSERT INTO tenants (id, organization_name, subdomain_slug, description, phone, address, status) VALUES
  ('b0000000-0000-0000-0000-000000000001', 'PropertiJaya Agency', 'propertijaya',
   'Agensi properti terpercaya sejak 2010. Melayani jual-beli dan sewa rumah, apartemen, dan ruko di Jabodetabek.',
   '0215551234', 'Jl. Jend. Sudirman No. 123, Jakarta Pusat 10220', 'active')
ON CONFLICT (subdomain_slug) DO NOTHING;

-- 6. users (demo)
-- Password: Admin@123 / Budi@123 / Andi@123 / Rina@123
INSERT INTO users (id, email, password_hash, name, phone, role, status) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'admin@propertyhub.id',  '$2a$12$dummyhash...', 'Super Admin',   '081100000001', 'platform_admin', 'active'),
  ('u0000000-0000-0000-0000-000000000001', 'budi@propertijaya.id',  '$2a$12$dummyhash...', 'Budi Santoso',  '081200000001', 'tenant_admin',   'active'),
  ('u0000000-0000-0000-0000-000000000002', 'andi@propertijaya.id',  '$2a$12$dummyhash...', 'Andi Pratama',  '081300000001', 'salesman',       'active'),
  ('u0000000-0000-0000-0000-000000000003', 'rina@email.com',        '$2a$12$dummyhash...', 'Rina Wijaya',   '081400000001', 'buyer',          'active')
ON CONFLICT (email) DO NOTHING;

-- 7. tenant_users
INSERT INTO tenant_users (id, tenant_id, user_id, tenant_role) VALUES
  ('tu000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000001', 'tenant_admin'),
  ('tu000000-0000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002', 'salesman')
ON CONFLICT (tenant_id, user_id) DO NOTHING;

-- 8. tenant_subscriptions
INSERT INTO tenant_subscriptions (id, tenant_id, plan_id, status) VALUES
  ('ts000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'p0000000-0000-0000-0000-000000000001', 'active')
ON CONFLICT (tenant_id) DO NOTHING;

-- 9. properties (3 demo: regular + bank_auction + company_asset)
INSERT INTO properties (id, tenant_id, salesman_id, property_type_id, location_id, title, description, price, listing_type, source_type, address, latitude, longitude, land_area, building_area, bedrooms, bathrooms, floors, certificate_type, status) VALUES
  ('pr000000-0000-0000-0000-000000000001',
   'b0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002',
   'pt000000-0000-0000-0000-000000000001', 'lc000000-0000-0000-0000-000000000003',
   'Rumah Minimalis 2 Lantai di Bintaro Sektor 7',
   'Rumah modern minimalis dengan desain terbuka, lingkungan asri dan tenang.',
   1850000000, 'sale', 'regular',
   'Jl. Bintaro Utama 7 No. 15', -6.2765000, 106.7183000,
   150, 200, 4, 3, 2, 'SHM', 'approved')
ON CONFLICT DO NOTHING;

INSERT INTO properties (id, tenant_id, salesman_id, property_type_id, location_id, title, description, price, listing_type, source_type, address, latitude, longitude, land_area, building_area, bedrooms, bathrooms, floors, certificate_type, status) VALUES
  ('pr000000-0000-0000-0000-000000000002',
   'b0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002',
   'pt000000-0000-0000-0000-000000000001', 'lc000000-0000-0000-0000-000000000006',
   'Rumah Lelang Murah di Depok — SHM',
   'Rumah lelang bank, kondisi butuh renovasi ringan. Harga di bawah pasar.',
   450000000, 'sale', 'bank_auction',
   'Jl. Margonda Raya No. 200', -6.4001000, 106.8186000,
   100, 120, 3, 2, 1, 'SHM', 'approved')
ON CONFLICT DO NOTHING;

INSERT INTO properties (id, tenant_id, salesman_id, property_type_id, location_id, title, description, price, listing_type, source_type, address, latitude, longitude, building_area, certificate_type, status) VALUES
  ('pr000000-0000-0000-0000-000000000003',
   'b0000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002',
   'pt000000-0000-0000-0000-000000000006', 'lc000000-0000-0000-0000-000000000001',
   'Ruang Kantor Eks Perusahaan — SCBD Lot 10',
   'Ruang kantor full floor di kawasan SCBD. View kota, fully furnished.',
   8500000000, 'sale', 'company_asset',
   'SCBD Lot 10 Lt. 25', -6.2243000, 106.8102000,
   350, 'SHGB', 'approved')
ON CONFLICT DO NOTHING;

-- 10. bank_auction_details
INSERT INTO bank_auction_details (id, property_id, bank_name, auction_number, auction_limit_price, auction_deposit, auction_date, auction_status) VALUES
  ('ba000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000002',
   'BankMaju', 'AUC-2026-001', 500000000, 50000000, '2026-07-15 10:00:00+07', 'open')
ON CONFLICT (property_id) DO NOTHING;

-- 11. company_asset_details
INSERT INTO company_asset_details (id, property_id, company_name, company_asset_code, disposal_type, asset_status, pic_name, pic_phone) VALUES
  ('ca000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000003',
   'CiptaGraha Corporindo', 'CG-ASSET-001', 'sale', 'available', 'Ratna Sari Dewi', '081200000005')
ON CONFLICT (property_id) DO NOTHING;

-- 12. property_images (placeholder URLs)
INSERT INTO property_images (id, property_id, image_url, watermarked_image_url, watermark_status, thumbnail_url, medium_url, is_primary, sort_order) VALUES
  ('pi000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000001',
   '/uploads/listings/pr001/original/abc.jpg', '/uploads/listings/pr001/watermarked/abc.jpg', 'skipped',
   '/uploads/listings/pr001/thumbnail/abc.jpg', '/uploads/listings/pr001/medium/abc.jpg', true, 0),
  ('pi000000-0000-0000-0000-000000000002', 'pr000000-0000-0000-0000-000000000002',
   '/uploads/listings/pr002/original/def.jpg', '/uploads/listings/pr002/watermarked/def.jpg', 'skipped',
   '/uploads/listings/pr002/thumbnail/def.jpg', '/uploads/listings/pr002/medium/def.jpg', true, 0),
  ('pi000000-0000-0000-0000-000000000003', 'pr000000-0000-0000-0000-000000000003',
   '/uploads/listings/pr003/original/ghi.jpg', '/uploads/listings/pr003/watermarked/ghi.jpg', 'skipped',
   '/uploads/listings/pr003/thumbnail/ghi.jpg', '/uploads/listings/pr003/medium/ghi.jpg', true, 0)
ON CONFLICT DO NOTHING;

-- 13. property_facilities
INSERT INTO property_facilities (id, property_id, facility_id, value) VALUES
  ('pf000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000001', 'fc000000-0000-0000-0000-000000000001', '2 mobil'),
  ('pf000000-0000-0000-0000-000000000002', 'pr000000-0000-0000-0000-000000000001', 'fc000000-0000-0000-0000-000000000003', 'luas'),
  ('pf000000-0000-0000-0000-000000000003', 'pr000000-0000-0000-0000-000000000001', 'fc000000-0000-0000-0000-000000000005', '24 jam'),
  ('pf000000-0000-0000-0000-000000000004', 'pr000000-0000-0000-0000-000000000002', 'fc000000-0000-0000-0000-000000000011', '900 Watt'),
  ('pf000000-0000-0000-0000-000000000005', 'pr000000-0000-0000-0000-000000000003', 'fc000000-0000-0000-0000-000000000015', '3 unit')
ON CONFLICT (property_id, facility_id) DO NOTHING;

-- 14. inquiries
INSERT INTO inquiries (id, property_id, buyer_id, message, status) VALUES
  ('iq000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000001',
   'u0000000-0000-0000-0000-000000000003', 'Apakah harga masih bisa nego?', 'unread')
ON CONFLICT DO NOTHING;

-- 15. favorites
INSERT INTO favorites (id, buyer_id, property_id) VALUES
  ('fv000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000003', 'pr000000-0000-0000-0000-000000000001')
ON CONFLICT (buyer_id, property_id) DO NOTHING;

-- 16. property_views
INSERT INTO property_views (id, property_id, user_id, ip_address) VALUES
  ('pv000000-0000-0000-0000-000000000001', 'pr000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000003', '127.0.0.1'),
  ('pv000000-0000-0000-0000-000000000002', 'pr000000-0000-0000-0000-000000000001', NULL, '192.168.1.1'),
  ('pv000000-0000-0000-0000-000000000003', 'pr000000-0000-0000-0000-000000000002', NULL, '192.168.1.2')
ON CONFLICT DO NOTHING;

-- 17. audit_logs
INSERT INTO audit_logs (id, user_id, tenant_id, action, module, reference_id, description) VALUES
  ('al000000-0000-0000-0000-000000000001', 'u0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001',
   'CREATE_PROPERTY', 'property', 'pr000000-0000-0000-0000-000000000001', 'Salesman membuat listing baru'),
  ('al000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
   'APPROVE_PROPERTY', 'property', 'pr000000-0000-0000-0000-000000000001', 'Admin menyetujui listing'),
  ('al000000-0000-0000-0000-000000000003', 'u0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001',
   'INVITE_SALESMAN', 'user', 'u0000000-0000-0000-0000-000000000002', 'Tenant admin mengundang salesman')
ON CONFLICT DO NOTHING;

COMMIT;

-- ============================================================
-- EOF
-- ============================================================
