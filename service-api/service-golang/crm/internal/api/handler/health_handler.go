// Health handlers expoem a disponibilidade minima do servico.
// Consultas de negocio nao devem entrar aqui.
package handler

import (
  "encoding/json"
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
)

func Live(writer http.ResponseWriter, request *http.Request) {
  respond(writer, dto.HealthResponse{
    Service: "crm",
    Status:  "live",
  })
}

func Ready(writer http.ResponseWriter, request *http.Request) {
  respond(writer, dto.HealthResponse{
    Service: "crm",
    Status:  "ready",
  })
}

func Details(writer http.ResponseWriter, request *http.Request) {
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(writer).Encode(dto.ReadinessResponse{
    Service: "crm",
    Status:  "ready",
    Dependencies: []dto.DependencyResponse{
      {
        Name:   "router",
        Status: "ready",
      },
      {
        Name:   "postgresql",
        Status: "pending-runtime-wiring",
      },
      {
        Name:   "edge",
        Status: "pending-runtime-wiring",
      },
    },
  })
}

func respond(writer http.ResponseWriter, response dto.HealthResponse) {
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(writer).Encode(response)
}
