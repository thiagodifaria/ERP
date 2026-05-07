CREATE TABLE IF NOT EXISTS catalog.item_versions (
  id BIGSERIAL PRIMARY KEY,
  item_id BIGINT NOT NULL REFERENCES catalog.items(id) ON DELETE CASCADE,
  version_number INTEGER NOT NULL,
  change_summary TEXT NOT NULL,
  snapshot_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ux_catalog_item_versions_item_version UNIQUE (item_id, version_number)
);

CREATE INDEX IF NOT EXISTS idx_catalog_item_versions_item_id
  ON catalog.item_versions (item_id, version_number DESC);

INSERT INTO catalog.item_versions (
  item_id,
  version_number,
  change_summary,
  snapshot_json
)
SELECT
  item.id,
  item.version_number,
  'initial_version',
  jsonb_build_object(
    'publicId', item.public_id,
    'sku', item.sku,
    'name', item.name,
    'itemType', item.item_type,
    'unitCode', item.unit_code,
    'priceBaseCents', item.price_base_cents,
    'currencyCode', item.currency_code,
    'active', item.active,
    'versionNumber', item.version_number,
    'attributes', item.attributes_json
  )
FROM catalog.items AS item
WHERE NOT EXISTS (
  SELECT 1
  FROM catalog.item_versions AS version
  WHERE version.item_id = item.id
    AND version.version_number = item.version_number
);
