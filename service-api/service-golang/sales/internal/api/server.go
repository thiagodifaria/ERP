// Server monta o servidor HTTP do servico.
// Middlewares, router e handlers de entrada ficam aqui.
package api

import (
	"net/http"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
)

func NewServer(
	cfg config.Config,
	logger *telemetry.Logger,
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouterWithRuntime(logger, opportunityRepository, proposalRepository, saleRepository, invoiceRepository, eventRepository, outboxRepository, cfg.RepositoryDriver),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
