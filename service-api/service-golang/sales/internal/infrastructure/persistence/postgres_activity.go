package persistence

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

type PostgresCommercialEventRepository struct {
	database *sql.DB
	tenantID int64
}

type PostgresOutboxEventRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresCommercialEventRepository(database *sql.DB, tenantSlug string) (*PostgresCommercialEventRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresCommercialEventRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func NewPostgresOutboxEventRepository(database *sql.DB, tenantSlug string) (*PostgresOutboxEventRepository, error) {
	tenantID, err := lookupTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresOutboxEventRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresCommercialEventRepository) ListByAggregate(aggregateType string, aggregatePublicID string) []entity.CommercialEvent {
	parsedAggregatePublicID, err := uuid.Parse(strings.TrimSpace(aggregatePublicID))
	if err != nil {
		return []entity.CommercialEvent{}
	}

	rows, err := repository.database.Query(
		`
      SELECT public_id::text, aggregate_type, aggregate_public_id::text, event_code, actor, summary, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM sales.commercial_events
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
		return []entity.CommercialEvent{}
	}
	defer rows.Close()

	response := make([]entity.CommercialEvent, 0)
	for rows.Next() {
		event, scanErr := scanCommercialEvent(rows)
		if scanErr == nil {
			response = append(response, event)
		}
	}

	return response
}

func (repository *PostgresCommercialEventRepository) Save(event entity.CommercialEvent) entity.CommercialEvent {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.commercial_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_code, actor, summary)
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

	saved, err := scanCommercialEvent(row)
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
      FROM sales.outbox_events
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
		event, scanErr := scanOutboxEvent(rows)
		if scanErr == nil {
			response = append(response, event)
		}
	}

	return response
}

func (repository *PostgresOutboxEventRepository) FindByPublicID(publicID string) *entity.OutboxEvent {
	parsedPublicID, err := uuid.Parse(strings.TrimSpace(publicID))
	if err != nil {
		return nil
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, aggregate_type, aggregate_public_id::text, event_type, payload, status, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), COALESCE(to_char(processed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
      FROM sales.outbox_events
      WHERE tenant_id = $1
        AND public_id = $2
      LIMIT 1
    `,
		repository.tenantID,
		parsedPublicID,
	)

	event, err := scanOutboxEvent(row)
	if err != nil {
		return nil
	}

	return &event
}

func (repository *PostgresOutboxEventRepository) Save(event entity.OutboxEvent) entity.OutboxEvent {
	row := repository.database.QueryRow(
		`
      INSERT INTO sales.outbox_events (tenant_id, public_id, aggregate_type, aggregate_public_id, event_type, payload, status, processed_at)
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

	saved, err := scanOutboxEvent(row)
	if err != nil {
		return event
	}

	return saved
}

func (repository *PostgresOutboxEventRepository) Update(event entity.OutboxEvent) entity.OutboxEvent {
	row := repository.database.QueryRow(
		`
      UPDATE sales.outbox_events
      SET status = $3,
          processed_at = NULLIF($4, '')::timestamptz
      WHERE tenant_id = $1
        AND public_id = $2
      RETURNING public_id::text, aggregate_type, aggregate_public_id::text, event_type, payload, status, to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), COALESCE(to_char(processed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
    `,
		repository.tenantID,
		uuid.MustParse(event.PublicID),
		event.Status,
		event.ProcessedAt,
	)

	updated, err := scanOutboxEvent(row)
	if err != nil {
		return event
	}

	return updated
}

func scanCommercialEvent(scanner scanner) (entity.CommercialEvent, error) {
	var publicID string
	var aggregateType string
	var aggregatePublicID string
	var eventCode string
	var actor string
	var summary string
	var createdAt string

	if err := scanner.Scan(&publicID, &aggregateType, &aggregatePublicID, &eventCode, &actor, &summary, &createdAt); err != nil {
		return entity.CommercialEvent{}, err
	}

	return entity.CommercialEvent{
		PublicID:          publicID,
		AggregateType:     aggregateType,
		AggregatePublicID: aggregatePublicID,
		EventCode:         eventCode,
		Actor:             actor,
		Summary:           summary,
		CreatedAt:         createdAt,
	}, nil
}

func scanOutboxEvent(scanner scanner) (entity.OutboxEvent, error) {
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
