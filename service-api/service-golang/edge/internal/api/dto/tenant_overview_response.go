// DTO do overview de tenant entregue pelo edge.
// O gateway devolve os blocos analiticos sem reformatar o shape de cada relatorio.
package dto

type TenantOverviewResponse struct {
	Service         string         `json:"service"`
	TenantSlug      string         `json:"tenantSlug"`
	GeneratedAt     string         `json:"generatedAt"`
	PipelineSummary map[string]any `json:"pipelineSummary"`
	ServicePulse    map[string]any `json:"servicePulse"`
	Tenant360       map[string]any `json:"tenant360"`
	AutomationBoard map[string]any `json:"automationBoard"`
}
