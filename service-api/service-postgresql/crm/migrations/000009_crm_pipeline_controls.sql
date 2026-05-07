ALTER TABLE crm.pipeline_configs
  ADD COLUMN IF NOT EXISTS territory_rules_json JSONB NOT NULL DEFAULT '[]'::jsonb,
  ADD COLUMN IF NOT EXISTS approval_policies_json JSONB NOT NULL DEFAULT '[]'::jsonb;
