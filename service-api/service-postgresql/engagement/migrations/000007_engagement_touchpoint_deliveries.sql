-- Entregas e tentativas operacionais de touchpoints do contexto engagement.
-- O historico de provider, template e falha operacional fica neste contexto.

CREATE TABLE IF NOT EXISTS engagement.touchpoint_deliveries (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  touchpoint_id BIGINT NOT NULL REFERENCES engagement.touchpoints(id) ON DELETE CASCADE,
  template_id BIGINT REFERENCES engagement.templates(id) ON DELETE SET NULL,
  public_id UUID NOT NULL,
  channel VARCHAR(40) NOT NULL,
  provider VARCHAR(40) NOT NULL,
  provider_message_id VARCHAR(160),
  status VARCHAR(40) NOT NULL DEFAULT 'queued',
  sent_by VARCHAR(120) NOT NULL,
  error_code VARCHAR(80),
  notes TEXT NOT NULL DEFAULT '',
  attempted_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_engagement_touchpoint_deliveries_public_id UNIQUE (public_id),
  CONSTRAINT ck_engagement_touchpoint_deliveries_channel CHECK (channel IN ('whatsapp', 'email', 'telegram', 'meta_ads', 'manual')),
  CONSTRAINT ck_engagement_touchpoint_deliveries_provider CHECK (provider IN ('resend', 'whatsapp_cloud', 'telegram_bot', 'manual')),
  CONSTRAINT ck_engagement_touchpoint_deliveries_status CHECK (status IN ('queued', 'sent', 'delivered', 'failed')),
  CONSTRAINT ck_engagement_touchpoint_deliveries_failure_payload CHECK (
    (status = 'failed' AND error_code IS NOT NULL)
    OR (status <> 'failed')
  )
);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoint_deliveries_tenant_id
  ON engagement.touchpoint_deliveries (tenant_id);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoint_deliveries_touchpoint_id
  ON engagement.touchpoint_deliveries (touchpoint_id);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoint_deliveries_template_id
  ON engagement.touchpoint_deliveries (template_id);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoint_deliveries_status
  ON engagement.touchpoint_deliveries (status);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoint_deliveries_provider
  ON engagement.touchpoint_deliveries (provider);
