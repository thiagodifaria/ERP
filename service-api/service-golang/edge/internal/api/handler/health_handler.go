// Health handlers expõem a disponibilidade minima do servico.
// Consultas de negocio nao devem entrar aqui.
package handler

import (
  "encoding/json"
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

func Live(writer http.ResponseWriter, request *http.Request) {
  respond(writer, dto.HealthResponse{
    Service: "edge",
    Status:  "live",
  })
}

func Ready(writer http.ResponseWriter, request *http.Request) {
  respond(writer, dto.HealthResponse{
    Service: "edge",
    Status:  "ready",
  })
}

func respond(writer http.ResponseWriter, response dto.HealthResponse) {
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(writer).Encode(response)
}
