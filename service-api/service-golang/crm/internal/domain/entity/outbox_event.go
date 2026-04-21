package entity

import (
	"strings"
	"time"
)

type OutboxEvent struct {
	PublicID          string
	AggregateType     string
	AggregatePublicID string
	EventType         string
	Payload           string
	Status            string
	CreatedAt         string
	ProcessedAt       string
}

func NewOutboxEvent(publicID string, aggregateType string, aggregatePublicID string, eventType string, payload string, createdAt time.Time) OutboxEvent {
	return OutboxEvent{
		PublicID:          strings.TrimSpace(publicID),
		AggregateType:     strings.TrimSpace(aggregateType),
		AggregatePublicID: strings.TrimSpace(aggregatePublicID),
		EventType:         strings.TrimSpace(eventType),
		Payload:           strings.TrimSpace(payload),
		Status:            "pending",
		CreatedAt:         createdAt.UTC().Format(time.RFC3339),
		ProcessedAt:       "",
	}
}

func RestoreOutboxEvent(publicID string, aggregateType string, aggregatePublicID string, eventType string, payload string, status string, createdAt string, processedAt string) OutboxEvent {
	return OutboxEvent{
		PublicID:          strings.TrimSpace(publicID),
		AggregateType:     strings.TrimSpace(aggregateType),
		AggregatePublicID: strings.TrimSpace(aggregatePublicID),
		EventType:         strings.TrimSpace(eventType),
		Payload:           strings.TrimSpace(payload),
		Status:            strings.TrimSpace(status),
		CreatedAt:         strings.TrimSpace(createdAt),
		ProcessedAt:       strings.TrimSpace(processedAt),
	}
}
