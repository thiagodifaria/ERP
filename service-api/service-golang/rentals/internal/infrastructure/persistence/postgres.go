package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/entity"
	repo "github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
)

type PostgresContractRepository struct {
	database            *sql.DB
	bootstrapTenantSlug string
	cachedTenantIDs     map[string]int64
}

func NewPostgresContractRepository(database *sql.DB, bootstrapTenantSlug string) *PostgresContractRepository {
	return &PostgresContractRepository{
		database:            database,
		bootstrapTenantSlug: strings.ToLower(strings.TrimSpace(bootstrapTenantSlug)),
		cachedTenantIDs:     map[string]int64{},
	}
}

func (repository *PostgresContractRepository) List(filters repo.ContractFilters) []entity.Contract {
	tenantID, tenantSlug, err := repository.resolveTenant(filters.TenantSlug)
	if err != nil {
		return []entity.Contract{}
	}

	query := `
      SELECT public_id::text, customer_public_id::text, title, property_code, currency_code, amount_cents, billing_day, starts_at, ends_at, status, terminated_at, termination_reason, created_at, updated_at
      FROM rentals.contracts
      WHERE tenant_id = $1
    `
	args := []any{tenantID}

	if normalizedStatus := strings.ToLower(strings.TrimSpace(filters.Status)); normalizedStatus != "" {
		query += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, normalizedStatus)
	}
	if normalizedCustomer := strings.TrimSpace(filters.CustomerPublicID); normalizedCustomer != "" {
		query += fmt.Sprintf(" AND customer_public_id = $%d::uuid", len(args)+1)
		args = append(args, normalizedCustomer)
	}

	query += " ORDER BY created_at, id"

	rows, err := repository.database.Query(query, args...)
	if err != nil {
		return []entity.Contract{}
	}
	defer rows.Close()

	response := make([]entity.Contract, 0)
	for rows.Next() {
		contract, ok := scanContract(rows, tenantSlug)
		if ok {
			response = append(response, contract)
		}
	}

	return response
}

func (repository *PostgresContractRepository) Summary(tenantSlug string) repo.ContractSummary {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return repo.ContractSummary{TenantSlug: strings.ToLower(strings.TrimSpace(tenantSlug))}
	}

	row := repository.database.QueryRow(
		`
      SELECT
        (SELECT count(*) FROM rentals.contracts WHERE tenant_id = $1),
        (SELECT count(*) FROM rentals.contracts WHERE tenant_id = $1 AND status = 'active'),
        (SELECT count(*) FROM rentals.contracts WHERE tenant_id = $1 AND status = 'terminated'),
        (SELECT count(*) FROM rentals.contract_charges WHERE tenant_id = $1 AND status = 'scheduled'),
        (SELECT count(*) FROM rentals.contract_charges WHERE tenant_id = $1 AND status = 'cancelled'),
        (SELECT count(*) FROM rentals.contract_adjustments WHERE tenant_id = $1),
        (SELECT count(*) FROM rentals.contract_events WHERE tenant_id = $1),
        (SELECT count(*) FROM rentals.outbox_events WHERE tenant_id = $1 AND status = 'pending'),
        (SELECT COALESCE(sum(amount_cents), 0) FROM rentals.contract_charges WHERE tenant_id = $1 AND status = 'scheduled'),
        (SELECT COALESCE(sum(amount_cents), 0) FROM rentals.contract_charges WHERE tenant_id = $1 AND status = 'cancelled')
    `,
		tenantID,
	)

	summary := repo.ContractSummary{TenantSlug: normalizedTenantSlug}
	if err := row.Scan(
		&summary.TotalContracts,
		&summary.ActiveContracts,
		&summary.TerminatedContracts,
		&summary.ScheduledCharges,
		&summary.CancelledCharges,
		&summary.Adjustments,
		&summary.HistoryEvents,
		&summary.PendingOutbox,
		&summary.ScheduledAmountCents,
		&summary.CancelledAmountCents,
	); err != nil {
		return repo.ContractSummary{TenantSlug: normalizedTenantSlug}
	}

	return summary
}

