// PostgresLeadRepository persists CRM leads in PostgreSQL for one bootstrap tenant.
package persistence

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type PostgresLeadRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresLeadRepository(database *sql.DB, tenantSlug string) (*PostgresLeadRepository, error) {
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

	return &PostgresLeadRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresLeadRepository) List() []entity.Lead {
	rows, err := repository.database.Query(
		`
      SELECT public_id::text, name, email::text, source, status, COALESCE(owner_user_public_id::text, '')
      FROM crm.leads
      WHERE tenant_id = $1
      ORDER BY created_at, id
    `,
		repository.tenantID,
	)
	if err != nil {
		return []entity.Lead{}
	}
	defer rows.Close()

	leads := make([]entity.Lead, 0)
	for rows.Next() {
		lead, scanErr := scanLead(rows)
		if scanErr == nil {
			leads = append(leads, lead)
		}
	}

	return leads
}

func (repository *PostgresLeadRepository) FindByPublicID(publicID string) *entity.Lead {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, name, email::text, source, status, COALESCE(owner_user_public_id::text, '')
      FROM crm.leads
      WHERE tenant_id = $1
        AND public_id = $2
    `,
		repository.tenantID,
		parsedPublicID,
	)

	lead, err := scanLead(row)
	if err != nil {
		return nil
	}

	return &lead
}

func (repository *PostgresLeadRepository) FindByEmail(email string) *entity.Lead {
	row := repository.database.QueryRow(
		`
      SELECT public_id::text, name, email::text, source, status, COALESCE(owner_user_public_id::text, '')
      FROM crm.leads
      WHERE tenant_id = $1
        AND email = $2
    `,
		repository.tenantID,
		strings.ToLower(strings.TrimSpace(email)),
	)

	lead, err := scanLead(row)
	if err != nil {
		return nil
	}

	return &lead
}

func (repository *PostgresLeadRepository) Save(lead entity.Lead) entity.Lead {
	publicID, ownerUserID := parseLeadIdentifiers(lead)

	row := repository.database.QueryRow(
		`
      INSERT INTO crm.leads (tenant_id, public_id, owner_user_public_id, name, email, source, status)
      VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING public_id::text, name, email::text, source, status, COALESCE(owner_user_public_id::text, '')
    `,
		repository.tenantID,
		publicID,
		ownerUserID,
		lead.Name,
		lead.Email,
		lead.Source,
		lead.Status,
	)

	savedLead, err := scanLead(row)
	if err != nil {
		return lead
	}

	return savedLead
}

func (repository *PostgresLeadRepository) Update(lead entity.Lead) entity.Lead {
	publicID, ownerUserID := parseLeadIdentifiers(lead)

	row := repository.database.QueryRow(
		`
      UPDATE crm.leads
      SET
        owner_user_public_id = $3,
        name = $4,
        email = $5,
        source = $6,
        status = $7
      WHERE tenant_id = $1
        AND public_id = $2
      RETURNING public_id::text, name, email::text, source, status, COALESCE(owner_user_public_id::text, '')
    `,
		repository.tenantID,
		publicID,
		ownerUserID,
		lead.Name,
		lead.Email,
		lead.Source,
		lead.Status,
	)

	updatedLead, err := scanLead(row)
	if err != nil {
		return lead
	}

	return updatedLead
}

func parseLeadIdentifiers(lead entity.Lead) (uuid.UUID, *uuid.UUID) {
	publicID := uuid.MustParse(lead.PublicID)

	if strings.TrimSpace(lead.OwnerUserID) == "" {
		return publicID, nil
	}

	ownerUserID := uuid.MustParse(lead.OwnerUserID)
	return publicID, &ownerUserID
}

type leadScanner interface {
	Scan(dest ...any) error
}

func scanLead(scanner leadScanner) (entity.Lead, error) {
	var publicID string
	var name string
	var email string
	var source string
	var status string
	var ownerUserID string

	if err := scanner.Scan(&publicID, &name, &email, &source, &status, &ownerUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Lead{}, err
		}

		return entity.Lead{}, err
	}

	lead, err := entity.NewLead(publicID, name, email, source, ownerUserID)
	if err != nil {
		return entity.Lead{}, err
	}

	return lead.TransitionTo(status)
}
