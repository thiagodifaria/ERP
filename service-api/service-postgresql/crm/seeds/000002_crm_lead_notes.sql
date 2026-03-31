-- Popula uma nota bootstrap por lead quando o tenant ainda nao tem historico operacional.
-- Este seed alinha a camada relacional com a nova trilha publica de notas do CRM.

INSERT INTO crm.lead_notes (tenant_id, lead_id, public_id, category, body, created_at)
SELECT
  lead.tenant_id,
  lead.id,
  gen_random_uuid(),
  'qualification',
  'Primeiro contato capturado e aguardando abordagem comercial.',
  timezone('utc', now())
FROM crm.leads AS lead
WHERE NOT EXISTS (
  SELECT 1
  FROM crm.lead_notes AS note
  WHERE note.lead_id = lead.id
);
