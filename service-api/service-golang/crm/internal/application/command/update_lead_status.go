// UpdateLeadStatus concentra a transicao minima de status do lead no bootstrap do CRM.
package command

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type UpdateLeadStatus struct {
	leadRepository repository.LeadRepository
}

type UpdateLeadStatusInput struct {
	PublicID string
	Status   string
}

type UpdateLeadStatusResult struct {
	Lead       *entity.Lead
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewUpdateLeadStatus(leadRepository repository.LeadRepository) UpdateLeadStatus {
	return UpdateLeadStatus{leadRepository: leadRepository}
}

func (useCase UpdateLeadStatus) Execute(input UpdateLeadStatusInput) UpdateLeadStatusResult {
	publicID := strings.TrimSpace(input.PublicID)
	lead := useCase.leadRepository.FindByPublicID(publicID)
	if lead == nil {
		return UpdateLeadStatusResult{
			ErrorCode: "lead_not_found",
			ErrorText: "Lead was not found.",
			NotFound:  true,
		}
	}

	transitionedLead, err := lead.TransitionTo(input.Status)
	if err != nil {
		switch err {
		case entity.ErrLeadStatusInvalid:
			return UpdateLeadStatusResult{
				ErrorCode:  "invalid_lead_status",
				ErrorText:  "Lead status is invalid.",
				BadRequest: true,
			}
		case entity.ErrLeadStatusTransitionInvalid:
			return UpdateLeadStatusResult{
				ErrorCode:  "invalid_lead_status_transition",
				ErrorText:  "Lead status transition is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateLeadStatusResult{
				ErrorCode:  "invalid_lead_status",
				ErrorText:  "Lead status is invalid.",
				BadRequest: true,
			}
		}
	}

	updatedLead := useCase.leadRepository.Update(transitionedLead)

	return UpdateLeadStatusResult{Lead: &updatedLead}
}
