CREATE TABLE IF NOT EXISTS engagement.provider_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants (id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  provider VARCHAR(40) NOT NULL,
  event_type VARCHAR(80) NOT NULL,
  direction VARCHAR(20) NOT NULL,
  external_event_id VARCHAR(120),
  lead_public_id UUID,
  touchpoint_public_id UUID,
  delivery_public_id UUID,
  workflow_run_public_id UUID,
  status VARCHAR(20) NOT NULL,
  payload_summary TEXT NOT NULL DEFAULT '',
  response_summary TEXT NOT NULL DEFAULT '',
  processed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT ck_engagement_provider_events_provider CHECK (provider IN ('resend', 'whatsapp_cloud', 'telegram_bot', 'meta_ads', 'manual')),
  CONSTRAINT ck_engagement_provider_events_direction CHECK (direction IN ('inbound', 'outbound')),
  CONSTRAINT ck_engagement_provider_events_status CHECK (status IN ('received', 'processed', 'failed'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_engagement_provider_events_external
  ON engagement.provider_events (tenant_id, provider, external_event_id)
  WHERE external_event_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_engagement_provider_events_status
  ON engagement.provider_events (tenant_id, status);

CREATE INDEX IF NOT EXISTS idx_engagement_provider_events_provider
  ON engagement.provider_events (tenant_id, provider);

CREATE INDEX IF NOT EXISTS idx_engagement_provider_events_touchpoint
  ON engagement.provider_events (tenant_id, touchpoint_public_id);
