// Handler que consolida a operacao de automacao em um unico payload.
// O edge cruza saude por fluxo, catalogo e entrega sem duplicar calculos analiticos.
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

type AutomationOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewAutomationOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) AutomationOverviewHandler {
	return AutomationOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler AutomationOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	automationBoard := map[string]any{}
	workflowDefinitionHealth := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/automation-board?tenant_slug="+encodedTenantSlug, &automationBoard); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/workflow-definition-health?tenant_slug="+encodedTenantSlug, &workflowDefinitionHealth); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.AutomationOverviewResponse{
		Service:                  handler.ServiceName,
		TenantSlug:               tenantSlug,
		GeneratedAt:              time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:         buildAutomationExecutiveSummary(servicePulse, automationBoard, workflowDefinitionHealth),
		ServicePulse:             servicePulse,
		AutomationBoard:          automationBoard,
		WorkflowDefinitionHealth: workflowDefinitionHealth,
	})
}

func buildAutomationExecutiveSummary(servicePulse map[string]any, automationBoard map[string]any, workflowDefinitionHealth map[string]any) dto.AutomationExecutiveSummary {
	stableDefinitions := readMapInt(workflowDefinitionHealth, "summary", "stable")
	attentionDefinitions := readMapInt(workflowDefinitionHealth, "summary", "attention")
	criticalDefinitions := readMapInt(workflowDefinitionHealth, "summary", "critical")
	status := "stable"

	if criticalDefinitions > 0 {
		status = "critical"
	} else if attentionDefinitions > 0 {
		status = "attention"
	}

	return dto.AutomationExecutiveSummary{
		Status:                     status,
		ActiveDefinitions:          readMapInt(automationBoard, "catalog", "definitionsActive"),
		StableDefinitions:          stableDefinitions,
		AttentionDefinitions:       attentionDefinitions,
		CriticalDefinitions:        criticalDefinitions,
		RunningControlRuns:         readMapInt(automationBoard, "control", "runningRuns"),
		CompletedRuntimeExecutions: readMapInt(automationBoard, "runtime", "completedExecutions"),
		ForwardedWebhookEvents:     readMapInt(servicePulse, "services", "webhookHub", "forwarded"),
	}
}
