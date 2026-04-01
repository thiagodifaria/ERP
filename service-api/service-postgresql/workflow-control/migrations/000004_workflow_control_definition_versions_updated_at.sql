-- Mantem updated_at coerente no historico de versoes do workflow-control.

DROP TRIGGER IF EXISTS trg_workflow_control_definition_versions_set_updated_at
  ON workflow_control.workflow_definition_versions;

CREATE TRIGGER trg_workflow_control_definition_versions_set_updated_at
BEFORE UPDATE ON workflow_control.workflow_definition_versions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