func (repository *PostgresContractRepository) FindByPublicID(tenantSlug string, publicID string) (entity.Contract, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Contract{}, false
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, customer_public_id::text, title, property_code, currency_code, amount_cents, billing_day, starts_at, ends_at, status, terminated_at, termination_reason, created_at, updated_at
      FROM rentals.contracts
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
    `,
		tenantID,
		strings.TrimSpace(publicID),
	)

	return scanContract(row, normalizedTenantSlug)
}

func (repository *PostgresContractRepository) Create(contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) entity.Contract {
	tenantID, _, err := repository.resolveTenant(contract.TenantSlug)
	if err != nil {
		return contract
	}

	transaction, err := repository.database.Begin()
	if err != nil {
		return contract
	}
	defer transaction.Rollback()

	var contractID int64
	err = transaction.QueryRow(
		`
      INSERT INTO rentals.contracts (
        tenant_id, public_id, customer_public_id, title, property_code, currency_code, amount_cents, billing_day, starts_at, ends_at, status, created_at, updated_at
      )
      VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
      RETURNING id
    `,
		tenantID,
		contract.PublicID,
		contract.CustomerPublicID,
		contract.Title,
		contract.PropertyCode,
		contract.CurrencyCode,
		contract.AmountCents,
		contract.BillingDay,
		contract.StartsAt,
		contract.EndsAt,
		contract.Status,
		contract.CreatedAt,
		contract.UpdatedAt,
	).Scan(&contractID)
	if err != nil {
		return contract
	}

	for _, charge := range charges {
		if _, err := transaction.Exec(
			`
        INSERT INTO rentals.contract_charges (
          tenant_id, contract_id, public_id, due_date, amount_cents, status, created_at, updated_at
        )
        VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, $8)
      `,
			tenantID,
			contractID,
			charge.PublicID,
			charge.DueDate,
			charge.AmountCents,
			charge.Status,
			charge.CreatedAt,
			charge.UpdatedAt,
		); err != nil {
			return contract
		}
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO rentals.contract_events (
        tenant_id, contract_id, public_id, event_code, summary, recorded_by, created_at
      )
      VALUES ($1, $2, $3::uuid, $4, $5, $6, $7)
    `,
		tenantID,
		contractID,
		event.PublicID,
		event.EventCode,
		event.Summary,
		event.RecordedBy,
		event.CreatedAt,
	); err != nil {
		return contract
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO rentals.outbox_events (
        tenant_id, public_id, aggregate_type, aggregate_public_id, event_type, payload, status, created_at
      )
      VALUES ($1, $2::uuid, $3, $4::uuid, $5, $6::jsonb, $7, $8)
    `,
		tenantID,
		outbox.PublicID,
		outbox.AggregateType,
		outbox.AggregatePublicID,
		outbox.EventType,
		outbox.Payload,
		outbox.Status,
		outbox.CreatedAt,
	); err != nil {
		return contract
	}

	if err := transaction.Commit(); err != nil {
		return contract
	}

	return contract
}

func (repository *PostgresContractRepository) ListCharges(tenantSlug string, contractPublicID string, status string) []entity.Charge {
	tenantID, _, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return []entity.Charge{}
	}

	query := `
      SELECT charge.public_id::text, charge.due_date, charge.amount_cents, charge.status, charge.created_at, charge.updated_at
      FROM rentals.contract_charges AS charge
      INNER JOIN rentals.contracts AS contract
        ON contract.id = charge.contract_id
      WHERE charge.tenant_id = $1
        AND contract.public_id = $2::uuid
    `
	args := []any{tenantID, strings.TrimSpace(contractPublicID)}
	if normalizedStatus := strings.ToLower(strings.TrimSpace(status)); normalizedStatus != "" {
		query += fmt.Sprintf(" AND charge.status = $%d", len(args)+1)
		args = append(args, normalizedStatus)
	}
	query += " ORDER BY charge.due_date, charge.id"

	rows, err := repository.database.Query(query, args...)
	if err != nil {
		return []entity.Charge{}
	}
	defer rows.Close()

	response := make([]entity.Charge, 0)
	for rows.Next() {
		var charge entity.Charge
		charge.ContractPublicID = strings.TrimSpace(contractPublicID)
		if err := rows.Scan(&charge.PublicID, &charge.DueDate, &charge.AmountCents, &charge.Status, &charge.CreatedAt, &charge.UpdatedAt); err == nil {
			response = append(response, charge)
		}
	}

	return response
}

func (repository *PostgresContractRepository) ListEvents(tenantSlug string, contractPublicID string) []entity.Event {
	tenantID, _, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return []entity.Event{}
	}

	rows, err := repository.database.Query(
		`
      SELECT event.public_id::text, event.event_code, event.summary, event.recorded_by, event.created_at
      FROM rentals.contract_events AS event
      INNER JOIN rentals.contracts AS contract
        ON contract.id = event.contract_id
      WHERE event.tenant_id = $1
        AND contract.public_id = $2::uuid
      ORDER BY event.created_at, event.id
    `,
		tenantID,
		strings.TrimSpace(contractPublicID),
	)
	if err != nil {
		return []entity.Event{}
	}
	defer rows.Close()

	response := make([]entity.Event, 0)
	for rows.Next() {
		var event entity.Event
		event.ContractPublicID = strings.TrimSpace(contractPublicID)
		if err := rows.Scan(&event.PublicID, &event.EventCode, &event.Summary, &event.RecordedBy, &event.CreatedAt); err == nil {
			response = append(response, event)
		}
	}

	return response
}

func (repository *PostgresContractRepository) ListAdjustments(tenantSlug string, contractPublicID string) []entity.Adjustment {
	tenantID, _, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return []entity.Adjustment{}
	}

	rows, err := repository.database.Query(
		`
      SELECT adjustment.public_id::text, adjustment.effective_at, adjustment.previous_amount_cents, adjustment.new_amount_cents, adjustment.reason, adjustment.recorded_by, adjustment.created_at
      FROM rentals.contract_adjustments AS adjustment
      INNER JOIN rentals.contracts AS contract
        ON contract.id = adjustment.contract_id
      WHERE adjustment.tenant_id = $1
        AND contract.public_id = $2::uuid
      ORDER BY adjustment.effective_at, adjustment.created_at, adjustment.id
    `,
		tenantID,
		strings.TrimSpace(contractPublicID),
	)
	if err != nil {
		return []entity.Adjustment{}
	}
	defer rows.Close()

	response := make([]entity.Adjustment, 0)
	for rows.Next() {
		var adjustment entity.Adjustment
		adjustment.ContractPublicID = strings.TrimSpace(contractPublicID)
		if err := rows.Scan(&adjustment.PublicID, &adjustment.EffectiveAt, &adjustment.PreviousAmountCents, &adjustment.NewAmountCents, &adjustment.Reason, &adjustment.RecordedBy, &adjustment.CreatedAt); err == nil {
			response = append(response, adjustment)
		}
	}

	return response
}

func (repository *PostgresContractRepository) SaveAdjustment(tenantSlug string, contract entity.Contract, adjustment entity.Adjustment, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool) {
	tenantID, _, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Contract{}, false
	}

	transaction, err := repository.database.Begin()
	if err != nil {
		return entity.Contract{}, false
	}
	defer transaction.Rollback()

	contractID, found := repository.lookupContractID(transaction, tenantID, contract.PublicID)
	if !found {
		return entity.Contract{}, false
	}

	if _, err := transaction.Exec(
		`
      UPDATE rentals.contracts
      SET amount_cents = $3, updated_at = $4
      WHERE tenant_id = $1
        AND public_id = $2::uuid
    `,
		tenantID,
		contract.PublicID,
		contract.AmountCents,
		contract.UpdatedAt,
	); err != nil {
		return entity.Contract{}, false
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO rentals.contract_adjustments (
        tenant_id, contract_id, public_id, effective_at, previous_amount_cents, new_amount_cents, reason, recorded_by, created_at
      )
      VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, $8, $9)
    `,
		tenantID,
		contractID,
		adjustment.PublicID,
		adjustment.EffectiveAt,
		adjustment.PreviousAmountCents,
		adjustment.NewAmountCents,
		adjustment.Reason,
		adjustment.RecordedBy,
		adjustment.CreatedAt,
	); err != nil {
		return entity.Contract{}, false
	}

	for _, charge := range charges {
		if _, err := transaction.Exec(
			`
        UPDATE rentals.contract_charges
        SET amount_cents = $4, status = $5, updated_at = $6
        WHERE tenant_id = $1
          AND contract_id = $2
          AND public_id = $3::uuid
      `,
			tenantID,
			contractID,
			charge.PublicID,
			charge.AmountCents,
			charge.Status,
			charge.UpdatedAt,
		); err != nil {
			return entity.Contract{}, false
		}
	}

	if err := repository.insertEventAndOutbox(transaction, tenantID, contractID, event, outbox); err != nil {
		return entity.Contract{}, false
	}

	if err := transaction.Commit(); err != nil {
		return entity.Contract{}, false
	}

	return contract, true
}

