-- Cria a tabela base de tenants do contexto de identidade.
-- O identificador publico sera preenchido pela aplicacao ate o contrato final de geracao ser fechado.

CREATE SCHEMA IF NOT EXISTS identity;

CREATE TABLE IF NOT EXISTS identity.tenants (
  id BIGSERIAL PRIMARY KEY,
  public_id UUID NOT NULL,
  slug VARCHAR(120) NOT NULL,
  display_name VARCHAR(180) NOT NULL,
  legal_name VARCHAR(220),
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_identity_tenants_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_tenants_slug UNIQUE (slug),
  CONSTRAINT ck_identity_tenants_status CHECK (status IN ('active', 'inactive', 'suspended', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_identity_tenants_status
  ON identity.tenants (status);

CREATE INDEX IF NOT EXISTS idx_identity_tenants_created_at
  ON identity.tenants (created_at DESC);
