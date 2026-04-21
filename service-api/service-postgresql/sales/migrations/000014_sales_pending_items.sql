-- Define pendencias operacionais ligadas a cada venda.
CREATE TABLE IF NOT EXISTS sales.pending_items (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  sale_id BIGINT NOT NULL REFERENCES sales.sales(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  code VARCHAR(80) NOT NULL,
  summary TEXT NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'open',
  resolved_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_pending_items_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_pending_items_status CHECK (status IN ('open', 'resolved', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_sales_pending_items_tenant_sale
  ON sales.pending_items (tenant_id, sale_id, created_at);

DROP TRIGGER IF EXISTS trg_sales_pending_items_updated_at ON sales.pending_items;
CREATE TRIGGER trg_sales_pending_items_updated_at
BEFORE UPDATE ON sales.pending_items
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
