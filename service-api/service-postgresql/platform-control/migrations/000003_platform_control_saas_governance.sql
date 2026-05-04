CREATE TABLE IF NOT EXISTS platform_control.quotas (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  metric_key TEXT NOT NULL,
  metric_unit TEXT NOT NULL,
  limit_value BIGINT NOT NULL DEFAULT 0,
  enforcement_mode TEXT NOT NULL DEFAULT 'soft',
  source TEXT NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_platform_control_quota_enforcement CHECK (enforcement_mode IN ('soft', 'hard'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_control_quotas_tenant_metric
  ON platform_control.quotas (tenant_id, metric_key);

CREATE TABLE IF NOT EXISTS platform_control.tenant_blocks (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  block_key TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT FALSE,
  reason TEXT NOT NULL DEFAULT '',
  scope TEXT NOT NULL DEFAULT 'tenant',
  source TEXT NOT NULL DEFAULT 'manual',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_control_blocks_tenant_key
  ON platform_control.tenant_blocks (tenant_id, block_key);

ALTER TABLE platform_control.lifecycle_jobs
  ADD COLUMN IF NOT EXISTS idempotency_key TEXT NULL,
  ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ NULL,
  ADD COLUMN IF NOT EXISTS failed_at TIMESTAMPTZ NULL,
  ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMPTZ NULL,
  ADD COLUMN IF NOT EXISTS failure_reason TEXT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_platform_control_jobs_tenant_type_idempotency
  ON platform_control.lifecycle_jobs (tenant_id, job_type, idempotency_key)
  WHERE idempotency_key IS NOT NULL;

CREATE TABLE IF NOT EXISTS platform_control.lifecycle_job_events (
  id BIGSERIAL PRIMARY KEY,
  job_id BIGINT NOT NULL REFERENCES platform_control.lifecycle_jobs(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  status TEXT NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
