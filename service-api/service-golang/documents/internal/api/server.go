package api

import (
	"net/http"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
)

func NewServer(cfg config.Config, logger *telemetry.Logger, attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouterWithRuntime(logger, attachmentRepository, uploadSessionRepository, cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
