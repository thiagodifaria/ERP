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
	ErrAttachmentFileSizeInvalid      = errors.New("attachment file size is invalid")
	ErrAttachmentVisibilityInvalid    = errors.New("attachment visibility is invalid")
)

type Attachment struct {
	PublicID       string
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
	CurrentVersion int
	VersionCount   int
	ArchiveReason  string
	ArchivedAt     *time.Time
	CreatedAt      time.Time
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
	fileSizeBytes int64,
	checksumSHA256 string,
	visibility string,
	retentionDays int,
	archiveReason string,
	archivedAt *time.Time,
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
	normalizedChecksumSHA256 := strings.ToLower(strings.TrimSpace(checksumSHA256))
	normalizedVisibility := strings.ToLower(strings.TrimSpace(visibility))
	normalizedArchiveReason := strings.TrimSpace(archiveReason)

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
	if fileSizeBytes < 0 {
		return Attachment{}, ErrAttachmentFileSizeInvalid
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
	if normalizedVisibility == "" {
		normalizedVisibility = "internal"
	}
	if normalizedVisibility != "internal" && normalizedVisibility != "restricted" && normalizedVisibility != "public" {
		return Attachment{}, ErrAttachmentVisibilityInvalid
	}
	if retentionDays <= 0 {
		retentionDays = 365
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	var normalizedArchivedAt *time.Time
	if archivedAt != nil && !archivedAt.IsZero() {
		value := archivedAt.UTC()
		normalizedArchivedAt = &value
	}

	return Attachment{
		PublicID:       normalizedPublicID,
		TenantSlug:     normalizedTenantSlug,
		OwnerType:      normalizedOwnerType,
		OwnerPublicID:  normalizedOwnerPublicID,
		FileName:       normalizedFileName,
		ContentType:    normalizedContentType,
		StorageKey:     normalizedStorageKey,
		StorageDriver:  normalizedStorageDriver,
		Source:         normalizedSource,
		UploadedBy:     normalizedUploadedBy,
		FileSizeBytes:  fileSizeBytes,
		ChecksumSHA256: normalizedChecksumSHA256,
		Visibility:     normalizedVisibility,
		RetentionDays:  retentionDays,
		CurrentVersion: 1,
		VersionCount:   1,
		ArchiveReason:  normalizedArchiveReason,
		ArchivedAt:     normalizedArchivedAt,
		CreatedAt:      createdAt.UTC(),
	}, nil
}

func (attachment Attachment) Archive(reason string, archivedAt time.Time) Attachment {
	if attachment.ArchivedAt != nil && !attachment.ArchivedAt.IsZero() {
		return attachment
	}

	normalizedArchivedAt := archivedAt.UTC()
	if normalizedArchivedAt.IsZero() {
		normalizedArchivedAt = time.Now().UTC()
	}

	attachment.ArchivedAt = &normalizedArchivedAt
	attachment.ArchiveReason = strings.TrimSpace(reason)
	if attachment.ArchiveReason == "" {
		attachment.ArchiveReason = "archived"
	}

	return attachment
}
