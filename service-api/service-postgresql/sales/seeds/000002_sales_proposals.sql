-- Popula uma proposta bootstrap por oportunidade quando ainda nao existe cobertura relacional.

INSERT INTO sales.proposals (tenant_id, opportunity_id, public_id, title, status, amount_cents)
SELECT
  opportunity.tenant_id,
  opportunity.id,
  gen_random_uuid(),
  concat(opportunity.title, ' Proposal'),
  'accepted',
  opportunity.amount_cents
FROM sales.opportunities AS opportunity
WHERE NOT EXISTS (
  SELECT 1
  FROM sales.proposals AS proposal
  WHERE proposal.opportunity_id = opportunity.id
);
