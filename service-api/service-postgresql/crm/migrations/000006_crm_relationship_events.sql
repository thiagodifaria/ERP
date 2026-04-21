CREATE TABLE IF NOT EXISTS crm.relationship_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  aggregate_type VARCHAR(40) NOT NULL,
  aggregate_public_id UUID NOT NULL,
  event_code VARCHAR(80) NOT NULL,
  actor VARCHAR(80) NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_crm_relationship_events_public_id UNIQUE (public_id)
);

CREATE INDEX IF NOT EXISTS idx_crm_relationship_events_tenant_aggregate
  ON crm.relationship_events (tenant_id, aggregate_type, aggregate_public_id, created_at, id);
