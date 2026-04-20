-- Cria a fundacao de convites, sessao, MFA e auditoria do contexto identity.
-- Esses registros sustentam autenticacao, bloqueio de acesso e trilha basica de seguranca.

CREATE TABLE IF NOT EXISTS identity.user_security_profiles (
  user_id BIGINT PRIMARY KEY,
  identity_provider_subject VARCHAR(180),
  mfa_enabled BOOLEAN NOT NULL DEFAULT false,
  mfa_secret VARCHAR(120),
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_user_security_profiles_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id)
    ON DELETE CASCADE,
  CONSTRAINT uq_identity_user_security_profiles_subject UNIQUE (identity_provider_subject)
);

CREATE TABLE IF NOT EXISTS identity.invites (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  tenant_slug VARCHAR(120) NOT NULL,
  user_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  invite_token VARCHAR(120) NOT NULL,
  email CITEXT NOT NULL,
  display_name VARCHAR(180),
  role_codes TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
  team_public_ids UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
  status VARCHAR(40) NOT NULL DEFAULT 'pending',
  expires_at TIMESTAMPTZ NOT NULL,
  accepted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_invites_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_identity_invites_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id)
    ON DELETE CASCADE,
  CONSTRAINT uq_identity_invites_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_invites_token UNIQUE (invite_token),
  CONSTRAINT ck_identity_invites_status CHECK (status IN ('pending', 'accepted', 'expired', 'revoked'))
);

CREATE INDEX IF NOT EXISTS idx_identity_invites_tenant_id
  ON identity.invites (tenant_id);

CREATE INDEX IF NOT EXISTS idx_identity_invites_email
  ON identity.invites (tenant_id, email);

CREATE TABLE IF NOT EXISTS identity.sessions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  session_token VARCHAR(120) NOT NULL,
  refresh_token VARCHAR(120) NOT NULL,
  identity_provider_subject VARCHAR(180),
  identity_provider_refresh_token TEXT,
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  expires_at TIMESTAMPTZ NOT NULL,
  refresh_expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  last_used_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  CONSTRAINT fk_identity_sessions_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_identity_sessions_user
    FOREIGN KEY (user_id)
    REFERENCES identity.users (id)
    ON DELETE CASCADE,
  CONSTRAINT uq_identity_sessions_public_id UNIQUE (public_id),
  CONSTRAINT uq_identity_sessions_session_token UNIQUE (session_token),
  CONSTRAINT uq_identity_sessions_refresh_token UNIQUE (refresh_token),
  CONSTRAINT ck_identity_sessions_status CHECK (status IN ('active', 'revoked', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_identity_sessions_tenant_user
  ON identity.sessions (tenant_id, user_id);

CREATE TABLE IF NOT EXISTS identity.security_audit_events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  actor_user_public_id UUID,
  subject_user_public_id UUID,
  event_code VARCHAR(80) NOT NULL,
  severity VARCHAR(20) NOT NULL,
  summary VARCHAR(240) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT fk_identity_security_audit_events_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES identity.tenants (id)
    ON DELETE CASCADE,
  CONSTRAINT uq_identity_security_audit_events_public_id UNIQUE (public_id),
  CONSTRAINT ck_identity_security_audit_events_severity CHECK (severity IN ('info', 'warning', 'critical'))
);

CREATE INDEX IF NOT EXISTS idx_identity_security_audit_events_tenant_created_at
  ON identity.security_audit_events (tenant_id, created_at DESC);
