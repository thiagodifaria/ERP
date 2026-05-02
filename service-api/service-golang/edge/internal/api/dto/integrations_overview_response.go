package dto

type IntegrationsExecutiveSummary struct {
	Status                  string `json:"status"`
	ConfiguredProviders     int    `json:"configuredProviders"`
	ActiveInboundProviders  int    `json:"activeInboundProviders"`
	ActiveOutboundProviders int    `json:"activeOutboundProviders"`
	InboundLeads            int    `json:"inboundLeads"`
	WorkflowDispatches      int    `json:"workflowDispatches"`
	FailedProviderEvents    int    `json:"failedProviderEvents"`
	DeadLetterEvents        int    `json:"deadLetterEvents"`
	OpenProviderRisks       int    `json:"openProviderRisks"`
}

type IntegrationsOverviewResponse struct {
	Service              string                       `json:"service"`
	TenantSlug           string                       `json:"tenantSlug"`
	GeneratedAt          string                       `json:"generatedAt"`
	ExecutiveSummary     IntegrationsExecutiveSummary `json:"executiveSummary"`
	ServicePulse         map[string]any               `json:"servicePulse"`
	EngagementOperations map[string]any               `json:"engagementOperations"`
	IntegrationReadiness map[string]any               `json:"integrationReadiness"`
}
