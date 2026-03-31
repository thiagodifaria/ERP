-- Cria a base de usuarios internos do tenant.
-- Identidade externa, seguranca e acessos convergem para este contexto.

CREATE TABLE IF NOT EXISTS identity.users (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  company_id BIGINT,
  public_id UUID NOT NULL,
  email CITEXT NOT NULL,
  display_name VARCHAR(180) NOT NULL,
  given_name VARCHAR(120),
  family_name VARCHAR(120),
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_users_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT fk_identity_users_company
    FOREIGN KEY (company_id)
    REFERENCES identity.companies (id),
  CONSTRAINT uq_identity_users_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_users_tenant_email UNIQUE (tenant_id, email),
  CONSTRAINT ck_identity_users_status CHECK (status IN ('active', 'inactive', 'invited', 'suspended', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_identity_users_tenant_id
  ON identity.users (tenant_id);

CREATE INDEX IF NOT EXISTS idx_identity_users_company_id
  ON identity.users (company_id);

CREATE INDEX IF NOT EXISTS idx_identity_users_status
  ON identity.users (status);
