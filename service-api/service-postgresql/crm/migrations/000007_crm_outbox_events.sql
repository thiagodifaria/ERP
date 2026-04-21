CREATE TABLE IF NOT EXISTS crm.outbox_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  aggregate_type VARCHAR(40) NOT NULL,
  aggregate_public_id UUID NOT NULL,
  event_type VARCHAR(120) NOT NULL,
  payload JSONB NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  processed_at TIMESTAMPTZ NULL,
  CONSTRAINT uq_crm_outbox_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_crm_outbox_events_status CHECK (status IN ('pending', 'processed'))
);

CREATE INDEX IF NOT EXISTS idx_crm_outbox_events_tenant_status
  ON crm.outbox_events (tenant_id, status, created_at, id);
