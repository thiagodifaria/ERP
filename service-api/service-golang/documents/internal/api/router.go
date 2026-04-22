package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository) http.Handler {
	return NewRouterWithRuntime(logger, attachmentRepository, "memory")
}

func NewRouterWithRuntime(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, repositoryDriver string) http.Handler {
	mux := http.NewServeMux()
	attachmentHandler := NewAttachmentHandler(attachmentRepository)

	mux.HandleFunc("/health/live", Live)
	mux.HandleFunc("/health/ready", Ready)
	mux.HandleFunc("/health/details", DetailsForRuntime(repositoryDriver))
	mux.HandleFunc("GET /api/documents/attachments", attachmentHandler.List)
	mux.HandleFunc("POST /api/documents/attachments", attachmentHandler.Create)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}", attachmentHandler.Get)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/archive", attachmentHandler.Archive)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/access-links", attachmentHandler.CreateAccessLink)

	return withCorrelation(logger, mux)
}
