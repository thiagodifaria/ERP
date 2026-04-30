package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository) http.Handler {
	return NewRouterWithRuntime(logger, attachmentRepository, uploadSessionRepository, "memory", "documents-local-secret")
}

func NewRouterWithRuntime(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository, repositoryDriver string, accessTokenSecret string) http.Handler {
	mux := http.NewServeMux()
	attachmentHandler := NewAttachmentHandler(attachmentRepository, uploadSessionRepository, accessTokenSecret)

	mux.HandleFunc("/health/live", Live)
	mux.HandleFunc("/health/ready", Ready)
	mux.HandleFunc("/health/details", DetailsForRuntime(repositoryDriver))
	mux.HandleFunc("GET /api/documents/attachments", attachmentHandler.List)
	mux.HandleFunc("POST /api/documents/attachments", attachmentHandler.Create)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}", attachmentHandler.Get)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}/download", attachmentHandler.Download)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/archive", attachmentHandler.Archive)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/access-links", attachmentHandler.CreateAccessLink)
	mux.HandleFunc("POST /api/documents/upload-sessions", attachmentHandler.CreateUploadSession)
	mux.HandleFunc("GET /api/documents/upload-sessions/{publicId}", attachmentHandler.GetUploadSession)
	mux.HandleFunc("POST /api/documents/upload-sessions/{publicId}/complete", attachmentHandler.CompleteUploadSession)

	return withCorrelation(logger, mux)
}
