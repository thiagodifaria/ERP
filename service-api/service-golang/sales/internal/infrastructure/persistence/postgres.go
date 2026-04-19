// Adapters relacionais do contexto sales.
// Cada tabela segue ownership do proprio contexto em PostgreSQL.
package persistence

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

type PostgresOpportunityRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresProposalRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresSaleRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresInvoiceRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresOpportunityRepository(database *sql.DB, tenantSlug string) (*PostgresOpportunityRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresOpportunityRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresOpportunityRepository) List() []entity.Opportunity {
	rows, err := repository.database.Query(
		`
      SELECT public_id::text, lead_public_id::text, title, stage, COALESCE(owner_user_public_id::text, ''), amount_cents
      FROM sales.opportunities
      WHERE tenant_id = $1
      ORDER BY created_at, id
    `,
		repository.tenantID,
	)
	if err != nil {
		return []entity.Opportunity{}
	}
	defer rows.Close()

	response := make([]entity.Opportunity, 0)
	for rows.Next() {
		opportunity, scanErr := scanOpportunity(rows)
		if scanErr == nil {
			response = append(response, opportunity)
		}
	}

	return response
}

func (repository *PostgresOpportunityRepository) FindByPublicID(publicID string) *entity.Opportunity {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, lead_public_id::text, title, stage, COALESCE(owner_user_public_id::text, ''), amount_cents
      FROM sales.opportunities
      WHERE tenant_id = $1
        AND public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	opportunity, err := scanOpportunity(row)
	if err != nil {
		return nil
	}

	return &opportunity
}

func (repository *PostgresOpportunityRepository) Save(opportunity entity.Opportunity) entity.Opportunity {
	publicID := uuid.MustParse(opportunity.PublicID)
	leadPublicID := uuid.MustParse(opportunity.LeadPublicID)
	var ownerUserID *uuid.UUID
	if strings.TrimSpace(opportunity.OwnerUserID) != "" {
		parsed := uuid.MustParse(opportunity.OwnerUserID)
		ownerUserID = &parsed
	}

	row := repository.database.QueryRow(
		`
      INSERT INTO sales.opportunities (tenant_id, public_id, lead_public_id, owner_user_public_id, title, stage, amount_cents)
      VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING public_id::text, lead_public_id::text, title, stage, COALESCE(owner_user_public_id::text, ''), amount_cents
    `,
		repository.tenantID,
		publicID,
		leadPublicID,
		ownerUserID,
		opportunity.Title,
		opportunity.Stage,
		opportunity.AmountCents,
	)

	saved, err := scanOpportunity(row)
	if err != nil {
		return opportunity
	}

	return saved
}

func (repository *PostgresOpportunityRepository) Update(opportunity entity.Opportunity) entity.Opportunity {
	publicID := uuid.MustParse(opportunity.PublicID)
	var ownerUserID *uuid.UUID
	if strings.TrimSpace(opportunity.OwnerUserID) != "" {
		parsed := uuid.MustParse(opportunity.OwnerUserID)
		ownerUserID = &parsed
	}

	row := repository.database.QueryRow(
		`
      UPDATE sales.opportunities
      SET
        title = $3,
        stage = $4,
        owner_user_public_id = $5,
        amount_cents = $6
      WHERE tenant_id = $1
        AND public_id = $2
      RETURNING public_id::text, lead_public_id::text, title, stage, COALESCE(owner_user_public_id::text, ''), amount_cents
    `,
		repository.tenantID,
		publicID,
		opportunity.Title,
		opportunity.Stage,
		ownerUserID,
		opportunity.AmountCents,
	)

	updated, err := scanOpportunity(row)
	if err != nil {
		return opportunity
	}

	return updated
}

func NewPostgresProposalRepository(database *sql.DB, tenantSlug string) (*PostgresProposalRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresProposalRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresProposalRepository) ListByOpportunityPublicID(opportunityPublicID string) []entity.Proposal {
	parsedOpportunityPublicID, err := uuid.Parse(strings.TrimSpace(opportunityPublicID))
	if err != nil {
		return []entity.Proposal{}
	}

	rows, err := repository.database.Query(
		`
      SELECT proposal.public_id::text, opportunity.public_id::text, proposal.title, proposal.status, proposal.amount_cents
      FROM sales.proposals AS proposal
      INNER JOIN sales.opportunities AS opportunity
        ON opportunity.id = proposal.opportunity_id
      WHERE proposal.tenant_id = $1
        AND opportunity.public_id = $2
      ORDER BY proposal.created_at, proposal.id
    `,
		repository.tenantID,
		parsedOpportunityPublicID,
	)
	if err != nil {
		return []entity.Proposal{}
	}
	defer rows.Close()

	response := make([]entity.Proposal, 0)
	for rows.Next() {
		proposal, scanErr := scanProposal(rows)
		if scanErr == nil {
			response = append(response, proposal)
		}
	}

	return response
}

