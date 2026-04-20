-- Registra uma trilha minima de historico e outbox bootstrap para o contexto sales.
INSERT INTO sales.commercial_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_code, actor, summary)
SELECT
  sale.tenant_id,
  '0195e7a0-7a9c-7c1f-8a44-4a6e71000005'::uuid,
  'sale',
  sale.public_id,
  'sale_bootstrapped',
  'sales',
  'Bootstrap sale loaded into the commercial ledger.'
FROM sales.sales AS sale
WHERE sale.public_id = '0195e7a0-7a9c-7c1f-8a44-4a6e71000003'::uuid
  AND NOT EXISTS (
    SELECT 1
    FROM sales.commercial_events AS event
    WHERE event.public_id = '0195e7a0-7a9c-7c1f-8a44-4a6e71000005'::uuid
  );

INSERT INTO sales.outbox_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_type, payload, status, processed_at)
SELECT
  sale.tenant_id,
  '0195e7a0-7a9c-7c1f-8a44-4a6e71000006'::uuid,
  'sale',
  sale.public_id,
  'sale.bootstrapped',
  jsonb_build_object(
    'salePublicId', sale.public_id::text,
    'status', sale.status,
    'amountCents', sale.amount_cents
  ),
  'processed',
  timezone('utc', now())
FROM sales.sales AS sale
WHERE sale.public_id = '0195e7a0-7a9c-7c1f-8a44-4a6e71000003'::uuid
  AND NOT EXISTS (
    SELECT 1
    FROM sales.outbox_events AS outbox_event
    WHERE outbox_event.public_id = '0195e7a0-7a9c-7c1f-8a44-4a6e71000006'::uuid
  );
