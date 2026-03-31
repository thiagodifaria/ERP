-- Mantem updated_at consistente no catalogo de workflow-control.

DROP TRIGGER IF EXISTS trg_workflow_control_definitions_set_updated_at
  ON workflow_control.workflow_definitions;

CREATE TRIGGER trg_workflow_control_definitions_set_updated_at
BEFORE UPDATE ON workflow_control.workflow_definitions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
