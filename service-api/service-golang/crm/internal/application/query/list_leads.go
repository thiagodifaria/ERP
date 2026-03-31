// ListLeads expõe a leitura minima de leads no bootstrap do CRM.
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

func (useCase ListLeads) Execute() []entity.Lead {
  return useCase.leadRepository.List()
}
