package api

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
)

func Live(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.HealthResponse{
		Service: "documents",
		Status:  "live",
	})
}

func Ready(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.HealthResponse{
		Service: "documents",
		Status:  "ready",
	})
}

func DetailsForRuntime(repositoryDriver string) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(dto.ReadinessResponse{
			Service: "documents",
			Status:  "ready",
			Dependencies: []dto.DependencyResponse{
				{Name: "http-api", Status: "ready"},
				{Name: repositoryDriver, Status: "ready"},
			},
		})
	}
}
