// ListLeadNotes entrega o historico operacional ligado a um lead especifico.
package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ListLeadNotes struct {
	leadNoteRepository repository.LeadNoteRepository
}

func NewListLeadNotes(leadNoteRepository repository.LeadNoteRepository) ListLeadNotes {
	return ListLeadNotes{leadNoteRepository: leadNoteRepository}
}

func (useCase ListLeadNotes) Execute(leadPublicID string) []entity.LeadNote {
	return useCase.leadNoteRepository.ListByLeadPublicID(leadPublicID)
}
