package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUploadSessionPublicIDInvalid        = errors.New("upload session public id is invalid")
	ErrUploadSessionOwnerTypeRequired      = errors.New("upload session owner type is required")
	ErrUploadSessionOwnerPublicIDInvalid   = errors.New("upload session owner public id is invalid")
	ErrUploadSessionFileNameRequired       = errors.New("upload session file name is required")
	ErrUploadSessionStorageKeyRequired     = errors.New("upload session storage key is required")
	ErrUploadSessionStatusInvalid          = errors.New("upload session status is invalid")
	ErrUploadSessionAttachmentIDInvalid    = errors.New("upload session attachment public id is invalid")
)

type UploadSession struct {
	PublicID          string
	TenantSlug        string
	OwnerType         string
	OwnerPublicID     string
	FileName          string
	ContentType       string
	StorageKey        string
	StorageDriver     string
	Source            string
	RequestedBy       string
	Visibility        string
	RetentionDays     int
	Status            string
	AttachmentPublicID string
	ExpiresAt         time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
}

func NewUploadSession(
	publicID string,
	tenantSlug string,
	ownerType string,
	ownerPublicID string,
	fileName string,
	contentType string,
	storageKey string,
	storageDriver string,
	source string,
	requestedBy string,
	visibility string,
	retentionDays int,
	status string,
	attachmentPublicID string,
	expiresAt time.Time,
	completedAt *time.Time,
	createdAt time.Time,
) (UploadSession, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedOwnerType := strings.TrimSpace(ownerType)
	normalizedOwnerPublicID := strings.TrimSpace(ownerPublicID)
	normalizedFileName := strings.TrimSpace(fileName)
	normalizedContentType := strings.TrimSpace(contentType)
	normalizedStorageKey := strings.TrimSpace(storageKey)
	normalizedStorageDriver := strings.TrimSpace(storageDriver)
	normalizedSource := strings.TrimSpace(source)
	normalizedRequestedBy := strings.TrimSpace(requestedBy)
	normalizedVisibility := strings.ToLower(strings.TrimSpace(visibility))
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	normalizedAttachmentPublicID := strings.TrimSpace(attachmentPublicID)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return UploadSession{}, ErrUploadSessionPublicIDInvalid
	}
	if normalizedOwnerType == "" {
		return UploadSession{}, ErrUploadSessionOwnerTypeRequired
	}
	if _, err := uuid.Parse(normalizedOwnerPublicID); err != nil {
		return UploadSession{}, ErrUploadSessionOwnerPublicIDInvalid
	}
	if normalizedFileName == "" {
		return UploadSession{}, ErrUploadSessionFileNameRequired
	}
	if normalizedStorageKey == "" {
		return UploadSession{}, ErrUploadSessionStorageKeyRequired
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
	if normalizedRequestedBy == "" {
		normalizedRequestedBy = "system"
	}
	if normalizedVisibility == "" {
		normalizedVisibility = "internal"
	}
	if retentionDays <= 0 {
		retentionDays = 365
	}
	if normalizedStatus == "" {
		normalizedStatus = "pending_upload"
	}
	if normalizedStatus != "pending_upload" && normalizedStatus != "completed" && normalizedStatus != "expired" {
		return UploadSession{}, ErrUploadSessionStatusInvalid
	}
	if normalizedAttachmentPublicID != "" {
		if _, err := uuid.Parse(normalizedAttachmentPublicID); err != nil {
			return UploadSession{}, ErrUploadSessionAttachmentIDInvalid
		}
	}
	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().Add(15 * time.Minute)
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	var normalizedCompletedAt *time.Time
	if completedAt != nil && !completedAt.IsZero() {
		value := completedAt.UTC()
		normalizedCompletedAt = &value
	}

	return UploadSession{
		PublicID:           normalizedPublicID,
		TenantSlug:         normalizedTenantSlug,
		OwnerType:          normalizedOwnerType,
		OwnerPublicID:      normalizedOwnerPublicID,
		FileName:           normalizedFileName,
		ContentType:        normalizedContentType,
		StorageKey:         normalizedStorageKey,
		StorageDriver:      normalizedStorageDriver,
		Source:             normalizedSource,
		RequestedBy:        normalizedRequestedBy,
		Visibility:         normalizedVisibility,
		RetentionDays:      retentionDays,
		Status:             normalizedStatus,
		AttachmentPublicID: normalizedAttachmentPublicID,
		ExpiresAt:          expiresAt.UTC(),
		CompletedAt:        normalizedCompletedAt,
		CreatedAt:          createdAt.UTC(),
	}, nil
}

func (session UploadSession) Complete(attachmentPublicID string, completedAt time.Time) UploadSession {
	if session.Status == "completed" {
		return session
	}

	value := completedAt.UTC()
	if value.IsZero() {
		value = time.Now().UTC()
	}

	session.Status = "completed"
	session.AttachmentPublicID = strings.TrimSpace(attachmentPublicID)
	session.CompletedAt = &value
	return session
}
