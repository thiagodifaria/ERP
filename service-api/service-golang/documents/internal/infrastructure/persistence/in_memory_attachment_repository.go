package persistence

import (
	"slices"
	"strings"
	"sync"

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

		response = append(response, attachment)
	}

	return slices.Clone(response)
}

func (repository *InMemoryAttachmentRepository) Save(attachment entity.Attachment) entity.Attachment {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.attachments = append(repository.attachments, attachment)
	return attachment
}
