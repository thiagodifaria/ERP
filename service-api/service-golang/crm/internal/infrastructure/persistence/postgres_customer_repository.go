package persistence

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type PostgresCustomerRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresCustomerRepository(database *sql.DB, tenantSlug string) (*PostgresCustomerRepository, error) {
	var tenantID int64
	if err := database.QueryRow(
		`
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
    `,
		strings.TrimSpace(tenantSlug),
	).Scan(&tenantID); err != nil {
		return nil, err
	}

	return &PostgresCustomerRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresCustomerRepository) List() []entity.Customer {
	rows, err := repository.database.Query(
		`
      SELECT customer.public_id::text, lead.public_id::text, customer.name, customer.email::text, customer.source, customer.status, COALESCE(customer.owner_user_public_id::text, '')
      FROM crm.customers AS customer
      INNER JOIN crm.leads AS lead
        ON lead.id = customer.lead_id
      WHERE customer.tenant_id = $1
      ORDER BY customer.created_at, customer.id
    `,
		repository.tenantID,
	)
	if err != nil {
		return []entity.Customer{}
	}
	defer rows.Close()

	response := make([]entity.Customer, 0)
	for rows.Next() {
		customer, scanErr := scanCustomer(rows)
		if scanErr == nil {
			response = append(response, customer)
		}
	}

	return response
}

func (repository *PostgresCustomerRepository) FindByPublicID(publicID string) *entity.Customer {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT customer.public_id::text, lead.public_id::text, customer.name, customer.email::text, customer.source, customer.status, COALESCE(customer.owner_user_public_id::text, '')
      FROM crm.customers AS customer
      INNER JOIN crm.leads AS lead
        ON lead.id = customer.lead_id
      WHERE customer.tenant_id = $1
        AND customer.public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	customer, err := scanCustomer(row)
	if err != nil {
		return nil
	}

	return &customer
}

func (repository *PostgresCustomerRepository) FindByEmail(email string) *entity.Customer {
	row := repository.database.QueryRow(
		`
      SELECT customer.public_id::text, lead.public_id::text, customer.name, customer.email::text, customer.source, customer.status, COALESCE(customer.owner_user_public_id::text, '')
      FROM crm.customers AS customer
      INNER JOIN crm.leads AS lead
        ON lead.id = customer.lead_id
      WHERE customer.tenant_id = $1
        AND customer.email = $2
    `,
		repository.tenantID,
		strings.ToLower(strings.TrimSpace(email)),
	)

	customer, err := scanCustomer(row)
	if err != nil {
		return nil
	}

	return &customer
}

func (repository *PostgresCustomerRepository) Save(customer entity.Customer) entity.Customer {
	publicID := uuid.MustParse(customer.PublicID)
	leadPublicID := uuid.MustParse(customer.LeadPublicID)
	var ownerUserID *uuid.UUID
	if strings.TrimSpace(customer.OwnerUserID) != "" {
		parsed := uuid.MustParse(customer.OwnerUserID)
		ownerUserID = &parsed
	}

	row := repository.database.QueryRow(
		`
      INSERT INTO crm.customers (tenant_id, lead_id, public_id, owner_user_public_id, name, email, source, status)
      SELECT
        $1,
        lead.id,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
      FROM crm.leads AS lead
      WHERE lead.tenant_id = $1
        AND lead.public_id = $2
      RETURNING
        public_id::text,
        $2::text,
        name,
        email::text,
        source,
        status,
        COALESCE(owner_user_public_id::text, '')
    `,
		repository.tenantID,
		leadPublicID,
		publicID,
		ownerUserID,
		customer.Name,
		customer.Email,
		customer.Source,
		customer.Status,
	)

	saved, err := scanCustomer(row)
	if err != nil {
		return customer
	}

	return saved
}

type customerScanner interface {
	Scan(dest ...any) error
}

func scanCustomer(scanner customerScanner) (entity.Customer, error) {
	var publicID string
	var leadPublicID string
	var name string
	var email string
	var source string
	var status string
	var ownerUserID string

	if err := scanner.Scan(&publicID, &leadPublicID, &name, &email, &source, &status, &ownerUserID); err != nil {
		return entity.Customer{}, err
	}

	return entity.RestoreCustomer(publicID, leadPublicID, name, email, source, status, ownerUserID)
}
