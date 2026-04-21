// InMemoryLeadRepository fornece leads de bootstrap ate a persistencia real entrar.
// A API publica nao deve depender da implementacao concreta.
package persistence

import (
	"fmt"
	"strings"
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type InMemoryLeadRepository struct {
	sync.Mutex
	leads []entity.Lead
}

const (
	BootstrapLeadPublicID      = "0195e7a0-7a9c-7c1f-8a44-4a6e70000001"
	BootstrapOwnerUserPublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000011"
)

func NewInMemoryLeadRepository(tenantSlug ...string) *InMemoryLeadRepository {
	repository := &InMemoryLeadRepository{}
	normalizedTenantSlug := normalizeCrmTenantSlug(firstTenantSlug(tenantSlug...))
	displayName := strings.ReplaceAll(normalizedTenantSlug, "-", " ")

	lead, _ := entity.NewLead(
		BootstrapLeadPublicID,
		fmt.Sprintf("%s Lead", displayName),
		fmt.Sprintf("lead@%s.local", normalizedTenantSlug),
		"manual",
		BootstrapOwnerUserPublicID,
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

func normalizeCrmTenantSlug(tenantSlug string) string {
	normalized := strings.ToLower(strings.TrimSpace(tenantSlug))
	if normalized == "" {
		return "bootstrap-ops"
	}

	return normalized
}

func firstTenantSlug(tenantSlug ...string) string {
	if len(tenantSlug) == 0 {
		return ""
	}

	return tenantSlug[0]
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
