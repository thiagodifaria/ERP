package dto

type GoLiveExecutiveSummary struct {
	Status            string `json:"status"`
	PlannedRollouts   int    `json:"plannedRollouts"`
	RunningRollouts   int    `json:"runningRollouts"`
	CompletedRollouts int    `json:"completedRollouts"`
	RolledBackRollouts int   `json:"rolledBackRollouts"`
	TrackedMetrics    int    `json:"trackedMetrics"`
	TotalQuantity     int    `json:"totalQuantity"`
	AdoptionPct       int    `json:"adoptionPct"`
	PendingAdjustments int   `json:"pendingAdjustments"`
	CriticalBottlenecks int  `json:"criticalBottlenecks"`
	RolloutReady      bool   `json:"rolloutReady"`
	MetricsObserved   bool   `json:"metricsObserved"`
}

type GoLiveOverviewResponse struct {
	Service        string                 `json:"service"`
	TenantSlug     string                 `json:"tenantSlug"`
	GeneratedAt    string                 `json:"generatedAt"`
	ExecutiveSummary GoLiveExecutiveSummary `json:"executiveSummary"`
	ServicePulse   map[string]any         `json:"servicePulse"`
	SaaSControl    map[string]any         `json:"saasControl"`
	GoLiveControl  map[string]any         `json:"goLiveControl"`
}
