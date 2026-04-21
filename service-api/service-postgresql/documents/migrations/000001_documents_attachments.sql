CREATE SCHEMA IF NOT EXISTS documents;

CREATE TABLE IF NOT EXISTS documents.attachments (
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
  uploaded_by VARCHAR(120) NOT NULL DEFAULT 'system',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_documents_attachments_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_documents_attachments_tenant_owner
  ON documents.attachments (tenant_id, owner_type, owner_public_id, created_at);

CREATE INDEX IF NOT EXISTS idx_documents_attachments_storage_key
  ON documents.attachments (storage_key);
