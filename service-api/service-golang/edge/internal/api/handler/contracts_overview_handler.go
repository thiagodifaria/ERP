package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
)

type ContractsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewContractsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) ContractsOverviewHandler {
	return ContractsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler ContractsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	integrationReadiness := map[string]any{}
	contractGovernance := map[string]any{}
	hardeningReview := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/integration-readiness?tenant_slug="+encodedTenantSlug, &integrationReadiness); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/contract-governance", &contractGovernance); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/hardening-review?tenant_slug="+encodedTenantSlug, &hardeningReview); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.ContractsOverviewResponse{
		Service:              handler.ServiceName,
		TenantSlug:           tenantSlug,
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:     buildContractsExecutiveSummary(contractGovernance),
		IntegrationReadiness: integrationReadiness,
		ContractGovernance:   contractGovernance,
		HardeningReview:      hardeningReview,
	})
}

func buildContractsExecutiveSummary(payload map[string]any) dto.ContractsExecutiveSummary {
	return dto.ContractsExecutiveSummary{
		Status:              readMapString(payload, "readiness", "status"),
		HTTPSpecs:           readMapInt(payload, "catalog", "httpSpecs"),
		EventSchemas:        readMapInt(payload, "catalog", "eventSchemas"),
		ADRCount:            readMapInt(payload, "catalog", "adrs"),
		ContractArtifacts:   readMapInt(payload, "catalog", "httpSpecs") + readMapInt(payload, "catalog", "eventSchemas") + readMapInt(payload, "catalog", "adrs"),
		NavigableAPIReady:   readMapBool(payload, "readiness", "navigableApiReady"),
		SchemaRegistryReady: readMapBool(payload, "readiness", "registryReady"),
	}
}
