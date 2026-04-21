package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"

type AttachmentFilters struct {
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
}

type AttachmentRepository interface {
	List(filters AttachmentFilters) []entity.Attachment
	Save(attachment entity.Attachment) entity.Attachment
}
