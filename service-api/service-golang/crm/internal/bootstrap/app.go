// App concentra o bootstrap do servico crm.
// Aqui ficam configuracao, logger e servidor HTTP.
package bootstrap

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

type App struct {
	Config   config.Config
	Logger   *telemetry.Logger
	Server   *http.Server
	Database *sql.DB
}

func NewApp() (*App, error) {
	cfg := config.Load()
	logger := telemetry.New(cfg.ServiceName)
	repositories, database, err := buildRepositories(cfg)
	if err != nil {
		return nil, err
	}
	attachmentGateway := buildAttachmentGateway(cfg)

	server := api.NewServer(cfg, logger, repositories, attachmentGateway)

	return &App{
		Config:   cfg,
		Logger:   logger,
		Server:   server,
		Database: database,
	}, nil
}

func (app *App) Run() error {
	app.Logger.Printf("starting %s on %s", app.Config.ServiceName, app.Config.HTTPAddress)
	return app.Server.ListenAndServe()
}

func buildRepositories(cfg config.Config) (repository.TenantRepositoryFactory, *sql.DB, error) {
	if cfg.RepositoryDriver != "postgres" {
		return persistence.NewInMemoryTenantRepositoryFactory(cfg.BootstrapTenantSlug), nil, nil
	}

	database, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, nil, err
	}

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, nil, err
	}

	return persistence.NewPostgresTenantRepositoryFactory(database, cfg.BootstrapTenantSlug), database, nil
}

func buildAttachmentGateway(cfg config.Config) repository.AttachmentGateway {
	if cfg.RepositoryDriver != "postgres" {
		return integration.NewInMemoryDocumentsGateway()
	}

	return integration.NewHTTPDocumentsGateway(cfg.DocumentsBaseURL)
}
