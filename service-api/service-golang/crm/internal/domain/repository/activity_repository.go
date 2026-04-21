package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"

type RelationshipEventRepository interface {
	ListByAggregate(aggregateType string, aggregatePublicID string) []entity.RelationshipEvent
	Save(event entity.RelationshipEvent) entity.RelationshipEvent
}

type OutboxEventRepository interface {
	ListPending(limit int) []entity.OutboxEvent
	Save(event entity.OutboxEvent) entity.OutboxEvent
}
