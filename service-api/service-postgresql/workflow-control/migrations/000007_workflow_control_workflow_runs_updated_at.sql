-- Mantem o controle de mutacao das execucoes de workflow.

DROP TRIGGER IF EXISTS trg_workflow_control_workflow_runs_set_updated_at
  ON workflow_control.workflow_runs;

CREATE TRIGGER trg_workflow_control_workflow_runs_set_updated_at
  BEFORE UPDATE ON workflow_control.workflow_runs
  FOR EACH ROW
  EXECUTE FUNCTION common_set_updated_at();

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_trigger_event
  ON workflow_control.workflow_runs (trigger_event);

CREATE INDEX IF NOT EXISTS idx_workflow_control_workflow_runs_subject_lookup
  ON workflow_control.workflow_runs (subject_type, subject_public_id);
