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

type IntegrationsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewIntegrationsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) IntegrationsOverviewHandler {
	return IntegrationsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler IntegrationsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	engagementOperations := map[string]any{}
	integrationReadiness := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/engagement-operations?tenant_slug="+encodedTenantSlug, &engagementOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/integration-readiness?tenant_slug="+encodedTenantSlug, &integrationReadiness); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.IntegrationsOverviewResponse{
		Service:              handler.ServiceName,
		TenantSlug:           tenantSlug,
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:     buildIntegrationsExecutiveSummary(integrationReadiness),
		ServicePulse:         servicePulse,
		EngagementOperations: engagementOperations,
		IntegrationReadiness: integrationReadiness,
	})
}

func buildIntegrationsExecutiveSummary(integrationReadiness map[string]any) dto.IntegrationsExecutiveSummary {
	return dto.IntegrationsExecutiveSummary{
		Status:                  readMapString(integrationReadiness, "readiness", "status"),
		ConfiguredProviders:     readMapInt(integrationReadiness, "providers", "configured"),
		ActiveInboundProviders:  readMapInt(integrationReadiness, "providers", "activeInboundProviders"),
		ActiveOutboundProviders: readMapInt(integrationReadiness, "providers", "activeOutboundProviders"),
		InboundLeads:            readMapInt(integrationReadiness, "flows", "inboundLeads"),
		WorkflowDispatches:      readMapInt(integrationReadiness, "flows", "workflowDispatches"),
		BusinessLinkedEvents:    readMapInt(integrationReadiness, "flows", "businessEntityLinkedEvents"),
		FailedProviderEvents:    readMapInt(integrationReadiness, "flows", "failedProviderEvents"),
		DeadLetterEvents:        readMapInt(integrationReadiness, "webhookHub", "deadLetterEvents"),
		OpenProviderRisks:       readMapInt(integrationReadiness, "readiness", "openProviderRisks"),
	}
}
