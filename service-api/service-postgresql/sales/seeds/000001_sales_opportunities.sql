-- Popula uma oportunidade bootstrap por tenant quando ainda nao existe trilha comercial em sales.
-- Este seed usa o lead bootstrap do CRM para alinhar o funil vertical.

INSERT INTO sales.opportunities (tenant_id, public_id, lead_public_id, customer_public_id, owner_user_public_id, title, stage, sale_type, amount_cents)
SELECT
  tenant.id,
  gen_random_uuid(),
  lead.public_id,
  customer.public_id,
  "user".public_id,
  concat(tenant.display_name, ' Opportunity'),
  'won',
  'new',
  125000
FROM identity.tenants AS tenant
INNER JOIN crm.leads AS lead
  ON lead.tenant_id = tenant.id
 AND lead.email = concat('lead@', tenant.slug, '.local')
INNER JOIN crm.customers AS customer
  ON customer.tenant_id = tenant.id
 AND customer.email = concat('lead@', tenant.slug, '.local')
LEFT JOIN identity.users AS "user"
  ON "user".tenant_id = tenant.id
 AND "user".email = concat('owner@', tenant.slug, '.local')
WHERE NOT EXISTS (
  SELECT 1
  FROM sales.opportunities AS opportunity
  WHERE opportunity.tenant_id = tenant.id
);
