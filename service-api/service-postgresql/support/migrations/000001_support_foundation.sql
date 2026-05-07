CREATE SCHEMA IF NOT EXISTS support;

CREATE TABLE IF NOT EXISTS support.queues (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  queue_key TEXT NOT NULL,
  name TEXT NOT NULL,
  sla_target_hours INTEGER NOT NULL DEFAULT 24,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_support_queues_tenant_key
  ON support.queues (tenant_id, queue_key);

CREATE TABLE IF NOT EXISTS support.cases (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  queue_id BIGINT NOT NULL REFERENCES support.queues(id) ON DELETE RESTRICT,
  public_id UUID NOT NULL UNIQUE,
  case_key TEXT NOT NULL,
  subject TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  priority TEXT NOT NULL DEFAULT 'medium',
  owner_user_id TEXT NULL,
  source_kind TEXT NOT NULL DEFAULT 'manual',
  entity_kind TEXT NULL,
  entity_public_id TEXT NULL,
  sla_due_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_support_case_status CHECK (status IN ('open', 'in_progress', 'waiting_customer', 'resolved', 'closed')),
  CONSTRAINT ck_support_case_priority CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_support_cases_tenant_case_key
  ON support.cases (tenant_id, case_key);

CREATE TABLE IF NOT EXISTS support.case_events (
  id BIGSERIAL PRIMARY KEY,
  case_id BIGINT NOT NULL REFERENCES support.cases(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  event_type TEXT NOT NULL,
  summary TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
