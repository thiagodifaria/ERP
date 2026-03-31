// InMemoryLeadRepository fornece leads de bootstrap ate a persistencia real entrar.
// A API publica nao deve depender da implementacao concreta.
package persistence

import (
  "strings"
  "sync"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type InMemoryLeadRepository struct {
  sync.Mutex
  leads []entity.Lead
}

func NewInMemoryLeadRepository() *InMemoryLeadRepository {
  repository := &InMemoryLeadRepository{}

  lead, _ := entity.NewLead(
    "lead-bootstrap-ops",
    "Bootstrap Ops Lead",
    "lead@bootstrap-ops.local",
    "manual",
    "owner-bootstrap-ops",
  )
  repository.leads = []entity.Lead{lead}

  return repository
}

func (repository *InMemoryLeadRepository) List() []entity.Lead {
  repository.Lock()
  defer repository.Unlock()

  copied := make([]entity.Lead, len(repository.leads))
  copy(copied, repository.leads)

  return copied
}

func (repository *InMemoryLeadRepository) FindByPublicID(publicID string) *entity.Lead {
  repository.Lock()
  defer repository.Unlock()

  for _, lead := range repository.leads {
    if lead.PublicID == publicID {
      copied := lead
      return &copied
    }
  }

  return nil
}

func (repository *InMemoryLeadRepository) FindByEmail(email string) *entity.Lead {
  repository.Lock()
  defer repository.Unlock()

  normalizedEmail := strings.ToLower(strings.TrimSpace(email))

  for _, lead := range repository.leads {
    if lead.Email == normalizedEmail {
      copied := lead
      return &copied
    }
  }

  return nil
}

func (repository *InMemoryLeadRepository) Save(lead entity.Lead) entity.Lead {
  repository.Lock()
  defer repository.Unlock()

  repository.leads = append(repository.leads, lead)
  return lead
}

func (repository *InMemoryLeadRepository) Update(lead entity.Lead) entity.Lead {
  repository.Lock()
  defer repository.Unlock()

  for index, currentLead := range repository.leads {
    if currentLead.PublicID == lead.PublicID {
      repository.leads[index] = lead
      return lead
    }
  }

  repository.leads = append(repository.leads, lead)
  return lead
}
