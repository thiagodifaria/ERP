-- Garante a manutencao automatica de updated_at nas execucoes do workflow-runtime.

DROP TRIGGER IF EXISTS trg_workflow_runtime_executions_set_updated_at ON workflow_runtime.executions;

CREATE TRIGGER trg_workflow_runtime_executions_set_updated_at
BEFORE UPDATE ON workflow_runtime.executions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
