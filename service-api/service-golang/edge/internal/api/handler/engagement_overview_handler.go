// Handler que consolida a operacao de engagement em um unico payload.
// O edge expoe saude de campanhas, templates e entregas sem duplicar calculos analiticos.
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

type EngagementOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewEngagementOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) EngagementOverviewHandler {
	return EngagementOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler EngagementOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	engagementOperations := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/engagement-operations?tenant_slug="+encodedTenantSlug, &engagementOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.EngagementOverviewResponse{
		Service:              handler.ServiceName,
		TenantSlug:           tenantSlug,
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:     buildEngagementExecutiveSummary(engagementOperations),
		ServicePulse:         servicePulse,
		Tenant360:            tenant360,
		EngagementOperations: engagementOperations,
	})
}

func buildEngagementExecutiveSummary(engagementOperations map[string]any) dto.EngagementExecutiveSummary {
	deliveryRate := readMapFloat(engagementOperations, "deliveries", "deliveryRate")
	failedDeliveries := readMapInt(engagementOperations, "deliveries", "failed")
	status := "stable"

	if failedDeliveries > 0 || deliveryRate < 0.8 {
		status = "attention"
	}
	if failedDeliveries > 2 || deliveryRate < 0.6 {
		status = "critical"
	}

	return dto.EngagementExecutiveSummary{
		Status:               status,
		Campaigns:            readMapInt(engagementOperations, "campaigns", "total"),
		ActiveCampaigns:      readMapInt(engagementOperations, "campaigns", "active"),
		Templates:            readMapInt(engagementOperations, "templates", "total"),
		Deliveries:           readMapInt(engagementOperations, "deliveries", "total"),
		DeliveredDeliveries:  readMapInt(engagementOperations, "deliveries", "delivered"),
		FailedDeliveries:     failedDeliveries,
		ConvertedTouchpoints: readMapInt(engagementOperations, "touchpoints", "converted"),
		DeliveryRate:         deliveryRate,
	}
}
