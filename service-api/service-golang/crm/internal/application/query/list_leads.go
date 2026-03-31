// ListLeads exposes operational lead listing for the CRM bootstrap.
package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ListLeads struct {
	leadRepository repository.LeadRepository
}

func NewListLeads(leadRepository repository.LeadRepository) ListLeads {
	return ListLeads{leadRepository: leadRepository}
}

func (useCase ListLeads) Execute(filters LeadFilters) []entity.Lead {
	return applyLeadFilters(useCase.leadRepository.List(), filters)
}
