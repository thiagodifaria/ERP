package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type GetUploadSession struct {
	repository repository.UploadSessionRepository
}

func NewGetUploadSession(repository repository.UploadSessionRepository) GetUploadSession {
	return GetUploadSession{repository: repository}
}

func (useCase GetUploadSession) Execute(tenantSlug string, publicID string) (*entity.UploadSession, bool) {
	session, ok := useCase.repository.FindByPublicID(strings.TrimSpace(tenantSlug), strings.TrimSpace(publicID))
	if !ok {
		return nil, false
	}

	return &session, true
}
