-- ============================================================
-- init.sql — PostgreSQL initialization
-- Creates extensions needed by the application
-- ============================================================

-- UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cryptography (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Earth distance (for geospatial nearby queries)
CREATE EXTENSION IF NOT EXISTS "earthdistance";
CREATE EXTENSION IF NOT EXISTS "cube";
