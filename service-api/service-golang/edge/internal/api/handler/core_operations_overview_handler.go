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

type CoreOperationsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewCoreOperationsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) CoreOperationsOverviewHandler {
	return CoreOperationsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler CoreOperationsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	coreOperations := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/core-operations?tenant_slug="+encodedTenantSlug, &coreOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.CoreOperationsOverviewResponse{
		Service:          handler.ServiceName,
		TenantSlug:       tenantSlug,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary: buildCoreOperationsExecutiveSummary(coreOperations),
		ServicePulse:     servicePulse,
		Tenant360:        tenant360,
		CoreOperations:   coreOperations,
	})
}

func buildCoreOperationsExecutiveSummary(payload map[string]any) dto.CoreOperationsExecutiveSummary {
	return dto.CoreOperationsExecutiveSummary{
		Status:                readMapString(payload, "readiness", "status"),
		CatalogItems:          readMapInt(payload, "summary", "catalogItems"),
		Suppliers:             readMapInt(payload, "summary", "suppliers"),
		SupportCases:          readMapInt(payload, "summary", "supportCases"),
		OverdueSupportCases:   readMapInt(payload, "support", "overdue"),
		UnreadNotifications:   readMapInt(payload, "notification", "unread"),
		CriticalNotifications: readMapInt(payload, "notification", "critical"),
	}
}
