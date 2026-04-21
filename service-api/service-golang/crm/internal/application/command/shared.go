package command

import (
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

func newPublicID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return uuid.Nil.String()
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	return uuid.UUID(raw).String()
}

func recordRelationshipEvent(eventRepository repository.RelationshipEventRepository, aggregateType string, aggregatePublicID string, eventCode string, actor string, summary string) {
	if eventRepository == nil {
		return
	}

	eventRepository.Save(entity.NewRelationshipEvent(newPublicID(), aggregateType, aggregatePublicID, eventCode, actor, summary, time.Now().UTC()))
}

func appendOutboxEvent(outboxRepository repository.OutboxEventRepository, aggregateType string, aggregatePublicID string, eventType string, payload map[string]any) {
	if outboxRepository == nil {
		return
	}

	serializedPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}

	outboxRepository.Save(entity.NewOutboxEvent(newPublicID(), aggregateType, aggregatePublicID, eventType, string(serializedPayload), time.Now().UTC()))
}
