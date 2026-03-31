// CreateLead concentra a criacao inicial de leads no bootstrap do CRM.
package command

import (
	"crypto/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type CreateLead struct {
	leadRepository repository.LeadRepository
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

func NewCreateLead(leadRepository repository.LeadRepository) CreateLead {
	return CreateLead{leadRepository: leadRepository}
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

	return CreateLeadResult{Lead: &createdLead}
}

func newPublicID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return uuid.Nil.String()
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	return uuid.UUID(raw).String()
}
