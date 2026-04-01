// Handler que agrega os principais relatorios de tenant do analytics.
// O edge passa a entregar um cockpit unico de leitura sem duplicar calculos de negocio.
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

type TenantOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewTenantOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) TenantOverviewHandler {
	return TenantOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler TenantOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	pipelineSummary := map[string]any{}
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	automationBoard := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/pipeline-summary?tenant_slug="+encodedTenantSlug, &pipelineSummary); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/automation-board?tenant_slug="+encodedTenantSlug, &automationBoard); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.TenantOverviewResponse{
		Service:         handler.ServiceName,
		TenantSlug:      tenantSlug,
		GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
		PipelineSummary: pipelineSummary,
		ServicePulse:    servicePulse,
		Tenant360:       tenant360,
		AutomationBoard: automationBoard,
	})
}
