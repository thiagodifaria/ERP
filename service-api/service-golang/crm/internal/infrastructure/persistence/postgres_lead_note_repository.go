// PostgresLeadNoteRepository persists CRM lead notes in PostgreSQL for one bootstrap tenant.
package persistence

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type PostgresLeadNoteRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresLeadNoteRepository(database *sql.DB, tenantSlug string) (*PostgresLeadNoteRepository, error) {
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

	return &PostgresLeadNoteRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresLeadNoteRepository) ListByLeadPublicID(leadPublicID string) []entity.LeadNote {
	parsedLeadPublicID, err := uuid.Parse(strings.TrimSpace(leadPublicID))
	if err != nil {
		return []entity.LeadNote{}
	}

	rows, err := repository.database.Query(
		`
      SELECT note.public_id::text, lead.public_id::text, note.body, note.category, note.created_at
      FROM crm.lead_notes AS note
      INNER JOIN crm.leads AS lead
        ON lead.id = note.lead_id
      WHERE note.tenant_id = $1
        AND lead.public_id = $2
      ORDER BY note.created_at, note.id
    `,
		repository.tenantID,
		parsedLeadPublicID,
	)
	if err != nil {
		return []entity.LeadNote{}
	}
	defer rows.Close()

	notes := make([]entity.LeadNote, 0)
	for rows.Next() {
		note, scanErr := scanLeadNote(rows)
		if scanErr == nil {
			notes = append(notes, note)
		}
	}

	return notes
}

func (repository *PostgresLeadNoteRepository) Save(note entity.LeadNote) entity.LeadNote {
	publicID := uuid.MustParse(note.PublicID)
	leadPublicID := uuid.MustParse(note.LeadPublicID)

	row := repository.database.QueryRow(
		`
      INSERT INTO crm.lead_notes (tenant_id, lead_id, public_id, category, body, created_at)
      SELECT $1, lead.id, $2, $4, $5, $6
      FROM crm.leads AS lead
      WHERE lead.tenant_id = $1
        AND lead.public_id = $3
      RETURNING
        public_id::text,
        $3::text,
        body,
        category,
        created_at
    `,
		repository.tenantID,
		publicID,
		leadPublicID,
		note.Category,
		note.Body,
		note.CreatedAt,
	)

	savedNote, err := scanLeadNote(row)
	if err != nil {
		return note
	}

	return savedNote
}

type leadNoteScanner interface {
	Scan(dest ...any) error
}

func scanLeadNote(scanner leadNoteScanner) (entity.LeadNote, error) {
	var publicID string
	var leadPublicID string
	var body string
	var category string
	var createdAt time.Time

	if err := scanner.Scan(&publicID, &leadPublicID, &body, &category, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.LeadNote{}, err
		}

		return entity.LeadNote{}, err
	}

	return entity.NewLeadNote(publicID, leadPublicID, body, category, createdAt)
}
