package entity

import (
	"strings"
	"time"
)

type RelationshipEvent struct {
	PublicID          string
	AggregateType     string
	AggregatePublicID string
	EventCode         string
	Actor             string
	Summary           string
	CreatedAt         string
}

func NewRelationshipEvent(publicID string, aggregateType string, aggregatePublicID string, eventCode string, actor string, summary string, createdAt time.Time) RelationshipEvent {
	return RelationshipEvent{
		PublicID:          strings.TrimSpace(publicID),
		AggregateType:     strings.TrimSpace(aggregateType),
		AggregatePublicID: strings.TrimSpace(aggregatePublicID),
		EventCode:         strings.TrimSpace(eventCode),
		Actor:             strings.TrimSpace(actor),
		Summary:           strings.TrimSpace(summary),
		CreatedAt:         createdAt.UTC().Format(time.RFC3339),
	}
}
