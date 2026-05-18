ALTER TABLE workflow_runtime.executions
  ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_executions_public_id_status_version
  ON workflow_runtime.executions (public_id, status, version);
