package dto

type OutboxEventResponse struct {
	PublicID          string `json:"publicId"`
	AggregateType     string `json:"aggregateType"`
	AggregatePublicID string `json:"aggregatePublicId"`
	EventType         string `json:"eventType"`
	Payload           string `json:"payload"`
	Status            string `json:"status"`
	CreatedAt         string `json:"createdAt"`
	ProcessedAt       string `json:"processedAt"`
}
