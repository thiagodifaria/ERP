-- Vincula o usuario owner padrao ao time Core do tenant.
-- O objetivo e manter banco e bootstrap HTTP alinhados durante a fase inicial.

INSERT INTO identity.team_memberships (tenant_id, team_id, user_id)
SELECT
  tenant.id,
  team.id,
  "user".id
FROM identity.tenants AS tenant
JOIN identity.teams AS team
  ON team.tenant_id = tenant.id
 AND team.name = 'Core'
JOIN identity.users AS "user"
  ON "user".tenant_id = tenant.id
 AND "user".email = concat('owner@', tenant.slug, '.local')
WHERE NOT EXISTS (
  SELECT 1
  FROM identity.team_memberships AS membership
  WHERE membership.team_id = team.id
    AND membership.user_id = "user".id
);
