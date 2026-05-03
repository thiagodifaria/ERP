package persistence

import (
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type InMemoryPipelineConfigRepository struct {
	mutex  sync.RWMutex
	config entity.PipelineConfig
}

func NewInMemoryPipelineConfigRepository() repository.PipelineConfigRepository {
	return &InMemoryPipelineConfigRepository{config: entity.DefaultPipelineConfig()}
}

func (repository *InMemoryPipelineConfigRepository) Get() *entity.PipelineConfig {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	config := repository.config
	return &config
}

func (repository *InMemoryPipelineConfigRepository) Save(config entity.PipelineConfig) entity.PipelineConfig {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.config = config
	return repository.config
}
