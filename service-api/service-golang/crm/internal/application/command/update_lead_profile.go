// UpdateLeadProfile handles partial profile changes for an existing CRM lead.
package command

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type UpdateLeadProfile struct {
	leadRepository   repository.LeadRepository
	eventRepository  repository.RelationshipEventRepository
	outboxRepository repository.OutboxEventRepository
}

type UpdateLeadProfileInput struct {
	PublicID string
	Name     *string
	Email    *string
	Source   *string
}

type UpdateLeadProfileResult struct {
	Lead       *entity.Lead
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	Conflict   bool
	NotFound   bool
}

func NewUpdateLeadProfile(
	leadRepository repository.LeadRepository,
	eventRepository repository.RelationshipEventRepository,
	outboxRepository repository.OutboxEventRepository,
) UpdateLeadProfile {
	return UpdateLeadProfile{
		leadRepository:   leadRepository,
		eventRepository:  eventRepository,
		outboxRepository: outboxRepository,
	}
}

func (useCase UpdateLeadProfile) Execute(input UpdateLeadProfileInput) UpdateLeadProfileResult {
	publicID := strings.TrimSpace(input.PublicID)
	lead := useCase.leadRepository.FindByPublicID(publicID)
	if lead == nil {
		return UpdateLeadProfileResult{
			ErrorCode: "lead_not_found",
			ErrorText: "Lead was not found.",
			NotFound:  true,
		}
	}

	name := lead.Name
	if input.Name != nil {
		name = *input.Name
	}

	email := lead.Email
	if input.Email != nil {
		email = *input.Email
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	if existing := useCase.leadRepository.FindByEmail(normalizedEmail); existing != nil && existing.PublicID != lead.PublicID {
		return UpdateLeadProfileResult{
			ErrorCode: "lead_email_conflict",
			ErrorText: "Lead email already exists.",
			Conflict:  true,
		}
	}

	source := lead.Source
	if input.Source != nil {
		source = *input.Source
	}

	revisedLead, err := lead.ReviseProfile(name, normalizedEmail, source)
	if err != nil {
		switch err {
		case entity.ErrLeadNameRequired:
			return UpdateLeadProfileResult{
				ErrorCode:  "invalid_lead_name",
				ErrorText:  "Lead name is required.",
				BadRequest: true,
			}
		case entity.ErrLeadEmailInvalid:
			return UpdateLeadProfileResult{
				ErrorCode:  "invalid_lead_email",
				ErrorText:  "Lead email is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateLeadProfileResult{
				ErrorCode:  "invalid_lead_profile",
				ErrorText:  "Lead profile is invalid.",
				BadRequest: true,
			}
		}
	}

	updatedLead := useCase.leadRepository.Update(revisedLead)
	recordRelationshipEvent(useCase.eventRepository, "lead", updatedLead.PublicID, "lead_profile_updated", "crm", "Lead profile updated.")
	appendOutboxEvent(useCase.outboxRepository, "lead", updatedLead.PublicID, "crm.lead.updated", map[string]any{
		"publicId":    updatedLead.PublicID,
		"name":        updatedLead.Name,
		"email":       updatedLead.Email,
		"source":      updatedLead.Source,
		"ownerUserId": updatedLead.OwnerUserID,
		"status":      updatedLead.Status,
	})

	return UpdateLeadProfileResult{Lead: &updatedLead}
}
