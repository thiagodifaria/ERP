package persistence

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type PostgresRelationshipEventRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresOutboxEventRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresRelationshipEventRepository(database *sql.DB, tenantSlug string) (*PostgresRelationshipEventRepository, error) {
	tenantID, err := lookupCrmTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresRelationshipEventRepository{database: database, tenantID: tenantID}, nil
}

func NewPostgresOutboxEventRepository(database *sql.DB, tenantSlug string) (*PostgresOutboxEventRepository, error) {
	tenantID, err := lookupCrmTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresOutboxEventRepository{database: database, tenantID: tenantID}, nil
}

func (repository *PostgresRelationshipEventRepository) ListByAggregate(aggregateType string, aggregatePublicID string) []entity.RelationshipEvent {
	parsedAggregatePublicID, err := uuid.Parse(strings.TrimSpace(aggregatePublicID))
	if err != nil {
		return []entity.RelationshipEvent{}
	}

	rows, err := repository.database.Query(
		`
      SELECT public_id::text, aggregate_type, aggregate_public_id::text, event_code, actor, summary, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM crm.relationship_events
      WHERE tenant_id = $1
        AND aggregate_type = $2
        AND aggregate_public_id = $3
      ORDER BY created_at, id
    `,
		repository.tenantID,
		strings.TrimSpace(aggregateType),
		parsedAggregatePublicID,
	)
	if err != nil {
		return []entity.RelationshipEvent{}
	}
	defer rows.Close()

	response := make([]entity.RelationshipEvent, 0)
	for rows.Next() {
		event, scanErr := scanRelationshipEvent(rows)
		if scanErr == nil {
			response = append(response, event)
		}
	}

	return response
}

func (repository *PostgresRelationshipEventRepository) Save(event entity.RelationshipEvent) entity.RelationshipEvent {
	row := repository.database.QueryRow(
		`
      INSERT INTO crm.relationship_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_code, actor, summary)
      VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING public_id::text, aggregate_type, aggregate_public_id::text, event_code, actor, summary, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
    `,
		repository.tenantID,
		uuid.MustParse(event.PublicID),
		event.AggregateType,
		uuid.MustParse(event.AggregatePublicID),
		event.EventCode,
		event.Actor,
		event.Summary,
	)

	saved, err := scanRelationshipEvent(row)
	if err != nil {
		return event
	}

	return saved
}

func (repository *PostgresOutboxEventRepository) ListPending(limit int) []entity.OutboxEvent {
	if limit <= 0 {
		limit = 100
	}

	rows, err := repository.database.Query(
		`
      SELECT public_id::text, aggregate_type, aggregate_public_id::text, event_type, payload, status, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), COALESCE(to_char(processed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM crm.outbox_events
      WHERE tenant_id = $1
        AND status = 'pending'
      ORDER BY created_at, id
      LIMIT $2
    `,
		repository.tenantID,
		limit,
	)
	if err != nil {
		return []entity.OutboxEvent{}
	}
	defer rows.Close()

	response := make([]entity.OutboxEvent, 0)
	for rows.Next() {
		event, scanErr := scanCrmOutboxEvent(rows)
		if scanErr == nil {
			response = append(response, event)
		}
	}

	return response
}

func (repository *PostgresOutboxEventRepository) Save(event entity.OutboxEvent) entity.OutboxEvent {
	row := repository.database.QueryRow(
		`
      INSERT INTO crm.outbox_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_type, payload, status, processed_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '')::timestamptz)
      RETURNING public_id::text, aggregate_type, aggregate_public_id::text, event_type, payload, status, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), COALESCE(to_char(processed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		uuid.MustParse(event.PublicID),
		event.AggregateType,
		uuid.MustParse(event.AggregatePublicID),
		event.EventType,
		event.Payload,
		event.Status,
		event.ProcessedAt,
	)

	saved, err := scanCrmOutboxEvent(row)
	if err != nil {
		return event
	}

	return saved
}

type activityScanner interface {
	Scan(dest ...any) error
}

func lookupCrmTenantID(database *sql.DB, tenantSlug string) (int64, error) {
	var tenantID int64
	err := database.QueryRow(
		`
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
    `,
		strings.TrimSpace(tenantSlug),
	).Scan(&tenantID)
	return tenantID, err
}

func scanRelationshipEvent(scanner activityScanner) (entity.RelationshipEvent, error) {
	var publicID string
	var aggregateType string
	var aggregatePublicID string
	var eventCode string
	var actor string
	var summary string
	var createdAt string

	if err := scanner.Scan(&publicID, &aggregateType, &aggregatePublicID, &eventCode, &actor, &summary, &createdAt); err != nil {
		return entity.RelationshipEvent{}, err
	}

	return entity.RelationshipEvent{
		PublicID:          publicID,
		AggregateType:     aggregateType,
		AggregatePublicID: aggregatePublicID,
		EventCode:         eventCode,
		Actor:             actor,
		Summary:           summary,
		CreatedAt:         createdAt,
	}, nil
}

func scanCrmOutboxEvent(scanner activityScanner) (entity.OutboxEvent, error) {
	var publicID string
	var aggregateType string
	var aggregatePublicID string
	var eventType string
	var payload string
	var status string
	var createdAt string
	var processedAt string

	if err := scanner.Scan(&publicID, &aggregateType, &aggregatePublicID, &eventType, &payload, &status, &createdAt, &processedAt); err != nil {
		return entity.OutboxEvent{}, err
	}

	return entity.RestoreOutboxEvent(publicID, aggregateType, aggregatePublicID, eventType, payload, status, createdAt, processedAt), nil
}
