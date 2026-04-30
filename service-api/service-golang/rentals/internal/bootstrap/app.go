package bootstrap

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/telemetry"
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
	contractRepository, database, err := buildRepository(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:   cfg,
		Logger:   logger,
		Server:   api.NewServer(cfg, logger, contractRepository, buildDocumentsGateway(cfg)),
		Database: database,
	}, nil
}

func (app *App) Run() error {
	app.Logger.Printf("starting %s on %s", app.Config.ServiceName, app.Config.HTTPAddress)
	return app.Server.ListenAndServe()
}

func buildRepository(cfg config.Config) (repository.ContractRepository, *sql.DB, error) {
	if cfg.RepositoryDriver != "postgres" {
		return persistence.NewInMemoryContractRepository(), nil, nil
	}

	database, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, nil, err
	}
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, nil, err
	}

	return persistence.NewPostgresContractRepository(database, cfg.BootstrapTenantSlug), database, nil
}

func buildDocumentsGateway(cfg config.Config) repository.AttachmentGateway {
	if cfg.DocumentsBaseURL == "" {
		return integration.NewInMemoryDocumentsGateway()
	}

	return integration.NewHTTPDocumentsGateway(cfg.DocumentsBaseURL)
}
