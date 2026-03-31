-- Ativa extensoes compartilhadas que serao reutilizadas por mais de um contexto.
-- Este arquivo nao deve carregar regra de negocio.

CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
