// DTO do cockpit comercial entregue pelo edge.
// O gateway agrega relatorios de analytics e expoe um resumo executivo unico.
package dto

type SalesExecutiveSummary struct {
	Status               string `json:"status"`
	LeadsCaptured        int    `json:"leadsCaptured"`
	Opportunities        int    `json:"opportunities"`
	Proposals            int    `json:"proposals"`
	SalesWon             int    `json:"salesWon"`
	BookedRevenueCents   int    `json:"bookedRevenueCents"`
	CompletedAutomations int    `json:"completedAutomations"`
}

type SalesOverviewResponse struct {
	Service          string                `json:"service"`
	TenantSlug       string                `json:"tenantSlug"`
	GeneratedAt      string                `json:"generatedAt"`
	ExecutiveSummary SalesExecutiveSummary `json:"executiveSummary"`
	PipelineSummary  map[string]any        `json:"pipelineSummary"`
	ServicePulse     map[string]any        `json:"servicePulse"`
	SalesJourney     map[string]any        `json:"salesJourney"`
}
