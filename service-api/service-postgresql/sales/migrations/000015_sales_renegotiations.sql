-- Define renegociacoes aplicadas ao ciclo comercial da venda.
CREATE TABLE IF NOT EXISTS sales.renegotiations (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  sale_id BIGINT NOT NULL REFERENCES sales.sales(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  reason TEXT NOT NULL,
  previous_amount_cents BIGINT NOT NULL,
  new_amount_cents BIGINT NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'applied',
  applied_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_renegotiations_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_renegotiations_previous_amount CHECK (previous_amount_cents > 0),
  CONSTRAINT ck_sales_renegotiations_new_amount CHECK (new_amount_cents > 0),
  CONSTRAINT ck_sales_renegotiations_status CHECK (status IN ('applied'))
);

CREATE INDEX IF NOT EXISTS idx_sales_renegotiations_tenant_sale
  ON sales.renegotiations (tenant_id, sale_id, created_at);
