package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type ListCommercialEventsByAggregate struct {
	eventRepository repository.CommercialEventRepository
}

type ListPendingOutboxEvents struct {
	outboxRepository repository.OutboxEventRepository
}

func NewListCommercialEventsByAggregate(eventRepository repository.CommercialEventRepository) ListCommercialEventsByAggregate {
	return ListCommercialEventsByAggregate{eventRepository: eventRepository}
}

func (useCase ListCommercialEventsByAggregate) Execute(aggregateType string, aggregatePublicID string) []entity.CommercialEvent {
	return useCase.eventRepository.ListByAggregate(aggregateType, aggregatePublicID)
}

func NewListPendingOutboxEvents(outboxRepository repository.OutboxEventRepository) ListPendingOutboxEvents {
	return ListPendingOutboxEvents{outboxRepository: outboxRepository}
}

func (useCase ListPendingOutboxEvents) Execute(limit int) []entity.OutboxEvent {
	return useCase.outboxRepository.ListPending(limit)
}
