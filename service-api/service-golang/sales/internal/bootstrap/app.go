// App concentra o bootstrap do servico sales.
// Aqui ficam configuracao, logger e servidor HTTP.
package bootstrap

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
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
	opportunityRepository, proposalRepository, saleRepository, invoiceRepository, installmentRepository, commissionRepository, pendingItemRepository, renegotiationRepository, eventRepository, outboxRepository, database, err := buildRepositories(cfg)
	if err != nil {
		return nil, err
	}

	server := api.NewServer(cfg, logger, opportunityRepository, proposalRepository, saleRepository, invoiceRepository, installmentRepository, commissionRepository, pendingItemRepository, renegotiationRepository, eventRepository, outboxRepository)

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

func buildRepositories(cfg config.Config) (repository.OpportunityRepository, repository.ProposalRepository, repository.SaleRepository, repository.InvoiceRepository, repository.InstallmentRepository, repository.CommissionRepository, repository.PendingItemRepository, repository.RenegotiationRepository, repository.CommercialEventRepository, repository.OutboxEventRepository, *sql.DB, error) {
	if cfg.RepositoryDriver != "postgres" {
		return persistence.NewInMemoryOpportunityRepository(), persistence.NewInMemoryProposalRepository(), persistence.NewInMemorySaleRepository(), persistence.NewInMemoryInvoiceRepository(), persistence.NewInMemoryInstallmentRepository(), persistence.NewInMemoryCommissionRepository(), persistence.NewInMemoryPendingItemRepository(), persistence.NewInMemoryRenegotiationRepository(), persistence.NewInMemoryCommercialEventRepository(), persistence.NewInMemoryOutboxEventRepository(), nil, nil
	}

	database, err := sql.Open("pgx", cfg.PostgresDSN())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	opportunityRepository, err := persistence.NewPostgresOpportunityRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	proposalRepository, err := persistence.NewPostgresProposalRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	saleRepository, err := persistence.NewPostgresSaleRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	invoiceRepository, err := persistence.NewPostgresInvoiceRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	installmentRepository, err := persistence.NewPostgresInstallmentRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	commissionRepository, err := persistence.NewPostgresCommissionRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	pendingItemRepository, err := persistence.NewPostgresPendingItemRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	renegotiationRepository, err := persistence.NewPostgresRenegotiationRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	eventRepository, err := persistence.NewPostgresCommercialEventRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	outboxRepository, err := persistence.NewPostgresOutboxEventRepository(database, cfg.BootstrapTenantSlug)
	if err != nil {
		_ = database.Close()
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return opportunityRepository, proposalRepository, saleRepository, invoiceRepository, installmentRepository, commissionRepository, pendingItemRepository, renegotiationRepository, eventRepository, outboxRepository, database, nil
}
