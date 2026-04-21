// CreateLeadNote registra novas notas operacionais vinculadas a um lead existente.
package command

import (
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type CreateLeadNote struct {
	leadRepository     repository.LeadRepository
	leadNoteRepository repository.LeadNoteRepository
	eventRepository    repository.RelationshipEventRepository
	outboxRepository   repository.OutboxEventRepository
}

type CreateLeadNoteInput struct {
	LeadPublicID string
	Body         string
	Category     string
}

type CreateLeadNoteResult struct {
	Note       *entity.LeadNote
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewCreateLeadNote(
	leadRepository repository.LeadRepository,
	leadNoteRepository repository.LeadNoteRepository,
	eventRepository repository.RelationshipEventRepository,
	outboxRepository repository.OutboxEventRepository,
) CreateLeadNote {
	return CreateLeadNote{
		leadRepository:     leadRepository,
		leadNoteRepository: leadNoteRepository,
		eventRepository:    eventRepository,
		outboxRepository:   outboxRepository,
	}
}

func (useCase CreateLeadNote) Execute(input CreateLeadNoteInput) CreateLeadNoteResult {
	lead := useCase.leadRepository.FindByPublicID(input.LeadPublicID)
	if lead == nil {
		return CreateLeadNoteResult{
			ErrorCode: "lead_not_found",
			ErrorText: "Lead was not found.",
			NotFound:  true,
		}
	}

	note, err := entity.NewLeadNote(
		newPublicID(),
		lead.PublicID,
		input.Body,
		input.Category,
		time.Now().UTC(),
	)
	if err != nil {
		switch err {
		case entity.ErrLeadNoteBodyRequired:
			return CreateLeadNoteResult{
				ErrorCode:  "invalid_note_body",
				ErrorText:  "Lead note body is required.",
				BadRequest: true,
			}
		default:
			return CreateLeadNoteResult{
				ErrorCode:  "invalid_lead_note",
				ErrorText:  "Lead note payload is invalid.",
				BadRequest: true,
			}
		}
	}

	createdNote := useCase.leadNoteRepository.Save(note)
	recordRelationshipEvent(useCase.eventRepository, "lead", lead.PublicID, "lead_interaction_recorded", "crm", "Lead interaction recorded in CRM.")
	appendOutboxEvent(useCase.outboxRepository, "lead", lead.PublicID, "crm.lead.interaction_recorded", map[string]any{
		"publicId":  lead.PublicID,
		"noteId":    createdNote.PublicID,
		"category":  createdNote.Category,
		"createdAt": createdNote.CreatedAt.UTC().Format(time.RFC3339),
	})
	return CreateLeadNoteResult{Note: &createdNote}
}
