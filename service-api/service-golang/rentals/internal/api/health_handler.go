package api

import (
	"encoding/json"
	"net/http"
)

func Live(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, HealthResponse{
		Service: "rentals",
		Status:  "live",
	})
}

func Ready(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, HealthResponse{
		Service: "rentals",
		Status:  "ready",
	})
}

func DetailsForRuntime(repositoryDriver string, documentsConfigured bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		dependencies := []DependencyResponse{
			{Name: "http-api", Status: "ready"},
			{Name: repositoryDriver, Status: "ready"},
		}
		if documentsConfigured {
			dependencies = append(dependencies, DependencyResponse{Name: "documents", Status: "ready"})
		} else {
			dependencies = append(dependencies, DependencyResponse{Name: "documents", Status: "unconfigured"})
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(ReadinessResponse{
			Service:      "rentals",
			Status:       "ready",
			Dependencies: dependencies,
		})
	}
}
