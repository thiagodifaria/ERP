package persistence

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

type PostgresInstallmentRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresCommissionRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresPendingItemRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresRenegotiationRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresInstallmentRepository(database *sql.DB, tenantSlug string) (*PostgresInstallmentRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}
	return &PostgresInstallmentRepository{database: database, tenantID: tenantID}, nil
}

func NewPostgresCommissionRepository(database *sql.DB, tenantSlug string) (*PostgresCommissionRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}
	return &PostgresCommissionRepository{database: database, tenantID: tenantID}, nil
}

func NewPostgresPendingItemRepository(database *sql.DB, tenantSlug string) (*PostgresPendingItemRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}
	return &PostgresPendingItemRepository{database: database, tenantID: tenantID}, nil
}

func NewPostgresRenegotiationRepository(database *sql.DB, tenantSlug string) (*PostgresRenegotiationRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}
	return &PostgresRenegotiationRepository{database: database, tenantID: tenantID}, nil
}

func (repository *PostgresInstallmentRepository) ListBySalePublicID(salePublicID string) []entity.Installment {
	rows, err := repository.database.Query(
		`
      SELECT installment.public_id::text, sale.public_id::text, installment.sequence_number, installment.amount_cents, to_char(installment.due_date, 'YYYY-MM-DD'), installment.status
      FROM sales.installments AS installment
      INNER JOIN sales.sales AS sale
        ON sale.id = installment.sale_id
      WHERE installment.tenant_id = $1
        AND sale.public_id = $2::uuid
      ORDER BY installment.sequence_number
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(salePublicID)),
	)
	if err != nil {
		return []entity.Installment{}
	}
	defer rows.Close()

	response := make([]entity.Installment, 0)
	for rows.Next() {
		installment, scanErr := scanInstallment(rows)
		if scanErr == nil {
			response = append(response, installment)
		}
	}

	return response
}

func (repository *PostgresInstallmentRepository) Save(installment entity.Installment) entity.Installment {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.installments (tenant_id, sale_id, public_id, sequence_number, amount_cents, due_date, status)
      SELECT $1, sale.id, $3, $4, $5, $6::date, $7
      FROM sales.sales AS sale
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2::uuid
      RETURNING $3::text, $2::text, sequence_number, amount_cents, to_char(due_date, 'YYYY-MM-DD'), status
    `,
		repository.tenantID,
		installment.SalePublicID,
		uuid.MustParse(installment.PublicID),
		installment.SequenceNumber,
		installment.AmountCents,
		installment.DueDate,
		installment.Status,
	)

	saved, err := scanInstallment(row)
	if err != nil {
		return installment
	}

	return saved
}

func (repository *PostgresInstallmentRepository) Update(installment entity.Installment) entity.Installment {
	row := repository.database.QueryRow(
		`
      UPDATE sales.installments AS installment
      SET status = $3
      FROM sales.sales AS sale
      WHERE installment.tenant_id = $1
        AND installment.public_id = $2::uuid
        AND sale.id = installment.sale_id
      RETURNING installment.public_id::text, sale.public_id::text, installment.sequence_number, installment.amount_cents, to_char(installment.due_date, 'YYYY-MM-DD'), installment.status
    `,
		repository.tenantID,
		installment.PublicID,
		installment.Status,
	)

	updated, err := scanInstallment(row)
	if err != nil {
		return installment
	}

	return updated
}

func (repository *PostgresCommissionRepository) ListBySalePublicID(salePublicID string) []entity.Commission {
	rows, err := repository.database.Query(
		`
      SELECT commission.public_id::text, sale.public_id::text, commission.recipient_user_public_id::text, commission.role_code, commission.rate_bps, commission.amount_cents, commission.status
      FROM sales.commissions AS commission
      INNER JOIN sales.sales AS sale
        ON sale.id = commission.sale_id
      WHERE commission.tenant_id = $1
        AND sale.public_id = $2::uuid
      ORDER BY commission.created_at, commission.id
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(salePublicID)),
	)
	if err != nil {
		return []entity.Commission{}
	}
	defer rows.Close()

	response := make([]entity.Commission, 0)
	for rows.Next() {
		commission, scanErr := scanCommission(rows)
		if scanErr == nil {
			response = append(response, commission)
		}
	}
	return response
}

func (repository *PostgresCommissionRepository) FindByPublicID(publicID string) *entity.Commission {
	row := repository.database.QueryRow(
		`
      SELECT commission.public_id::text, sale.public_id::text, commission.recipient_user_public_id::text, commission.role_code, commission.rate_bps, commission.amount_cents, commission.status
      FROM sales.commissions AS commission
      INNER JOIN sales.sales AS sale
        ON sale.id = commission.sale_id
      WHERE commission.tenant_id = $1
        AND commission.public_id = $2::uuid
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(publicID)),
	)

	commission, err := scanCommission(row)
	if err != nil {
		return nil
	}
	return &commission
}

