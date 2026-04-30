package command

import (
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type CompleteUploadSession struct {
	uploadSessionRepository repository.UploadSessionRepository
	attachmentRepository    repository.AttachmentRepository
}

type CompleteUploadSessionInput struct {
	TenantSlug      string
	PublicID        string
	UploadedBy      string
	FileSizeBytes   int64
	ChecksumSHA256  string
}

type CompleteUploadSessionResult struct {
	UploadSession *entity.UploadSession
	Attachment    *entity.Attachment
	ErrorCode     string
	ErrorText     string
	BadRequest    bool
	NotFound      bool
	Conflict      bool
}

func NewCompleteUploadSession(uploadSessionRepository repository.UploadSessionRepository, attachmentRepository repository.AttachmentRepository) CompleteUploadSession {
	return CompleteUploadSession{
		uploadSessionRepository: uploadSessionRepository,
		attachmentRepository:    attachmentRepository,
	}
}

func (useCase CompleteUploadSession) Execute(input CompleteUploadSessionInput) CompleteUploadSessionResult {
	session, ok := useCase.uploadSessionRepository.FindByPublicID(strings.TrimSpace(input.TenantSlug), strings.TrimSpace(input.PublicID))
	if !ok {
		return CompleteUploadSessionResult{NotFound: true}
	}
	if session.Status == "completed" {
		return CompleteUploadSessionResult{
			ErrorCode:  "upload_session_completed",
			ErrorText:  "Upload session was already completed.",
			Conflict:   true,
		}
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		return CompleteUploadSessionResult{
			ErrorCode:  "upload_session_expired",
			ErrorText:  "Upload session is expired.",
			BadRequest: true,
		}
	}

	attachment, err := entity.NewAttachment(
		newAttachmentPublicID(),
		session.TenantSlug,
		session.OwnerType,
		session.OwnerPublicID,
		session.FileName,
		session.ContentType,
		session.StorageKey,
		session.StorageDriver,
		session.Source,
		input.UploadedBy,
		input.FileSizeBytes,
		input.ChecksumSHA256,
		session.Visibility,
		session.RetentionDays,
		"",
		nil,
		time.Time{},
	)
	if err != nil {
		return CompleteUploadSessionResult{
			ErrorCode:  "invalid_attachment",
			ErrorText:  "Attachment payload is invalid.",
			BadRequest: true,
		}
	}
	if strings.TrimSpace(attachment.UploadedBy) == "" {
		attachment.UploadedBy = session.RequestedBy
	}

	createdAttachment := useCase.attachmentRepository.Save(attachment)
	completedSession, ok := useCase.uploadSessionRepository.Complete(session.TenantSlug, session.PublicID, createdAttachment.PublicID, time.Now().UTC())
	if !ok {
		return CompleteUploadSessionResult{NotFound: true}
	}

	return CompleteUploadSessionResult{
		UploadSession: &completedSession,
		Attachment:    &createdAttachment,
	}
}
