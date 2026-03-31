-- Semeia o catalogo inicial de workflow-control para tenants bootstrap.

INSERT INTO workflow_control.workflow_definitions (
  tenant_id,
  public_id,
  key,
  name,
  description,
  status,
  trigger
)
SELECT
  tenant.id,
  gen_random_uuid(),
  'lead-follow-up',
  'Lead Follow-Up',
  'Orquestra o acompanhamento inicial de novos leads do CRM.',
  'active',
  'lead.created'
FROM identity.tenants AS tenant
WHERE tenant.slug IN ('bootstrap-ops', 'northwind-group', 'smoke-identity-bootstrap')
ON CONFLICT (tenant_id, key) DO NOTHING;
