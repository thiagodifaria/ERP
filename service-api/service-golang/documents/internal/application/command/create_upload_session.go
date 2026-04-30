package command

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type CreateUploadSession struct {
	repository repository.UploadSessionRepository
}

type CreateUploadSessionInput struct {
	TenantSlug      string
	OwnerType       string
	OwnerPublicID   string
	FileName        string
	ContentType     string
	StorageKey      string
	StorageDriver   string
	Source          string
	RequestedBy     string
	Visibility      string
	RetentionDays   int
	ExpiresInSeconds int
}

type CreateUploadSessionResult struct {
	UploadSession *entity.UploadSession
	ErrorCode     string
	ErrorText     string
	BadRequest    bool
}

func NewCreateUploadSession(repository repository.UploadSessionRepository) CreateUploadSession {
	return CreateUploadSession{repository: repository}
}

func (useCase CreateUploadSession) Execute(input CreateUploadSessionInput) CreateUploadSessionResult {
	publicID := newUploadSessionPublicID()
	expiresAt := time.Now().UTC().Add(15 * time.Minute)
	if input.ExpiresInSeconds > 0 {
		expiresAt = time.Now().UTC().Add(time.Duration(input.ExpiresInSeconds) * time.Second)
	}

	storageKey := strings.TrimSpace(input.StorageKey)
	if storageKey == "" {
		fileName := sanitizeFileName(input.FileName)
		storageKey = fmt.Sprintf(
			"%s/%s/%s/%s",
			strings.ToLower(strings.TrimSpace(input.TenantSlug)),
			strings.ReplaceAll(strings.TrimSpace(input.OwnerType), ".", "/"),
			strings.TrimSpace(input.OwnerPublicID),
			fileName,
		)
	}

	session, err := entity.NewUploadSession(
		publicID,
		input.TenantSlug,
		input.OwnerType,
		input.OwnerPublicID,
		input.FileName,
		input.ContentType,
		storageKey,
		input.StorageDriver,
		input.Source,
		input.RequestedBy,
		input.Visibility,
		input.RetentionDays,
		"pending_upload",
		"",
		expiresAt,
		nil,
		time.Time{},
	)
	if err != nil {
		return CreateUploadSessionResult{
			ErrorCode:  "invalid_upload_session",
			ErrorText:  "Upload session payload is invalid.",
			BadRequest: true,
		}
	}

	created := useCase.repository.Save(session)
	return CreateUploadSessionResult{UploadSession: &created}
}

func newUploadSessionPublicID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return uuid.Nil.String()
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80
	return uuid.UUID(raw).String()
}

func sanitizeFileName(value string) string {
	base := strings.TrimSpace(filepath.Base(value))
	base = strings.ReplaceAll(base, " ", "-")
	base = strings.ReplaceAll(base, "..", "-")
	if base == "" {
		return "upload.bin"
	}

	return base
}
