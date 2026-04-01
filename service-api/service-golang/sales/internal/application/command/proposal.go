// Commands do ciclo de propostas do contexto sales.
package command

import (
	"errors"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type CreateProposal struct {
	opportunityRepository repository.OpportunityRepository
	proposalRepository    repository.ProposalRepository
}

type CreateProposalInput struct {
	OpportunityPublicID string
	Title               string
	AmountCents         int64
}

type CreateProposalResult struct {
	Proposal   *entity.Proposal
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewCreateProposal(opportunityRepository repository.OpportunityRepository, proposalRepository repository.ProposalRepository) CreateProposal {
	return CreateProposal{
		opportunityRepository: opportunityRepository,
		proposalRepository:    proposalRepository,
	}
}

func (useCase CreateProposal) Execute(input CreateProposalInput) CreateProposalResult {
	opportunity := useCase.opportunityRepository.FindByPublicID(input.OpportunityPublicID)
	if opportunity == nil {
		return CreateProposalResult{
			ErrorCode: "opportunity_not_found",
			ErrorText: "Opportunity was not found.",
			NotFound:  true,
		}
	}

	proposal, err := entity.NewProposal(newPublicID(), input.OpportunityPublicID, input.Title, input.AmountCents)
	if err != nil {
		return mapProposalValidationError(err)
	}

	if opportunity.Stage == "qualified" {
		updatedOpportunity, transitionErr := opportunity.TransitionTo("proposal")
		if transitionErr == nil {
			useCase.opportunityRepository.Update(updatedOpportunity)
		}
	}

	created := useCase.proposalRepository.Save(proposal)
	return CreateProposalResult{Proposal: &created}
}

type UpdateProposalStatus struct {
	proposalRepository repository.ProposalRepository
}

type UpdateProposalStatusInput struct {
	PublicID string
	Status   string
}

type UpdateProposalStatusResult struct {
	Proposal   *entity.Proposal
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewUpdateProposalStatus(proposalRepository repository.ProposalRepository) UpdateProposalStatus {
	return UpdateProposalStatus{proposalRepository: proposalRepository}
}

func (useCase UpdateProposalStatus) Execute(input UpdateProposalStatusInput) UpdateProposalStatusResult {
	proposal := useCase.proposalRepository.FindByPublicID(input.PublicID)
	if proposal == nil {
		return UpdateProposalStatusResult{
			ErrorCode: "proposal_not_found",
			ErrorText: "Proposal was not found.",
			NotFound:  true,
		}
	}

	updatedProposal, err := proposal.TransitionTo(input.Status)
	if err != nil {
		switch err {
		case entity.ErrProposalStatusInvalid:
			return UpdateProposalStatusResult{
				ErrorCode:  "invalid_proposal_status",
				ErrorText:  "Proposal status is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateProposalStatusResult{
				ErrorCode:  "invalid_proposal_status_transition",
				ErrorText:  "Proposal status transition is invalid.",
				BadRequest: true,
			}
		}
	}

	saved := useCase.proposalRepository.Update(updatedProposal)
	return UpdateProposalStatusResult{Proposal: &saved}
}

func mapProposalValidationError(err error) CreateProposalResult {
	switch {
	case errors.Is(err, entity.ErrProposalTitleRequired):
		return CreateProposalResult{
			ErrorCode:  "invalid_proposal_title",
			ErrorText:  "Proposal title is required.",
			BadRequest: true,
		}
	case errors.Is(err, entity.ErrProposalAmountCentsInvalid):
		return CreateProposalResult{
			ErrorCode:  "invalid_proposal_amount",
			ErrorText:  "Proposal amount cents must be greater than zero.",
			BadRequest: true,
		}
	default:
		return CreateProposalResult{
			ErrorCode:  "invalid_proposal",
			ErrorText:  "Proposal payload is invalid.",
			BadRequest: true,
		}
	}
}
