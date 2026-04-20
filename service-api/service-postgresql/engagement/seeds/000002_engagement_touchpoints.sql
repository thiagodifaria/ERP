-- Popula um touchpoint bootstrap por tenant quando ainda nao existe massa inicial de engagement.
-- O seed se alinha ao lead bootstrap do CRM e ao workflow bootstrap do workflow-control.

INSERT INTO engagement.touchpoints (
  tenant_id,
  campaign_id,
  public_id,
  lead_public_id,
  channel,
  contact_value,
  source,
  status,
  workflow_definition_key,
  last_workflow_run_public_id,
  created_by,
  notes
)
SELECT
  tenant.id,
  campaign.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000e101'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000e201'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000e301'::uuid
    ELSE gen_random_uuid()
  END,
  lead.public_id,
  'whatsapp',
  concat('+5531', lpad((tenant.id + 999999999)::text, 10, '0')),
  'crm',
  'responded',
  campaign.workflow_definition_key,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-000000000301'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-000000000302'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-000000000303'::uuid
    ELSE NULL
  END,
  'bootstrap-seed',
  'Lead respondeu ao primeiro contato omnichannel bootstrap.'
FROM identity.tenants AS tenant
INNER JOIN engagement.campaigns AS campaign
  ON campaign.tenant_id = tenant.id
 AND campaign.key = 'lead-follow-up-campaign'
INNER JOIN crm.leads AS lead
  ON lead.tenant_id = tenant.id
 AND lead.email = concat('lead@', tenant.slug, '.local')
WHERE NOT EXISTS (
  SELECT 1
  FROM engagement.touchpoints AS touchpoint
  WHERE touchpoint.tenant_id = tenant.id
    AND touchpoint.campaign_id = campaign.id
    AND touchpoint.lead_public_id = lead.public_id
);
