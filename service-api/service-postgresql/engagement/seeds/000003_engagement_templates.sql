-- Popula templates bootstrap por tenant quando ainda nao existe catalogo inicial de mensagem.

INSERT INTO engagement.templates (
  tenant_id,
  public_id,
  key,
  name,
  channel,
  status,
  provider,
  subject,
  body
)
SELECT
  tenant.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000f101'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000f201'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000f301'::uuid
    ELSE gen_random_uuid()
  END,
  'lead-follow-up-whatsapp',
  concat(tenant.display_name, ' Lead Follow-Up Template'),
  'whatsapp',
  'active',
  'manual',
  NULL,
  'Ola {{firstName}}, queremos continuar seu atendimento com a proxima acao do processo.'
FROM identity.tenants AS tenant
WHERE NOT EXISTS (
  SELECT 1
  FROM engagement.templates AS template
  WHERE template.tenant_id = tenant.id
    AND template.key = 'lead-follow-up-whatsapp'
);
