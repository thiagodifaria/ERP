package dto

type CoreOperationsExecutiveSummary struct {
	Status                string `json:"status"`
	CatalogItems          int    `json:"catalogItems"`
	Suppliers             int    `json:"suppliers"`
	SupportCases          int    `json:"supportCases"`
	OverdueSupportCases   int    `json:"overdueSupportCases"`
	UnreadNotifications   int    `json:"unreadNotifications"`
	CriticalNotifications int    `json:"criticalNotifications"`
}

type CoreOperationsOverviewResponse struct {
	Service          string                         `json:"service"`
	TenantSlug       string                         `json:"tenantSlug"`
	GeneratedAt      string                         `json:"generatedAt"`
	ExecutiveSummary CoreOperationsExecutiveSummary `json:"executiveSummary"`
	ServicePulse     map[string]any                 `json:"servicePulse"`
	Tenant360        map[string]any                 `json:"tenant360"`
	CoreOperations   map[string]any                 `json:"coreOperations"`
}
