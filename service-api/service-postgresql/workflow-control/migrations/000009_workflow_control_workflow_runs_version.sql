-- Adiciona controle monotônico de versão para proteger transições concorrentes.

ALTER TABLE workflow_control.workflow_runs
  ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_version
  ON workflow_control.workflow_runs (version);
