-- Popula o usuario owner padrao de cada tenant para bootstrap local e contratos iniciais.
-- O usuario pode ser substituido depois por fluxo real de convite e autenticacao.

INSERT INTO identity.users (tenant_id, company_id, public_id, email, display_name, given_name, family_name, status)
SELECT
  tenant.id,
  company.id,
  gen_random_uuid(),
  concat('owner@', tenant.slug, '.local'),
  concat(tenant.display_name, ' Owner'),
  NULL,
  NULL,
  'active'
FROM identity.tenants AS tenant
LEFT JOIN LATERAL (
  SELECT company.id
  FROM identity.companies AS company
  WHERE company.tenant_id = tenant.id
  ORDER BY company.id
  LIMIT 1
) AS company ON TRUE
WHERE NOT EXISTS (
  SELECT 1
  FROM identity.users AS "user"
  WHERE "user".tenant_id = tenant.id
    AND "user".email = concat('owner@', tenant.slug, '.local')
);
