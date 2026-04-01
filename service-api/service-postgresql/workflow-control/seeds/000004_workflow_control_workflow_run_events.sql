-- Semeia um evento bootstrap por execucao inicial de workflow.

INSERT INTO workflow_control.workflow_run_events (
  tenant_id,
  public_id,
  workflow_run_id,
  category,
  body,
  created_by
)
SELECT
  tenant.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-000000000501'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-000000000502'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-000000000503'::uuid
  END,
  workflow_run.id,
  'note',
  'Execucao bootstrap criada para acompanhamento inicial do fluxo.',
  'bootstrap-seed'
FROM identity.tenants AS tenant
INNER JOIN workflow_control.workflow_runs AS workflow_run
  ON workflow_run.tenant_id = tenant.id
 AND workflow_run.public_id = CASE tenant.slug
   WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-000000000301'::uuid
   WHEN 'northwind-group' THEN '00000000-0000-0000-0000-000000000302'::uuid
   WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-000000000303'::uuid
 END
WHERE tenant.slug IN ('bootstrap-ops', 'northwind-group', 'smoke-identity-bootstrap')
ON CONFLICT (public_id) DO NOTHING;
