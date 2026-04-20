-- Acrescenta o plano de acoes na definicao corrente e no snapshot publicado.
-- O runtime passa a consumir esta trilha para delays, retries e compensacoes.

ALTER TABLE workflow_control.workflow_definitions
  ADD COLUMN IF NOT EXISTS actions JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE workflow_control.workflow_definition_versions
  ADD COLUMN IF NOT EXISTS snapshot_actions JSONB NOT NULL DEFAULT '[]'::jsonb;
