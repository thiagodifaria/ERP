-- Cria a estrutura inicial de invoices do contexto sales.
-- O faturamento inicial permanece junto do comercial ate o dominio billing amadurecer.

CREATE TABLE IF NOT EXISTS sales.invoices (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  sale_id BIGINT NOT NULL REFERENCES sales.sales(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  invoice_number VARCHAR(60) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'draft',
  amount_cents BIGINT NOT NULL,
  due_date DATE NOT NULL,
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_invoices_public_id UNIQUE (public_id),
  CONSTRAINT uq_sales_invoices_sale_id UNIQUE (sale_id),
  CONSTRAINT uq_sales_invoices_number_per_tenant UNIQUE (tenant_id, invoice_number),
  CONSTRAINT ck_sales_invoices_status CHECK (status IN ('draft', 'sent', 'paid', 'cancelled')),
  CONSTRAINT ck_sales_invoices_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_sales_invoices_tenant_id
  ON sales.invoices (tenant_id);

CREATE INDEX IF NOT EXISTS idx_sales_invoices_status
  ON sales.invoices (status);

CREATE INDEX IF NOT EXISTS idx_sales_invoices_due_date
  ON sales.invoices (due_date);