func (repository *PostgresProposalRepository) FindByPublicID(publicID string) *entity.Proposal {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT proposal.public_id::text, opportunity.public_id::text, proposal.title, proposal.status, proposal.amount_cents
      FROM sales.proposals AS proposal
      INNER JOIN sales.opportunities AS opportunity
        ON opportunity.id = proposal.opportunity_id
      WHERE proposal.tenant_id = $1
        AND proposal.public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	proposal, err := scanProposal(row)
	if err != nil {
		return nil
	}

	return &proposal
}

func (repository *PostgresProposalRepository) Save(proposal entity.Proposal) entity.Proposal {
	opportunityPublicID := uuid.MustParse(proposal.OpportunityPublicID)

	row := repository.database.QueryRow(
		`
      INSERT INTO sales.proposals (tenant_id, opportunity_id, public_id, title, status, amount_cents)
      SELECT
        $1,
        opportunity.id,
        $3,
        $4,
        $5,
        $6
      FROM sales.opportunities AS opportunity
      WHERE opportunity.tenant_id = $1
        AND opportunity.public_id = $2
      RETURNING public_id::text, $2::text, title, status, amount_cents
    `,
		repository.tenantID,
		opportunityPublicID,
		uuid.MustParse(proposal.PublicID),
		proposal.Title,
		proposal.Status,
		proposal.AmountCents,
	)

	saved, err := scanProposal(row)
	if err != nil {
		return proposal
	}

	return saved
}

func (repository *PostgresProposalRepository) Update(proposal entity.Proposal) entity.Proposal {
	publicID := uuid.MustParse(proposal.PublicID)

	row := repository.database.QueryRow(
		`
      UPDATE sales.proposals AS proposal
      SET
        title = $3,
        status = $4,
        amount_cents = $5
      FROM sales.opportunities AS opportunity
      WHERE proposal.tenant_id = $1
        AND proposal.public_id = $2
        AND opportunity.id = proposal.opportunity_id
      RETURNING proposal.public_id::text, opportunity.public_id::text, proposal.title, proposal.status, proposal.amount_cents
    `,
		repository.tenantID,
		publicID,
		proposal.Title,
		proposal.Status,
		proposal.AmountCents,
	)

	updated, err := scanProposal(row)
	if err != nil {
		return proposal
	}

	return updated
}

