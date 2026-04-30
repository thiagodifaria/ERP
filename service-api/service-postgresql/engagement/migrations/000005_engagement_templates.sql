-- Catalogo relacional de templates omnichannel do contexto engagement.
-- Ownership de conteudo reutilizavel e provider padrao fica neste contexto.

CREATE TABLE IF NOT EXISTS engagement.templates (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  key VARCHAR(120) NOT NULL,
  name VARCHAR(160) NOT NULL,
  channel VARCHAR(40) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'draft',
  provider VARCHAR(40) NOT NULL,
  subject VARCHAR(160),
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_engagement_templates_public_id UNIQUE (public_id),
  CONSTRAINT uq_engagement_templates_tenant_key UNIQUE (tenant_id, key),
  CONSTRAINT ck_engagement_templates_channel CHECK (channel IN ('whatsapp', 'email', 'telegram', 'meta_ads', 'manual')),
  CONSTRAINT ck_engagement_templates_status CHECK (status IN ('draft', 'active', 'archived')),
  CONSTRAINT ck_engagement_templates_provider CHECK (provider IN ('resend', 'whatsapp_cloud', 'telegram_bot', 'manual'))
);

CREATE INDEX IF NOT EXISTS idx_engagement_templates_tenant_id
  ON engagement.templates (tenant_id);

CREATE INDEX IF NOT EXISTS idx_engagement_templates_channel
  ON engagement.templates (channel);

CREATE INDEX IF NOT EXISTS idx_engagement_templates_status
  ON engagement.templates (status);

CREATE INDEX IF NOT EXISTS idx_engagement_templates_provider
  ON engagement.templates (provider);
