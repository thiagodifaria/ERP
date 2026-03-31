-- Popula o time Core padrao de cada tenant para bootstrap local e contratos iniciais.
-- O primeiro time organiza o owner inicial e serve de base para contratos de acesso.

INSERT INTO identity.teams (tenant_id, company_id, public_id, name, status)
SELECT
  tenant.id,
  company.id,
  gen_random_uuid(),
  'Core',
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
  FROM identity.teams AS team
  WHERE team.tenant_id = tenant.id
    AND team.name = 'Core'
);
