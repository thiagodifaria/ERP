-- Cria a tabela principal de intake do contexto webhook-hub.
-- O servico dono do webhook deve persistir seu buffer e seu estado aqui.

CREATE SCHEMA IF NOT EXISTS webhook_hub;

CREATE TABLE IF NOT EXISTS webhook_hub.webhook_events (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL,
  provider VARCHAR(80) NOT NULL,
  event_type VARCHAR(160) NOT NULL,
  external_id VARCHAR(160) NOT NULL,
  payload_summary TEXT,
  status VARCHAR(40) NOT NULL,
  received_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_webhook_hub_webhook_events_public_id UNIQUE (public_id),
  CONSTRAINT uq_webhook_hub_webhook_events_provider_external_id UNIQUE (provider, external_id),
  CONSTRAINT ck_webhook_hub_webhook_events_status CHECK (status IN ('received', 'validated', 'queued', 'processing', 'forwarded', 'failed', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_events_provider
  ON webhook_hub.webhook_events (provider);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_events_status
  ON webhook_hub.webhook_events (status);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_webhook_events_received_at
  ON webhook_hub.webhook_events (received_at DESC);
