package dto

type SaaSExecutiveSummary struct {
	Status              string `json:"status"`
	EntitlementsTotal   int    `json:"entitlementsTotal"`
	EnabledEntitlements int    `json:"enabledEntitlements"`
	ActiveQuotas        int    `json:"activeQuotas"`
	ActiveBlocks        int    `json:"activeBlocks"`
	TrackedMetrics      int    `json:"trackedMetrics"`
	QueuedJobs          int    `json:"queuedJobs"`
	CompletedJobs       int    `json:"completedJobs"`
}

type SaaSOverviewResponse struct {
	Service          string               `json:"service"`
	TenantSlug       string               `json:"tenantSlug"`
	GeneratedAt      string               `json:"generatedAt"`
	ExecutiveSummary SaaSExecutiveSummary `json:"executiveSummary"`
	ServicePulse     map[string]any       `json:"servicePulse"`
	Tenant360        map[string]any       `json:"tenant360"`
	SaaSControl      map[string]any       `json:"saasControl"`
}
