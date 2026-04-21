-- Cria a estrutura inicial de vendas do contexto sales.
-- Cada venda referencia a oportunidade e a proposta que originaram o fechamento.

CREATE TABLE IF NOT EXISTS sales.sales (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  opportunity_id BIGINT NOT NULL REFERENCES sales.opportunities(id),
  proposal_id BIGINT NOT NULL REFERENCES sales.proposals(id),
  public_id UUID NOT NULL,
  customer_public_id UUID NOT NULL,
  owner_user_public_id UUID,
  sale_type VARCHAR(40) NOT NULL DEFAULT 'new',
  status VARCHAR(40) NOT NULL DEFAULT 'active',
  amount_cents BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_sales_public_id UNIQUE (public_id),
  CONSTRAINT uq_sales_sales_proposal_id UNIQUE (proposal_id),
  CONSTRAINT ck_sales_sales_sale_type CHECK (sale_type IN ('new', 'upsell', 'renewal', 'cross_sell')),
  CONSTRAINT ck_sales_sales_status CHECK (status IN ('active', 'invoiced', 'cancelled')),
  CONSTRAINT ck_sales_sales_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_sales_sales_tenant_id
  ON sales.sales (tenant_id);

CREATE INDEX IF NOT EXISTS idx_sales_sales_status
  ON sales.sales (status);

CREATE INDEX IF NOT EXISTS idx_sales_sales_customer_public_id
  ON sales.sales (customer_public_id);
