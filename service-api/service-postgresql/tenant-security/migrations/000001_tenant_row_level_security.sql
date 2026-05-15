DO $$
DECLARE
  target_table record;
  policy_name text;
BEGIN
  FOR target_table IN
    SELECT table_schema, table_name
    FROM information_schema.columns
    WHERE column_name = 'tenant_id'
      AND table_schema NOT IN ('pg_catalog', 'information_schema')
    GROUP BY table_schema, table_name
  LOOP
    policy_name := format('%I_tenant_isolation', target_table.table_name);

    EXECUTE format(
      'ALTER TABLE %I.%I ENABLE ROW LEVEL SECURITY',
      target_table.table_schema,
      target_table.table_name
    );

    IF NOT EXISTS (
      SELECT 1
      FROM pg_policies
      WHERE schemaname = target_table.table_schema
        AND tablename = target_table.table_name
        AND policyname = policy_name
    ) THEN
      EXECUTE format(
        'CREATE POLICY %I ON %I.%I USING (tenant_id::text = current_setting(''erp.tenant_id'', true)) WITH CHECK (tenant_id::text = current_setting(''erp.tenant_id'', true))',
        policy_name,
        target_table.table_schema,
        target_table.table_name
      );
    END IF;
  END LOOP;
END $$;
