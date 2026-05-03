package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type IntelligenceHandler struct {
	repositories repository.TenantRepositoryFactory
}

func NewIntelligenceHandler(repositories repository.TenantRepositoryFactory) *IntelligenceHandler {
	return &IntelligenceHandler{repositories: repositories}
}

func (handler *IntelligenceHandler) GetPipelineConfig(response http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	payload, err := query.GetPipelineConfig(handler.repositories, tenantSlug)
	if err != nil {
		writeBadRequest(response, "tenant_not_found", "Tenant was not found.")
		return
	}

	writeJSON(response, http.StatusOK, payload)
}

func (handler *IntelligenceHandler) UpsertPipelineConfig(response http.ResponseWriter, request *http.Request) {
	var payload dto.UpsertPipelineConfigRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(response, "invalid_json", "Payload is invalid.")
		return
	}

	stages := make([]entity.PipelineStage, 0, len(payload.Stages))
	for _, stage := range payload.Stages {
		stages = append(stages, entity.PipelineStage{
			Key:              stage.Key,
			Name:             stage.Name,
			RequiresApproval: stage.RequiresApproval,
		})
	}

	config, err := command.UpsertPipelineConfig(handler.repositories, payload.TenantSlug, payload.Name, stages, payload.AutoScoring)
	if err != nil {
		writeBadRequest(response, "pipeline_config_invalid", "Pipeline configuration is invalid.")
		return
	}

	writeJSON(response, http.StatusOK, config)
}

func (handler *IntelligenceHandler) Summary(response http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	payload, err := query.GetLeadIntelligenceSummary(handler.repositories, tenantSlug)
	if err != nil {
		writeBadRequest(response, "tenant_not_found", "Tenant was not found.")
		return
	}

	writeJSON(response, http.StatusOK, payload)
}
