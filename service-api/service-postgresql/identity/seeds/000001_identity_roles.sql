-- Popula os papeis basicos do tenant para bootstrap local e contratos iniciais.
-- Estes registros podem evoluir depois para fluxo administravel por API.

INSERT INTO identity.roles (tenant_id, public_id, code, display_name, status)
SELECT
  tenant.id,
  gen_random_uuid(),
  seed.code,
  seed.display_name,
  'active'
FROM identity.tenants AS tenant
CROSS JOIN (
  VALUES
    ('owner', 'Owner'),
    ('admin', 'Administrator'),
    ('manager', 'Manager'),
    ('operator', 'Operator'),
    ('viewer', 'Viewer')
) AS seed(code, display_name)
WHERE NOT EXISTS (
  SELECT 1
  FROM identity.roles AS role
  WHERE role.tenant_id = tenant.id
    AND role.code = seed.code
);