func (repository *PostgresCommissionRepository) Save(commission entity.Commission) entity.Commission {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.commissions (tenant_id, sale_id, public_id, recipient_user_public_id, role_code, rate_bps, amount_cents, status)
      SELECT $1, sale.id, $3, $4::uuid, $5, $6, $7, $8
      FROM sales.sales AS sale
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2::uuid
      RETURNING $3::text, $2::text, recipient_user_public_id::text, role_code, rate_bps, amount_cents, status
    `,
		repository.tenantID,
		commission.SalePublicID,
		commission.PublicID,
		commission.RecipientUserID,
		commission.RoleCode,
		commission.RateBps,
		commission.AmountCents,
		commission.Status,
	)

	saved, err := scanCommission(row)
	if err != nil {
		return commission
	}
	return saved
}

func (repository *PostgresCommissionRepository) Update(commission entity.Commission) entity.Commission {
	row := repository.database.QueryRow(
		`
      UPDATE sales.commissions
      SET status = $3
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING public_id::text, $4::text, recipient_user_public_id::text, role_code, rate_bps, amount_cents, status
    `,
		repository.tenantID,
		commission.PublicID,
		commission.Status,
		commission.SalePublicID,
	)

	updated, err := scanCommission(row)
	if err != nil {
		return commission
	}
	return updated
}

func (repository *PostgresPendingItemRepository) ListBySalePublicID(salePublicID string) []entity.PendingItem {
	rows, err := repository.database.Query(
		`
      SELECT item.public_id::text, sale.public_id::text, item.code, item.summary, item.status, COALESCE(to_char(item.resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.pending_items AS item
      INNER JOIN sales.sales AS sale
        ON sale.id = item.sale_id
      WHERE item.tenant_id = $1
        AND sale.public_id = $2::uuid
      ORDER BY item.created_at, item.id
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(salePublicID)),
	)
	if err != nil {
		return []entity.PendingItem{}
	}
	defer rows.Close()

	response := make([]entity.PendingItem, 0)
	for rows.Next() {
		item, scanErr := scanPendingItem(rows)
		if scanErr == nil {
			response = append(response, item)
		}
	}
	return response
}

