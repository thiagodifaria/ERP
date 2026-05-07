CREATE SCHEMA IF NOT EXISTS fiscal;

CREATE TABLE IF NOT EXISTS fiscal.company_profiles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  tax_regime TEXT NOT NULL DEFAULT 'simples_nacional',
  cnae TEXT NOT NULL DEFAULT '',
  state_registration TEXT NOT NULL DEFAULT '',
  municipal_registration TEXT NOT NULL DEFAULT '',
  certificate_mode TEXT NOT NULL DEFAULT 'a1',
  certificate_label TEXT NOT NULL DEFAULT 'local-fallback',
  environment_mode TEXT NOT NULL DEFAULT 'homologation',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fiscal_company_profiles_company
  ON fiscal.company_profiles (tenant_id, company_public_id);

CREATE TABLE IF NOT EXISTS fiscal.retention_policies (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  data_domain TEXT NOT NULL,
  classification TEXT NOT NULL DEFAULT 'internal',
  retention_days INTEGER NOT NULL DEFAULT 365,
  anonymize_after_days INTEGER NOT NULL DEFAULT 730,
  source TEXT NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fiscal_retention_policies_domain
  ON fiscal.retention_policies (tenant_id, company_public_id, data_domain);

CREATE TABLE IF NOT EXISTS fiscal.documents (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  document_kind TEXT NOT NULL,
  series_code TEXT NOT NULL,
  number_code TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'issued',
  customer_public_id TEXT NULL,
  amount_cents BIGINT NOT NULL DEFAULT 0,
  provider_key TEXT NOT NULL DEFAULT 'local',
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  cancelled_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_fiscal_document_kind CHECK (document_kind IN ('nfe', 'nfse')),
  CONSTRAINT ck_fiscal_document_status CHECK (status IN ('issued', 'cancelled', 'corrected'))
);

CREATE TABLE IF NOT EXISTS fiscal.privacy_requests (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  request_type TEXT NOT NULL,
  subject_kind TEXT NOT NULL,
  subject_public_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'received',
  requested_by TEXT NOT NULL,
  consent_reference TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_fiscal_privacy_request_type CHECK (request_type IN ('access', 'portability', 'anonymization', 'deletion', 'consent_revoke')),
  CONSTRAINT ck_fiscal_privacy_request_status CHECK (status IN ('received', 'processing', 'completed', 'denied'))
);

CREATE TABLE IF NOT EXISTS fiscal.audit_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  category TEXT NOT NULL,
  summary TEXT NOT NULL,
  actor TEXT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
