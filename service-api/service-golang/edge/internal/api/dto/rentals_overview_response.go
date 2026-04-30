// DTO do cockpit de locacoes entregue pelo edge.
// O gateway cruza operacao contratual, pulso do tenant e footprint documental.
package dto

type RentalsExecutiveSummary struct {
	Status                 string `json:"status"`
	Contracts              int    `json:"contracts"`
	ActiveContracts        int    `json:"activeContracts"`
	ScheduledCharges       int    `json:"scheduledCharges"`
	PaidCharges            int    `json:"paidCharges"`
	CancelledCharges       int    `json:"cancelledCharges"`
	OutstandingAmountCents int    `json:"outstandingAmountCents"`
	OverdueCharges         int    `json:"overdueCharges"`
	Attachments            int    `json:"attachments"`
}

type RentalsOverviewResponse struct {
	Service          string                  `json:"service"`
	TenantSlug       string                  `json:"tenantSlug"`
	GeneratedAt      string                  `json:"generatedAt"`
	ExecutiveSummary RentalsExecutiveSummary `json:"executiveSummary"`
	ServicePulse     map[string]any          `json:"servicePulse"`
	Tenant360        map[string]any          `json:"tenant360"`
	RentalOperations map[string]any          `json:"rentalOperations"`
}
