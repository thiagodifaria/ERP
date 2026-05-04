CREATE TABLE IF NOT EXISTS platform_control.provider_defaults (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  capability_key TEXT NOT NULL,
  provider_key TEXT NOT NULL,
  provider_type TEXT NOT NULL,
  mode TEXT NOT NULL DEFAULT 'unconfigured',
  critical BOOLEAN NOT NULL DEFAULT FALSE,
  fallback_allowed BOOLEAN NOT NULL DEFAULT FALSE,
  env_key TEXT NULL,
  source TEXT NOT NULL DEFAULT 'tenant-default',
  configured BOOLEAN NOT NULL DEFAULT FALSE,
  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_platform_control_provider_default_mode
    CHECK (mode IN ('configured', 'fallback', 'manual', 'unconfigured', 'disabled'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_control_provider_defaults_tenant_capability
  ON platform_control.provider_defaults (tenant_id, capability_key);
