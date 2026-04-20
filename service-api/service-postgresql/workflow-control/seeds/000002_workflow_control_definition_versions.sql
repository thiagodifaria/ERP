-- Semeia a versao inicial do catalogo bootstrap de workflow-control.

INSERT INTO workflow_control.workflow_definition_versions (
  tenant_id,
  workflow_definition_id,
  version_number,
  snapshot_name,
  snapshot_description,
  snapshot_status,
  snapshot_trigger,
  snapshot_actions
)
SELECT
  definition.tenant_id,
  definition.id,
  1,
  definition.name,
  definition.description,
  definition.status,
  definition.trigger,
  definition.actions
FROM workflow_control.workflow_definitions AS definition
ON CONFLICT (workflow_definition_id, version_number) DO UPDATE
SET
  snapshot_name = EXCLUDED.snapshot_name,
  snapshot_description = EXCLUDED.snapshot_description,
  snapshot_status = EXCLUDED.snapshot_status,
  snapshot_trigger = EXCLUDED.snapshot_trigger,
  snapshot_actions = EXCLUDED.snapshot_actions;
