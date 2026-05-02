ALTER TABLE engagement.touchpoints
  ADD COLUMN IF NOT EXISTS business_entity_type VARCHAR(80),
  ADD COLUMN IF NOT EXISTS business_entity_public_id UUID;

UPDATE engagement.touchpoints
SET business_entity_type = 'crm.lead',
    business_entity_public_id = lead_public_id
WHERE lead_public_id IS NOT NULL
  AND business_entity_public_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_business_entity
  ON engagement.touchpoints (tenant_id, business_entity_type, business_entity_public_id);

ALTER TABLE engagement.provider_events
  ADD COLUMN IF NOT EXISTS business_entity_type VARCHAR(80),
  ADD COLUMN IF NOT EXISTS business_entity_public_id UUID;

UPDATE engagement.provider_events
SET business_entity_type = 'crm.lead',
    business_entity_public_id = lead_public_id
WHERE lead_public_id IS NOT NULL
  AND business_entity_public_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_engagement_provider_events_business_entity
  ON engagement.provider_events (tenant_id, business_entity_type, business_entity_public_id);
