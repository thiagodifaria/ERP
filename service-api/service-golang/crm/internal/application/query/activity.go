package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ListRelationshipEventsByAggregate struct {
	eventRepository repository.RelationshipEventRepository
}

func NewListRelationshipEventsByAggregate(eventRepository repository.RelationshipEventRepository) ListRelationshipEventsByAggregate {
	return ListRelationshipEventsByAggregate{eventRepository: eventRepository}
}

func (useCase ListRelationshipEventsByAggregate) Execute(aggregateType string, aggregatePublicID string) []entity.RelationshipEvent {
	if useCase.eventRepository == nil {
		return []entity.RelationshipEvent{}
	}

	return useCase.eventRepository.ListByAggregate(aggregateType, aggregatePublicID)
}

type ListPendingOutboxEvents struct {
	outboxRepository repository.OutboxEventRepository
}

func NewListPendingOutboxEvents(outboxRepository repository.OutboxEventRepository) ListPendingOutboxEvents {
	return ListPendingOutboxEvents{outboxRepository: outboxRepository}
}

func (useCase ListPendingOutboxEvents) Execute(limit int) []entity.OutboxEvent {
	if useCase.outboxRepository == nil {
		return []entity.OutboxEvent{}
	}

	return useCase.outboxRepository.ListPending(limit)
}
