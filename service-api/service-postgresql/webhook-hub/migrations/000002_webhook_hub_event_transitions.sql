-- Registra a trilha cronologica de transicoes do webhook-hub.
-- Cada mudanca de status gera uma linha auditavel neste ledger.

CREATE TABLE IF NOT EXISTS webhook_hub.webhook_event_transitions (
  id BIGSERIAL PRIMARY KEY,
  webhook_event_id BIGINT NOT NULL REFERENCES webhook_hub.webhook_events(id) ON DELETE CASCADE,
  status VARCHAR(40) NOT NULL,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT ck_webhook_hub_webhook_event_transitions_status CHECK (status IN ('received', 'validated', 'queued', 'processing', 'forwarded', 'failed', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_event_transitions_event_id
  ON webhook_hub.webhook_event_transitions (webhook_event_id);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_event_transitions_changed_at
  ON webhook_hub.webhook_event_transitions (changed_at DESC);
