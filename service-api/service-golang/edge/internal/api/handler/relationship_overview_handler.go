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

type RelationshipOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewRelationshipOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) RelationshipOverviewHandler {
	return RelationshipOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler RelationshipOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	pipelineSummary := map[string]any{}
	relationshipIntelligence := map[string]any{}
	tenant360 := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/pipeline-summary?tenant_slug="+encodedTenantSlug, &pipelineSummary); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/relationship-intelligence?tenant_slug="+encodedTenantSlug, &relationshipIntelligence); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.RelationshipOverviewResponse{
		Service:                  handler.ServiceName,
		TenantSlug:               tenantSlug,
		GeneratedAt:              time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:         buildRelationshipExecutiveSummary(relationshipIntelligence),
		PipelineSummary:          pipelineSummary,
		RelationshipIntelligence: relationshipIntelligence,
		Tenant360:                tenant360,
	})
}

func buildRelationshipExecutiveSummary(payload map[string]any) dto.RelationshipExecutiveSummary {
	return dto.RelationshipExecutiveSummary{
		Status:                readMapString(payload, "readiness", "status"),
		AverageLeadScore:      readMapInt(payload, "scoring", "average"),
		HotLeads:              readMapInt(payload, "scoring", "hot"),
		PipelineConfigs:       readMapInt(payload, "pipeline", "configs"),
		PipelineStages:        readMapInt(payload, "pipeline", "stages"),
		OpenSupportCases:      readMapInt(payload, "support", "openCases"),
		OverdueSupportCases:   readMapInt(payload, "support", "overdueCases"),
		WeightedPipelineCents: readMapInt(payload, "forecast", "weightedPipelineCents"),
		BookedRevenueCents:    readMapInt(payload, "forecast", "bookedRevenueCents"),
	}
}
