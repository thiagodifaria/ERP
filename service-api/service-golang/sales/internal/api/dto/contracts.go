// DTOs da API descrevem contratos de entrada e saida.
// Objetos de dominio nao devem ser expostos diretamente aqui.
package dto

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type ReadinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []DependencyResponse `json:"dependencies"`
}

type DependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OpportunityResponse struct {
	PublicID     string `json:"publicId"`
	LeadPublicID string `json:"leadPublicId"`
	Title        string `json:"title"`
	Stage        string `json:"stage"`
	OwnerUserID  string `json:"ownerUserId"`
	AmountCents  int64  `json:"amountCents"`
}

type OpportunitySummaryResponse struct {
	Total            int            `json:"total"`
	TotalAmountCents int64          `json:"totalAmountCents"`
	ByStage          map[string]int `json:"byStage"`
}

type CreateOpportunityRequest struct {
	LeadPublicID string `json:"leadPublicId"`
	Title        string `json:"title"`
	OwnerUserID  string `json:"ownerUserId"`
	AmountCents  int64  `json:"amountCents"`
}

type UpdateOpportunityRequest struct {
	Title       string `json:"title"`
	OwnerUserID string `json:"ownerUserId"`
	AmountCents int64  `json:"amountCents"`
}

type UpdateOpportunityStageRequest struct {
	Stage string `json:"stage"`
}

type ProposalResponse struct {
	PublicID            string `json:"publicId"`
	OpportunityPublicID string `json:"opportunityPublicId"`
	Title               string `json:"title"`
	Status              string `json:"status"`
	AmountCents         int64  `json:"amountCents"`
}

type CreateProposalRequest struct {
	Title       string `json:"title"`
	AmountCents int64  `json:"amountCents"`
}

type UpdateProposalStatusRequest struct {
	Status string `json:"status"`
}

type SaleResponse struct {
	PublicID            string `json:"publicId"`
	OpportunityPublicID string `json:"opportunityPublicId"`
	ProposalPublicID    string `json:"proposalPublicId"`
	Status              string `json:"status"`
	AmountCents         int64  `json:"amountCents"`
}

type SaleSummaryResponse struct {
	Total              int            `json:"total"`
	BookedRevenueCents int64          `json:"bookedRevenueCents"`
	ByStatus           map[string]int `json:"byStatus"`
}

type UpdateSaleStatusRequest struct {
	Status string `json:"status"`
}

type InvoiceResponse struct {
	PublicID     string `json:"publicId"`
	SalePublicID string `json:"salePublicId"`
	Number       string `json:"number"`
	Status       string `json:"status"`
	AmountCents  int64  `json:"amountCents"`
	DueDate      string `json:"dueDate"`
	PaidAt       string `json:"paidAt"`
}

type InvoiceSummaryResponse struct {
	Total              int            `json:"total"`
	OpenAmountCents    int64          `json:"openAmountCents"`
	PaidAmountCents    int64          `json:"paidAmountCents"`
	OverdueAmountCents int64          `json:"overdueAmountCents"`
	OverdueCount       int            `json:"overdueCount"`
	ByStatus           map[string]int `json:"byStatus"`
}

type CreateInvoiceRequest struct {
	Number  string `json:"number"`
	DueDate string `json:"dueDate"`
}

type UpdateInvoiceStatusRequest struct {
	Status string `json:"status"`
}
