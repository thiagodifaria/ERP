-- Cria a estrutura inicial de definicoes do contexto workflow-control.
-- O catalogo fica versionavel no banco desde o primeiro runtime real.

CREATE SCHEMA IF NOT EXISTS workflow_control;

CREATE TABLE IF NOT EXISTS workflow_control.workflow_definitions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  key VARCHAR(120) NOT NULL,
  name VARCHAR(180) NOT NULL,
  description TEXT,
  status VARCHAR(40) NOT NULL DEFAULT 'draft',
  trigger VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_control_definitions_public_id UNIQUE (public_id),
  CONSTRAINT uq_workflow_control_definitions_tenant_key UNIQUE (tenant_id, key),
  CONSTRAINT ck_workflow_control_definitions_status CHECK (status IN ('draft', 'active', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_control_definitions_tenant_id
  ON workflow_control.workflow_definitions (tenant_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_definitions_status
  ON workflow_control.workflow_definitions (status);

CREATE INDEX IF NOT EXISTS idx_workflow_control_definitions_trigger
  ON workflow_control.workflow_definitions (trigger);
