// Router define as rotas publicas do servico edge.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/middleware"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, healthHandler handler.HealthHandler) http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("/health/live", healthHandler.Live)
  mux.HandleFunc("/health/ready", healthHandler.Ready)
  mux.HandleFunc("/health/details", healthHandler.Details)

  return middleware.WithCorrelation(logger, mux)
}
