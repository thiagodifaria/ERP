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

func NewServer(cfg config.Config, logger *telemetry.Logger, leadRepository repository.LeadRepository) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           NewRouter(logger, leadRepository),
		ReadHeaderTimeout: 5 * time.Second,
	}
}
