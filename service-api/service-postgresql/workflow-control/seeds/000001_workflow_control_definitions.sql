-- Semeia o catalogo inicial de workflow-control para tenants bootstrap.

INSERT INTO workflow_control.workflow_definitions (
  tenant_id,
  public_id,
  key,
  name,
  description,
  status,
  trigger,
  actions
)
SELECT
  tenant.id,
  gen_random_uuid(),
  'lead-follow-up',
  'Lead Follow-Up',
  'Orquestra o acompanhamento inicial de novos leads do CRM.',
  'active',
  'lead.created',
  jsonb_build_array(
    jsonb_build_object(
      'stepId', 'create-task',
      'actionKey', 'task.create',
      'label', 'Criar tarefa comercial inicial',
      'delaySeconds', NULL,
      'compensationActionKey', 'task.create'
    ),
    jsonb_build_object(
      'stepId', 'cooldown',
      'actionKey', 'delay.wait',
      'label', 'Aguardar janela curta de acompanhamento',
      'delaySeconds', 1,
      'compensationActionKey', NULL
    ),
    jsonb_build_object(
      'stepId', 'notify-webhook',
      'actionKey', 'integration.webhook',
      'label', 'Emitir webhook operacional',
      'delaySeconds', NULL,
      'compensationActionKey', 'integration.webhook'
    )
  )
FROM identity.tenants AS tenant
WHERE tenant.slug IN ('bootstrap-ops', 'northwind-group', 'smoke-identity-bootstrap')
ON CONFLICT (tenant_id, key) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  status = EXCLUDED.status,
  trigger = EXCLUDED.trigger,
  actions = EXCLUDED.actions;
