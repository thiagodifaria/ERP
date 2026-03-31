// CreateLead concentra a criacao inicial de leads no bootstrap do CRM.
package command

import (
  "crypto/rand"
  "encoding/hex"
  "strings"

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
        ErrorCode: "invalid_lead_name",
        ErrorText: "Lead name is required.",
        BadRequest: true,
      }
    case entity.ErrLeadEmailInvalid:
      return CreateLeadResult{
        ErrorCode: "invalid_lead_email",
        ErrorText: "Lead email is invalid.",
        BadRequest: true,
      }
    default:
      return CreateLeadResult{
        ErrorCode: "invalid_lead",
        ErrorText: "Lead payload is invalid.",
        BadRequest: true,
      }
    }
  }

  createdLead := useCase.leadRepository.Save(lead)

  return CreateLeadResult{Lead: &createdLead}
}

func newPublicID() string {
  raw := make([]byte, 8)
  if _, err := rand.Read(raw); err != nil {
    return "lead-fallback"
  }

  return "lead-" + hex.EncodeToString(raw)
}
