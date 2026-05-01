package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type ListAttachmentVersions struct {
	repository repository.AttachmentRepository
}

func NewListAttachmentVersions(repository repository.AttachmentRepository) ListAttachmentVersions {
	return ListAttachmentVersions{repository: repository}
}

func (useCase ListAttachmentVersions) Execute(tenantSlug string, attachmentPublicID string) []entity.AttachmentVersion {
	return useCase.repository.ListVersions(strings.TrimSpace(tenantSlug), strings.TrimSpace(attachmentPublicID))
}
