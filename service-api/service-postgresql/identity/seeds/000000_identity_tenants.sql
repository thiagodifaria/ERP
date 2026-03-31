-- Garante os tenants bootstrap historicos usados pelos services em memoria.
-- Isso alinha o bootstrap relacional ao comportamento atual da API de identidade.

INSERT INTO identity.tenants (public_id, slug, display_name, status)
VALUES
  ('0195e7a0-7a9c-7c1f-8a44-4a6e50000001', 'bootstrap-ops', 'Bootstrap Ops', 'active'),
  ('0195e7a0-7a9c-7c1f-8a44-4a6e50000002', 'northwind-group', 'Northwind Group', 'active')
ON CONFLICT (slug) DO NOTHING;
