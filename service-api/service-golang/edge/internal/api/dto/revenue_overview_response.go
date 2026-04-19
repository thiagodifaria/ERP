// DTO do cockpit de receita entregue pelo edge.
// O gateway agrega comercial e cobranca em um resumo executivo unico.
package dto

type RevenueExecutiveSummary struct {
	Status            string `json:"status"`
	SalesWon          int    `json:"salesWon"`
	Invoices          int    `json:"invoices"`
	PaidInvoices      int    `json:"paidInvoices"`
	OpenAmountCents   int    `json:"openAmountCents"`
	PaidAmountCents   int    `json:"paidAmountCents"`
	OverdueInvoices   int    `json:"overdueInvoices"`
	CollectionRateBps int    `json:"collectionRateBps"`
}

type RevenueOverviewResponse struct {
	Service           string                  `json:"service"`
	TenantSlug        string                  `json:"tenantSlug"`
	GeneratedAt       string                  `json:"generatedAt"`
	ExecutiveSummary  RevenueExecutiveSummary `json:"executiveSummary"`
	ServicePulse      map[string]any          `json:"servicePulse"`
	SalesJourney      map[string]any          `json:"salesJourney"`
	RevenueOperations map[string]any          `json:"revenueOperations"`
}
