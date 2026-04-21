-- Define a trilha de comissoes operacionais ligadas ao comercial.
CREATE TABLE IF NOT EXISTS sales.commissions (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  sale_id BIGINT NOT NULL REFERENCES sales.sales(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  recipient_user_public_id UUID NOT NULL,
  role_code VARCHAR(60) NOT NULL,
  rate_bps INTEGER NOT NULL,
  amount_cents BIGINT NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_commissions_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_commissions_rate_bps CHECK (rate_bps > 0 AND rate_bps <= 10000),
  CONSTRAINT ck_sales_commissions_amount CHECK (amount_cents > 0),
  CONSTRAINT ck_sales_commissions_status CHECK (status IN ('pending', 'blocked', 'released'))
);

CREATE INDEX IF NOT EXISTS idx_sales_commissions_tenant_sale
  ON sales.commissions (tenant_id, sale_id, created_at);

DROP TRIGGER IF EXISTS trg_sales_commissions_updated_at ON sales.commissions;
CREATE TRIGGER trg_sales_commissions_updated_at
BEFORE UPDATE ON sales.commissions
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
