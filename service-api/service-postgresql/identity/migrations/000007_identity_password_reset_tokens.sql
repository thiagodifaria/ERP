-- Adiciona tokens de recuperacao de senha ao contexto identity.
-- O objetivo e permitir reset controlado com expiracao e auditoria.

CREATE TABLE IF NOT EXISTS identity.password_reset_tokens (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  reset_token VARCHAR(120) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'pending',
  expires_at TIMESTAMPTZ NOT NULL,
  consumed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_password_reset_tokens_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_identity_password_reset_tokens_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id)
    ON DELETE CASCADE,
  CONSTRAINT uq_identity_password_reset_tokens_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_password_reset_tokens_reset_token UNIQUE (reset_token),
  CONSTRAINT ck_identity_password_reset_tokens_status CHECK (status IN ('pending', 'consumed', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_identity_password_reset_tokens_tenant_user
  ON identity.password_reset_tokens (tenant_id, user_id, created_at DESC);
