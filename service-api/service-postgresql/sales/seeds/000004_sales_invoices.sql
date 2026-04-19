-- Popula uma invoice bootstrap por venda quando ainda nao existe cobranca relacional.

INSERT INTO sales.invoices (tenant_id, sale_id, public_id, invoice_number, status, amount_cents, due_date)
SELECT
  sale.tenant_id,
  sale.id,
  gen_random_uuid(),
  concat(upper(replace(tenant.slug, '-', '')), '-INV-', lpad(sale.id::text, 4, '0')),
  'sent',
  sale.amount_cents,
  timezone('utc', now())::date + INTERVAL '15 days'
FROM sales.sales AS sale
INNER JOIN identity.tenants AS tenant
  ON tenant.id = sale.tenant_id
WHERE NOT EXISTS (
  SELECT 1
  FROM sales.invoices AS invoice
  WHERE invoice.sale_id = sale.id
);
