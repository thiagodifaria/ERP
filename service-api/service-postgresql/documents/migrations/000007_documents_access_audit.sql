CREATE TABLE IF NOT EXISTS documents.access_link_revocations (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  attachment_id BIGINT NOT NULL REFERENCES documents.attachments(id) ON DELETE CASCADE,
  token_hash CHAR(64) NOT NULL,
  reason VARCHAR(160) NOT NULL DEFAULT 'manual_revocation',
  actor VARCHAR(160) NOT NULL DEFAULT 'system',
  correlation_id VARCHAR(120) NOT NULL DEFAULT '',
  revoked_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_documents_access_link_revocations_token_hash UNIQUE (token_hash)
);

CREATE INDEX IF NOT EXISTS idx_documents_access_link_revocations_attachment
  ON documents.access_link_revocations (tenant_id, attachment_id, revoked_at DESC);

CREATE TABLE IF NOT EXISTS documents.audit_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  attachment_public_id UUID,
  event_code VARCHAR(80) NOT NULL,
  actor VARCHAR(160) NOT NULL DEFAULT 'system',
  reason TEXT NOT NULL DEFAULT '',
  correlation_id VARCHAR(120) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_documents_audit_events_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_documents_audit_events_tenant_attachment
  ON documents.audit_events (tenant_id, attachment_public_id, created_at DESC);
