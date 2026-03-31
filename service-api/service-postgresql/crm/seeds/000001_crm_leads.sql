-- Popula um lead bootstrap por tenant quando ainda nao existe massa inicial de CRM.
-- Este seed alinha o banco com o bootstrap em memoria atual do servico CRM.

INSERT INTO crm.leads (tenant_id, public_id, owner_user_public_id, name, email, source, status)
SELECT
  tenant.id,
  gen_random_uuid(),
  "user".public_id,
  concat(tenant.display_name, ' Lead'),
  concat('lead@', tenant.slug, '.local'),
  'manual',
  'captured'
FROM identity.tenants AS tenant
LEFT JOIN identity.users AS "user"
  ON "user".tenant_id = tenant.id
 AND "user".email = concat('owner@', tenant.slug, '.local')
WHERE NOT EXISTS (
  SELECT 1
  FROM crm.leads AS lead
  WHERE lead.tenant_id = tenant.id
    AND lead.email = concat('lead@', tenant.slug, '.local')
);
