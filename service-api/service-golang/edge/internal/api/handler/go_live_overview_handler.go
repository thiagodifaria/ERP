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

type GoLiveOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewGoLiveOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) GoLiveOverviewHandler {
	return GoLiveOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler GoLiveOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	saasControl := map[string]any{}
	goLiveControl := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/saas-control?tenant_slug="+encodedTenantSlug, &saasControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/go-live-control?tenant_slug="+encodedTenantSlug, &goLiveControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.GoLiveOverviewResponse{
		Service:          handler.ServiceName,
		TenantSlug:       tenantSlug,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary: buildGoLiveExecutiveSummary(goLiveControl),
		ServicePulse:     servicePulse,
		SaaSControl:      saasControl,
		GoLiveControl:    goLiveControl,
	})
}

func buildGoLiveExecutiveSummary(payload map[string]any) dto.GoLiveExecutiveSummary {
	return dto.GoLiveExecutiveSummary{
		Status:             readMapString(payload, "readiness", "status"),
		PlannedRollouts:    readMapInt(payload, "rollouts", "planned"),
		RunningRollouts:    readMapInt(payload, "rollouts", "running"),
		CompletedRollouts:  readMapInt(payload, "rollouts", "completed"),
		RolledBackRollouts: readMapInt(payload, "rollouts", "rolledBack"),
		TrackedMetrics:     readMapInt(payload, "adoption", "trackedMetrics"),
		TotalQuantity:      readMapInt(payload, "adoption", "totalQuantity"),
		AdoptionPct:        readMapInt(payload, "adoption", "adoptionPct"),
		PendingAdjustments: readMapInt(payload, "adjustments", "recommended"),
		CriticalBottlenecks: readMapInt(payload, "bottlenecks", "critical"),
		RolloutReady:       readMapBool(payload, "readiness", "rolloutReady"),
		MetricsObserved:    readMapBool(payload, "readiness", "metricsObserved"),
	}
}
