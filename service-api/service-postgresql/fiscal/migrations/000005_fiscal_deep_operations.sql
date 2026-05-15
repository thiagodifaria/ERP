CREATE TABLE IF NOT EXISTS fiscal.deep_operations (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL UNIQUE,
  collection TEXT NOT NULL,
  record_key TEXT NOT NULL,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'queued',
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_fiscal_deep_operations_tenant_collection_key
  ON fiscal.deep_operations (tenant_id, collection, record_key);

CREATE INDEX IF NOT EXISTS ix_fiscal_deep_operations_tenant_collection_status
  ON fiscal.deep_operations (tenant_id, collection, status);

DROP TRIGGER IF EXISTS trg_fiscal_deep_operations_updated_at ON fiscal.deep_operations;
CREATE TRIGGER trg_fiscal_deep_operations_updated_at
BEFORE UPDATE ON fiscal.deep_operations
FOR EACH ROW EXECUTE FUNCTION common.set_updated_at();
