CREATE TABLE IF NOT EXISTS fiscal.document_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  document_public_id UUID NOT NULL,
  event_type TEXT NOT NULL,
  summary TEXT NOT NULL,
  actor TEXT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_fiscal_document_events_document
  ON fiscal.document_events (tenant_id, document_public_id, created_at DESC);

CREATE TABLE IF NOT EXISTS fiscal.consents (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  company_public_id TEXT NOT NULL,
  subject_kind TEXT NOT NULL,
  subject_public_id TEXT NOT NULL,
  purpose_key TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'granted',
  source TEXT NOT NULL DEFAULT 'ops',
  granted_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_fiscal_consent_status CHECK (status IN ('granted', 'revoked'))
);

CREATE INDEX IF NOT EXISTS ix_fiscal_consents_subject
  ON fiscal.consents (tenant_id, subject_kind, subject_public_id, updated_at DESC);
