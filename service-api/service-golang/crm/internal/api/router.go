// Router define as rotas publicas do servico crm.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/handler"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/middleware"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger) http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("/health/live", handler.Live)
  mux.HandleFunc("/health/ready", handler.Ready)
  mux.HandleFunc("/health/details", handler.Details)

  return middleware.WithCorrelation(logger, mux)
}
