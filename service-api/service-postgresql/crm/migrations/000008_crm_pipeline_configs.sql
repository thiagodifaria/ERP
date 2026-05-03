CREATE TABLE IF NOT EXISTS crm.pipeline_configs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  config_name TEXT NOT NULL,
  stages_json JSONB NOT NULL,
  auto_scoring BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_crm_pipeline_configs_tenant
  ON crm.pipeline_configs (tenant_id);

DROP TRIGGER IF EXISTS trg_crm_pipeline_configs_updated_at ON crm.pipeline_configs;
CREATE TRIGGER trg_crm_pipeline_configs_updated_at
BEFORE UPDATE ON crm.pipeline_configs
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
