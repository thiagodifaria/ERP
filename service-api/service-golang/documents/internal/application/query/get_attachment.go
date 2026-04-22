package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type GetAttachment struct {
	repository repository.AttachmentRepository
}

func NewGetAttachment(repository repository.AttachmentRepository) GetAttachment {
	return GetAttachment{repository: repository}
}

func (useCase GetAttachment) Execute(tenantSlug string, publicID string) (*entity.Attachment, bool) {
	attachment, ok := useCase.repository.FindByPublicID(strings.TrimSpace(tenantSlug), strings.TrimSpace(publicID))
	if !ok {
		return nil, false
	}

	return &attachment, true
}