func NewPostgresSaleRepository(database *sql.DB, tenantSlug string) (*PostgresSaleRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresSaleRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func NewPostgresInvoiceRepository(database *sql.DB, tenantSlug string) (*PostgresInvoiceRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresInvoiceRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresSaleRepository) List() []entity.Sale {
	rows, err := repository.database.Query(
		`
      SELECT sale.public_id::text, opportunity.public_id::text, proposal.public_id::text, sale.status, sale.amount_cents
      FROM sales.sales AS sale
      INNER JOIN sales.opportunities AS opportunity
        ON opportunity.id = sale.opportunity_id
      INNER JOIN sales.proposals AS proposal
        ON proposal.id = sale.proposal_id
      WHERE sale.tenant_id = $1
      ORDER BY sale.created_at, sale.id
    `,
		repository.tenantID,
	)
	if err != nil {
		return []entity.Sale{}
	}
	defer rows.Close()

	response := make([]entity.Sale, 0)
	for rows.Next() {
		sale, scanErr := scanSale(rows)
		if scanErr == nil {
			response = append(response, sale)
		}
	}

	return response
}

func (repository *PostgresSaleRepository) FindByPublicID(publicID string) *entity.Sale {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT sale.public_id::text, opportunity.public_id::text, proposal.public_id::text, sale.status, sale.amount_cents
      FROM sales.sales AS sale
      INNER JOIN sales.opportunities AS opportunity
        ON opportunity.id = sale.opportunity_id
      INNER JOIN sales.proposals AS proposal
        ON proposal.id = sale.proposal_id
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	sale, err := scanSale(row)
	if err != nil {
		return nil
	}

	return &sale
}

func (repository *PostgresSaleRepository) FindByProposalPublicID(proposalPublicID string) *entity.Sale {
	parsedProposalPublicID, err := uuid.Parse(strings.TrimSpace(proposalPublicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT sale.public_id::text, opportunity.public_id::text, proposal.public_id::text, sale.status, sale.amount_cents
      FROM sales.sales AS sale
      INNER JOIN sales.opportunities AS opportunity
        ON opportunity.id = sale.opportunity_id
      INNER JOIN sales.proposals AS proposal
        ON proposal.id = sale.proposal_id
      WHERE sale.tenant_id = $1
        AND proposal.public_id = $2
    `,
		repository.tenantID,
		parsedProposalPublicID,
	)

	sale, err := scanSale(row)
	if err != nil {
		return nil
	}

	return &sale
}

func (repository *PostgresSaleRepository) Save(sale entity.Sale) entity.Sale {
	opportunityPublicID := uuid.MustParse(sale.OpportunityPublicID)
	proposalPublicID := uuid.MustParse(sale.ProposalPublicID)

	row := repository.database.QueryRow(
		`
      INSERT INTO sales.sales (tenant_id, opportunity_id, proposal_id, public_id, status, amount_cents)
      SELECT
        $1,
        opportunity.id,
        proposal.id,
        $4,
        $5,
        $6
      FROM sales.opportunities AS opportunity
      INNER JOIN sales.proposals AS proposal
        ON proposal.opportunity_id = opportunity.id
      WHERE opportunity.tenant_id = $1
        AND opportunity.public_id = $2
        AND proposal.public_id = $3
      RETURNING $4::text, $2::text, $3::text, status, amount_cents
    `,
		repository.tenantID,
		opportunityPublicID,
		proposalPublicID,
		uuid.MustParse(sale.PublicID),
		sale.Status,
		sale.AmountCents,
	)

	saved, err := scanSale(row)
	if err != nil {
		return sale
	}

	return saved
}

func (repository *PostgresSaleRepository) Update(sale entity.Sale) entity.Sale {
	publicID := uuid.MustParse(sale.PublicID)

	row := repository.database.QueryRow(
		`
      UPDATE sales.sales AS sale
      SET
        status = $3,
        amount_cents = $4
      FROM sales.opportunities AS opportunity,
           sales.proposals AS proposal
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2
        AND opportunity.id = sale.opportunity_id
        AND proposal.id = sale.proposal_id
      RETURNING sale.public_id::text, opportunity.public_id::text, proposal.public_id::text, sale.status, sale.amount_cents
    `,
		repository.tenantID,
		publicID,
		sale.Status,
		sale.AmountCents,
	)

	updated, err := scanSale(row)
	if err != nil {
		return sale
	}

	return updated
}

func (repository *PostgresInvoiceRepository) List() []entity.Invoice {
	rows, err := repository.database.Query(
		`
      SELECT invoice.public_id::text, sale.public_id::text, invoice.invoice_number, invoice.status, invoice.amount_cents, to_char(invoice.due_date, 'YYYY-MM-DD'), COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.invoices AS invoice
      INNER JOIN sales.sales AS sale
        ON sale.id = invoice.sale_id
      WHERE invoice.tenant_id = $1
      ORDER BY invoice.created_at, invoice.id
    `,
		repository.tenantID,
	)
	if err != nil {
		return []entity.Invoice{}
	}
	defer rows.Close()

	response := make([]entity.Invoice, 0)
	for rows.Next() {
		invoice, scanErr := scanInvoice(rows)
		if scanErr == nil {
			response = append(response, invoice)
		}
	}

	return response
}

func (repository *PostgresInvoiceRepository) FindByPublicID(publicID string) *entity.Invoice {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT invoice.public_id::text, sale.public_id::text, invoice.invoice_number, invoice.status, invoice.amount_cents, to_char(invoice.due_date, 'YYYY-MM-DD'), COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.invoices AS invoice
      INNER JOIN sales.sales AS sale
        ON sale.id = invoice.sale_id
      WHERE invoice.tenant_id = $1
        AND invoice.public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	invoice, err := scanInvoice(row)
	if err != nil {
		return nil
	}

	return &invoice
}

func (repository *PostgresInvoiceRepository) FindBySalePublicID(salePublicID string) *entity.Invoice {
	parsedSalePublicID, err := uuid.Parse(strings.TrimSpace(salePublicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT invoice.public_id::text, sale.public_id::text, invoice.invoice_number, invoice.status, invoice.amount_cents, to_char(invoice.due_date, 'YYYY-MM-DD'), COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.invoices AS invoice
      INNER JOIN sales.sales AS sale
        ON sale.id = invoice.sale_id
      WHERE invoice.tenant_id = $1
        AND sale.public_id = $2
    `,
		repository.tenantID,
		parsedSalePublicID,
	)

	invoice, err := scanInvoice(row)
	if err != nil {
		return nil
	}

	return &invoice
}

func (repository *PostgresInvoiceRepository) Save(invoice entity.Invoice) entity.Invoice {
	salePublicID := uuid.MustParse(invoice.SalePublicID)
	publicID := uuid.MustParse(invoice.PublicID)
	var paidAt *time.Time
	if strings.TrimSpace(invoice.PaidAt) != "" {
		parsedPaidAt := mustParseRFC3339(invoice.PaidAt)
		paidAt = &parsedPaidAt
	}

	row := repository.database.QueryRow(
		`
      INSERT INTO sales.invoices (tenant_id, sale_id, public_id, invoice_number, status, amount_cents, due_date, paid_at)
      SELECT
        $1,
        sale.id,
        $3,
        $4,
        $5,
        $6,
        $7::date,
        $8
      FROM sales.sales AS sale
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2
      RETURNING $3::text, $2::text, invoice_number, status, amount_cents, to_char(due_date, 'YYYY-MM-DD'), COALESCE(to_char(paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		salePublicID,
		publicID,
		invoice.Number,
		invoice.Status,
		invoice.AmountCents,
		invoice.DueDate,
		paidAt,
	)

	saved, err := scanInvoice(row)
	if err != nil {
		return invoice
	}

	return saved
}

func (repository *PostgresInvoiceRepository) Update(invoice entity.Invoice) entity.Invoice {
	publicID := uuid.MustParse(invoice.PublicID)
	var paidAt *time.Time
	if strings.TrimSpace(invoice.PaidAt) != "" {
		parsedPaidAt := mustParseRFC3339(invoice.PaidAt)
		paidAt = &parsedPaidAt
	}

	row := repository.database.QueryRow(
		`
      UPDATE sales.invoices AS invoice
      SET
        invoice_number = $3,
        status = $4,
        amount_cents = $5,
        due_date = $6::date,
        paid_at = $7
      FROM sales.sales AS sale
      WHERE invoice.tenant_id = $1
        AND invoice.public_id = $2
        AND sale.id = invoice.sale_id
      RETURNING invoice.public_id::text, sale.public_id::text, invoice.invoice_number, invoice.status, invoice.amount_cents, to_char(invoice.due_date, 'YYYY-MM-DD'), COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		publicID,
		invoice.Number,
		invoice.Status,
		invoice.AmountCents,
		invoice.DueDate,
		paidAt,
	)

	updated, err := scanInvoice(row)
	if err != nil {
		return invoice
	}

	return updated
}

func lookupTenantID(database *sql.DB, tenantSlug string) (int64, error) {
	var tenantID int64
	if err := database.QueryRow(
		`
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
    `,
		strings.TrimSpace(tenantSlug),
	).Scan(&tenantID); err != nil {
		return 0, err
	}

	return tenantID, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanOpportunity(scanner scanner) (entity.Opportunity, error) {
	var publicID string
	var leadPublicID string
	var title string
	var stage string
	var ownerUserID string
	var amountCents int64

	if err := scanner.Scan(&publicID, &leadPublicID, &title, &stage, &ownerUserID, &amountCents); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Opportunity{}, err
		}

		return entity.Opportunity{}, err
	}

	return entity.RestoreOpportunity(publicID, leadPublicID, title, ownerUserID, amountCents, stage)
}

func scanProposal(scanner scanner) (entity.Proposal, error) {
	var publicID string
	var opportunityPublicID string
	var title string
	var status string
	var amountCents int64

	if err := scanner.Scan(&publicID, &opportunityPublicID, &title, &status, &amountCents); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Proposal{}, err
		}

		return entity.Proposal{}, err
	}

	return entity.RestoreProposal(publicID, opportunityPublicID, title, amountCents, status)
}

func scanSale(scanner scanner) (entity.Sale, error) {
	var publicID string
	var opportunityPublicID string
	var proposalPublicID string
	var status string
	var amountCents int64

	if err := scanner.Scan(&publicID, &opportunityPublicID, &proposalPublicID, &status, &amountCents); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Sale{}, err
		}

		return entity.Sale{}, err
	}

	return entity.RestoreSale(publicID, opportunityPublicID, proposalPublicID, amountCents, status)
}

func scanInvoice(scanner scanner) (entity.Invoice, error) {
	var publicID string
	var salePublicID string
	var number string
	var status string
	var amountCents int64
	var dueDate string
	var paidAt string

	if err := scanner.Scan(&publicID, &salePublicID, &number, &status, &amountCents, &dueDate, &paidAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Invoice{}, err
		}

		return entity.Invoice{}, err
	}

	return entity.RestoreInvoice(publicID, salePublicID, number, amountCents, dueDate, status, paidAt)
}

func mustParseRFC3339(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}

	return parsed
}
