CREATE TABLE IF NOT EXISTS platform_control.incidents (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  title TEXT NOT NULL,
  service_key TEXT NOT NULL,
  severity TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  impact TEXT NOT NULL DEFAULT 'under investigation',
  owner TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  resolved_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ck_platform_control_incident_severity CHECK (severity IN ('sev1', 'sev2', 'sev3', 'sev4')),
  CONSTRAINT ck_platform_control_incident_status CHECK (status IN ('open', 'investigating', 'mitigating', 'resolved', 'cancelled'))
);

CREATE TABLE IF NOT EXISTS platform_control.incident_timeline_events (
  id BIGSERIAL PRIMARY KEY,
  incident_id BIGINT NOT NULL REFERENCES platform_control.incidents(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  event_type TEXT NOT NULL,
  summary TEXT NOT NULL,
  actor TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS platform_control.incident_actions (
  id BIGSERIAL PRIMARY KEY,
  incident_id BIGINT NOT NULL REFERENCES platform_control.incidents(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  title TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open',
  owner TEXT NOT NULL,
  due_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS platform_control.postmortems (
  id BIGSERIAL PRIMARY KEY,
  incident_id BIGINT NOT NULL REFERENCES platform_control.incidents(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  root_cause TEXT NOT NULL,
  impact_summary TEXT NOT NULL,
  preventive_actions_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  evidence_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
