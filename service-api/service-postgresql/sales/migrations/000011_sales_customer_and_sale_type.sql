-- Evolui oportunidades e vendas com vinculo explicito de cliente e tipo comercial.
ALTER TABLE sales.opportunities
  ADD COLUMN IF NOT EXISTS customer_public_id UUID,
  ADD COLUMN IF NOT EXISTS sale_type VARCHAR(40) NOT NULL DEFAULT 'new';

UPDATE sales.opportunities AS opportunity
SET customer_public_id = customer.public_id
FROM crm.customers AS customer
INNER JOIN crm.leads AS lead
  ON lead.id = customer.lead_id
WHERE customer.tenant_id = opportunity.tenant_id
  AND lead.public_id = opportunity.lead_public_id
  AND opportunity.customer_public_id IS NULL;

ALTER TABLE sales.opportunities
  ALTER COLUMN customer_public_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'ck_sales_opportunities_sale_type'
  ) THEN
    ALTER TABLE sales.opportunities
      ADD CONSTRAINT ck_sales_opportunities_sale_type
      CHECK (sale_type IN ('new', 'upsell', 'renewal', 'cross_sell'));
  END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_sales_opportunities_customer_public_id
  ON sales.opportunities (customer_public_id);

ALTER TABLE sales.sales
  ADD COLUMN IF NOT EXISTS customer_public_id UUID,
  ADD COLUMN IF NOT EXISTS owner_user_public_id UUID,
  ADD COLUMN IF NOT EXISTS sale_type VARCHAR(40) NOT NULL DEFAULT 'new';

UPDATE sales.sales AS sale
SET
  customer_public_id = opportunity.customer_public_id,
  owner_user_public_id = opportunity.owner_user_public_id,
  sale_type = opportunity.sale_type
FROM sales.opportunities AS opportunity
WHERE opportunity.id = sale.opportunity_id
  AND sale.customer_public_id IS NULL;

ALTER TABLE sales.sales
  ALTER COLUMN customer_public_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'ck_sales_sales_sale_type'
  ) THEN
    ALTER TABLE sales.sales
      ADD CONSTRAINT ck_sales_sales_sale_type
      CHECK (sale_type IN ('new', 'upsell', 'renewal', 'cross_sell'));
  END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_sales_sales_customer_public_id
  ON sales.sales (customer_public_id);
