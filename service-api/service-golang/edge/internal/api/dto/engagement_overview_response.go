// DTO do cockpit de engagement entregue pelo edge.
// O gateway cruza operacao de campanhas, templates e entregas sem replicar logica analitica.
package dto

type EngagementExecutiveSummary struct {
	Status               string  `json:"status"`
	Campaigns            int     `json:"campaigns"`
	ActiveCampaigns      int     `json:"activeCampaigns"`
	Templates            int     `json:"templates"`
	Deliveries           int     `json:"deliveries"`
	DeliveredDeliveries  int     `json:"deliveredDeliveries"`
	FailedDeliveries     int     `json:"failedDeliveries"`
	ConvertedTouchpoints int     `json:"convertedTouchpoints"`
	BusinessLinked       int     `json:"businessLinked"`
	ProviderLinkedEvents int     `json:"providerLinkedEvents"`
	DeliveryRate         float64 `json:"deliveryRate"`
}

type EngagementOverviewResponse struct {
	Service              string                     `json:"service"`
	TenantSlug           string                     `json:"tenantSlug"`
	GeneratedAt          string                     `json:"generatedAt"`
	ExecutiveSummary     EngagementExecutiveSummary `json:"executiveSummary"`
	ServicePulse         map[string]any             `json:"servicePulse"`
	Tenant360            map[string]any             `json:"tenant360"`
	EngagementOperations map[string]any             `json:"engagementOperations"`
}
