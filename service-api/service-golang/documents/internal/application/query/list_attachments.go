package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type ListAttachments struct {
	repository repository.AttachmentRepository
}

func NewListAttachments(repository repository.AttachmentRepository) ListAttachments {
	return ListAttachments{repository: repository}
}

func (useCase ListAttachments) Execute(filters repository.AttachmentFilters) []entity.Attachment {
	return useCase.repository.List(filters)
}
