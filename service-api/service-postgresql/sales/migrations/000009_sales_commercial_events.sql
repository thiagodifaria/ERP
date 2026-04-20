-- Cria o ledger auditavel das transicoes comerciais do contexto sales.
CREATE SCHEMA IF NOT EXISTS sales;

CREATE TABLE IF NOT EXISTS sales.commercial_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  aggregate_type VARCHAR(40) NOT NULL,
  aggregate_public_id UUID NOT NULL,
  event_code VARCHAR(80) NOT NULL,
  actor VARCHAR(80) NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_commercial_events_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_sales_commercial_events_tenant_aggregate
  ON sales.commercial_events (tenant_id, aggregate_type, aggregate_public_id, created_at);
