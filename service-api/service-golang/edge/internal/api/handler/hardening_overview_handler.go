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

type HardeningOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewHardeningOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) HardeningOverviewHandler {
	return HardeningOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler HardeningOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	platformReliability := map[string]any{}
	hardeningReview := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/platform-reliability?tenant_slug="+encodedTenantSlug, &platformReliability); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/hardening-review?tenant_slug="+encodedTenantSlug, &hardeningReview); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.HardeningOverviewResponse{
		Service:             handler.ServiceName,
		TenantSlug:          tenantSlug,
		GeneratedAt:         time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:    buildHardeningExecutiveSummary(platformReliability, hardeningReview),
		ServicePulse:        servicePulse,
		PlatformReliability: platformReliability,
		HardeningReview:     hardeningReview,
	})
}

func buildHardeningExecutiveSummary(platformReliability map[string]any, hardeningReview map[string]any) dto.HardeningExecutiveSummary {
	return dto.HardeningExecutiveSummary{
		Status:                 readMapString(hardeningReview, "summary", "status"),
		StableChecks:           readMapInt(hardeningReview, "summary", "stableChecks"),
		AttentionChecks:        readMapInt(hardeningReview, "summary", "attentionChecks"),
		CriticalChecks:         readMapInt(hardeningReview, "summary", "criticalChecks"),
		DeadLetterEvents:       readMapInt(platformReliability, "stability", "deadLetterEvents"),
		FailedPaymentAttempts:  readMapInt(platformReliability, "stability", "failedPaymentAttempts"),
		OpenSecurityAlerts:     readMapInt(hardeningReview, "reviews", "security", "auditEvents"),
		LatestBenchmarkStatus:  readMapString(hardeningReview, "reviews", "performance", "latestBenchmarkStatus"),
		BackupRestoreValidated: readMapBool(hardeningReview, "reviews", "backupRestore", "validated"),
	}
}
