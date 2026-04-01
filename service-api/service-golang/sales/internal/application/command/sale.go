// Commands do ciclo de fechamento comercial do contexto sales.
package command

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type ConvertProposalToSale struct {
	opportunityRepository repository.OpportunityRepository
	proposalRepository    repository.ProposalRepository
	saleRepository        repository.SaleRepository
}

type ConvertProposalToSaleInput struct {
	ProposalPublicID string
}

type ConvertProposalToSaleResult struct {
	Sale       *entity.Sale
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
	Conflict   bool
}

func NewConvertProposalToSale(
	opportunityRepository repository.OpportunityRepository,
	proposalRepository repository.ProposalRepository,
	saleRepository repository.SaleRepository,
) ConvertProposalToSale {
	return ConvertProposalToSale{
		opportunityRepository: opportunityRepository,
		proposalRepository:    proposalRepository,
		saleRepository:        saleRepository,
	}
}

func (useCase ConvertProposalToSale) Execute(input ConvertProposalToSaleInput) ConvertProposalToSaleResult {
	proposal := useCase.proposalRepository.FindByPublicID(input.ProposalPublicID)
	if proposal == nil {
		return ConvertProposalToSaleResult{
			ErrorCode: "proposal_not_found",
			ErrorText: "Proposal was not found.",
			NotFound:  true,
		}
	}

	if existing := useCase.saleRepository.FindByProposalPublicID(proposal.PublicID); existing != nil {
		return ConvertProposalToSaleResult{
			ErrorCode: "sale_already_exists_for_proposal",
			ErrorText: "Proposal was already converted to sale.",
			Conflict:  true,
		}
	}

	if proposal.Status == "draft" || proposal.Status == "rejected" {
		return ConvertProposalToSaleResult{
			ErrorCode:  "proposal_not_convertible",
			ErrorText:  "Proposal must be sent or accepted before conversion.",
			BadRequest: true,
		}
	}

	opportunity := useCase.opportunityRepository.FindByPublicID(proposal.OpportunityPublicID)
	if opportunity == nil {
		return ConvertProposalToSaleResult{
			ErrorCode: "opportunity_not_found",
			ErrorText: "Opportunity was not found.",
			NotFound:  true,
		}
	}

	if proposal.Status == "sent" {
		acceptedProposal, err := proposal.TransitionTo("accepted")
		if err != nil {
			return ConvertProposalToSaleResult{
				ErrorCode:  "proposal_not_convertible",
				ErrorText:  "Proposal must be sent or accepted before conversion.",
				BadRequest: true,
			}
		}

		updatedProposal := useCase.proposalRepository.Update(acceptedProposal)
		proposal = &updatedProposal
	}

	targetOpportunity := *opportunity
	var err error
	switch opportunity.Stage {
	case "qualified":
		targetOpportunity, err = targetOpportunity.TransitionTo("proposal")
		if err == nil {
			targetOpportunity, err = targetOpportunity.TransitionTo("won")
		}
	case "proposal", "negotiation":
		targetOpportunity, err = targetOpportunity.TransitionTo("won")
	case "won":
		err = nil
	default:
		err = entity.ErrOpportunityStageTransitionInvalid
	}

	if err != nil {
		return ConvertProposalToSaleResult{
			ErrorCode:  "opportunity_not_convertible",
			ErrorText:  "Opportunity stage is not eligible for conversion.",
			BadRequest: true,
		}
	}

	useCase.opportunityRepository.Update(targetOpportunity)

	sale, err := entity.NewSale(newPublicID(), targetOpportunity.PublicID, proposal.PublicID, proposal.AmountCents)
	if err != nil {
		return ConvertProposalToSaleResult{
			ErrorCode:  "invalid_sale",
			ErrorText:  "Sale payload is invalid.",
			BadRequest: true,
		}
	}

	created := useCase.saleRepository.Save(sale)
	return ConvertProposalToSaleResult{Sale: &created}
}

type UpdateSaleStatus struct {
	saleRepository repository.SaleRepository
}

type UpdateSaleStatusInput struct {
	PublicID string
	Status   string
}

type UpdateSaleStatusResult struct {
	Sale       *entity.Sale
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewUpdateSaleStatus(saleRepository repository.SaleRepository) UpdateSaleStatus {
	return UpdateSaleStatus{saleRepository: saleRepository}
}

func (useCase UpdateSaleStatus) Execute(input UpdateSaleStatusInput) UpdateSaleStatusResult {
	sale := useCase.saleRepository.FindByPublicID(input.PublicID)
	if sale == nil {
		return UpdateSaleStatusResult{
			ErrorCode: "sale_not_found",
			ErrorText: "Sale was not found.",
			NotFound:  true,
		}
	}

	updatedSale, err := sale.TransitionTo(input.Status)
	if err != nil {
		switch err {
		case entity.ErrSaleStatusInvalid:
			return UpdateSaleStatusResult{
				ErrorCode:  "invalid_sale_status",
				ErrorText:  "Sale status is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateSaleStatusResult{
				ErrorCode:  "invalid_sale_status_transition",
				ErrorText:  "Sale status transition is invalid.",
				BadRequest: true,
			}
		}
	}

	saved := useCase.saleRepository.Update(updatedSale)
	return UpdateSaleStatusResult{Sale: &saved}
}
