package persistence

import (
	"strings"
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

const BootstrapCustomerPublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000081"

type InMemoryCustomerRepository struct {
	sync.Mutex
	customers []entity.Customer
}

func NewInMemoryCustomerRepository() *InMemoryCustomerRepository {
	lead, _ := entity.RestoreLead(
		BootstrapLeadPublicID,
		"Bootstrap Lead",
		"lead@bootstrap-ops.local",
		"manual",
		"qualified",
		BootstrapOwnerUserPublicID,
	)
	customer, _ := entity.NewCustomerFromLead(BootstrapCustomerPublicID, lead)

	return &InMemoryCustomerRepository{
		customers: []entity.Customer{customer},
	}
}

func (repository *InMemoryCustomerRepository) List() []entity.Customer {
	repository.Lock()
	defer repository.Unlock()

	copied := make([]entity.Customer, len(repository.customers))
	copy(copied, repository.customers)
	return copied
}

func (repository *InMemoryCustomerRepository) FindByPublicID(publicID string) *entity.Customer {
	repository.Lock()
	defer repository.Unlock()

	for _, customer := range repository.customers {
		if customer.PublicID == strings.TrimSpace(publicID) {
			copied := customer
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryCustomerRepository) FindByEmail(email string) *entity.Customer {
	repository.Lock()
	defer repository.Unlock()

	target := strings.ToLower(strings.TrimSpace(email))
	for _, customer := range repository.customers {
		if customer.Email == target {
			copied := customer
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryCustomerRepository) Save(customer entity.Customer) entity.Customer {
	repository.Lock()
	defer repository.Unlock()

	repository.customers = append(repository.customers, customer)
	return customer
}
