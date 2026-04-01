-- Ledger cronologico das transicoes de execucao.
-- Cada mudanca de estado produz uma linha auditavel neste historico.

CREATE TABLE IF NOT EXISTS workflow_runtime.execution_transitions (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL,
  execution_id BIGINT NOT NULL REFERENCES workflow_runtime.executions(id) ON DELETE CASCADE,
  status VARCHAR(40) NOT NULL,
  changed_by VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_runtime_execution_transitions_public_id UNIQUE (public_id),
  CONSTRAINT ck_workflow_runtime_execution_transitions_status CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_execution_transitions_execution_id
  ON workflow_runtime.execution_transitions (execution_id);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_execution_transitions_status
  ON workflow_runtime.execution_transitions (status);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_execution_transitions_created_at
  ON workflow_runtime.execution_transitions (created_at DESC);
