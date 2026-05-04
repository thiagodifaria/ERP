CREATE TABLE IF NOT EXISTS webhook_hub.outbound_endpoints (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id),
  public_id UUID NOT NULL,
  event_type VARCHAR(160) NOT NULL,
  target_url TEXT NOT NULL,
  signing_mode VARCHAR(40) NOT NULL DEFAULT 'hmac_sha256',
  secret_hint VARCHAR(120) NOT NULL DEFAULT '',
  retry_policy VARCHAR(40) NOT NULL DEFAULT 'exponential',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  dlq_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_webhook_hub_outbound_endpoints_public_id UNIQUE (public_id),
  CONSTRAINT ck_webhook_hub_outbound_endpoints_signing_mode CHECK (signing_mode IN ('hmac_sha256', 'none')),
  CONSTRAINT ck_webhook_hub_outbound_endpoints_retry_policy CHECK (retry_policy IN ('exponential', 'linear', 'manual'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_outbound_endpoints_tenant_event
  ON webhook_hub.outbound_endpoints (tenant_id, event_type, active);

DROP TRIGGER IF EXISTS trg_webhook_hub_outbound_endpoints_updated_at ON webhook_hub.outbound_endpoints;
CREATE TRIGGER trg_webhook_hub_outbound_endpoints_updated_at
BEFORE UPDATE ON webhook_hub.outbound_endpoints
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

CREATE TABLE IF NOT EXISTS webhook_hub.outbound_deliveries (
  id BIGSERIAL PRIMARY KEY,
  outbound_endpoint_id BIGINT NOT NULL REFERENCES webhook_hub.outbound_endpoints (id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  webhook_event_public_id UUID,
  status VARCHAR(40) NOT NULL,
  attempt_number INT NOT NULL DEFAULT 1,
  response_code INT,
  last_error_code VARCHAR(120),
  last_error_message TEXT,
  delivered_at TIMESTAMPTZ,
  dead_lettered_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_webhook_hub_outbound_deliveries_public_id UNIQUE (public_id),
  CONSTRAINT ck_webhook_hub_outbound_deliveries_status CHECK (status IN ('queued', 'delivered', 'failed', 'dead_letter'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_outbound_deliveries_endpoint
  ON webhook_hub.outbound_deliveries (outbound_endpoint_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_webhook_hub_outbound_deliveries_dead_letter
  ON webhook_hub.outbound_deliveries (dead_lettered_at DESC)
  WHERE dead_lettered_at IS NOT NULL;

DROP TRIGGER IF EXISTS trg_webhook_hub_outbound_deliveries_updated_at ON webhook_hub.outbound_deliveries;
CREATE TRIGGER trg_webhook_hub_outbound_deliveries_updated_at
BEFORE UPDATE ON webhook_hub.outbound_deliveries
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