func (repository *PostgresContractRepository) SaveTermination(tenantSlug string, contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool) {
	tenantID, _, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Contract{}, false
	}

	transaction, err := repository.database.Begin()
	if err != nil {
		return entity.Contract{}, false
	}
	defer transaction.Rollback()

	contractID, found := repository.lookupContractID(transaction, tenantID, contract.PublicID)
	if !found {
		return entity.Contract{}, false
	}

	if _, err := transaction.Exec(
		`
      UPDATE rentals.contracts
      SET status = $3, ends_at = $4, terminated_at = $5, termination_reason = $6, updated_at = $7
      WHERE tenant_id = $1
        AND public_id = $2::uuid
    `,
		tenantID,
		contract.PublicID,
		contract.Status,
		contract.EndsAt,
		contract.TerminatedAt,
		contract.TerminationReason,
		contract.UpdatedAt,
	); err != nil {
		return entity.Contract{}, false
	}

	for _, charge := range charges {
		if _, err := transaction.Exec(
			`
        UPDATE rentals.contract_charges
        SET amount_cents = $4, status = $5, updated_at = $6
        WHERE tenant_id = $1
          AND contract_id = $2
          AND public_id = $3::uuid
      `,
			tenantID,
			contractID,
			charge.PublicID,
			charge.AmountCents,
			charge.Status,
			charge.UpdatedAt,
		); err != nil {
			return entity.Contract{}, false
		}
	}

	if err := repository.insertEventAndOutbox(transaction, tenantID, contractID, event, outbox); err != nil {
		return entity.Contract{}, false
	}

	if err := transaction.Commit(); err != nil {
		return entity.Contract{}, false
	}

	return contract, true
}

