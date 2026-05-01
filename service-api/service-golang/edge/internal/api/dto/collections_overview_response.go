// DTO do cockpit de collections entregue pelo edge.
// O gateway consolida cobranca, recuperacao e risco de receita sem duplicar calculos do analytics.
package dto

type CollectionsExecutiveSummary struct {
	Status                string `json:"status"`
	CasesTotal            int    `json:"casesTotal"`
	CriticalCases         int    `json:"criticalCases"`
	InvoicesInRecovery    int    `json:"invoicesInRecovery"`
	OpenAmountCents       int    `json:"openAmountCents"`
	RecoveredAmountCents  int    `json:"recoveredAmountCents"`
	FailedPaymentAttempts int    `json:"failedPaymentAttempts"`
	ActivePromises        int    `json:"activePromises"`
	NextActionsDue        int    `json:"nextActionsDue"`
	RecoveryRateBps       int    `json:"recoveryRateBps"`
}

type CollectionsOverviewResponse struct {
	Service            string                      `json:"service"`
	TenantSlug         string                      `json:"tenantSlug"`
	GeneratedAt        string                      `json:"generatedAt"`
	ExecutiveSummary   CollectionsExecutiveSummary `json:"executiveSummary"`
	ServicePulse       map[string]any              `json:"servicePulse"`
	Tenant360          map[string]any              `json:"tenant360"`
	FinanceControl     map[string]any              `json:"financeControl"`
	CollectionsControl map[string]any              `json:"collectionsControl"`
}
