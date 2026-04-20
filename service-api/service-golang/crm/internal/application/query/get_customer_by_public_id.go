package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type GetCustomerByPublicID struct {
	customerRepository repository.CustomerRepository
}

func NewGetCustomerByPublicID(customerRepository repository.CustomerRepository) GetCustomerByPublicID {
	return GetCustomerByPublicID{customerRepository: customerRepository}
}

func (useCase GetCustomerByPublicID) Execute(publicID string) *entity.Customer {
	return useCase.customerRepository.FindByPublicID(publicID)
}
