// Router define as rotas publicas do servico crm.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/handler"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/middleware"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, leadRepository repository.LeadRepository) http.Handler {
  mux := http.NewServeMux()
  leadHandler := handler.NewLeadHandler(
    query.NewListLeads(leadRepository),
    command.NewCreateLead(leadRepository),
  )

  mux.HandleFunc("/health/live", handler.Live)
  mux.HandleFunc("/health/ready", handler.Ready)
  mux.HandleFunc("/health/details", handler.Details)
  mux.HandleFunc("GET /api/crm/leads", leadHandler.List)
  mux.HandleFunc("POST /api/crm/leads", leadHandler.Create)

  return middleware.WithCorrelation(logger, mux)
}
