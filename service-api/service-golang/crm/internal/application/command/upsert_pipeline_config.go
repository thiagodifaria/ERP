package command

import (
	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

func UpsertPipelineConfig(factory repository.TenantRepositoryFactory, tenantSlug string, name string, stages []entity.PipelineStage, autoScoring bool) (entity.PipelineConfig, error) {
	tenantRepositories, err := factory.ForTenant(tenantSlug)
	if err != nil {
		return entity.PipelineConfig{}, err
	}

	current := tenantRepositories.PipelineConfigRepository.Get()
	publicID := uuid.NewString()
	if current != nil && current.PublicID != "" {
		publicID = current.PublicID
	}

	config, err := entity.NewPipelineConfig(publicID, name, stages, autoScoring)
	if err != nil {
		return entity.PipelineConfig{}, err
	}

	return tenantRepositories.PipelineConfigRepository.Save(config), nil
}
