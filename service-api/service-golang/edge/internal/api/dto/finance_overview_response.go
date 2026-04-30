// DTO do cockpit financeiro entregue pelo edge.
// O gateway agrega tesouraria, billing e risco operacional em um resumo unico.
package dto

type FinanceExecutiveSummary struct {
	Status                       string `json:"status"`
	CurrentBalanceCents          int    `json:"currentBalanceCents"`
	MonthlyRecurringRevenueCents int    `json:"monthlyRecurringRevenueCents"`
	ReceivablesPaidCents         int    `json:"receivablesPaidCents"`
	PayablesPaidCents            int    `json:"payablesPaidCents"`
	FailedPaymentAttempts        int    `json:"failedPaymentAttempts"`
	ActiveSubscriptions          int    `json:"activeSubscriptions"`
	PeriodClosures               int    `json:"periodClosures"`
	NetOperationalMarginCents    int    `json:"netOperationalMarginCents"`
}

type FinanceOverviewResponse struct {
	Service           string                  `json:"service"`
	TenantSlug        string                  `json:"tenantSlug"`
	GeneratedAt       string                  `json:"generatedAt"`
	ExecutiveSummary  FinanceExecutiveSummary `json:"executiveSummary"`
	ServicePulse      map[string]any          `json:"servicePulse"`
	Tenant360         map[string]any          `json:"tenant360"`
	FinanceControl    map[string]any          `json:"financeControl"`
	RevenueOperations map[string]any          `json:"revenueOperations"`
}
