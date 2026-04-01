-- Popula uma venda bootstrap por proposta quando ainda nao existe fechamento relacional.

INSERT INTO sales.sales (tenant_id, opportunity_id, proposal_id, public_id, status, amount_cents)
SELECT
  proposal.tenant_id,
  opportunity.id,
  proposal.id,
  gen_random_uuid(),
  'active',
  proposal.amount_cents
FROM sales.proposals AS proposal
INNER JOIN sales.opportunities AS opportunity
  ON opportunity.id = proposal.opportunity_id
WHERE NOT EXISTS (
  SELECT 1
  FROM sales.sales AS sale
  WHERE sale.proposal_id = proposal.id
);
