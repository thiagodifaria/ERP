-- Popula uma entrega bootstrap por tenant para dar contexto operacional ao engagement.

INSERT INTO engagement.touchpoint_deliveries (
  tenant_id,
  touchpoint_id,
  template_id,
  public_id,
  channel,
  provider,
  provider_message_id,
  status,
  sent_by,
  error_code,
  notes,
  attempted_at
)
SELECT
  tenant.id,
  touchpoint.id,
  template.id,
  CASE tenant.slug
    WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000d101'::uuid
    WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000d201'::uuid
    WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000d301'::uuid
    ELSE gen_random_uuid()
  END,
  'whatsapp',
  'manual',
  concat('bootstrap-', tenant.slug, '-delivery-001'),
  'delivered',
  'bootstrap-seed',
  NULL,
  'Entrega bootstrap confirmada para o primeiro touchpoint do tenant.',
  timezone('utc', now())
FROM identity.tenants AS tenant
INNER JOIN engagement.touchpoints AS touchpoint
  ON touchpoint.tenant_id = tenant.id
INNER JOIN engagement.templates AS template
  ON template.tenant_id = tenant.id
 AND template.key = 'lead-follow-up-whatsapp'
WHERE touchpoint.public_id = CASE tenant.slug
  WHEN 'bootstrap-ops' THEN '00000000-0000-0000-0000-00000000e101'::uuid
  WHEN 'northwind-group' THEN '00000000-0000-0000-0000-00000000e201'::uuid
  WHEN 'smoke-identity-bootstrap' THEN '00000000-0000-0000-0000-00000000e301'::uuid
  ELSE (
    SELECT bootstrap_touchpoint.public_id
    FROM engagement.touchpoints AS bootstrap_touchpoint
    WHERE bootstrap_touchpoint.tenant_id = tenant.id
    ORDER BY bootstrap_touchpoint.id
    LIMIT 1
  )
END
AND NOT EXISTS (
  SELECT 1
  FROM engagement.touchpoint_deliveries AS delivery
  WHERE delivery.tenant_id = tenant.id
    AND delivery.touchpoint_id = touchpoint.id
    AND delivery.provider_message_id = concat('bootstrap-', tenant.slug, '-delivery-001')
);
