package dto

type PlatformReliabilityExecutiveSummary struct {
	Status                   string `json:"status"`
	PendingWebhookEvents     int    `json:"pendingWebhookEvents"`
	DeadLetterEvents         int    `json:"deadLetterEvents"`
	FailedWorkflowExecutions int    `json:"failedWorkflowExecutions"`
	CriticalRecoveryCases    int    `json:"criticalRecoveryCases"`
	FailedPaymentAttempts    int    `json:"failedPaymentAttempts"`
	WebhookForwardingRateBps int    `json:"webhookForwardingRateBps"`
	WorkflowSuccessRateBps   int    `json:"workflowSuccessRateBps"`
	BillingRecoveryRateBps   int    `json:"billingRecoveryRateBps"`
	OpenCriticalRisks        int    `json:"openCriticalRisks"`
}

type PlatformReliabilityOverviewResponse struct {
	Service             string                              `json:"service"`
	TenantSlug          string                              `json:"tenantSlug"`
	GeneratedAt         string                              `json:"generatedAt"`
	ExecutiveSummary    PlatformReliabilityExecutiveSummary `json:"executiveSummary"`
	ServicePulse        map[string]any                      `json:"servicePulse"`
	DeliveryReliability map[string]any                      `json:"deliveryReliability"`
	PlatformReliability map[string]any                      `json:"platformReliability"`
}
