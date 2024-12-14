CREATE TYPE user_role AS ENUM ('operator', 'cashier', 'finance', 'administrator');
CREATE TYPE data_status AS ENUM ('active', 'inactive');
CREATE TYPE journal_method AS ENUM ('automatic', 'manual');
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX products_name_trgm_idx ON products USING gin (name gin_trgm_ops);
CREATE INDEX products_desc_trgm_idx ON products USING gin (description gin_trgm_ops);