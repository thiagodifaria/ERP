CREATE SCHEMA IF NOT EXISTS supplier;

CREATE TABLE IF NOT EXISTS supplier.categories (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  category_key TEXT NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_supplier_categories_tenant_key
  ON supplier.categories (tenant_id, category_key);

CREATE TABLE IF NOT EXISTS supplier.suppliers (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  category_id BIGINT NULL REFERENCES supplier.categories(id) ON DELETE SET NULL,
  public_id UUID NOT NULL UNIQUE,
  company_name TEXT NOT NULL,
  trade_name TEXT NOT NULL,
  tax_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  payable_term_days INTEGER NOT NULL DEFAULT 30,
  bank_name TEXT NOT NULL DEFAULT '',
  pix_key TEXT NOT NULL DEFAULT '',
  contact_email TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_supplier_status CHECK (status IN ('active', 'watchlist', 'inactive'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_supplier_suppliers_tenant_tax_id
  ON supplier.suppliers (tenant_id, tax_id);
