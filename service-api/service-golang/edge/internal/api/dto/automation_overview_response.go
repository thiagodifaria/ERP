// DTO do cockpit de automacao entregue pelo edge.
// O gateway agrega relatorios de analytics e expoe um resumo executivo unico.
package dto

type AutomationExecutiveSummary struct {
	Status                     string `json:"status"`
	ActiveDefinitions          int    `json:"activeDefinitions"`
	StableDefinitions          int    `json:"stableDefinitions"`
	AttentionDefinitions       int    `json:"attentionDefinitions"`
	CriticalDefinitions        int    `json:"criticalDefinitions"`
	RunningControlRuns         int    `json:"runningControlRuns"`
	CompletedRuntimeExecutions int    `json:"completedRuntimeExecutions"`
	ForwardedWebhookEvents     int    `json:"forwardedWebhookEvents"`
}

type AutomationOverviewResponse struct {
	Service                  string                     `json:"service"`
	TenantSlug               string                     `json:"tenantSlug"`
	GeneratedAt              string                     `json:"generatedAt"`
	ExecutiveSummary         AutomationExecutiveSummary `json:"executiveSummary"`
	ServicePulse             map[string]any             `json:"servicePulse"`
	AutomationBoard          map[string]any             `json:"automationBoard"`
	WorkflowDefinitionHealth map[string]any             `json:"workflowDefinitionHealth"`
}
