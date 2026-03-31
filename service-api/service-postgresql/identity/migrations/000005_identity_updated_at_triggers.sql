-- Liga a funcao compartilhada de `updated_at` nas tabelas operacionais de identidade.
-- O objetivo aqui e padronizar auditoria tecnica minima de alteracao.

DROP TRIGGER IF EXISTS trg_identity_tenants_set_updated_at ON identity.tenants;
CREATE TRIGGER trg_identity_tenants_set_updated_at
BEFORE UPDATE ON identity.tenants
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_identity_companies_set_updated_at ON identity.companies;
CREATE TRIGGER trg_identity_companies_set_updated_at
BEFORE UPDATE ON identity.companies
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_identity_users_set_updated_at ON identity.users;
CREATE TRIGGER trg_identity_users_set_updated_at
BEFORE UPDATE ON identity.users
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_identity_teams_set_updated_at ON identity.teams;
CREATE TRIGGER trg_identity_teams_set_updated_at
BEFORE UPDATE ON identity.teams
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();

DROP TRIGGER IF EXISTS trg_identity_roles_set_updated_at ON identity.roles;
CREATE TRIGGER trg_identity_roles_set_updated_at
BEFORE UPDATE ON identity.roles
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
