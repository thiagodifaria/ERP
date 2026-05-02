-- Expande o webhook-hub com DLQ e trilha de erro operacional.
-- O objetivo e permitir requeue seguro e troubleshooting previsivel.

ALTER TABLE webhook_hub.webhook_events
  ADD COLUMN IF NOT EXISTS retry_count INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS last_error_code VARCHAR(120),
  ADD COLUMN IF NOT EXISTS last_error_message TEXT,
  ADD COLUMN IF NOT EXISTS dead_lettered_at TIMESTAMPTZ;

ALTER TABLE webhook_hub.webhook_event_transitions
  DROP CONSTRAINT IF EXISTS ck_webhook_hub_webhook_event_transitions_status;

ALTER TABLE webhook_hub.webhook_events
  DROP CONSTRAINT IF EXISTS ck_webhook_hub_webhook_events_status;

ALTER TABLE webhook_hub.webhook_events
  ADD CONSTRAINT ck_webhook_hub_webhook_events_status
  CHECK (status IN ('received', 'validated', 'queued', 'processing', 'forwarded', 'failed', 'rejected', 'dead_letter'));

ALTER TABLE webhook_hub.webhook_event_transitions
  ADD CONSTRAINT ck_webhook_hub_webhook_event_transitions_status
  CHECK (status IN ('received', 'validated', 'queued', 'processing', 'forwarded', 'failed', 'rejected', 'dead_letter'));

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_events_dead_lettered_at
  ON webhook_hub.webhook_events (dead_lettered_at DESC)
  WHERE dead_lettered_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_events_retry_count
  ON webhook_hub.webhook_events (retry_count DESC);
