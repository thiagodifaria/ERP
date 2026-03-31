// GetLeadByPublicID expõe a leitura individual de leads em bootstrap.
package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type GetLeadByPublicID struct {
	leadRepository repository.LeadRepository
}

func NewGetLeadByPublicID(leadRepository repository.LeadRepository) GetLeadByPublicID {
	return GetLeadByPublicID{leadRepository: leadRepository}
}

func (useCase GetLeadByPublicID) Execute(publicID string) *entity.Lead {
	return useCase.leadRepository.FindByPublicID(publicID)
}
