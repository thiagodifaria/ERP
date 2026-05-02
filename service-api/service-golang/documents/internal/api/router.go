package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
)

func NewRouter(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository) http.Handler {
	return NewRouterWithRuntime(
		logger,
		attachmentRepository,
		uploadSessionRepository,
		config.Config{
			RepositoryDriver:  "memory",
			AccessTokenSecret: "documents-local-secret",
			StorageDriver:     "local",
		},
	)
}

func NewRouterWithRuntime(logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository, cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	attachmentHandler := NewAttachmentHandler(attachmentRepository, uploadSessionRepository, cfg.AccessTokenSecret)

	mux.HandleFunc("/health/live", Live)
	mux.HandleFunc("/health/ready", Ready)
	mux.HandleFunc("/health/details", DetailsForRuntime(cfg.RepositoryDriver, cfg))
	mux.HandleFunc("GET /api/documents/storage/capabilities", StorageCapabilitiesForRuntime(cfg))
	mux.HandleFunc("GET /api/documents/storage/capabilities/{provider}", StorageCapabilitiesForRuntime(cfg))
	mux.HandleFunc("GET /api/documents/attachments", attachmentHandler.List)
	mux.HandleFunc("POST /api/documents/attachments", attachmentHandler.Create)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}", attachmentHandler.Get)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}/versions", attachmentHandler.ListVersions)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/versions", attachmentHandler.CreateVersion)
	mux.HandleFunc("GET /api/documents/attachments/{publicId}/download", attachmentHandler.Download)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/archive", attachmentHandler.Archive)
	mux.HandleFunc("POST /api/documents/attachments/{publicId}/access-links", attachmentHandler.CreateAccessLink)
	mux.HandleFunc("POST /api/documents/upload-sessions", attachmentHandler.CreateUploadSession)
	mux.HandleFunc("GET /api/documents/upload-sessions/{publicId}", attachmentHandler.GetUploadSession)
	mux.HandleFunc("POST /api/documents/upload-sessions/{publicId}/complete", attachmentHandler.CompleteUploadSession)

	return withCorrelation(logger, mux)
}
