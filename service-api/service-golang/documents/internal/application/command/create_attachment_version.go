package command

import (
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type CreateAttachmentVersion struct {
	repository repository.AttachmentRepository
}

type CreateAttachmentVersionInput struct {
	TenantSlug         string
	AttachmentPublicID string
	FileName           string
	ContentType        string
	StorageKey         string
	StorageDriver      string
	Source             string
	UploadedBy         string
	FileSizeBytes      int64
	ChecksumSHA256     string
}

type CreateAttachmentVersionResult struct {
	Attachment *entity.Attachment
	Version    *entity.AttachmentVersion
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewCreateAttachmentVersion(repository repository.AttachmentRepository) CreateAttachmentVersion {
	return CreateAttachmentVersion{repository: repository}
}

func (useCase CreateAttachmentVersion) Execute(input CreateAttachmentVersionInput) CreateAttachmentVersionResult {
	attachment, found := useCase.repository.FindByPublicID(strings.TrimSpace(input.TenantSlug), strings.TrimSpace(input.AttachmentPublicID))
	if !found {
		return CreateAttachmentVersionResult{NotFound: true}
	}

	version, err := entity.NewAttachmentVersion(
		newAttachmentPublicID(),
		attachment.TenantSlug,
		attachment.PublicID,
		attachment.CurrentVersion+1,
		input.FileName,
		input.ContentType,
		input.StorageKey,
		input.StorageDriver,
		input.Source,
		input.UploadedBy,
		input.FileSizeBytes,
		input.ChecksumSHA256,
		time.Time{},
	)
	if err != nil {
		return CreateAttachmentVersionResult{
			ErrorCode:  "invalid_attachment_version",
			ErrorText:  "Attachment version payload is invalid.",
			BadRequest: true,
		}
	}

	updatedAttachment, createdVersion, ok := useCase.repository.CreateVersion(attachment.TenantSlug, attachment.PublicID, version)
	if !ok {
		return CreateAttachmentVersionResult{NotFound: true}
	}

	return CreateAttachmentVersionResult{
		Attachment: &updatedAttachment,
		Version:    &createdVersion,
	}
}
