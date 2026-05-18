CREATE SCHEMA IF NOT EXISTS platform_control;

CREATE TABLE IF NOT EXISTS platform_control.provider_activation_runs (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    tenant_slug TEXT NOT NULL,
    provider_key TEXT NOT NULL,
    domain TEXT NOT NULL,
    action TEXT NOT NULL,
    credential_key TEXT NOT NULL,
    credential_configured BOOLEAN NOT NULL DEFAULT false,
    status TEXT NOT NULL,
    status_code INTEGER,
    result JSONB NOT NULL DEFAULT '{}'::jsonb,
    secret_value_exposed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_provider_activation_runs_tenant ON platform_control.provider_activation_runs (tenant_slug, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_provider_activation_runs_provider ON platform_control.provider_activation_runs (provider_key, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_provider_activation_runs_domain ON platform_control.provider_activation_runs (tenant_slug, domain, created_at DESC);
