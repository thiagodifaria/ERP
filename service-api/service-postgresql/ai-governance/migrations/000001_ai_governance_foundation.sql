CREATE SCHEMA IF NOT EXISTS ai_governance;

CREATE TABLE IF NOT EXISTS ai_governance.tool_registry (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  tool_key TEXT NOT NULL UNIQUE,
  service_key TEXT NOT NULL,
  mode TEXT NOT NULL DEFAULT 'read',
  capability_key TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_governance.ai_policies (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL UNIQUE,
  policy_key TEXT NOT NULL UNIQUE,
  effect TEXT NOT NULL,
  mode TEXT NOT NULL,
  requires_tenant BOOLEAN NOT NULL DEFAULT TRUE,
  requires_audit BOOLEAN NOT NULL DEFAULT TRUE,
  mutation_allowed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_governance.assistant_runs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  actor TEXT NOT NULL,
  mode TEXT NOT NULL DEFAULT 'read-only',
  status TEXT NOT NULL,
  prompt_text TEXT NOT NULL,
  answer_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  redaction_findings_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_governance.assistant_run_actions (
  id BIGSERIAL PRIMARY KEY,
  run_id BIGINT NOT NULL REFERENCES ai_governance.assistant_runs(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  tool_key TEXT NOT NULL,
  mode TEXT NOT NULL,
  status TEXT NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_governance.prompt_audit_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  actor TEXT NOT NULL,
  run_public_id UUID NOT NULL,
  redaction_findings_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  tool_count INTEGER NOT NULL DEFAULT 0,
  denied_tool_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_governance.redaction_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  findings_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

