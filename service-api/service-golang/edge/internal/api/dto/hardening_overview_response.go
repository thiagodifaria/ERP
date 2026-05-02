package dto

type HardeningExecutiveSummary struct {
	Status                 string `json:"status"`
	StableChecks           int    `json:"stableChecks"`
	AttentionChecks        int    `json:"attentionChecks"`
	CriticalChecks         int    `json:"criticalChecks"`
	DeadLetterEvents       int    `json:"deadLetterEvents"`
	FailedPaymentAttempts  int    `json:"failedPaymentAttempts"`
	OpenSecurityAlerts     int    `json:"openSecurityAlerts"`
	LatestBenchmarkStatus  string `json:"latestBenchmarkStatus"`
	BackupRestoreValidated bool   `json:"backupRestoreValidated"`
	CriticalProviderGaps   int    `json:"criticalProviderGaps"`
	HttpSpecs             int    `json:"httpSpecs"`
}

type HardeningOverviewResponse struct {
	Service             string                    `json:"service"`
	TenantSlug          string                    `json:"tenantSlug"`
	GeneratedAt         string                    `json:"generatedAt"`
	ExecutiveSummary    HardeningExecutiveSummary `json:"executiveSummary"`
	ServicePulse        map[string]any            `json:"servicePulse"`
	PlatformReliability map[string]any            `json:"platformReliability"`
	HardeningReview     map[string]any            `json:"hardeningReview"`
}
