CREATE TABLE IF NOT EXISTS documents.upload_sessions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  owner_type VARCHAR(80) NOT NULL,
  owner_public_id UUID NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  content_type VARCHAR(120) NOT NULL DEFAULT 'application/octet-stream',
  storage_key TEXT NOT NULL,
  storage_driver VARCHAR(80) NOT NULL DEFAULT 'manual',
  source VARCHAR(80) NOT NULL DEFAULT 'manual',
  requested_by VARCHAR(120) NOT NULL DEFAULT 'system',
  visibility VARCHAR(32) NOT NULL DEFAULT 'internal',
  retention_days INTEGER NOT NULL DEFAULT 365,
  status VARCHAR(32) NOT NULL DEFAULT 'pending_upload',
  attachment_public_id UUID NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_documents_upload_sessions_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_documents_upload_sessions_tenant_status
  ON documents.upload_sessions (tenant_id, status, created_at);

CREATE INDEX IF NOT EXISTS idx_documents_upload_sessions_owner
  ON documents.upload_sessions (tenant_id, owner_type, owner_public_id, created_at);
