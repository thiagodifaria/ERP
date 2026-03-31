-- Vincula o papel owner ao usuario owner padrao do tenant.
-- Este seed garante coerencia entre contracts HTTP e bootstrap relacional inicial.

INSERT INTO identity.user_roles (tenant_id, user_id, role_id)
SELECT
  tenant.id,
  "user".id,
  role.id
FROM identity.tenants AS tenant
JOIN identity.users AS "user"
  ON "user".tenant_id = tenant.id
 AND "user".email = concat('owner@', tenant.slug, '.local')
JOIN identity.roles AS role
  ON role.tenant_id = tenant.id
 AND role.code = 'owner'
WHERE NOT EXISTS (
  SELECT 1
  FROM identity.user_roles AS user_role
  WHERE user_role.user_id = "user".id
    AND user_role.role_id = role.id
);
