-- Cria a trilha relacional inicial de notas vinculadas a leads do CRM.
-- Historico operacional e follow-ups do relacionamento passam a viver aqui.

CREATE TABLE IF NOT EXISTS crm.lead_notes (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL,
  lead_id BIGINT NOT NULL,
  public_id UUID NOT NULL,
  category VARCHAR(80) NOT NULL DEFAULT 'internal',
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc', now()),
  CONSTRAINT uq_crm_lead_notes_public_id UNIQUE (public_id),
  CONSTRAINT fk_crm_lead_notes_lead_id
    FOREIGN KEY (lead_id)
    REFERENCES crm.leads (id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_crm_lead_notes_tenant_id
  ON crm.lead_notes (tenant_id);

CREATE INDEX IF NOT EXISTS idx_crm_lead_notes_lead_id
  ON crm.lead_notes (lead_id);
