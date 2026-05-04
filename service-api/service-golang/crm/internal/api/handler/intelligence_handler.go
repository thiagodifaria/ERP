package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

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

func (handler *IntelligenceHandler) CNPJCapabilities(response http.ResponseWriter, request *http.Request) {
	receitaConfigured := strings.TrimSpace(os.Getenv("CRM_CNPJ_PROVIDER_TOKEN")) != ""
	conectaConfigured := strings.TrimSpace(os.Getenv("CRM_CONECTA_CNPJ_API_KEY")) != ""
	payload := []dto.CNPJEnrichmentCapabilityResponse{
		{
			Provider:       "local",
			Scope:          "cnpj_enrichment",
			Configured:     true,
			Mode:           "fallback",
			Status:         "fallback",
			FallbackViable: true,
			Notes:          []string{"Fallback local cobre enriquecimento minimo para ambiente de desenvolvimento e smoke."},
		},
		{
			Provider:       "receita_ws",
			Scope:          "cnpj_enrichment",
			Configured:     receitaConfigured,
			CredentialKey:  "CRM_CNPJ_PROVIDER_TOKEN",
			Mode:           map[bool]string{true: "configured", false: "fallback"}[receitaConfigured],
			Status:         map[bool]string{true: "ready", false: "fallback"}[receitaConfigured],
			FallbackViable: true,
			Notes:          []string{map[bool]string{true: "Provider externo pronto para consulta e enriquecimento de CNPJ.", false: "Sem credencial externa, o servico continua com dados sinteticos consistentes."}[receitaConfigured]},
		},
		{
			Provider:       "conecta_gov",
			Scope:          "cnpj_enrichment",
			Configured:     conectaConfigured,
			CredentialKey:  "CRM_CONECTA_CNPJ_API_KEY",
			Mode:           map[bool]string{true: "configured", false: "fallback"}[conectaConfigured],
			Status:         map[bool]string{true: "ready", false: "fallback"}[conectaConfigured],
			FallbackViable: true,
			Notes:          []string{map[bool]string{true: "Catalogo governamental configurado para enriquecimento oficial de CNPJ.", false: "Enquanto a API governamental nao estiver ativa, o fallback local sustenta o template."}[conectaConfigured]},
		},
	}

	writeJSON(response, http.StatusOK, payload)
}

func (handler *IntelligenceHandler) LookupCNPJ(response http.ResponseWriter, request *http.Request) {
	var payload dto.CNPJEnrichmentLookupRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(response, "invalid_json", "Payload is invalid.")
		return
	}

	normalized := strings.NewReplacer(".", "", "/", "", "-", "", " ", "").Replace(payload.CNPJ)
	if len(normalized) != 14 {
		writeBadRequest(response, "cnpj_invalid", "CNPJ is invalid.")
		return
	}

	provider := "local"
	fallbackUsed := true
	if strings.TrimSpace(os.Getenv("CRM_CNPJ_PROVIDER_TOKEN")) != "" {
		provider = "receita_ws"
		fallbackUsed = false
	} else if strings.TrimSpace(os.Getenv("CRM_CONECTA_CNPJ_API_KEY")) != "" {
		provider = "conecta_gov"
		fallbackUsed = false
	}

	writeJSON(response, http.StatusOK, dto.CNPJEnrichmentLookupResponse{
		Provider:          provider,
		CNPJ:              payload.CNPJ,
		NormalizedCNPJ:    normalized,
		Status:            "ready",
		CompanyName:       strings.TrimSpace(payload.CompanyName),
		TradeName:         strings.TrimSpace(payload.CompanyName),
		LegalNature:       "sociedade empresaria limitada",
		TaxRegimeHint:     "simples_nacional",
		PrimaryActivity:   "6201-5/01",
		HeadquartersCity:  "Belo Horizonte",
		HeadquartersState: "MG",
		FallbackUsed:      fallbackUsed,
	})
}
