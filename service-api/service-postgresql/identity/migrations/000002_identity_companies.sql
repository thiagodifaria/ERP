-- Cria a estrutura inicial de empresas ou unidades vinculadas a um tenant.
-- Ownership de identidade e estrutura organizacional basica fica neste contexto.

CREATE TABLE IF NOT EXISTS identity.companies (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  display_name VARCHAR(180) NOT NULL,
  legal_name VARCHAR(220),
  tax_id VARCHAR(32),
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_companies_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id),
  CONSTRAINT uq_identity_companies_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_companies_tenant_name UNIQUE (tenant_id, display_name),
  CONSTRAINT ck_identity_companies_status CHECK (status IN ('active', 'inactive', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_identity_companies_tenant_id
  ON identity.companies (tenant_id);

CREATE INDEX IF NOT EXISTS idx_identity_companies_status
  ON identity.companies (status);
