-- Cria a base inicial de projecoes financeiras derivadas do comercial.
CREATE SCHEMA IF NOT EXISTS finance;

CREATE TABLE IF NOT EXISTS finance.receivable_projections (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL REFERENCES identity.tenants(id) ON DELETE CASCADE,
  public_id UUID NOT NULL,
  source_event_public_id UUID NOT NULL,
  projection_kind VARCHAR(40) NOT NULL,
  sale_public_id UUID NOT NULL,
  invoice_public_id UUID NULL,
  status VARCHAR(30) NOT NULL,
  amount_cents BIGINT NOT NULL,
  due_date DATE NULL,
  snapshot_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_finance_receivable_projections_public_id UNIQUE (public_id),
  CONSTRAINT uq_finance_receivable_projections_source_event UNIQUE (source_event_public_id),
  CONSTRAINT ck_finance_receivable_projections_kind CHECK (projection_kind IN ('sale-booking', 'invoice')),
  CONSTRAINT ck_finance_receivable_projections_status CHECK (status IN ('forecast', 'open', 'paid', 'cancelled')),
  CONSTRAINT ck_finance_receivable_projections_amount CHECK (amount_cents > 0)
);

CREATE INDEX IF NOT EXISTS idx_finance_receivable_projections_tenant_status
  ON finance.receivable_projections (tenant_id, status, created_at);

CREATE INDEX IF NOT EXISTS idx_finance_receivable_projections_sale_public_id
  ON finance.receivable_projections (sale_public_id);

CREATE INDEX IF NOT EXISTS idx_finance_receivable_projections_invoice_public_id
  ON finance.receivable_projections (invoice_public_id);
