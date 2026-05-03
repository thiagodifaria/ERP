CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE IF NOT EXISTS catalog.categories (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  category_key TEXT NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_categories_tenant_key
  ON catalog.categories (tenant_id, category_key);

CREATE TABLE IF NOT EXISTS catalog.items (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  category_id BIGINT REFERENCES catalog.categories(id) ON DELETE SET NULL,
  public_id UUID NOT NULL UNIQUE,
  sku TEXT NOT NULL,
  name TEXT NOT NULL,
  item_type TEXT NOT NULL,
  unit_code TEXT NOT NULL,
  price_base_cents BIGINT NOT NULL DEFAULT 0,
  currency_code TEXT NOT NULL DEFAULT 'BRL',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  version_number INTEGER NOT NULL DEFAULT 1,
  attributes_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_catalog_items_item_type CHECK (item_type IN ('product', 'service'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_catalog_items_tenant_sku
  ON catalog.items (tenant_id, sku);
