package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/config"
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

func buildStorageCapabilities(cfg config.Config) []dto.StorageCapabilityResponse {
	s3Configured := cfg.StorageDriver == "s3" || cfg.StorageDriver == "s3_compatible" || (cfg.StorageBucket != "" && cfg.StorageEndpoint != "")
	r2Configured := cfg.StorageDriver == "r2" || (cfg.R2AccountID != "" && cfg.R2Bucket != "")

	return []dto.StorageCapabilityResponse{
		{
			Provider:       "local",
			Scope:          "storage",
			Configured:     true,
			Mode:           "local",
			Status:         "manual",
			FallbackViable: true,
			SupportsLinks:  true,
			SupportsUpload: true,
			Notes:          []string{"Storage local permanece disponivel para runtime de desenvolvimento e smoke."},
		},
		{
			Provider:       "s3_compatible",
			Scope:          "storage",
			Configured:     s3Configured,
			CredentialKey:  "DOCUMENTS_STORAGE_BUCKET",
			Mode:           map[bool]string{true: "configured", false: "fallback"}[s3Configured],
			Status:         map[bool]string{true: "ready", false: "fallback"}[s3Configured],
			FallbackViable: true,
			SupportsLinks:  true,
			SupportsUpload: true,
			Notes:          []string{map[bool]string{true: "Storage S3-compativel pronto para anexos, links temporarios e uploads.", false: "Sem bucket/endpoint remoto, o servico opera com storage local sem quebrar o fluxo."}[s3Configured]},
		},
		{
			Provider:       "cloudflare_r2",
			Scope:          "storage",
			Configured:     r2Configured,
			CredentialKey:  "DOCUMENTS_R2_ACCOUNT_ID",
			Mode:           map[bool]string{true: "configured", false: "fallback"}[r2Configured],
			Status:         map[bool]string{true: "ready", false: "fallback"}[r2Configured],
			FallbackViable: true,
			SupportsLinks:  true,
			SupportsUpload: true,
			Notes:          []string{map[bool]string{true: "Cloudflare R2 pronto como backend remoto de storage.", false: "Fallback local cobre runtime sem conta R2 configurada."}[r2Configured]},
		},
	}
}

func DetailsForRuntime(repositoryDriver string, cfg config.Config) http.HandlerFunc {
	capabilities := buildStorageCapabilities(cfg)
	signingCapabilities := buildSigningCapabilities(cfg)

	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		dependencies := []dto.DependencyResponse{
			{Name: "http-api", Status: "ready"},
			{Name: repositoryDriver, Status: "ready"},
		}
		for _, capability := range capabilities {
			dependencies = append(dependencies, dto.DependencyResponse{
				Name:   "storage:" + capability.Provider,
				Status: capability.Status,
			})
		}
		for _, capability := range signingCapabilities {
			dependencies = append(dependencies, dto.DependencyResponse{
				Name:   "signing:" + capability.Provider,
				Status: capability.Status,
			})
		}

		_ = json.NewEncoder(writer).Encode(dto.ReadinessResponse{
			Service:      "documents",
			Status:       "ready",
			Dependencies: dependencies,
		})
	}
}

func StorageCapabilitiesForRuntime(cfg config.Config) http.HandlerFunc {
	capabilities := buildStorageCapabilities(cfg)

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		provider := strings.TrimSpace(request.PathValue("provider"))
		if provider == "" {
			writer.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(writer).Encode(capabilities)
			return
		}

		for _, capability := range capabilities {
			if capability.Provider == provider {
				writer.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(writer).Encode(capability)
				return
			}
		}

		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(map[string]string{
			"code":    "storage_capability_not_found",
			"message": "Storage capability was not found.",
		})
	}
}
