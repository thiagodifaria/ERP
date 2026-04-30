package api

import (
	"net/http"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/telemetry"
)

func NewServer(cfg config.Config, logger *telemetry.Logger, contractRepository repository.ContractRepository, attachmentGateway repository.AttachmentGateway) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouterWithRuntime(logger, contractRepository, attachmentGateway, cfg.RepositoryDriver, cfg.DocumentsBaseURL != ""),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
