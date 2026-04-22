package command

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type ArchiveAttachment struct {
	repository repository.AttachmentRepository
}

type ArchiveAttachmentInput struct {
	TenantSlug string
	PublicID   string
	Reason     string
}

type ArchiveAttachmentResult struct {
	Attachment *entity.Attachment
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewArchiveAttachment(repository repository.AttachmentRepository) ArchiveAttachment {
	return ArchiveAttachment{repository: repository}
}

func (useCase ArchiveAttachment) Execute(input ArchiveAttachmentInput) ArchiveAttachmentResult {
	if strings.TrimSpace(input.PublicID) == "" {
		return ArchiveAttachmentResult{
			ErrorCode:  "attachment_public_id_required",
			ErrorText:  "Attachment public id is required.",
			BadRequest: true,
		}
	}

	attachment, ok := useCase.repository.Archive(strings.TrimSpace(input.TenantSlug), strings.TrimSpace(input.PublicID), strings.TrimSpace(input.Reason))
	if !ok {
		return ArchiveAttachmentResult{NotFound: true}
	}

	return ArchiveAttachmentResult{Attachment: &attachment}
}
