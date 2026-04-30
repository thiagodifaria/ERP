package api

import (
	"net/http"
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, contractRepository repository.ContractRepository, attachmentGateway repository.AttachmentGateway) http.Handler {
	return NewRouterWithRuntime(logger, contractRepository, attachmentGateway, "memory", false)
}

func NewRouterWithRuntime(logger *telemetry.Logger, contractRepository repository.ContractRepository, attachmentGateway repository.AttachmentGateway, repositoryDriver string, documentsConfigured bool) http.Handler {
	mux := http.NewServeMux()
	handler := NewContractHandler(contractRepository, attachmentGateway)

	mux.HandleFunc("/health/live", Live)
	mux.HandleFunc("/health/ready", Ready)
	mux.HandleFunc("/health/details", DetailsForRuntime(repositoryDriver, documentsConfigured))
	mux.HandleFunc("GET /api/rentals/contracts", handler.ListContracts)
	mux.HandleFunc("GET /api/rentals/contracts/summary", handler.GetSummary)
	mux.HandleFunc("POST /api/rentals/contracts", handler.CreateContract)
	mux.HandleFunc("GET /api/rentals/contracts/{publicId}", handler.GetContract)
	mux.HandleFunc("GET /api/rentals/contracts/{publicId}/charges", handler.ListCharges)
	mux.HandleFunc("GET /api/rentals/contracts/{publicId}/history", handler.ListHistory)
	mux.HandleFunc("GET /api/rentals/contracts/{publicId}/adjustments", handler.ListAdjustments)
	mux.HandleFunc("POST /api/rentals/contracts/{publicId}/adjustments", handler.CreateAdjustment)
	mux.HandleFunc("POST /api/rentals/contracts/{publicId}/terminate", handler.TerminateContract)
	mux.HandleFunc("GET /api/rentals/contracts/{publicId}/attachments", handler.ListAttachments)
	mux.HandleFunc("POST /api/rentals/contracts/{publicId}/attachments", handler.CreateAttachment)

	return withCorrelation(logger, mux)
}

func hasDocumentsGatewayConfigured(documentsConfigured bool, attachmentGateway repository.AttachmentGateway) bool {
	return documentsConfigured || attachmentGateway != nil || strings.TrimSpace("") != ""
}
