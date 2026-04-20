-- Cria o outbox transacional inicial para publicacao de eventos comerciais.
CREATE SCHEMA IF NOT EXISTS sales;

CREATE TABLE IF NOT EXISTS sales.outbox_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  aggregate_type VARCHAR(40) NOT NULL,
  aggregate_public_id UUID NOT NULL,
  event_type VARCHAR(120) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  processed_at TIMESTAMPTZ NULL,
  CONSTRAINT uq_sales_outbox_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_outbox_events_status CHECK (status IN ('pending', 'processed'))
);

CREATE INDEX IF NOT EXISTS idx_sales_outbox_events_tenant_status
  ON sales.outbox_events (tenant_id, status, created_at);
