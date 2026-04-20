-- Popula campanhas bootstrap por tenant quando ainda nao existe trilha inicial de engagement.

INSERT INTO engagement.campaigns (
  tenant_id,
  public_id,
  key,
  name,
  description,
  channel,
  status,
  touchpoint_goal,
  workflow_definition_key,
  budget_cents
)
SELECT
  tenant.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000c101'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000c201'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000c301'::uuid
    ELSE gen_random_uuid()
  END,
  'lead-follow-up-campaign',
  concat(tenant.display_name, ' Lead Follow-Up'),
  'Campanha bootstrap para o primeiro contato omnichannel com novos leads.',
  'whatsapp',
  'active',
  'book-meeting',
  'lead-follow-up',
  95000
FROM identity.tenants AS tenant
WHERE NOT EXISTS (
  SELECT 1
  FROM engagement.campaigns AS campaign
  WHERE campaign.tenant_id = tenant.id
    AND campaign.key = 'lead-follow-up-campaign'
);

INSERT INTO engagement.campaigns (
  tenant_id,
  public_id,
  key,
  name,
  description,
  channel,
  status,
  touchpoint_goal,
  workflow_definition_key,
  budget_cents
)
SELECT
  tenant.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000c102'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000c202'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000c302'::uuid
    ELSE gen_random_uuid()
  END,
  'proposal-nurture-email',
  concat(tenant.display_name, ' Proposal Nurture'),
  'Campanha bootstrap para aquecer propostas abertas por email.',
  'email',
  'paused',
  'proposal-reminder',
  'proposal-reminder',
  35000
FROM identity.tenants AS tenant
WHERE NOT EXISTS (
  SELECT 1
  FROM engagement.campaigns AS campaign
  WHERE campaign.tenant_id = tenant.id
    AND campaign.key = 'proposal-nurture-email'
);
