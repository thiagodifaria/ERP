CREATE TABLE IF NOT EXISTS simulation.load_benchmark_runs (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT REFERENCES identity.tenants(id) ON DELETE SET NULL,
  public_id UUID NOT NULL UNIQUE,
  benchmark_key TEXT NOT NULL,
  input_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  output_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now())
);

CREATE INDEX IF NOT EXISTS idx_simulation_load_benchmark_runs_tenant_created
  ON simulation.load_benchmark_runs (tenant_id, created_at DESC);
