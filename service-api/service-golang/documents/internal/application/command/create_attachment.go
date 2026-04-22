package command

import (
	"crypto/rand"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type CreateAttachment struct {
	repository repository.AttachmentRepository
}

type CreateAttachmentInput struct {
	TenantSlug     string
	OwnerType      string
	OwnerPublicID  string
	FileName       string
	ContentType    string
	StorageKey     string
	StorageDriver  string
	Source         string
	UploadedBy     string
	FileSizeBytes  int64
	ChecksumSHA256 string
	Visibility     string
	RetentionDays  int
}

type CreateAttachmentResult struct {
	Attachment *entity.Attachment
	ErrorCode  string
	ErrorText  string
	BadRequest bool
}

func NewCreateAttachment(repository repository.AttachmentRepository) CreateAttachment {
	return CreateAttachment{repository: repository}
}

func (useCase CreateAttachment) Execute(input CreateAttachmentInput) CreateAttachmentResult {
	attachment, err := entity.NewAttachment(
		newAttachmentPublicID(),
		input.TenantSlug,
		input.OwnerType,
		input.OwnerPublicID,
		input.FileName,
		input.ContentType,
		input.StorageKey,
		input.StorageDriver,
		input.Source,
		input.UploadedBy,
		input.FileSizeBytes,
		input.ChecksumSHA256,
		input.Visibility,
		input.RetentionDays,
		"",
		nil,
		time.Time{},
	)
	if err != nil {
		return CreateAttachmentResult{
			ErrorCode:  "invalid_attachment",
			ErrorText:  "Attachment payload is invalid.",
			BadRequest: true,
		}
	}

	created := useCase.repository.Save(attachment)
	return CreateAttachmentResult{Attachment: &created}
}

func newAttachmentPublicID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return uuid.Nil.String()
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	return uuid.UUID(raw).String()
}
