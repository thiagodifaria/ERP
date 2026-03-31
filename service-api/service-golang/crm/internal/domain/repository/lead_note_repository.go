// LeadNoteRepository define a persistencia minima das notas ligadas a um lead.
package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"

type LeadNoteRepository interface {
	ListByLeadPublicID(leadPublicID string) []entity.LeadNote
	Save(note entity.LeadNote) entity.LeadNote
}
