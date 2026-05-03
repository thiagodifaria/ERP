package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"

type PipelineConfigRepository interface {
	Get() *entity.PipelineConfig
	Save(config entity.PipelineConfig) entity.PipelineConfig
}
