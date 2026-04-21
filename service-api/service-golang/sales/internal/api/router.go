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
	installmentRepository repository.InstallmentRepository,
	commissionRepository repository.CommissionRepository,
	pendingItemRepository repository.PendingItemRepository,
	renegotiationRepository repository.RenegotiationRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) http.Handler {
	return NewRouterWithRuntime(logger, opportunityRepository, proposalRepository, saleRepository, invoiceRepository, installmentRepository, commissionRepository, pendingItemRepository, renegotiationRepository, eventRepository, outboxRepository, "memory")
}

func NewRouterWithRuntime(
	logger *telemetry.Logger,
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	installmentRepository repository.InstallmentRepository,
	commissionRepository repository.CommissionRepository,
	pendingItemRepository repository.PendingItemRepository,
	renegotiationRepository repository.RenegotiationRepository,
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
	operationsHandler := handler.NewOperationsHandler(
		query.NewListInstallmentsBySale(installmentRepository),
		command.NewCreateInstallmentSchedule(saleRepository, installmentRepository, eventRepository, outboxRepository),
		query.NewListCommissionsBySale(commissionRepository),
		command.NewCreateCommission(saleRepository, commissionRepository, eventRepository),
		command.NewUpdateCommissionStatus(commissionRepository, eventRepository),
		query.NewListPendingItemsBySale(pendingItemRepository),
		command.NewCreatePendingItem(saleRepository, pendingItemRepository, eventRepository),
		command.NewResolvePendingItem(pendingItemRepository, eventRepository),
		query.NewListRenegotiationsBySale(renegotiationRepository),
		command.NewApplyRenegotiation(saleRepository, invoiceRepository, installmentRepository, commissionRepository, renegotiationRepository, eventRepository, outboxRepository),
		command.NewCancelSale(saleRepository, invoiceRepository, installmentRepository, commissionRepository, pendingItemRepository, eventRepository, outboxRepository),
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
	mux.HandleFunc("POST /api/sales/sales/{publicId}/cancel", operationsHandler.CancelSale)
	mux.HandleFunc("GET /api/sales/sales/{publicId}/installments", operationsHandler.ListInstallments)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/installments", operationsHandler.CreateInstallments)
	mux.HandleFunc("GET /api/sales/sales/{publicId}/commissions", operationsHandler.ListCommissions)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/commissions", operationsHandler.CreateCommission)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/commissions/{commissionPublicId}/block", operationsHandler.BlockCommission)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/commissions/{commissionPublicId}/release", operationsHandler.ReleaseCommission)
	mux.HandleFunc("GET /api/sales/sales/{publicId}/pending-items", operationsHandler.ListPendingItems)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/pending-items", operationsHandler.CreatePendingItem)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/pending-items/{pendingItemPublicId}/resolve", operationsHandler.ResolvePendingItem)
	mux.HandleFunc("GET /api/sales/sales/{publicId}/renegotiations", operationsHandler.ListRenegotiations)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/renegotiations", operationsHandler.ApplyRenegotiation)
	mux.HandleFunc("POST /api/sales/sales/{publicId}/invoice", invoiceHandler.Create)
	mux.HandleFunc("GET /api/sales/invoices", invoiceHandler.List)
	mux.HandleFunc("GET /api/sales/invoices/summary", invoiceHandler.Summary)
	mux.HandleFunc("GET /api/sales/invoices/{publicId}", invoiceHandler.GetByPublicID)
	mux.HandleFunc("GET /api/sales/invoices/{publicId}/history", activityHandler.ListInvoiceHistory)
	mux.HandleFunc("PATCH /api/sales/invoices/{publicId}/status", invoiceHandler.UpdateStatus)
	mux.HandleFunc("GET /api/sales/outbox/pending", activityHandler.ListPendingOutbox)

	return middleware.WithCorrelation(logger, mux)
}
