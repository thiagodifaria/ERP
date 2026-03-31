-- Popula a empresa padrao de cada tenant para bootstrap local e contratos iniciais.
-- O registro nasce alinhado ao display name do tenant quando ainda nao existe empresa base.

INSERT INTO identity.companies (tenant_id, public_id, display_name, legal_name, tax_id, status)
SELECT
  tenant.id,
  gen_random_uuid(),
  tenant.display_name,
  tenant.display_name,
  NULL,
  'active'
FROM identity.tenants AS tenant
WHERE NOT EXISTS (
  SELECT 1
  FROM identity.companies AS company
  WHERE company.tenant_id = tenant.id
    AND company.display_name = tenant.display_name
);
