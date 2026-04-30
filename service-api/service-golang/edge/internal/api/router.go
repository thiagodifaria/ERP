// Router define as rotas publicas do servico edge.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/middleware"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func NewRouter(
	logger *telemetry.Logger,
	healthHandler handler.HealthHandler,
	opsHandler handler.OpsHandler,
	tenantOverviewHandler handler.TenantOverviewHandler,
	automationOverviewHandler handler.AutomationOverviewHandler,
	engagementOverviewHandler handler.EngagementOverviewHandler,
	salesOverviewHandler handler.SalesOverviewHandler,
	revenueOverviewHandler handler.RevenueOverviewHandler,
	rentalsOverviewHandler handler.RentalsOverviewHandler,
	identityBaseURL string,
	accessResolver integration.TenantAccessResolver,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", healthHandler.Live)
	mux.HandleFunc("/health/ready", healthHandler.Ready)
	mux.HandleFunc("/health/details", healthHandler.Details)
	mux.HandleFunc("/api/edge/ops/health", opsHandler.Health)
	mux.Handle("/api/edge/ops/tenant-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(tenantOverviewHandler.Overview)))
	mux.Handle("/api/edge/ops/automation-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(automationOverviewHandler.Overview)))
	mux.Handle("/api/edge/ops/engagement-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(engagementOverviewHandler.Overview)))
	mux.Handle("/api/edge/ops/sales-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(salesOverviewHandler.Overview)))
	mux.Handle("/api/edge/ops/revenue-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(revenueOverviewHandler.Overview)))
	mux.Handle("/api/edge/ops/rentals-overview", middleware.WithTenantAccess(identityBaseURL, accessResolver, http.HandlerFunc(rentalsOverviewHandler.Overview)))

	return middleware.WithCorrelation(logger, mux)
}
