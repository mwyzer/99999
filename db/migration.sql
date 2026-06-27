-- ============================================================
-- migration.sql — Full Database Migration
-- Multi-Tenant Property Information System — MVP
-- PostgreSQL 16
-- Run: psql -U propertyhub -d propertyhub -f migration.sql
-- ============================================================

BEGIN;

-- ============================================================
-- EXTENSIONS
-- ============================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "cube";
CREATE EXTENSION IF NOT EXISTS "earthdistance";

-- ============================================================
-- ENUM TYPES
-- ============================================================

-- Tenant status
DO $$ BEGIN
    CREATE TYPE tenant_status AS ENUM ('active', 'suspended');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- User role
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('buyer', 'salesman', 'tenant_admin', 'platform_admin');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- User status
DO $$ BEGIN
    CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Plan type
DO $$ BEGIN
    CREATE TYPE plan_type AS ENUM ('free', 'premium');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Listing type
DO $$ BEGIN
    CREATE TYPE listing_type AS ENUM ('sale', 'rent');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Property type
DO $$ BEGIN
    CREATE TYPE property_type AS ENUM (
        'house', 'land', 'apartment', 'shophouse', 'warehouse', 'office', 'villa'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Source type
DO $$ BEGIN
    CREATE TYPE source_type AS ENUM ('regular', 'bank_auction', 'company_asset');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Certificate type
DO $$ BEGIN
    CREATE TYPE certificate_type AS ENUM ('SHM', 'SHGB', 'Girik', 'Lainnya');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Listing status
DO $$ BEGIN
    CREATE TYPE listing_status AS ENUM (
        'draft', 'pending', 'approved', 'rejected',
        'sold', 'rented', 'inactive', 'deleted'
    );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- ============================================================
-- TABLES
-- ============================================================

-- 1. tenants
CREATE TABLE IF NOT EXISTS tenants (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_name VARCHAR(200) NOT NULL,
    subdomain_slug    VARCHAR(100) NOT NULL,
    logo_url          VARCHAR(500),
    description       TEXT,
    phone             VARCHAR(20),
    address           TEXT,
    status            tenant_status NOT NULL DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    CONSTRAINT uq_tenants_subdomain_slug UNIQUE (subdomain_slug)
);

-- 2. users
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(200) NOT NULL,
    phone         VARCHAR(20),
    photo_url     VARCHAR(500),
    role          user_role NOT NULL,
    status        user_status NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL
);

-- 3. subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                UUID NOT NULL,
    plan_type                plan_type NOT NULL DEFAULT 'free',
    max_salesmen             INTEGER NOT NULL DEFAULT 5,
    max_listings_per_salesman INTEGER NOT NULL DEFAULT 5,
    start_date               TIMESTAMPTZ NOT NULL DEFAULT now(),
    end_date                 TIMESTAMPTZ,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ,
    CONSTRAINT uq_subscriptions_tenant UNIQUE (tenant_id),
    CONSTRAINT fk_subscriptions_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- 4. property_listings
CREATE TABLE IF NOT EXISTS property_listings (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL,
    salesman_id      UUID NOT NULL,
    title            VARCHAR(300) NOT NULL,
    description      TEXT,
    price            DECIMAL(16,2) NOT NULL,
    listing_type     listing_type NOT NULL,
    property_type    property_type NOT NULL,
    source_type      source_type NOT NULL DEFAULT 'regular',
    address          TEXT,
    city             VARCHAR(100),
    province         VARCHAR(100),
    latitude         DECIMAL(10,7),
    longitude        DECIMAL(10,7),
    land_area        DECIMAL(12,2),
    building_area    DECIMAL(12,2),
    bedrooms         INTEGER,
    bathrooms        INTEGER,
    floors           INTEGER,
    certificate_type certificate_type,
    facilities       JSONB DEFAULT '{}'::jsonb,
    status           listing_status NOT NULL DEFAULT 'draft',
    reject_reason    TEXT,
    approved_by      UUID,
    approved_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ,
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT fk_listings_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_listings_salesman FOREIGN KEY (salesman_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT fk_listings_approver FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 5. property_photos
CREATE TABLE IF NOT EXISTS property_photos (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id      UUID NOT NULL,
    file_name       VARCHAR(255),
    original_url    VARCHAR(500) NOT NULL,
    thumbnail_url   VARCHAR(500),
    medium_url      VARCHAR(500),
    watermarked_url VARCHAR(500) NOT NULL,
    sort_order      INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_photos_listing FOREIGN KEY (listing_id) REFERENCES property_listings(id) ON DELETE CASCADE
);

-- 6. saved_properties
CREATE TABLE IF NOT EXISTS saved_properties (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    buyer_id   UUID NOT NULL,
    listing_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_saved_buyer_listing UNIQUE (buyer_id, listing_id),
    CONSTRAINT fk_saved_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_saved_listing FOREIGN KEY (listing_id) REFERENCES property_listings(id) ON DELETE CASCADE
);

-- 7. audit_logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID,
    user_role   VARCHAR(20) NOT NULL,
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   VARCHAR(36) NOT NULL,
    old_values  JSONB,
    new_values  JSONB,
    ip_address  VARCHAR(45),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================
-- INDEXES
-- ============================================================

-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_listings_tenant_id ON property_listings(tenant_id);
CREATE INDEX IF NOT EXISTS idx_listings_salesman_id ON property_listings(salesman_id);
CREATE INDEX IF NOT EXISTS idx_photos_listing_id ON property_photos(listing_id);
CREATE INDEX IF NOT EXISTS idx_saved_buyer_id ON saved_properties(buyer_id);
CREATE INDEX IF NOT EXISTS idx_saved_listing_id ON saved_properties(listing_id);

-- Query performance indexes — Listings
CREATE INDEX IF NOT EXISTS idx_listings_tenant_status ON property_listings(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_listings_status_city ON property_listings(status, city);
CREATE INDEX IF NOT EXISTS idx_listings_status_property_type ON property_listings(status, property_type);
CREATE INDEX IF NOT EXISTS idx_listings_status_source_type ON property_listings(status, source_type);
CREATE INDEX IF NOT EXISTS idx_listings_status_price ON property_listings(status, price);
CREATE INDEX IF NOT EXISTS idx_listings_status_created ON property_listings(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_listings_city_status ON property_listings(city, status);
CREATE INDEX IF NOT EXISTS idx_listings_tenant_salesman_status ON property_listings(tenant_id, salesman_id, status);

-- Geospatial index
CREATE INDEX IF NOT EXISTS idx_listings_geo ON property_listings
    USING gist (ll_to_earth(latitude, longitude))
    WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Query performance indexes — Users
CREATE INDEX IF NOT EXISTS idx_users_tenant_role ON users(tenant_id, role);
CREATE INDEX IF NOT EXISTS idx_users_role_status ON users(role, status);

-- Query performance indexes — Audit
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id, created_at DESC);

-- Full-text search index
CREATE INDEX IF NOT EXISTS idx_listings_fts ON property_listings
    USING gin (to_tsvector('indonesian', coalesce(title, '') || ' ' || coalesce(description, '')));

COMMIT;