func (repository *PostgresPendingItemRepository) FindByPublicID(publicID string) *entity.PendingItem {
	row := repository.database.QueryRow(
		`
      SELECT item.public_id::text, sale.public_id::text, item.code, item.summary, item.status, COALESCE(to_char(item.resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.pending_items AS item
      INNER JOIN sales.sales AS sale
        ON sale.id = item.sale_id
      WHERE item.tenant_id = $1
        AND item.public_id = $2::uuid
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(publicID)),
	)

	item, err := scanPendingItem(row)
	if err != nil {
		return nil
	}
	return &item
}

func (repository *PostgresPendingItemRepository) Save(item entity.PendingItem) entity.PendingItem {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.pending_items (tenant_id, sale_id, public_id, code, summary, status, resolved_at)
      SELECT $1, sale.id, $3, $4, $5, $6, NULLIF($7, '')::timestamptz
      FROM sales.sales AS sale
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2::uuid
      RETURNING $3::text, $2::text, code, summary, status, COALESCE(to_char(resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		item.SalePublicID,
		item.PublicID,
		item.Code,
		item.Summary,
		item.Status,
		item.ResolvedAt,
	)

	saved, err := scanPendingItem(row)
	if err != nil {
		return item
	}
	return saved
}

func (repository *PostgresPendingItemRepository) Update(item entity.PendingItem) entity.PendingItem {
	row := repository.database.QueryRow(
		`
      UPDATE sales.pending_items
      SET status = $3,
          resolved_at = NULLIF($4, '')::timestamptz
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING public_id::text, $5::text, code, summary, status, COALESCE(to_char(resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		item.PublicID,
		item.Status,
		item.ResolvedAt,
		item.SalePublicID,
	)

	updated, err := scanPendingItem(row)
	if err != nil {
		return item
	}
	return updated
}

func (repository *PostgresRenegotiationRepository) ListBySalePublicID(salePublicID string) []entity.Renegotiation {
	rows, err := repository.database.Query(
		`
      SELECT renegotiation.public_id::text, sale.public_id::text, renegotiation.reason, renegotiation.previous_amount_cents, renegotiation.new_amount_cents, renegotiation.status, to_char(renegotiation.applied_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM sales.renegotiations AS renegotiation
      INNER JOIN sales.sales AS sale
        ON sale.id = renegotiation.sale_id
      WHERE renegotiation.tenant_id = $1
        AND sale.public_id = $2::uuid
      ORDER BY renegotiation.created_at, renegotiation.id
    `,
		repository.tenantID,
		uuid.MustParse(strings.TrimSpace(salePublicID)),
	)
	if err != nil {
		return []entity.Renegotiation{}
	}
	defer rows.Close()

	response := make([]entity.Renegotiation, 0)
	for rows.Next() {
		renegotiation, scanErr := scanRenegotiation(rows)
		if scanErr == nil {
			response = append(response, renegotiation)
		}
	}
	return response
}

func (repository *PostgresRenegotiationRepository) Save(renegotiation entity.Renegotiation) entity.Renegotiation {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.renegotiations (tenant_id, sale_id, public_id, reason, previous_amount_cents, new_amount_cents, status, applied_at)
      SELECT $1, sale.id, $3, $4, $5, $6, $7, $8::timestamptz
      FROM sales.sales AS sale
      WHERE sale.tenant_id = $1
        AND sale.public_id = $2::uuid
      RETURNING $3::text, $2::text, reason, previous_amount_cents, new_amount_cents, status, to_char(applied_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
    `,
		repository.tenantID,
		renegotiation.SalePublicID,
		renegotiation.PublicID,
		renegotiation.Reason,
		renegotiation.PreviousAmountCents,
		renegotiation.NewAmountCents,
		renegotiation.Status,
		renegotiation.AppliedAt,
	)

	saved, err := scanRenegotiation(row)
	if err != nil {
		return renegotiation
	}
	return saved
}

func scanInstallment(scanner scanner) (entity.Installment, error) {
	var publicID string
	var salePublicID string
	var sequenceNumber int
	var amountCents int64
	var dueDate string
	var status string
	if err := scanner.Scan(&publicID, &salePublicID, &sequenceNumber, &amountCents, &dueDate, &status); err != nil {
		return entity.Installment{}, err
	}
	return entity.RestoreInstallment(publicID, salePublicID, sequenceNumber, amountCents, dueDate, status)
}

func scanCommission(scanner scanner) (entity.Commission, error) {
	var publicID string
	var salePublicID string
	var recipientUserID string
	var roleCode string
	var rateBps int
	var amountCents int64
	var status string
	if err := scanner.Scan(&publicID, &salePublicID, &recipientUserID, &roleCode, &rateBps, &amountCents, &status); err != nil {
		return entity.Commission{}, err
	}
	return entity.RestoreCommission(publicID, salePublicID, recipientUserID, roleCode, rateBps, amountCents, status)
}

func scanPendingItem(scanner scanner) (entity.PendingItem, error) {
	var publicID string
	var salePublicID string
	var code string
	var summary string
	var status string
	var resolvedAt string
	if err := scanner.Scan(&publicID, &salePublicID, &code, &summary, &status, &resolvedAt); err != nil {
		return entity.PendingItem{}, err
	}
	return entity.RestorePendingItem(publicID, salePublicID, code, summary, status, resolvedAt)
}

func scanRenegotiation(scanner scanner) (entity.Renegotiation, error) {
	var publicID string
	var salePublicID string
	var reason string
	var previousAmountCents int64
	var newAmountCents int64
	var status string
	var appliedAt string
	if err := scanner.Scan(&publicID, &salePublicID, &reason, &previousAmountCents, &newAmountCents, &status, &appliedAt); err != nil {
		return entity.Renegotiation{}, err
	}
	return entity.RestoreRenegotiation(publicID, salePublicID, reason, previousAmountCents, newAmountCents, status, appliedAt)
}
