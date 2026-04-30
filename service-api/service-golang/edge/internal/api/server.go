// Server monta o servidor HTTP do servico.
// Middlewares, router e handlers de entrada ficam aqui.
package api

import (
	"net/http"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func NewServer(cfg config.Config, logger *telemetry.Logger) *http.Server {
	checker := integration.NewHTTPHealthChecker(cfg.DownstreamTimeout)
	dependencies := []integration.ServiceEndpoint{
		{Name: "identity", BaseURL: cfg.IdentityBaseURL},
		{Name: "crm", BaseURL: cfg.CRMBaseURL},
		{Name: "workflow-control", BaseURL: cfg.WorkflowControlBaseURL},
		{Name: "workflow-runtime", BaseURL: cfg.WorkflowRuntimeBaseURL},
		{Name: "analytics", BaseURL: cfg.AnalyticsBaseURL},
		{Name: "webhook-hub", BaseURL: cfg.WebhookHubBaseURL},
		{Name: "sales", BaseURL: cfg.SalesBaseURL},
	}
	healthHandler := handler.NewHealthHandler(cfg.ServiceName, checker, dependencies)
	opsHandler := handler.NewOpsHandler(cfg.ServiceName, checker, dependencies)
	tenantOverviewHandler := handler.NewTenantOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	automationOverviewHandler := handler.NewAutomationOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	engagementOverviewHandler := handler.NewEngagementOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	salesOverviewHandler := handler.NewSalesOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	revenueOverviewHandler := handler.NewRevenueOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	financeOverviewHandler := handler.NewFinanceOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	rentalsOverviewHandler := handler.NewRentalsOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)
	accessResolver := integration.NewHTTPIdentityAccessResolver(cfg.DownstreamTimeout)

	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouter(logger, healthHandler, opsHandler, tenantOverviewHandler, automationOverviewHandler, engagementOverviewHandler, salesOverviewHandler, revenueOverviewHandler, financeOverviewHandler, rentalsOverviewHandler, cfg.IdentityBaseURL, accessResolver),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
