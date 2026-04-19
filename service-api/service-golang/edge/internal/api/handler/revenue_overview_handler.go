// Handler que consolida comercial fechado e cobranca em um unico payload.
// O edge expoe risco de recebimento sem duplicar calculos do plano analitico.
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

type RevenueOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewRevenueOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) RevenueOverviewHandler {
	return RevenueOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler RevenueOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	salesJourney := map[string]any{}
	revenueOperations := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/sales-journey?tenant_slug="+encodedTenantSlug, &salesJourney); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/revenue-operations?tenant_slug="+encodedTenantSlug, &revenueOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.RevenueOverviewResponse{
		Service:           handler.ServiceName,
		TenantSlug:        tenantSlug,
		GeneratedAt:       time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:  buildRevenueExecutiveSummary(salesJourney, revenueOperations),
		ServicePulse:      servicePulse,
		SalesJourney:      salesJourney,
		RevenueOperations: revenueOperations,
	})
}

func buildRevenueExecutiveSummary(salesJourney map[string]any, revenueOperations map[string]any) dto.RevenueExecutiveSummary {
	overdueInvoices := readMapInt(revenueOperations, "risk", "overdueInvoices")
	openAmountCents := readMapInt(revenueOperations, "invoices", "openAmountCents")
	status := "stable"

	if overdueInvoices > 0 {
		status = "attention"
	}

	if overdueInvoices > 2 || openAmountCents > 500000 {
		status = "critical"
	}

	return dto.RevenueExecutiveSummary{
		Status:            status,
		SalesWon:          readMapInt(salesJourney, "sales", "total"),
		Invoices:          readMapInt(revenueOperations, "invoices", "total"),
		PaidInvoices:      readMapInt(revenueOperations, "invoices", "byStatus", "paid"),
		OpenAmountCents:   openAmountCents,
		PaidAmountCents:   readMapInt(revenueOperations, "invoices", "paidAmountCents"),
		OverdueInvoices:   overdueInvoices,
		CollectionRateBps: readMapRateBps(revenueOperations, "collections", "collectionRate"),
	}
}

func readMapRateBps(payload map[string]any, path ...string) int {
	return int(readMapFloat(payload, path...) * 10000)
}

func readMapFloat(payload map[string]any, path ...string) float64 {
	var current any = payload

	for _, key := range path {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return 0
		}

		nextValue, ok := currentMap[key]
		if !ok {
			return 0
		}

		current = nextValue
	}

	switch value := current.(type) {
	case float64:
		return value
	case float32:
		return float64(value)
	case int:
		return float64(value)
	case int64:
		return float64(value)
	case json.Number:
		parsed, err := value.Float64()
		if err == nil {
			return parsed
		}
	}

	return 0
}