func (repository *PostgresContractRepository) resolveTenant(tenantSlug string) (int64, string, error) {
	slug := strings.ToLower(strings.TrimSpace(tenantSlug))
	if slug == "" {
		slug = repository.bootstrapTenantSlug
	}

	if tenantID, ok := repository.cachedTenantIDs[slug]; ok {
		return tenantID, slug, nil
	}

	var tenantID int64
	if err := repository.database.QueryRow(
		`
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
      LIMIT 1
    `,
		slug,
	).Scan(&tenantID); err != nil {
		return 0, slug, err
	}

	repository.cachedTenantIDs[slug] = tenantID
	return tenantID, slug, nil
}

func (repository *PostgresContractRepository) lookupContractID(scanner interface {
	QueryRow(query string, args ...any) *sql.Row
}, tenantID int64, contractPublicID string) (int64, bool) {
	var contractID int64
	if err := scanner.QueryRow(
		`
      SELECT id
      FROM rentals.contracts
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
    `,
		tenantID,
		contractPublicID,
	).Scan(&contractID); err != nil {
		return 0, false
	}

	return contractID, true
}

func (repository *PostgresContractRepository) insertEventAndOutbox(transaction *sql.Tx, tenantID int64, contractID int64, event entity.Event, outbox entity.OutboxEvent) error {
	if _, err := transaction.Exec(
		`
      INSERT INTO rentals.contract_events (
        tenant_id, contract_id, public_id, event_code, summary, recorded_by, created_at
      )
      VALUES ($1, $2, $3::uuid, $4, $5, $6, $7)
    `,
		tenantID,
		contractID,
		event.PublicID,
		event.EventCode,
		event.Summary,
		event.RecordedBy,
		event.CreatedAt,
	); err != nil {
		return err
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO rentals.outbox_events (
        tenant_id, public_id, aggregate_type, aggregate_public_id, event_type, payload, status, created_at
      )
      VALUES ($1, $2::uuid, $3, $4::uuid, $5, $6::jsonb, $7, $8)
    `,
		tenantID,
		outbox.PublicID,
		outbox.AggregateType,
		outbox.AggregatePublicID,
		outbox.EventType,
		outbox.Payload,
		outbox.Status,
		outbox.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

type contractScanner interface {
	Scan(dest ...any) error
}

func scanContract(scanner contractScanner, tenantSlug string) (entity.Contract, bool) {
	var publicID string
	var customerPublicID string
	var title string
	var propertyCode string
	var currencyCode string
	var amountCents int64
	var billingDay int
	var startsAt time.Time
	var endsAt time.Time
	var status string
	var terminatedAt sql.NullTime
	var terminationReason string
	var createdAt time.Time
	var updatedAt time.Time

	if err := scanner.Scan(&publicID, &customerPublicID, &title, &propertyCode, &currencyCode, &amountCents, &billingDay, &startsAt, &endsAt, &status, &terminatedAt, &terminationReason, &createdAt, &updatedAt); err != nil {
		return entity.Contract{}, false
	}

	var terminatedAtValue *time.Time
	if terminatedAt.Valid {
		value := terminatedAt.Time.UTC()
		terminatedAtValue = &value
	}

	contract, err := entity.NewContract(
		publicID,
		tenantSlug,
		customerPublicID,
		title,
		propertyCode,
		currencyCode,
		amountCents,
		billingDay,
		startsAt,
		endsAt,
		status,
		terminatedAtValue,
		terminationReason,
		createdAt,
		updatedAt,
	)
	if err != nil {
		return entity.Contract{}, false
	}

	return contract, true
}
