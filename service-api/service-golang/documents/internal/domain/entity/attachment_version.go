package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAttachmentVersionAttachmentPublicIDInvalid = errors.New("attachment version attachment public id is invalid")
	ErrAttachmentVersionNumberInvalid             = errors.New("attachment version number is invalid")
)

type AttachmentVersion struct {
	PublicID           string
	TenantSlug         string
	AttachmentPublicID string
	VersionNumber      int
	FileName           string
	ContentType        string
	StorageKey         string
	StorageDriver      string
	Source             string
	UploadedBy         string
	FileSizeBytes      int64
	ChecksumSHA256     string
	CreatedAt          time.Time
}

func NewAttachmentVersion(
	publicID string,
	tenantSlug string,
	attachmentPublicID string,
	versionNumber int,
	fileName string,
	contentType string,
	storageKey string,
	storageDriver string,
	source string,
	uploadedBy string,
	fileSizeBytes int64,
	checksumSHA256 string,
	createdAt time.Time,
) (AttachmentVersion, error) {
	if _, err := uuid.Parse(strings.TrimSpace(attachmentPublicID)); err != nil {
		return AttachmentVersion{}, ErrAttachmentVersionAttachmentPublicIDInvalid
	}
	if versionNumber <= 0 {
		return AttachmentVersion{}, ErrAttachmentVersionNumberInvalid
	}

	attachment, err := NewAttachment(
		publicID,
		tenantSlug,
		"documents.attachment",
		attachmentPublicID,
		fileName,
		contentType,
		storageKey,
		storageDriver,
		source,
		uploadedBy,
		fileSizeBytes,
		checksumSHA256,
		"internal",
		365,
		"",
		nil,
		createdAt,
	)
	if err != nil {
		return AttachmentVersion{}, err
	}

	return AttachmentVersion{
		PublicID:           attachment.PublicID,
		TenantSlug:         attachment.TenantSlug,
		AttachmentPublicID: attachmentPublicID,
		VersionNumber:      versionNumber,
		FileName:           attachment.FileName,
		ContentType:        attachment.ContentType,
		StorageKey:         attachment.StorageKey,
		StorageDriver:      attachment.StorageDriver,
		Source:             attachment.Source,
		UploadedBy:         attachment.UploadedBy,
		FileSizeBytes:      attachment.FileSizeBytes,
		ChecksumSHA256:     attachment.ChecksumSHA256,
		CreatedAt:          attachment.CreatedAt,
	}, nil
}
