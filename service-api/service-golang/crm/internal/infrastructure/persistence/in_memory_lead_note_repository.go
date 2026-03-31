// InMemoryLeadNoteRepository entrega o primeiro historico operacional enquanto a persistencia relacional amadurece.
package persistence

import (
	"slices"
	"sync"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

const BootstrapLeadNotePublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000032"

type InMemoryLeadNoteRepository struct {
	sync.Mutex
	notesByLeadPublicID map[string][]entity.LeadNote
}

func NewInMemoryLeadNoteRepository() *InMemoryLeadNoteRepository {
	repository := &InMemoryLeadNoteRepository{
		notesByLeadPublicID: map[string][]entity.LeadNote{},
	}

	bootstrapNote, _ := entity.NewLeadNote(
		BootstrapLeadNotePublicID,
		BootstrapLeadPublicID,
		"Primeiro contato capturado e aguardando abordagem comercial.",
		"qualification",
		time.Date(2026, time.March, 31, 14, 0, 0, 0, time.UTC),
	)
	repository.notesByLeadPublicID[BootstrapLeadPublicID] = []entity.LeadNote{bootstrapNote}

	return repository
}

func (repository *InMemoryLeadNoteRepository) ListByLeadPublicID(leadPublicID string) []entity.LeadNote {
	repository.Lock()
	defer repository.Unlock()

	notes := repository.notesByLeadPublicID[leadPublicID]
	return slices.Clone(notes)
}

func (repository *InMemoryLeadNoteRepository) Save(note entity.LeadNote) entity.LeadNote {
	repository.Lock()
	defer repository.Unlock()

	repository.notesByLeadPublicID[note.LeadPublicID] = append(
		repository.notesByLeadPublicID[note.LeadPublicID],
		note,
	)

	return note
}
