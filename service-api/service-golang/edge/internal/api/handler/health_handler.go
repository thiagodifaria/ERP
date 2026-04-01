// Health handlers expõem a disponibilidade minima do servico.
// Consultas de negocio nao devem entrar aqui.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
)

type HealthHandler struct {
	ServiceName  string
	Checker      integration.HealthChecker
	Dependencies []integration.ServiceEndpoint
}

func NewHealthHandler(serviceName string, checker integration.HealthChecker, dependencies []integration.ServiceEndpoint) HealthHandler {
	return HealthHandler{
		ServiceName:  serviceName,
		Checker:      checker,
		Dependencies: dependencies,
	}
}

func (handler HealthHandler) Live(writer http.ResponseWriter, request *http.Request) {
	respond(writer, dto.HealthResponse{
		Service: handler.ServiceName,
		Status:  "live",
	})
}

func (handler HealthHandler) Ready(writer http.ResponseWriter, request *http.Request) {
	respond(writer, dto.HealthResponse{
		Service: handler.ServiceName,
		Status:  "ready",
	})
}

func (handler HealthHandler) Details(writer http.ResponseWriter, request *http.Request) {
	status := "ready"
	dependencies := []dto.DependencyResponse{
		{
			Name:   "router",
			Status: "ready",
		},
	}

	for _, dependency := range handler.Dependencies {
		health := handler.Checker.Check(context.Background(), dependency)
		if health.Status != "ready" {
			status = "degraded"
		}

		dependencies = append(dependencies, health)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.ReadinessResponse{
		Service:      handler.ServiceName,
		Status:       status,
		Dependencies: dependencies,
	})
}

func respond(writer http.ResponseWriter, response dto.HealthResponse) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}
