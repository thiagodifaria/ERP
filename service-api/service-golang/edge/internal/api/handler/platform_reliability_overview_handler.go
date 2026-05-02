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

type PlatformReliabilityOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewPlatformReliabilityOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) PlatformReliabilityOverviewHandler {
	return PlatformReliabilityOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler PlatformReliabilityOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	deliveryReliability := map[string]any{}
	platformReliability := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/delivery-reliability", &deliveryReliability); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/platform-reliability?tenant_slug="+encodedTenantSlug, &platformReliability); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.PlatformReliabilityOverviewResponse{
		Service:             handler.ServiceName,
		TenantSlug:          tenantSlug,
		GeneratedAt:         time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:    buildPlatformReliabilityExecutiveSummary(platformReliability),
		ServicePulse:        servicePulse,
		DeliveryReliability: deliveryReliability,
		PlatformReliability: platformReliability,
	})
}

func buildPlatformReliabilityExecutiveSummary(platformReliability map[string]any) dto.PlatformReliabilityExecutiveSummary {
	return dto.PlatformReliabilityExecutiveSummary{
		Status:                   readMapString(platformReliability, "stability", "status"),
		PendingWebhookEvents:     readMapInt(platformReliability, "stability", "pendingWebhookEvents"),
		DeadLetterEvents:         readMapInt(platformReliability, "stability", "deadLetterEvents"),
		FailedWorkflowExecutions: readMapInt(platformReliability, "stability", "failedWorkflowExecutions"),
		CriticalRecoveryCases:    readMapInt(platformReliability, "stability", "criticalRecoveryCases"),
		FailedPaymentAttempts:    readMapInt(platformReliability, "stability", "failedPaymentAttempts"),
		WebhookForwardingRateBps: int(readMapFloat(platformReliability, "serviceLevelObjectives", "webhookForwardingRate") * 10000),
		WorkflowSuccessRateBps:   int(readMapFloat(platformReliability, "serviceLevelObjectives", "workflowSuccessRate") * 10000),
		BillingRecoveryRateBps:   int(readMapFloat(platformReliability, "serviceLevelObjectives", "billingRecoveryRate") * 10000),
		OpenCriticalRisks:        readMapInt(platformReliability, "safeguards", "openCriticalRisks"),
	}
}

func readMapString(payload map[string]any, path ...string) string {
	var current any = payload

	for _, key := range path {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return ""
		}

		nextValue, ok := currentMap[key]
		if !ok {
			return ""
		}

		current = nextValue
	}

	value, ok := current.(string)
	if !ok {
		return ""
	}

	return value
}
