// Server monta o servidor HTTP do servico.
// Middlewares, router e handlers de entrada ficam aqui.
package api

import (
	"net/http"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

func NewServer(
	cfg config.Config,
	logger *telemetry.Logger,
	leadRepository repository.LeadRepository,
	leadNoteRepository repository.LeadNoteRepository,
) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouterWithRuntime(logger, leadRepository, leadNoteRepository, cfg.RepositoryDriver),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
