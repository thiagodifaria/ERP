-- Padroniza o outbox de CRM com metadados de retry, lock e diagnostico.

ALTER TABLE crm.outbox_events
  ADD COLUMN IF NOT EXISTS attempts INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS locked_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS next_attempt_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS last_error TEXT;

ALTER TABLE crm.outbox_events
  DROP CONSTRAINT IF EXISTS ck_crm_outbox_events_status;

ALTER TABLE crm.outbox_events
  ADD CONSTRAINT ck_crm_outbox_events_status
  CHECK (status IN ('pending', 'processing', 'processed', 'failed', 'dead_letter'));

CREATE INDEX IF NOT EXISTS idx_crm_outbox_events_retry
  ON crm.outbox_events (status, next_attempt_at, locked_at, created_at);
