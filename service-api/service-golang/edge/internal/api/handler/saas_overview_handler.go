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

type SaaSOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewSaaSOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) SaaSOverviewHandler {
	return SaaSOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler SaaSOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	saasControl := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/saas-control?tenant_slug="+encodedTenantSlug, &saasControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.SaaSOverviewResponse{
		Service:          handler.ServiceName,
		TenantSlug:       tenantSlug,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary: buildSaaSExecutiveSummary(saasControl),
		ServicePulse:     servicePulse,
		Tenant360:        tenant360,
		SaaSControl:      saasControl,
	})
}

func buildSaaSExecutiveSummary(payload map[string]any) dto.SaaSExecutiveSummary {
	return dto.SaaSExecutiveSummary{
		Status:              readMapString(payload, "readiness", "status"),
		EntitlementsTotal:   readMapInt(payload, "entitlements", "total"),
		EnabledEntitlements: readMapInt(payload, "entitlements", "enabled"),
		ActiveQuotas:        readMapInt(payload, "quotas", "active"),
		ActiveBlocks:        readMapInt(payload, "blocks", "active"),
		TrackedMetrics:      readMapInt(payload, "metering", "trackedMetrics"),
		QueuedJobs:          readMapInt(payload, "lifecycle", "queued"),
		CompletedJobs:       readMapInt(payload, "lifecycle", "completed"),
	}
}
