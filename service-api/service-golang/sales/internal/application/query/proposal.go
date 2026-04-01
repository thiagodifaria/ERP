// Queries de leitura para propostas do contexto sales.
package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type ListProposalsByOpportunity struct {
	proposalRepository repository.ProposalRepository
}

type GetProposalByPublicID struct {
	proposalRepository repository.ProposalRepository
}

func NewListProposalsByOpportunity(proposalRepository repository.ProposalRepository) ListProposalsByOpportunity {
	return ListProposalsByOpportunity{proposalRepository: proposalRepository}
}

func (useCase ListProposalsByOpportunity) Execute(opportunityPublicID string) []entity.Proposal {
	return useCase.proposalRepository.ListByOpportunityPublicID(opportunityPublicID)
}

func NewGetProposalByPublicID(proposalRepository repository.ProposalRepository) GetProposalByPublicID {
	return GetProposalByPublicID{proposalRepository: proposalRepository}
}

func (useCase GetProposalByPublicID) Execute(publicID string) *entity.Proposal {
	return useCase.proposalRepository.FindByPublicID(publicID)
}
