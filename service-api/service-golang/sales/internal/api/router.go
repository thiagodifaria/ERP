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
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) http.Handler {
	return NewRouterWithRuntime(logger, opportunityRepository, proposalRepository, saleRepository, invoiceRepository, eventRepository, outboxRepository, "memory")
}

func NewRouterWithRuntime(
	logger *telemetry.Logger,
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
	repositoryDriver string,
) http.Handler {
	mux := http.NewServeMux()
	opportunityHandler := handler.NewOpportunityHandler(
		query.NewListOpportunities(opportunityRepository),
		query.NewGetOpportunitySummary(opportunityRepository),
		query.NewGetOpportunityByPublicID(opportunityRepository),
		command.NewCreateOpportunity(opportunityRepository, eventRepository),
		command.NewUpdateOpportunityProfile(opportunityRepository, eventRepository),
		command.NewUpdateOpportunityStage(opportunityRepository, eventRepository),
	)
	proposalHandler := handler.NewProposalHandler(
		query.NewListProposalsByOpportunity(proposalRepository),
		query.NewGetProposalByPublicID(proposalRepository),
		command.NewCreateProposal(opportunityRepository, proposalRepository, eventRepository),
		command.NewUpdateProposalStatus(proposalRepository, eventRepository),
		command.NewConvertProposalToSale(opportunityRepository, proposalRepository, saleRepository, eventRepository, outboxRepository),
	)
	saleHandler := handler.NewSaleHandler(
		query.NewListSales(saleRepository),
		query.NewGetSaleSummary(saleRepository),
		query.NewGetSaleByPublicID(saleRepository),
		command.NewUpdateSaleStatus(saleRepository, eventRepository, outboxRepository),
	)
	invoiceHandler := handler.NewInvoiceHandler(
		query.NewListInvoices(invoiceRepository),
		query.NewGetInvoiceSummary(invoiceRepository),
		query.NewGetInvoiceByPublicID(invoiceRepository),
		command.NewCreateInvoice(saleRepository, invoiceRepository, eventRepository, outboxRepository),
		command.NewUpdateInvoiceStatus(invoiceRepository, eventRepository, outboxRepository),
	)
	activityHandler := handler.NewActivityHandler(
		query.NewListCommercialEventsByAggregate(eventRepository),
		query.NewListPendingOutboxEvents(outboxRepository),
	)

	mux.HandleFunc("/health/live", handler.Live)
	mux.HandleFunc("/health/ready", handler.Ready)
	mux.HandleFunc("/health/details", handler.DetailsForRuntime(repositoryDriver))
	mux.HandleFunc("GET /api/sales/opportunities", opportunityHandler.List)
	mux.HandleFunc("GET /api/sales/opportunities/summary", opportunityHandler.Summary)
	mux.HandleFunc("POST /api/sales/opportunities", opportunityHandler.Create)
	mux.HandleFunc("GET /api/sales/opportunities/{publicId}", opportunityHandler.GetByPublicID)
	mux.HandleFunc("GET /api/sales/opportunities/{publicId}/history", activityHandler.ListOpportunityHistory)
	mux.HandleFunc("PATCH /api/sales/opportunities/{publicId}", opportunityHandler.Update)
	mux.HandleFunc("PATCH /api/sales/opportunities/{publicId}/stage", opportunityHandler.UpdateStage)
	mux.HandleFunc("GET /api/sales/opportunities/{publicId}/proposals", proposalHandler.ListByOpportunity)
	mux.HandleFunc("POST /api/sales/opportunities/{publicId}/proposals", proposalHandler.Create)
	mux.HandleFunc("GET /api/sales/proposals/{publicId}", proposalHandler.GetByPublicID)
	mux.HandleFunc("GET /api/sales/proposals/{publicId}/history", activityHandler.ListProposalHistory)
	mux.HandleFunc("PATCH /api/sales/proposals/{publicId}/status", proposalHandler.UpdateStatus)
	mux.HandleFunc("POST /api/sales/proposals/{publicId}/convert", proposalHandler.Convert)
	mux.HandleFunc("GET /api/sales/sales", saleHandler.List)
	mux.HandleFunc("GET /api/sales/sales/summary", saleHandler.Summary)
	mux.HandleFunc("GET /api/sales/sales/{publicId}", saleHandler.GetByPublicID)
	mux.HandleFunc("GET /api/sales/sales/{publicId}/history", activityHandler.ListSaleHistory)
	mux.HandleFunc("PATCH /api/sales/sales/{publicId}/status", saleHandler.UpdateStatus)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/invoice", invoiceHandler.Create)
	mux.HandleFunc("GET /api/sales/invoices", invoiceHandler.List)
	mux.HandleFunc("GET /api/sales/invoices/summary", invoiceHandler.Summary)
	mux.HandleFunc("GET /api/sales/invoices/{publicId}", invoiceHandler.GetByPublicID)
	mux.HandleFunc("GET /api/sales/invoices/{publicId}/history", activityHandler.ListInvoiceHistory)
	mux.HandleFunc("PATCH /api/sales/invoices/{publicId}/status", invoiceHandler.UpdateStatus)
	mux.HandleFunc("GET /api/sales/outbox/pending", activityHandler.ListPendingOutbox)

	return middleware.WithCorrelation(logger, mux)
}
