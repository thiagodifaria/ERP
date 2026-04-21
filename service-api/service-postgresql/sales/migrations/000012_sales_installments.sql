-- Define a estrutura de parcelamento comercial das vendas.
CREATE TABLE IF NOT EXISTS sales.installments (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  sale_id BIGINT NOT NULL REFERENCES sales.sales(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  sequence_number INTEGER NOT NULL,
  amount_cents BIGINT NOT NULL,
  due_date DATE NOT NULL,
  status VARCHAR(30) NOT NULL DEFAULT 'scheduled',
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_installments_public_id UNIQUE (public_id),
  CONSTRAINT uq_sales_installments_sale_sequence UNIQUE (sale_id, sequence_number),
  CONSTRAINT ck_sales_installments_status CHECK (status IN ('scheduled', 'paid', 'cancelled')),
  CONSTRAINT ck_sales_installments_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_sales_installments_tenant_sale
  ON sales.installments (tenant_id, sale_id, sequence_number);

DROP TRIGGER IF EXISTS trg_sales_installments_updated_at ON sales.installments;
CREATE TRIGGER trg_sales_installments_updated_at
BEFORE UPDATE ON sales.installments
FOR EACH ROW
EXECUTE FUNCTION common_set_updated_at();
