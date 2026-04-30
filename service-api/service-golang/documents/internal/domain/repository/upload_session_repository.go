package repository

import (
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
)

type UploadSessionRepository interface {
	Save(session entity.UploadSession) entity.UploadSession
	FindByPublicID(tenantSlug string, publicID string) (entity.UploadSession, bool)
	Complete(tenantSlug string, publicID string, attachmentPublicID string, completedAt time.Time) (entity.UploadSession, bool)
}
