package command

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type CreateLead struct {
	leadRepository   repository.LeadRepository
	eventRepository  repository.RelationshipEventRepository
	outboxRepository repository.OutboxEventRepository
}

type CreateLeadInput struct {
	Name        string
	Email       string
	Source      string
	OwnerUserID string
}

type CreateLeadResult struct {
	Lead       *entity.Lead
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	Conflict   bool
}

func NewCreateLead(
	leadRepository repository.LeadRepository,
	eventRepository repository.RelationshipEventRepository,
	outboxRepository repository.OutboxEventRepository,
) CreateLead {
	return CreateLead{
		leadRepository:   leadRepository,
		eventRepository:  eventRepository,
		outboxRepository: outboxRepository,
	}
}

func (useCase CreateLead) Execute(input CreateLeadInput) CreateLeadResult {
	normalizedEmail := strings.ToLower(strings.TrimSpace(input.Email))

	if existing := useCase.leadRepository.FindByEmail(normalizedEmail); existing != nil {
		return CreateLeadResult{
			ErrorCode: "lead_email_conflict",
			ErrorText: "Lead email already exists.",
			Conflict:  true,
		}
	}

	lead, err := entity.NewLead(
		newPublicID(),
		input.Name,
		normalizedEmail,
		input.Source,
		input.OwnerUserID,
	)
	if err != nil {
		switch err {
		case entity.ErrLeadNameRequired:
			return CreateLeadResult{
				ErrorCode:  "invalid_lead_name",
				ErrorText:  "Lead name is required.",
				BadRequest: true,
			}
		case entity.ErrLeadEmailInvalid:
			return CreateLeadResult{
				ErrorCode:  "invalid_lead_email",
				ErrorText:  "Lead email is invalid.",
				BadRequest: true,
			}
		case entity.ErrLeadOwnerUserIDInvalid:
			return CreateLeadResult{
				ErrorCode:  "invalid_owner_user_id",
				ErrorText:  "Lead owner user id is invalid.",
				BadRequest: true,
			}
		default:
			return CreateLeadResult{
				ErrorCode:  "invalid_lead",
				ErrorText:  "Lead payload is invalid.",
				BadRequest: true,
			}
		}
	}

	createdLead := useCase.leadRepository.Save(lead)
	recordRelationshipEvent(useCase.eventRepository, "lead", createdLead.PublicID, "lead_created", "crm", "Lead created in CRM.")
	appendOutboxEvent(useCase.outboxRepository, "lead", createdLead.PublicID, "crm.lead.created", map[string]any{
		"publicId":    createdLead.PublicID,
		"email":       createdLead.Email,
		"source":      createdLead.Source,
		"status":      createdLead.Status,
		"ownerUserId": createdLead.OwnerUserID,
	})

	return CreateLeadResult{Lead: &createdLead}
}
