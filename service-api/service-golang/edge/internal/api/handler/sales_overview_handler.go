// Handler que consolida a operacao comercial em um unico payload.
// O edge cruza funil, pulso de servicos e jornada de vendas sem duplicar calculos analiticos.
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

type SalesOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewSalesOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) SalesOverviewHandler {
	return SalesOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler SalesOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	pipelineSummary := map[string]any{}
	servicePulse := map[string]any{}
	salesJourney := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/pipeline-summary?tenant_slug="+encodedTenantSlug, &pipelineSummary); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/sales-journey?tenant_slug="+encodedTenantSlug, &salesJourney); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.SalesOverviewResponse{
		Service:          handler.ServiceName,
		TenantSlug:       tenantSlug,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary: buildSalesExecutiveSummary(pipelineSummary, salesJourney),
		PipelineSummary:  pipelineSummary,
		ServicePulse:     servicePulse,
		SalesJourney:     salesJourney,
	})
}

func buildSalesExecutiveSummary(pipelineSummary map[string]any, salesJourney map[string]any) dto.SalesExecutiveSummary {
	leadsCaptured := readMapInt(pipelineSummary, "metrics", "leadsCaptured")
	opportunities := readMapInt(salesJourney, "opportunities", "total")
	proposals := readMapInt(salesJourney, "proposals", "total")
	salesWon := readMapInt(salesJourney, "sales", "total")
	bookedRevenueCents := readMapInt(salesJourney, "sales", "bookedRevenueCents")
	completedAutomations := readMapInt(salesJourney, "automation", "runtimeCompleted")
	status := "stable"

	if salesWon == 0 && (opportunities > 0 || proposals > 0) {
		status = "critical"
	} else if proposals > salesWon {
		status = "attention"
	}

	return dto.SalesExecutiveSummary{
		Status:               status,
		LeadsCaptured:        leadsCaptured,
		Opportunities:        opportunities,
		Proposals:            proposals,
		SalesWon:             salesWon,
		BookedRevenueCents:   bookedRevenueCents,
		CompletedAutomations: completedAutomations,
	}
}
