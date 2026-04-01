// Handlers operacionais do edge para leitura do estado da plataforma.
// O gateway consolida diagnostico rapido sem reimplementar logica de dominio.
package handler

import (
  "context"
  "encoding/json"
  "net/http"
  "time"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
)

type OpsHandler struct {
  ServiceName  string
  Checker      integration.HealthChecker
  Dependencies []integration.ServiceEndpoint
}

func NewOpsHandler(serviceName string, checker integration.HealthChecker, dependencies []integration.ServiceEndpoint) OpsHandler {
  return OpsHandler{
    ServiceName:  serviceName,
    Checker:      checker,
    Dependencies: dependencies,
  }
}

func (handler OpsHandler) Health(writer http.ResponseWriter, request *http.Request) {
  snapshots := make([]dto.ServiceHealthSnapshot, 0, len(handler.Dependencies))
  summary := dto.OpsHealthSummary{
    Total: len(handler.Dependencies),
  }
  status := "ready"

  for _, dependency := range handler.Dependencies {
    snapshot := handler.Checker.Details(context.Background(), dependency)
    if snapshot.Status == "ready" {
      summary.Ready++
    } else {
      summary.Degraded++
      status = "degraded"
    }

    snapshots = append(snapshots, snapshot)
  }

  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(writer).Encode(dto.OpsHealthResponse{
    Service:     handler.ServiceName,
    Status:      status,
    GeneratedAt: time.Now().UTC().Format(time.RFC3339),
    Summary:     summary,
    Services:    snapshots,
  })
}
