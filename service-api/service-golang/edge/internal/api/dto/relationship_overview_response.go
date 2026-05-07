package dto

type RelationshipExecutiveSummary struct {
	Status               string `json:"status"`
	AverageLeadScore     int    `json:"averageLeadScore"`
	HotLeads             int    `json:"hotLeads"`
	PipelineConfigs      int    `json:"pipelineConfigs"`
	PipelineStages       int    `json:"pipelineStages"`
	OpenSupportCases     int    `json:"openSupportCases"`
	OverdueSupportCases  int    `json:"overdueSupportCases"`
	WeightedPipelineCents int   `json:"weightedPipelineCents"`
	BookedRevenueCents   int    `json:"bookedRevenueCents"`
}

type RelationshipOverviewResponse struct {
	Service                  string                      `json:"service"`
	TenantSlug               string                      `json:"tenantSlug"`
	GeneratedAt              string                      `json:"generatedAt"`
	ExecutiveSummary         RelationshipExecutiveSummary `json:"executiveSummary"`
	PipelineSummary          map[string]any              `json:"pipelineSummary"`
	RelationshipIntelligence map[string]any              `json:"relationshipIntelligence"`
	Tenant360                map[string]any              `json:"tenant360"`
}
