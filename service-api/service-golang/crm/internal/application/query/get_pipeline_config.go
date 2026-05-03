package query

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"

func GetPipelineConfig(factory repository.TenantRepositoryFactory, tenantSlug string) (map[string]any, error) {
	tenantRepositories, err := factory.ForTenant(tenantSlug)
	if err != nil {
		return nil, err
	}

	config := tenantRepositories.PipelineConfigRepository.Get()
	if config == nil {
		return map[string]any{}, nil
	}

	return map[string]any{
		"publicId":    config.PublicID,
		"name":        config.Name,
		"autoScoring": config.AutoScoring,
		"stages":      config.Stages,
	}, nil
}
