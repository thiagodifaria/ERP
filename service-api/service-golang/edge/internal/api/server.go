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
  }
  healthHandler := handler.NewHealthHandler(cfg.ServiceName, checker, dependencies)
  opsHandler := handler.NewOpsHandler(cfg.ServiceName, checker, dependencies)
  tenantOverviewHandler := handler.NewTenantOverviewHandler(cfg.ServiceName, cfg.AnalyticsBaseURL, checker)

  return &http.Server{
    Addr:              cfg.HTTPAddress,
    Handler:           NewRouter(logger, healthHandler, opsHandler, tenantOverviewHandler),
    ReadHeaderTimeout: 5 * time.Second,
  }
}
