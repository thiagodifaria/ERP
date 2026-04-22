ALTER TABLE documents.attachments
  ADD COLUMN IF NOT EXISTS file_size_bytes BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS checksum_sha256 VARCHAR(128) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS visibility VARCHAR(32) NOT NULL DEFAULT 'internal',
  ADD COLUMN IF NOT EXISTS retention_days INTEGER NOT NULL DEFAULT 365,
  ADD COLUMN IF NOT EXISTS archive_reason TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS archived_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_documents_attachments_archived
  ON documents.attachments (tenant_id, archived_at, created_at);

CREATE INDEX IF NOT EXISTS idx_documents_attachments_visibility
  ON documents.attachments (tenant_id, visibility, created_at);
