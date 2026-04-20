-- Cria o catalogo relacional de campanhas omnichannel do contexto engagement.
-- Ownership de campanhas e cadencia de contato fica neste contexto.

CREATE SCHEMA IF NOT EXISTS engagement;

CREATE TABLE IF NOT EXISTS engagement.campaigns (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  key VARCHAR(120) NOT NULL,
  name VARCHAR(160) NOT NULL,
  description TEXT NOT NULL,
  channel VARCHAR(40) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'draft',
  touchpoint_goal VARCHAR(80) NOT NULL,
  workflow_definition_key VARCHAR(120),
  budget_cents BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_engagement_campaigns_public_id UNIQUE (public_id),
  CONSTRAINT uq_engagement_campaigns_tenant_key UNIQUE (tenant_id, key),
  CONSTRAINT ck_engagement_campaigns_channel CHECK (channel IN ('whatsapp', 'email', 'telegram', 'meta_ads', 'manual')),
  CONSTRAINT ck_engagement_campaigns_status CHECK (status IN ('draft', 'active', 'paused', 'archived')),
  CONSTRAINT ck_engagement_campaigns_budget CHECK (budget_cents >= 0)
);

CREATE INDEX IF NOT EXISTS idx_engagement_campaigns_tenant_id
  ON engagement.campaigns (tenant_id);

CREATE INDEX IF NOT EXISTS idx_engagement_campaigns_channel
  ON engagement.campaigns (channel);

CREATE INDEX IF NOT EXISTS idx_engagement_campaigns_status
  ON engagement.campaigns (status);
