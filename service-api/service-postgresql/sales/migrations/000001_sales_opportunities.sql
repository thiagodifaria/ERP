-- Cria a estrutura inicial de oportunidades do contexto sales.
-- Ownership do pipeline comercial avancado fica neste contexto.

CREATE SCHEMA IF NOT EXISTS sales;

CREATE TABLE IF NOT EXISTS sales.opportunities (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  public_id UUID NOT NULL,
  lead_public_id UUID NOT NULL,
  customer_public_id UUID NOT NULL,
  owner_user_public_id UUID,
  title VARCHAR(180) NOT NULL,
  stage VARCHAR(40) NOT NULL DEFAULT 'qualified',
  sale_type VARCHAR(40) NOT NULL DEFAULT 'new',
  amount_cents BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_opportunities_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_opportunities_stage CHECK (stage IN ('qualified', 'proposal', 'negotiation', 'won', 'lost')),
  CONSTRAINT ck_sales_opportunities_sale_type CHECK (sale_type IN ('new', 'upsell', 'renewal', 'cross_sell')),
  CONSTRAINT ck_sales_opportunities_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_sales_opportunities_tenant_id
  ON sales.opportunities (tenant_id);

CREATE INDEX IF NOT EXISTS idx_sales_opportunities_stage
  ON sales.opportunities (stage);

CREATE INDEX IF NOT EXISTS idx_sales_opportunities_lead_public_id
  ON sales.opportunities (lead_public_id);

CREATE INDEX IF NOT EXISTS idx_sales_opportunities_customer_public_id
  ON sales.opportunities (customer_public_id);
