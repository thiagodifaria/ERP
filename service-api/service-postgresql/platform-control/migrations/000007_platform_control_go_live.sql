CREATE TABLE IF NOT EXISTS platform_control.go_live_rollouts (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  target_env TEXT NOT NULL DEFAULT 'production',
  wave_key TEXT NOT NULL DEFAULT 'wave-1',
  status TEXT NOT NULL DEFAULT 'planned',
  requested_by TEXT NOT NULL,
  rollback_playbook TEXT NOT NULL,
  adoption_target_pct INTEGER NOT NULL DEFAULT 70,
  readiness_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ NULL,
  completed_at TIMESTAMPTZ NULL,
  rolled_back_at TIMESTAMPTZ NULL,
  CONSTRAINT ck_platform_control_go_live_status CHECK (status IN ('planned', 'running', 'completed', 'rolled_back'))
);

CREATE TABLE IF NOT EXISTS platform_control.go_live_rollout_events (
  id BIGSERIAL PRIMARY KEY,
  rollout_id BIGINT NOT NULL REFERENCES platform_control.go_live_rollouts(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  status TEXT NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
