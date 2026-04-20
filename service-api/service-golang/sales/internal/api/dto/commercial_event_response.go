package dto

type CommercialEventResponse struct {
	PublicID          string `json:"publicId"`
	AggregateType     string `json:"aggregateType"`
	AggregatePublicID string `json:"aggregatePublicId"`
	EventCode         string `json:"eventCode"`
	Actor             string `json:"actor"`
	Summary           string `json:"summary"`
	CreatedAt         string `json:"createdAt"`
}
