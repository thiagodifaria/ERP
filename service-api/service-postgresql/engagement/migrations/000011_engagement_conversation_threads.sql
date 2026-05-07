ALTER TABLE engagement.touchpoints
  ADD COLUMN IF NOT EXISTS thread_public_id UUID,
  ADD COLUMN IF NOT EXISTS participant_kind VARCHAR(80),
  ADD COLUMN IF NOT EXISTS participant_public_id UUID;

UPDATE engagement.touchpoints
SET thread_public_id = lead_public_id
WHERE lead_public_id IS NOT NULL
  AND thread_public_id IS NULL;

UPDATE engagement.touchpoints
SET participant_kind = COALESCE(business_entity_type, 'crm.lead'),
    participant_public_id = COALESCE(business_entity_public_id, lead_public_id)
WHERE participant_public_id IS NULL;

ALTER TABLE engagement.touchpoints
  ALTER COLUMN thread_public_id SET NOT NULL,
  ALTER COLUMN participant_kind SET NOT NULL,
  ALTER COLUMN participant_public_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_thread
  ON engagement.touchpoints (tenant_id, thread_public_id, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_engagement_touchpoints_participant
  ON engagement.touchpoints (tenant_id, participant_kind, participant_public_id, updated_at DESC);
