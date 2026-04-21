package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"

type CustomerRepository interface {
	List() []entity.Customer
	FindByPublicID(publicID string) *entity.Customer
	FindByEmail(email string) *entity.Customer
	Save(customer entity.Customer) entity.Customer
}

type TenantRepositorySet struct {
	LeadRepository     LeadRepository
	LeadNoteRepository LeadNoteRepository
	CustomerRepository CustomerRepository
}

type TenantRepositoryFactory interface {
	BootstrapTenantSlug() string
	ForTenant(tenantSlug string) (TenantRepositorySet, error)
}
