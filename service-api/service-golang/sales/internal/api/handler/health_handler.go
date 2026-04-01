// Health handlers expoem a disponibilidade minima do servico.
// Consultas de negocio nao devem entrar aqui.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
)

func Live(writer http.ResponseWriter, request *http.Request) {
	respond(writer, dto.HealthResponse{
		Service: "sales",
		Status:  "live",
	})
}

func Ready(writer http.ResponseWriter, request *http.Request) {
	respond(writer, dto.HealthResponse{
		Service: "sales",
		Status:  "ready",
	})
}

func Details(writer http.ResponseWriter, request *http.Request) {
	DetailsForRuntime("memory")(writer, request)
}

func DetailsForRuntime(repositoryDriver string) http.HandlerFunc {
	readiness := buildReadiness(repositoryDriver)

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(readiness)
	}
}

func buildReadiness(repositoryDriver string) dto.ReadinessResponse {
	postgresqlStatus := "pending-runtime-wiring"
	if repositoryDriver == "postgres" {
		postgresqlStatus = "ready"
	}

	return dto.ReadinessResponse{
		Service: "sales",
		Status:  "ready",
		Dependencies: []dto.DependencyResponse{
			{Name: "router", Status: "ready"},
			{Name: "postgresql", Status: postgresqlStatus},
			{Name: "crm", Status: "pending-runtime-wiring"},
		},
	}
}

func respond(writer http.ResponseWriter, response dto.HealthResponse) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}
