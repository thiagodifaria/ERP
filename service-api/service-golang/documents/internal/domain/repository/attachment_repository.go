package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"

type AttachmentFilters struct {
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
	Source        string
	Visibility    string
	Archived      string
}

type AttachmentRepository interface {
	List(filters AttachmentFilters) []entity.Attachment
	FindByPublicID(tenantSlug string, publicID string) (entity.Attachment, bool)
	Save(attachment entity.Attachment) entity.Attachment
	Archive(tenantSlug string, publicID string, reason string) (entity.Attachment, bool)
	ListVersions(tenantSlug string, attachmentPublicID string) []entity.AttachmentVersion
	CreateVersion(tenantSlug string, attachmentPublicID string, version entity.AttachmentVersion) (entity.Attachment, entity.AttachmentVersion, bool)
}
