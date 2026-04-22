package persistence

import (
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type InMemoryAttachmentRepository struct {
	mutex       sync.Mutex
	attachments []entity.Attachment
}

func NewInMemoryAttachmentRepository() *InMemoryAttachmentRepository {
	return &InMemoryAttachmentRepository{
		attachments: []entity.Attachment{},
	}
}

func (repository *InMemoryAttachmentRepository) List(filters repository.AttachmentFilters) []entity.Attachment {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	response := make([]entity.Attachment, 0, len(repository.attachments))
	for _, attachment := range repository.attachments {
		if filters.TenantSlug != "" && attachment.TenantSlug != strings.ToLower(strings.TrimSpace(filters.TenantSlug)) {
			continue
		}
		if filters.OwnerType != "" && attachment.OwnerType != strings.TrimSpace(filters.OwnerType) {
			continue
		}
		if filters.OwnerPublicID != "" && attachment.OwnerPublicID != strings.TrimSpace(filters.OwnerPublicID) {
			continue
		}
		if filters.Source != "" && attachment.Source != strings.TrimSpace(filters.Source) {
			continue
		}
		if filters.Visibility != "" && attachment.Visibility != strings.ToLower(strings.TrimSpace(filters.Visibility)) {
			continue
		}
		if filters.Archived == "true" && attachment.ArchivedAt == nil {
			continue
		}
		if filters.Archived == "false" && attachment.ArchivedAt != nil {
			continue
		}

		response = append(response, attachment)
	}

	return slices.Clone(response)
}

func (repository *InMemoryAttachmentRepository) FindByPublicID(tenantSlug string, publicID string) (entity.Attachment, bool) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedPublicID := strings.TrimSpace(publicID)
	for _, attachment := range repository.attachments {
		if normalizedTenantSlug != "" && attachment.TenantSlug != normalizedTenantSlug {
			continue
		}
		if attachment.PublicID == normalizedPublicID {
			return attachment, true
		}
	}

	return entity.Attachment{}, false
}

func (repository *InMemoryAttachmentRepository) Save(attachment entity.Attachment) entity.Attachment {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.attachments = append(repository.attachments, attachment)
	return attachment
}

func (repository *InMemoryAttachmentRepository) Archive(tenantSlug string, publicID string, reason string) (entity.Attachment, bool) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedPublicID := strings.TrimSpace(publicID)
	for index, attachment := range repository.attachments {
		if normalizedTenantSlug != "" && attachment.TenantSlug != normalizedTenantSlug {
			continue
		}
		if attachment.PublicID != normalizedPublicID {
			continue
		}

		updated := attachment.Archive(reason, time.Now().UTC())
		repository.attachments[index] = updated
		return updated, true
	}

	return entity.Attachment{}, false
}
