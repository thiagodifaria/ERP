package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAttachmentPublicIDInvalid      = errors.New("attachment public id is invalid")
	ErrAttachmentOwnerTypeRequired    = errors.New("attachment owner type is required")
	ErrAttachmentOwnerPublicIDInvalid = errors.New("attachment owner public id is invalid")
	ErrAttachmentFileNameRequired     = errors.New("attachment file name is required")
	ErrAttachmentStorageKeyRequired   = errors.New("attachment storage key is required")
)

type Attachment struct {
	PublicID      string
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
	FileName      string
	ContentType   string
	StorageKey    string
	StorageDriver string
	Source        string
	UploadedBy    string
	CreatedAt     time.Time
}

func NewAttachment(
	publicID string,
	tenantSlug string,
	ownerType string,
	ownerPublicID string,
	fileName string,
	contentType string,
	storageKey string,
	storageDriver string,
	source string,
	uploadedBy string,
	createdAt time.Time,
) (Attachment, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedOwnerType := strings.TrimSpace(ownerType)
	normalizedOwnerPublicID := strings.TrimSpace(ownerPublicID)
	normalizedFileName := strings.TrimSpace(fileName)
	normalizedContentType := strings.TrimSpace(contentType)
	normalizedStorageKey := strings.TrimSpace(storageKey)
	normalizedStorageDriver := strings.TrimSpace(storageDriver)
	normalizedSource := strings.TrimSpace(source)
	normalizedUploadedBy := strings.TrimSpace(uploadedBy)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Attachment{}, ErrAttachmentPublicIDInvalid
	}
	if normalizedOwnerType == "" {
		return Attachment{}, ErrAttachmentOwnerTypeRequired
	}
	if _, err := uuid.Parse(normalizedOwnerPublicID); err != nil {
		return Attachment{}, ErrAttachmentOwnerPublicIDInvalid
	}
	if normalizedFileName == "" {
		return Attachment{}, ErrAttachmentFileNameRequired
	}
	if normalizedStorageKey == "" {
		return Attachment{}, ErrAttachmentStorageKeyRequired
	}
	if normalizedTenantSlug == "" {
		normalizedTenantSlug = "bootstrap-ops"
	}
	if normalizedContentType == "" {
		normalizedContentType = "application/octet-stream"
	}
	if normalizedStorageDriver == "" {
		normalizedStorageDriver = "manual"
	}
	if normalizedSource == "" {
		normalizedSource = "manual"
	}
	if normalizedUploadedBy == "" {
		normalizedUploadedBy = "system"
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	return Attachment{
		PublicID:      normalizedPublicID,
		TenantSlug:    normalizedTenantSlug,
		OwnerType:     normalizedOwnerType,
		OwnerPublicID: normalizedOwnerPublicID,
		FileName:      normalizedFileName,
		ContentType:   normalizedContentType,
		StorageKey:    normalizedStorageKey,
		StorageDriver: normalizedStorageDriver,
		Source:        normalizedSource,
		UploadedBy:    normalizedUploadedBy,
		CreatedAt:     createdAt.UTC(),
	}, nil
}
