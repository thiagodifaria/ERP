package persistence

import (
	"slices"
	"sync"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type InMemoryRelationshipEventRepository struct {
	sync.Mutex
	events []entity.RelationshipEvent
}

type InMemoryOutboxEventRepository struct {
	sync.Mutex
	events []entity.OutboxEvent
}

func NewInMemoryRelationshipEventRepository(_ ...string) *InMemoryRelationshipEventRepository {
	return &InMemoryRelationshipEventRepository{
		events: []entity.RelationshipEvent{
			entity.NewRelationshipEvent(
				"0195e7a0-7a9c-7c1f-8a44-4a6e70000091",
				"lead",
				BootstrapLeadPublicID,
				"lead_created",
				"crm",
				"Bootstrap lead loaded for CRM runtime.",
				time.Date(2026, time.March, 31, 13, 58, 0, 0, time.UTC),
			),
		},
	}
}

func (repository *InMemoryRelationshipEventRepository) ListByAggregate(aggregateType string, aggregatePublicID string) []entity.RelationshipEvent {
	repository.Lock()
	defer repository.Unlock()

	response := make([]entity.RelationshipEvent, 0)
	for _, event := range repository.events {
		if event.AggregateType == aggregateType && event.AggregatePublicID == aggregatePublicID {
			response = append(response, event)
		}
	}

	return slices.Clone(response)
}

func (repository *InMemoryRelationshipEventRepository) Save(event entity.RelationshipEvent) entity.RelationshipEvent {
	repository.Lock()
	defer repository.Unlock()

	repository.events = append(repository.events, event)
	return event
}

func NewInMemoryOutboxEventRepository(_ ...string) *InMemoryOutboxEventRepository {
	return &InMemoryOutboxEventRepository{
		events: []entity.OutboxEvent{},
	}
}

func (repository *InMemoryOutboxEventRepository) ListPending(limit int) []entity.OutboxEvent {
	repository.Lock()
	defer repository.Unlock()

	if limit <= 0 {
		limit = 100
	}

	response := make([]entity.OutboxEvent, 0, limit)
	for _, event := range repository.events {
		if event.Status != "pending" {
			continue
		}

		response = append(response, event)
		if len(response) == limit {
			break
		}
	}

	return slices.Clone(response)
}

func (repository *InMemoryOutboxEventRepository) Save(event entity.OutboxEvent) entity.OutboxEvent {
	repository.Lock()
	defer repository.Unlock()

	repository.events = append(repository.events, event)
	return event
}
