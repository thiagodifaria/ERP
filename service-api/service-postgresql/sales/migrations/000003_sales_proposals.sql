-- Cria a estrutura inicial de propostas do contexto sales.
-- Cada proposta fica vinculada a uma oportunidade do proprio tenant.

CREATE TABLE IF NOT EXISTS sales.proposals (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id),
  opportunity_id BIGINT NOT NULL REFERENCES sales.opportunities(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  title VARCHAR(180) NOT NULL,
  status VARCHAR(40) NOT NULL DEFAULT 'draft',
  amount_cents BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_sales_proposals_public_id UNIQUE (public_id),
  CONSTRAINT ck_sales_proposals_status CHECK (status IN ('draft', 'sent', 'accepted', 'rejected')),
  CONSTRAINT ck_sales_proposals_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_sales_proposals_tenant_id
  ON sales.proposals (tenant_id);

CREATE INDEX IF NOT EXISTS idx_sales_proposals_opportunity_id
  ON sales.proposals (opportunity_id);

CREATE INDEX IF NOT EXISTS idx_sales_proposals_status
  ON sales.proposals (status);
