-- Estrutura base do runtime de execucao de workflows.
-- Cada execucao pertence a um tenant e mantém timestamps do ciclo de vida operacional.

CREATE SCHEMA IF NOT EXISTS workflow_runtime;

CREATE TABLE IF NOT EXISTS workflow_runtime.executions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  workflow_definition_key VARCHAR(120) NOT NULL,
  subject_type VARCHAR(120) NOT NULL,
  subject_public_id UUID NOT NULL,
  initiated_by VARCHAR(120) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'pending',
  retry_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  failed_at TIMESTAMPTZ,
  cancelled_at TIMESTAMPTZ,
  CONSTRAINT uq_workflow_runtime_executions_public_id UNIQUE (public_id),
  CONSTRAINT ck_workflow_runtime_executions_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_executions_tenant_id
  ON workflow_runtime.executions (tenant_id);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_executions_status
  ON workflow_runtime.executions (status);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_executions_workflow_definition_key
  ON workflow_runtime.executions (workflow_definition_key);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_executions_subject_type
  ON workflow_runtime.executions (subject_type);
