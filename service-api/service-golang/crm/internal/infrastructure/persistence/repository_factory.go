package persistence

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type InMemoryTenantRepositoryFactory struct {
	bootstrapTenantSlug string
	mutex               sync.Mutex
	bundles             map[string]repository.TenantRepositorySet
}

func NewInMemoryTenantRepositoryFactory(bootstrapTenantSlug string) *InMemoryTenantRepositoryFactory {
	return &InMemoryTenantRepositoryFactory{
		bootstrapTenantSlug: normalizeCrmTenantSlug(bootstrapTenantSlug),
		bundles:             map[string]repository.TenantRepositorySet{},
	}
}

func (factory *InMemoryTenantRepositoryFactory) BootstrapTenantSlug() string {
	return factory.bootstrapTenantSlug
}

func (factory *InMemoryTenantRepositoryFactory) ForTenant(tenantSlug string) (repository.TenantRepositorySet, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	slug := normalizeCrmTenantSlug(tenantSlug)
	if slug == "" {
		slug = factory.bootstrapTenantSlug
	}

	if bundle, ok := factory.bundles[slug]; ok {
		return bundle, nil
	}

	bundle := repository.TenantRepositorySet{
		LeadRepository:              NewInMemoryLeadRepository(slug),
		LeadNoteRepository:          NewInMemoryLeadNoteRepository(slug),
		CustomerRepository:          NewInMemoryCustomerRepository(slug),
		RelationshipEventRepository: NewInMemoryRelationshipEventRepository(slug),
		OutboxEventRepository:       NewInMemoryOutboxEventRepository(slug),
	}
	factory.bundles[slug] = bundle
	return bundle, nil
}

type PostgresTenantRepositoryFactory struct {
	database            *sql.DB
	bootstrapTenantSlug string
}

func NewPostgresTenantRepositoryFactory(database *sql.DB, bootstrapTenantSlug string) *PostgresTenantRepositoryFactory {
	return &PostgresTenantRepositoryFactory{
		database:            database,
		bootstrapTenantSlug: normalizeCrmTenantSlug(bootstrapTenantSlug),
	}
}

func (factory *PostgresTenantRepositoryFactory) BootstrapTenantSlug() string {
	return factory.bootstrapTenantSlug
}

func (factory *PostgresTenantRepositoryFactory) ForTenant(tenantSlug string) (repository.TenantRepositorySet, error) {
	slug := normalizeCrmTenantSlug(tenantSlug)
	if strings.TrimSpace(slug) == "" {
		slug = factory.bootstrapTenantSlug
	}

	leadRepository, err := NewPostgresLeadRepository(factory.database, slug)
	if err != nil {
		return repository.TenantRepositorySet{}, err
	}

	leadNoteRepository, err := NewPostgresLeadNoteRepository(factory.database, slug)
	if err != nil {
		return repository.TenantRepositorySet{}, err
	}

	customerRepository, err := NewPostgresCustomerRepository(factory.database, slug)
	if err != nil {
		return repository.TenantRepositorySet{}, err
	}

	eventRepository, err := NewPostgresRelationshipEventRepository(factory.database, slug)
	if err != nil {
		return repository.TenantRepositorySet{}, err
	}

	outboxRepository, err := NewPostgresOutboxEventRepository(factory.database, slug)
	if err != nil {
		return repository.TenantRepositorySet{}, err
	}

	return repository.TenantRepositorySet{
		LeadRepository:              leadRepository,
		LeadNoteRepository:          leadNoteRepository,
		CustomerRepository:          customerRepository,
		RelationshipEventRepository: eventRepository,
		OutboxEventRepository:       outboxRepository,
	}, nil
}
