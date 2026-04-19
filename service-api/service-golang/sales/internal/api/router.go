// Router define as rotas publicas do servico sales.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/middleware"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
)

func NewRouter(
	logger *telemetry.Logger,
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
) http.Handler {
	return NewRouterWithRuntime(logger, opportunityRepository, proposalRepository, saleRepository, invoiceRepository, "memory")
}

func NewRouterWithRuntime(
	logger *telemetry.Logger,
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	repositoryDriver string,
) http.Handler {
	mux := http.NewServeMux()
	opportunityHandler := handler.NewOpportunityHandler(
		query.NewListOpportunities(opportunityRepository),
		query.NewGetOpportunitySummary(opportunityRepository),
		query.NewGetOpportunityByPublicID(opportunityRepository),
		command.NewCreateOpportunity(opportunityRepository),
		command.NewUpdateOpportunityProfile(opportunityRepository),
		command.NewUpdateOpportunityStage(opportunityRepository),
	)
	proposalHandler := handler.NewProposalHandler(
		query.NewListProposalsByOpportunity(proposalRepository),
		query.NewGetProposalByPublicID(proposalRepository),
		command.NewCreateProposal(opportunityRepository, proposalRepository),
		command.NewUpdateProposalStatus(proposalRepository),
		command.NewConvertProposalToSale(opportunityRepository, proposalRepository, saleRepository),
	)
	saleHandler := handler.NewSaleHandler(
		query.NewListSales(saleRepository),
		query.NewGetSaleSummary(saleRepository),
		query.NewGetSaleByPublicID(saleRepository),
		command.NewUpdateSaleStatus(saleRepository),
	)
	invoiceHandler := handler.NewInvoiceHandler(
		query.NewListInvoices(invoiceRepository),
		query.NewGetInvoiceSummary(invoiceRepository),
		query.NewGetInvoiceByPublicID(invoiceRepository),
		command.NewCreateInvoice(saleRepository, invoiceRepository),
		command.NewUpdateInvoiceStatus(invoiceRepository),
	)

	mux.HandleFunc("/health/live", handler.Live)
	mux.HandleFunc("/health/ready", handler.Ready)
	mux.HandleFunc("/health/details", handler.DetailsForRuntime(repositoryDriver))
	mux.HandleFunc("GET /api/sales/opportunities", opportunityHandler.List)
	mux.HandleFunc("GET /api/sales/opportunities/summary", opportunityHandler.Summary)
	mux.HandleFunc("POST /api/sales/opportunities", opportunityHandler.Create)
	mux.HandleFunc("GET /api/sales/opportunities/{publicId}", opportunityHandler.GetByPublicID)
	mux.HandleFunc("PATCH /api/sales/opportunities/{publicId}", opportunityHandler.Update)
	mux.HandleFunc("PATCH /api/sales/opportunities/{publicId}/stage", opportunityHandler.UpdateStage)
	mux.HandleFunc("GET /api/sales/opportunities/{publicId}/proposals", proposalHandler.ListByOpportunity)
	mux.HandleFunc("POST /api/sales/opportunities/{publicId}/proposals", proposalHandler.Create)
	mux.HandleFunc("GET /api/sales/proposals/{publicId}", proposalHandler.GetByPublicID)
	mux.HandleFunc("PATCH /api/sales/proposals/{publicId}/status", proposalHandler.UpdateStatus)
	mux.HandleFunc("POST /api/sales/proposals/{publicId}/convert", proposalHandler.Convert)
	mux.HandleFunc("GET /api/sales/sales", saleHandler.List)
	mux.HandleFunc("GET /api/sales/sales/summary", saleHandler.Summary)
	mux.HandleFunc("GET /api/sales/sales/{publicId}", saleHandler.GetByPublicID)
	mux.HandleFunc("PATCH /api/sales/sales/{publicId}/status", saleHandler.UpdateStatus)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/invoice", invoiceHandler.Create)
	mux.HandleFunc("GET /api/sales/invoices", invoiceHandler.List)
	mux.HandleFunc("GET /api/sales/invoices/summary", invoiceHandler.Summary)
	mux.HandleFunc("GET /api/sales/invoices/{publicId}", invoiceHandler.GetByPublicID)
	mux.HandleFunc("PATCH /api/sales/invoices/{publicId}/status", invoiceHandler.UpdateStatus)

	return middleware.WithCorrelation(logger, mux)
}
