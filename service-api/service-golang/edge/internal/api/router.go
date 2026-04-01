// Router define as rotas publicas do servico edge.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/middleware"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func NewRouter(
	logger *telemetry.Logger,
	healthHandler handler.HealthHandler,
	opsHandler handler.OpsHandler,
	tenantOverviewHandler handler.TenantOverviewHandler,
	automationOverviewHandler handler.AutomationOverviewHandler,
	salesOverviewHandler handler.SalesOverviewHandler,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health/live", healthHandler.Live)
	mux.HandleFunc("/health/ready", healthHandler.Ready)
	mux.HandleFunc("/health/details", healthHandler.Details)
	mux.HandleFunc("/api/edge/ops/health", opsHandler.Health)
	mux.HandleFunc("/api/edge/ops/tenant-overview", tenantOverviewHandler.Overview)
	mux.HandleFunc("/api/edge/ops/automation-overview", automationOverviewHandler.Overview)
	mux.HandleFunc("/api/edge/ops/sales-overview", salesOverviewHandler.Overview)

	return middleware.WithCorrelation(logger, mux)
}
