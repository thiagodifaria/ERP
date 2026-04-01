-- Cria a estrutura base de execucoes de workflows publicadas.
-- Cada run aponta para a definicao e para a versao usada no momento do disparo.

CREATE TABLE IF NOT EXISTS workflow_control.workflow_runs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  workflow_definition_id BIGINT NOT NULL REFERENCES workflow_control.workflow_definitions(id),
  workflow_definition_version_id BIGINT NOT NULL REFERENCES workflow_control.workflow_definition_versions(id),
  status VARCHAR(40) NOT NULL DEFAULT 'pending',
  trigger_event VARCHAR(120) NOT NULL,
  subject_type VARCHAR(120) NOT NULL,
  subject_public_id UUID NOT NULL,
  initiated_by VARCHAR(120) NOT NULL,
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  failed_at TIMESTAMPTZ,
  cancelled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_control_workflow_runs_public_id UNIQUE (public_id),
  CONSTRAINT ck_workflow_control_workflow_runs_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_tenant_id
  ON workflow_control.workflow_runs (tenant_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_definition_id
  ON workflow_control.workflow_runs (workflow_definition_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_version_id
  ON workflow_control.workflow_runs (workflow_definition_version_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_status
  ON workflow_control.workflow_runs (status);
