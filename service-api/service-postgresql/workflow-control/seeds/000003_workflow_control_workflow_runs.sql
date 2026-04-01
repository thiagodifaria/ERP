-- Semeia execucoes bootstrap do workflow lead-follow-up para tenants operacionais.

INSERT INTO workflow_control.workflow_runs (
  tenant_id,
  public_id,
  workflow_definition_id,
  workflow_definition_version_id,
  status,
  trigger_event,
  subject_type,
  subject_public_id,
  initiated_by,
  started_at
)
SELECT
  tenant.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-000000000301'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-000000000302'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-000000000303'::uuid
  END,
  definition.id,
  version.id,
  'running',
  'lead.created',
  'crm.lead',
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-000000000401'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-000000000402'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-000000000403'::uuid
  END,
  'bootstrap-seed',
  timezone('utc', now())
FROM identity.tenants AS tenant
INNER JOIN workflow_control.workflow_definitions AS definition
  ON definition.tenant_id = tenant.id
 AND definition.key = 'lead-follow-up'
INNER JOIN workflow_control.workflow_definition_versions AS version
  ON version.tenant_id = tenant.id
 AND version.workflow_definition_id = definition.id
 AND version.version_number = 1
WHERE tenant.slug IN ('bootstrap-ops', 'northwind-group', 'smoke-identity-bootstrap')
ON CONFLICT (public_id) DO NOTHING;
