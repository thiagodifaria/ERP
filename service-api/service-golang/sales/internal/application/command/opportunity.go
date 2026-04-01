// Commands do ciclo de oportunidade do contexto sales.
package command

import (
	"errors"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type CreateOpportunity struct {
	opportunityRepository repository.OpportunityRepository
}

type CreateOpportunityInput struct {
	LeadPublicID string
	Title        string
	OwnerUserID  string
	AmountCents  int64
}

type CreateOpportunityResult struct {
	Opportunity *entity.Opportunity
	ErrorCode   string
	ErrorText   string
	BadRequest  bool
}

func NewCreateOpportunity(opportunityRepository repository.OpportunityRepository) CreateOpportunity {
	return CreateOpportunity{opportunityRepository: opportunityRepository}
}

func (useCase CreateOpportunity) Execute(input CreateOpportunityInput) CreateOpportunityResult {
	opportunity, err := entity.NewOpportunity(newPublicID(), input.LeadPublicID, input.Title, input.OwnerUserID, input.AmountCents)
	if err != nil {
		return mapOpportunityValidationError(err)
	}

	created := useCase.opportunityRepository.Save(opportunity)
	return CreateOpportunityResult{Opportunity: &created}
}

type UpdateOpportunityProfile struct {
	opportunityRepository repository.OpportunityRepository
}

type UpdateOpportunityProfileInput struct {
	PublicID    string
	Title       string
	OwnerUserID string
	AmountCents int64
}

type UpdateOpportunityProfileResult struct {
	Opportunity *entity.Opportunity
	ErrorCode   string
	ErrorText   string
	BadRequest  bool
	NotFound    bool
}

func NewUpdateOpportunityProfile(opportunityRepository repository.OpportunityRepository) UpdateOpportunityProfile {
	return UpdateOpportunityProfile{opportunityRepository: opportunityRepository}
}

func (useCase UpdateOpportunityProfile) Execute(input UpdateOpportunityProfileInput) UpdateOpportunityProfileResult {
	opportunity := useCase.opportunityRepository.FindByPublicID(input.PublicID)
	if opportunity == nil {
		return UpdateOpportunityProfileResult{
			ErrorCode: "opportunity_not_found",
			ErrorText: "Opportunity was not found.",
			NotFound:  true,
		}
	}

	revised, err := opportunity.Revise(input.Title, input.OwnerUserID, input.AmountCents)
	if err != nil {
		validation := mapOpportunityValidationError(err)
		return UpdateOpportunityProfileResult{
			ErrorCode:  validation.ErrorCode,
			ErrorText:  validation.ErrorText,
			BadRequest: validation.BadRequest,
		}
	}

	updated := useCase.opportunityRepository.Update(revised)
	return UpdateOpportunityProfileResult{Opportunity: &updated}
}

type UpdateOpportunityStage struct {
	opportunityRepository repository.OpportunityRepository
}

type UpdateOpportunityStageInput struct {
	PublicID string
	Stage    string
}

type UpdateOpportunityStageResult struct {
	Opportunity *entity.Opportunity
	ErrorCode   string
	ErrorText   string
	BadRequest  bool
	NotFound    bool
}

func NewUpdateOpportunityStage(opportunityRepository repository.OpportunityRepository) UpdateOpportunityStage {
	return UpdateOpportunityStage{opportunityRepository: opportunityRepository}
}

func (useCase UpdateOpportunityStage) Execute(input UpdateOpportunityStageInput) UpdateOpportunityStageResult {
	opportunity := useCase.opportunityRepository.FindByPublicID(input.PublicID)
	if opportunity == nil {
		return UpdateOpportunityStageResult{
			ErrorCode: "opportunity_not_found",
			ErrorText: "Opportunity was not found.",
			NotFound:  true,
		}
	}

	updatedOpportunity, err := opportunity.TransitionTo(input.Stage)
	if err != nil {
		switch err {
		case entity.ErrOpportunityStageInvalid:
			return UpdateOpportunityStageResult{
				ErrorCode:  "invalid_opportunity_stage",
				ErrorText:  "Opportunity stage is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateOpportunityStageResult{
				ErrorCode:  "invalid_opportunity_stage_transition",
				ErrorText:  "Opportunity stage transition is invalid.",
				BadRequest: true,
			}
		}
	}

	saved := useCase.opportunityRepository.Update(updatedOpportunity)
	return UpdateOpportunityStageResult{Opportunity: &saved}
}

func mapOpportunityValidationError(err error) CreateOpportunityResult {
	switch {
	case errors.Is(err, entity.ErrOpportunityLeadPublicIDInvalid):
		return CreateOpportunityResult{
			ErrorCode:  "invalid_lead_public_id",
			ErrorText:  "Lead public id is invalid.",
			BadRequest: true,
		}
	case errors.Is(err, entity.ErrOpportunityTitleRequired):
		return CreateOpportunityResult{
			ErrorCode:  "invalid_opportunity_title",
			ErrorText:  "Opportunity title is required.",
			BadRequest: true,
		}
	case errors.Is(err, entity.ErrOpportunityOwnerUserIDInvalid):
		return CreateOpportunityResult{
			ErrorCode:  "invalid_owner_user_id",
			ErrorText:  "Opportunity owner user id is invalid.",
			BadRequest: true,
		}
	case errors.Is(err, entity.ErrOpportunityAmountCentsInvalid):
		return CreateOpportunityResult{
			ErrorCode:  "invalid_opportunity_amount",
			ErrorText:  "Opportunity amount cents must be greater than zero.",
			BadRequest: true,
		}
	default:
		return CreateOpportunityResult{
			ErrorCode:  "invalid_opportunity",
			ErrorText:  "Opportunity payload is invalid.",
			BadRequest: true,
		}
	}
}
