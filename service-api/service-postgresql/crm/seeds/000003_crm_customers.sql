-- Converte o lead bootstrap em cliente quando o tenant ainda nao possui base de clientes.
INSERT INTO crm.customers (tenant_id, lead_id, public_id, owner_user_public_id, name, email, source, status)
SELECT
  lead.tenant_id,
  lead.id,
  gen_random_uuid(),
  lead.owner_user_public_id,
  lead.name,
  lead.email,
  lead.source,
  'active'
FROM crm.leads AS lead
WHERE lead.email = concat('lead@', (
  SELECT tenant.slug
  FROM identity.tenants AS tenant
  WHERE tenant.id = lead.tenant_id
), '.local')
  AND NOT EXISTS (
    SELECT 1
    FROM crm.customers AS customer
    WHERE customer.tenant_id = lead.tenant_id
      AND customer.lead_id = lead.id
  );
