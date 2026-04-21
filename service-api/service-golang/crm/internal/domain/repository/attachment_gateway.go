package repository

import "time"

type AttachmentRecord struct {
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

type CreateAttachmentInput struct {
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
	FileName      string
	ContentType   string
	StorageKey    string
	StorageDriver string
	Source        string
	UploadedBy    string
}

type AttachmentGateway interface {
	List(tenantSlug string, ownerType string, ownerPublicID string) ([]AttachmentRecord, error)
	Create(input CreateAttachmentInput) (*AttachmentRecord, error)
}
