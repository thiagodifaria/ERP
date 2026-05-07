package dto

type ComplianceExecutiveSummary struct {
	Status                 string `json:"status"`
	FiscalDocuments        int    `json:"fiscalDocuments"`
	CancelledDocuments     int    `json:"cancelledDocuments"`
	DocumentEvents         int    `json:"documentEvents"`
	PendingPrivacyRequests int    `json:"pendingPrivacyRequests"`
	GrantedConsents        int    `json:"grantedConsents"`
	RestrictedDocuments    int    `json:"restrictedDocuments"`
	AuditEvents            int    `json:"auditEvents"`
}

type ComplianceOverviewResponse struct {
	Service           string                      `json:"service"`
	TenantSlug        string                      `json:"tenantSlug"`
	GeneratedAt       string                      `json:"generatedAt"`
	ExecutiveSummary  ComplianceExecutiveSummary  `json:"executiveSummary"`
	ServicePulse      map[string]any              `json:"servicePulse"`
	DocumentGovernance map[string]any             `json:"documentGovernance"`
	ComplianceControl map[string]any              `json:"complianceControl"`
}
