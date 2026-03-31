// Server monta o servidor HTTP do servico.
// Middlewares, router e handlers de entrada ficam aqui.
package api

import (
  "net/http"
  "time"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/config"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func NewServer(cfg config.Config, logger *telemetry.Logger) *http.Server {
  return &http.Server{
    Addr:              cfg.HTTPAddress,
    Handler:           NewRouter(logger),
    ReadHeaderTimeout: 5 * time.Second,
  }
}
