-- Cria o catalogo de versoes publicadas de workflow definitions.
-- Cada publicacao guarda um snapshot do estado funcional da definicao.

CREATE TABLE IF NOT EXISTS workflow_control.workflow_definition_versions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  workflow_definition_id BIGINT NOT NULL REFERENCES workflow_control.workflow_definitions(id),
  version_number INTEGER NOT NULL,
  snapshot_name VARCHAR(180) NOT NULL,
  snapshot_description TEXT,
  snapshot_status VARCHAR(40) NOT NULL,
  snapshot_trigger VARCHAR(120) NOT NULL,
  published_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_workflow_control_definition_versions UNIQUE (workflow_definition_id, version_number),
  CONSTRAINT ck_workflow_control_definition_versions_status CHECK (snapshot_status IN ('draft', 'active', 'archived'))
);

CREATE INDEX IF NOT EXISTS idx_workflow_control_definition_versions_tenant_id
  ON workflow_control.workflow_definition_versions (tenant_id);

CREATE INDEX IF NOT EXISTS idx_workflow_control_definition_versions_definition_id
  ON workflow_control.workflow_definition_versions (workflow_definition_id);
