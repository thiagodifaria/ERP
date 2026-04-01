-- Registra eventos operacionais e notas sobre o ciclo de vida das execucoes.

CREATE TABLE IF NOT EXISTS workflow_control.workflow_run_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  workflow_run_id BIGINT NOT NULL REFERENCES workflow_control.workflow_runs(id),
  category VARCHAR(60) NOT NULL,
  body TEXT NOT NULL,
  created_by VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_control_workflow_run_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_workflow_control_workflow_run_events_category CHECK (category IN ('status', 'note'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_run_events_tenant_id
  ON workflow_control.workflow_run_events (tenant_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_run_events_run_id
  ON workflow_control.workflow_run_events (workflow_run_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_run_events_category
  ON workflow_control.workflow_run_events (category);
