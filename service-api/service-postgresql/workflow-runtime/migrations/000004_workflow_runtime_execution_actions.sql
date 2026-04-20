-- Expande o runtime para rastrear plano publicado, delays e ledger de acoes.
-- Cada execucao passa a conhecer a versao publicada consumida pelo runtime.

ALTER TABLE workflow_runtime.executions
  ADD COLUMN IF NOT EXISTS workflow_definition_version_number INTEGER NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS current_action_index INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS waiting_until TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS workflow_runtime.execution_actions (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL,
  execution_id BIGINT NOT NULL REFERENCES workflow_runtime.executions(id) ON DELETE CASCADE,
  step_id VARCHAR(160) NOT NULL,
  action_key VARCHAR(120) NOT NULL,
  label VARCHAR(180) NOT NULL,
  status VARCHAR(40) NOT NULL,
  delay_seconds INTEGER,
  compensation_action_key VARCHAR(120),
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_runtime_execution_actions_public_id UNIQUE (public_id),
  CONSTRAINT ck_workflow_runtime_execution_actions_status CHECK (status IN ('waiting', 'completed', 'compensated'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_execution_actions_execution_id
  ON workflow_runtime.execution_actions (execution_id);

CREATE INDEX IF NOT EXISTS idx_workflow_runtime_execution_actions_status
  ON workflow_runtime.execution_actions (status);
