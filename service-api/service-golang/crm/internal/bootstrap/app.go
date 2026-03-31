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
	leadRepository, leadNoteRepository, database, err := buildRepositories(cfg)
	if err != nil {
		return nil, err
	}

	server := api.NewServer(cfg, logger, leadRepository, leadNoteRepository)

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

func buildRepositories(cfg config.Config) (repository.LeadRepository, repository.LeadNoteRepository, *sql.DB, error) {
	if cfg.RepositoryDriver != "postgres" {
		return persistence.NewInMemoryLeadRepository(), persistence.NewInMemoryLeadNoteRepository(), nil, nil
	}

	database, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, nil, nil, err
	}

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, nil, nil, err
	}

	leadRepository, err := persistence.NewPostgresLeadRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, err
	}

	return leadRepository, persistence.NewInMemoryLeadNoteRepository(), database, nil
}
