CREATE SCHEMA IF NOT EXISTS platform_control;

CREATE TABLE IF NOT EXISTS platform_control.entitlements (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  capability_key TEXT NOT NULL,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  plan_code TEXT NOT NULL DEFAULT 'custom',
  limit_value BIGINT NOT NULL DEFAULT 0,
  source TEXT NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_control_entitlements_tenant_key
  ON platform_control.entitlements (tenant_id, capability_key);

CREATE TABLE IF NOT EXISTS platform_control.usage_snapshots (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  metric_key TEXT NOT NULL,
  metric_unit TEXT NOT NULL,
  quantity BIGINT NOT NULL DEFAULT 0,
  source TEXT NOT NULL DEFAULT 'manual',
  captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS platform_control.lifecycle_jobs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  job_type TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'queued',
  requested_by TEXT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ NULL,
  CONSTRAINT ck_platform_control_job_type CHECK (job_type IN ('onboarding', 'offboarding')),
  CONSTRAINT ck_platform_control_job_status CHECK (status IN ('queued', 'running', 'completed', 'failed', 'cancelled'))
);
