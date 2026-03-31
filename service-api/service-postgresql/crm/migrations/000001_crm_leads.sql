-- Cria a estrutura inicial de leads do contexto CRM.
-- Ownership de relacionamento comercial basico fica neste contexto.

CREATE SCHEMA IF NOT EXISTS crm;

CREATE TABLE IF NOT EXISTS crm.leads (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  owner_user_public_id UUID,
  name VARCHAR(180) NOT NULL,
  email CITEXT NOT NULL,
  source VARCHAR(80) NOT NULL DEFAULT 'manual',
  status VARCHAR(40) NOT NULL DEFAULT 'captured',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_crm_leads_public_id UNIQUE (public_id),
  CONSTRAINT uq_crm_leads_tenant_email UNIQUE (tenant_id, email),
  CONSTRAINT ck_crm_leads_status CHECK (status IN ('captured', 'contacted', 'qualified', 'disqualified'))
);

CREATE INDEX IF NOT EXISTS idx_crm_leads_tenant_id
  ON crm.leads (tenant_id);

CREATE INDEX IF NOT EXISTS idx_crm_leads_status
  ON crm.leads (status);

CREATE INDEX IF NOT EXISTS idx_crm_leads_owner_user_public_id
  ON crm.leads (owner_user_public_id);
