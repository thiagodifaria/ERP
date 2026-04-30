package persistence

import (
	"strings"
	"sync"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
)

type InMemoryUploadSessionRepository struct {
	mutex    sync.Mutex
	sessions []entity.UploadSession
}

func NewInMemoryUploadSessionRepository() *InMemoryUploadSessionRepository {
	return &InMemoryUploadSessionRepository{
		sessions: []entity.UploadSession{},
	}
}

func (repository *InMemoryUploadSessionRepository) Save(session entity.UploadSession) entity.UploadSession {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.sessions = append(repository.sessions, session)
	return session
}

func (repository *InMemoryUploadSessionRepository) FindByPublicID(tenantSlug string, publicID string) (entity.UploadSession, bool) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedPublicID := strings.TrimSpace(publicID)
	for _, session := range repository.sessions {
		if normalizedTenantSlug != "" && session.TenantSlug != normalizedTenantSlug {
			continue
		}
		if session.PublicID == normalizedPublicID {
			return session, true
		}
	}

	return entity.UploadSession{}, false
}

func (repository *InMemoryUploadSessionRepository) Complete(tenantSlug string, publicID string, attachmentPublicID string, completedAt time.Time) (entity.UploadSession, bool) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedPublicID := strings.TrimSpace(publicID)
	for index, session := range repository.sessions {
		if normalizedTenantSlug != "" && session.TenantSlug != normalizedTenantSlug {
			continue
		}
		if session.PublicID != normalizedPublicID {
			continue
		}

		updated := session.Complete(strings.TrimSpace(attachmentPublicID), completedAt)
		repository.sessions[index] = updated
		return updated, true
	}

	return entity.UploadSession{}, false
}
