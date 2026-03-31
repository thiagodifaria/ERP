-- Cria estruturas basicas de times, papeis e memberships.
-- Autorizacao fina externa pode crescer depois sem perder ownership do dominio.

CREATE TABLE IF NOT EXISTS identity.teams (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  company_id BIGINT,
  public_id UUID NOT NULL,
  name VARCHAR(160) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_teams_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT fk_identity_teams_company
    FOREIGN KEY (company_id)
    REFERENCES identity.companies (id),
  CONSTRAINT uq_identity_teams_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_teams_tenant_name UNIQUE (tenant_id, name),
  CONSTRAINT ck_identity_teams_status CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE TABLE IF NOT EXISTS identity.roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  code VARCHAR(80) NOT NULL,
  display_name VARCHAR(160) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_roles_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT uq_identity_roles_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_roles_tenant_code UNIQUE (tenant_id, code),
  CONSTRAINT ck_identity_roles_status CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE TABLE IF NOT EXISTS identity.team_memberships (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  team_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_team_memberships_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT fk_identity_team_memberships_team
    FOREIGN KEY (team_id)
    REFERENCES identity.teams (id),
  CONSTRAINT fk_identity_team_memberships_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id),
  CONSTRAINT uq_identity_team_memberships UNIQUE (team_id, user_id)
);

CREATE TABLE IF NOT EXISTS identity.user_roles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_user_roles_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT fk_identity_user_roles_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id),
  CONSTRAINT fk_identity_user_roles_role
    FOREIGN KEY (role_id)
    REFERENCES identity.roles (id),
  CONSTRAINT uq_identity_user_roles UNIQUE (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_identity_teams_tenant_id
  ON identity.teams (tenant_id);

CREATE INDEX IF NOT EXISTS idx_identity_roles_tenant_id
  ON identity.roles (tenant_id);

CREATE INDEX IF NOT EXISTS idx_identity_team_memberships_user_id
  ON identity.team_memberships (user_id);

CREATE INDEX IF NOT EXISTS idx_identity_user_roles_user_id
  ON identity.user_roles (user_id);
