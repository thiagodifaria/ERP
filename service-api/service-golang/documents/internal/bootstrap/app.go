package bootstrap

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
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
	attachmentRepository, uploadSessionRepository, database, err := buildRepository(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:   cfg,
		Logger:   logger,
		Server:   api.NewServer(cfg, logger, attachmentRepository, uploadSessionRepository),
		Database: database,
	}, nil
}

func (app *App) Run() error {
	app.Logger.Printf("starting %s on %s", app.Config.ServiceName, app.Config.HTTPAddress)
	return app.Server.ListenAndServe()
}

func buildRepository(cfg config.Config) (repository.AttachmentRepository, repository.UploadSessionRepository, *sql.DB, error) {
	if cfg.RepositoryDriver != "postgres" {
		return persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository(), nil, nil
	}

	database, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, nil, nil, err
	}

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, nil, nil, err
	}

	return persistence.NewPostgresAttachmentRepository(database, cfg.BootstrapTenantSlug), persistence.NewPostgresUploadSessionRepository(database, cfg.BootstrapTenantSlug), database, nil
}
