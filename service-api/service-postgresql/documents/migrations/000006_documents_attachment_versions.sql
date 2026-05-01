ALTER TABLE documents.attachments
  ADD COLUMN IF NOT EXISTS current_version_number INTEGER NOT NULL DEFAULT 1;

CREATE TABLE IF NOT EXISTS documents.attachment_versions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  attachment_id BIGINT NOT NULL REFERENCES documents.attachments(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  version_number INTEGER NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  content_type VARCHAR(120) NOT NULL DEFAULT 'application/octet-stream',
  storage_key TEXT NOT NULL,
  storage_driver VARCHAR(80) NOT NULL DEFAULT 'manual',
  source VARCHAR(80) NOT NULL DEFAULT 'manual',
  uploaded_by VARCHAR(120) NOT NULL DEFAULT 'system',
  file_size_bytes BIGINT NOT NULL DEFAULT 0,
  checksum_sha256 VARCHAR(128) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_documents_attachment_versions_public_id UNIQUE (public_id),
  CONSTRAINT uq_documents_attachment_versions_attachment_version UNIQUE (attachment_id, version_number),
  CONSTRAINT ck_documents_attachment_versions_file_size CHECK (file_size_bytes >= 0),
  CONSTRAINT ck_documents_attachment_versions_version_number CHECK (version_number > 0)
);

CREATE INDEX IF NOT EXISTS idx_documents_attachment_versions_attachment
  ON documents.attachment_versions (attachment_id, version_number DESC, created_at DESC);
