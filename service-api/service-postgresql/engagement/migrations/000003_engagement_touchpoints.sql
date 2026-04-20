-- Cria a trilha relacional de touchpoints e follow-ups do contexto engagement.
-- Cada touchpoint pode apontar para um lead do CRM e para um workflow associado.

CREATE TABLE IF NOT EXISTS engagement.touchpoints (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  campaign_id BIGINT NOT NULL REFERENCES engagement.campaigns(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  lead_public_id UUID NOT NULL,
  channel VARCHAR(40) NOT NULL,
  contact_value VARCHAR(160) NOT NULL,
  source VARCHAR(60) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'queued',
  workflow_definition_key VARCHAR(120),
  last_workflow_run_public_id UUID,
  created_by VARCHAR(120) NOT NULL,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_engagement_touchpoints_public_id UNIQUE (public_id),
  CONSTRAINT ck_engagement_touchpoints_channel CHECK (channel IN ('whatsapp', 'email', 'telegram', 'meta_ads', 'manual')),
  CONSTRAINT ck_engagement_touchpoints_status CHECK (status IN ('queued', 'sent', 'delivered', 'responded', 'converted', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_tenant_id
  ON engagement.touchpoints (tenant_id);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_campaign_id
  ON engagement.touchpoints (campaign_id);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_status
  ON engagement.touchpoints (status);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_channel
  ON engagement.touchpoints (channel);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_lead_public_id
  ON engagement.touchpoints (lead_public_id);
