-- Cria a estrutura inicial de clientes convertidos a partir do CRM.
CREATE SCHEMA IF NOT EXISTS crm;

CREATE TABLE IF NOT EXISTS crm.customers (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  lead_id BIGINT NOT NULL REFERENCES crm.leads(id) ON DELETE RESTRICT,
  public_id UUID NOT NULL,
  owner_user_public_id UUID NULL,
  name VARCHAR(160) NOT NULL,
  email CITEXT NOT NULL,
  source VARCHAR(80) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_crm_customers_public_id UNIQUE (public_id),
  CONSTRAINT uq_crm_customers_tenant_email UNIQUE (tenant_id, email),
  CONSTRAINT uq_crm_customers_tenant_lead_id UNIQUE (tenant_id, lead_id),
  CONSTRAINT ck_crm_customers_status CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX IF NOT EXISTS idx_crm_customers_tenant_id
  ON crm.customers (tenant_id);

CREATE INDEX IF NOT EXISTS idx_crm_customers_owner_user_public_id
  ON crm.customers (owner_user_public_id);
