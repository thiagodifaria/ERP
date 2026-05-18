CREATE TABLE IF NOT EXISTS platform_control.policy_decisions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  policy_key TEXT NOT NULL,
  policy_version TEXT NOT NULL,
  domain TEXT NOT NULL,
  action TEXT NOT NULL,
  effect TEXT NOT NULL,
  decision TEXT NOT NULL,
  actor TEXT NOT NULL,
  reason TEXT NOT NULL,
  context_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  evaluated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_platform_control_policy_decisions_tenant_action
  ON platform_control.policy_decisions (tenant_id, action, evaluated_at DESC);

CREATE TABLE IF NOT EXISTS platform_control.timeline_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  source_service TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_public_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  severity TEXT NOT NULL DEFAULT 'info',
  actor TEXT NOT NULL,
  correlation_id TEXT NOT NULL,
  summary TEXT NOT NULL,
  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_platform_control_timeline_tenant_entity
  ON platform_control.timeline_events (tenant_id, entity_type, entity_public_id, created_at DESC);

CREATE TABLE IF NOT EXISTS platform_control.approval_requests (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  command_type TEXT NOT NULL,
  domain TEXT NOT NULL,
  status TEXT NOT NULL,
  requested_by TEXT NOT NULL,
  approved_by TEXT,
  rejected_by TEXT,
  executed_by TEXT,
  justification TEXT NOT NULL,
  policy_decision_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  command_payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  decided_at TIMESTAMPTZ,
  executed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS ix_platform_control_approval_requests_tenant_status
  ON platform_control.approval_requests (tenant_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS platform_control.runbook_runs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  runbook_key TEXT NOT NULL,
  title TEXT NOT NULL,
  domain TEXT NOT NULL,
  status TEXT NOT NULL,
  requested_by TEXT NOT NULL,
  approval_public_id UUID,
  context_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS platform_control.runbook_steps (
  id BIGSERIAL PRIMARY KEY,
  run_id BIGINT NOT NULL REFERENCES platform_control.runbook_runs(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  sequence INTEGER NOT NULL,
  step_key TEXT NOT NULL,
  status TEXT NOT NULL,
  summary TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS platform_control.evidence_records (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  source_service TEXT NOT NULL,
  evidence_type TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_public_id TEXT NOT NULL,
  actor TEXT NOT NULL,
  classification TEXT NOT NULL DEFAULT 'internal',
  retention TEXT NOT NULL DEFAULT 'p2y',
  payload_hash TEXT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_platform_control_evidence_tenant_entity
  ON platform_control.evidence_records (tenant_id, entity_type, entity_public_id, created_at DESC);
