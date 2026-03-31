-- Cria a funcao compartilhada para manter `updated_at` coerente em tabelas operacionais.
-- Regras de negocio nao devem ser embutidas neste gatilho.

CREATE OR REPLACE FUNCTION common_set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = timezone('utc', now());
  RETURN NEW;
END;
$$;
