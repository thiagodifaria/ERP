// UpdateLeadOwner handles direct lead ownership changes in the CRM bootstrap.
package command

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type UpdateLeadOwner struct {
	leadRepository repository.LeadRepository
}

type UpdateLeadOwnerInput struct {
	PublicID    string
	OwnerUserID string
}

type UpdateLeadOwnerResult struct {
	Lead      *entity.Lead
	ErrorCode string
	ErrorText string
	NotFound  bool
}

func NewUpdateLeadOwner(leadRepository repository.LeadRepository) UpdateLeadOwner {
	return UpdateLeadOwner{leadRepository: leadRepository}
}

func (useCase UpdateLeadOwner) Execute(input UpdateLeadOwnerInput) UpdateLeadOwnerResult {
	publicID := strings.TrimSpace(input.PublicID)
	lead := useCase.leadRepository.FindByPublicID(publicID)
	if lead == nil {
		return UpdateLeadOwnerResult{
			ErrorCode: "lead_not_found",
			ErrorText: "Lead was not found.",
			NotFound:  true,
		}
	}

	updatedLead := useCase.leadRepository.Update(lead.AssignOwner(input.OwnerUserID))

	return UpdateLeadOwnerResult{Lead: &updatedLead}
}
