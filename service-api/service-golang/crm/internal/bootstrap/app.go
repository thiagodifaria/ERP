// App concentra o bootstrap do servico crm.
// Aqui ficam configuracao, logger e servidor HTTP.
package bootstrap

import (
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/config"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

type App struct {
  Config config.Config
  Logger *telemetry.Logger
  Server *http.Server
}

func NewApp() (*App, error) {
  cfg := config.Load()
  logger := telemetry.New(cfg.ServiceName)
  server := api.NewServer(cfg, logger)

  return &App{
    Config: cfg,
    Logger: logger,
    Server: server,
  }, nil
}

func (app *App) Run() error {
  app.Logger.Printf("starting %s on %s", app.Config.ServiceName, app.Config.HTTPAddress)
  return app.Server.ListenAndServe()
}
