package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ListCustomers struct {
	customerRepository repository.CustomerRepository
}

func NewListCustomers(customerRepository repository.CustomerRepository) ListCustomers {
	return ListCustomers{customerRepository: customerRepository}
}

func (useCase ListCustomers) Execute() []entity.Customer {
	return useCase.customerRepository.List()
}
