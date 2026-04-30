// Handler que consolida saude financeira, tesouraria e billing em um unico payload.
// O edge expoe visao executiva do caixa sem duplicar a logica do plano analitico.
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

type FinanceOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewFinanceOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) FinanceOverviewHandler {
	return FinanceOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler FinanceOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	financeControl := map[string]any{}
	revenueOperations := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/finance-control?tenant_slug="+encodedTenantSlug, &financeControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/revenue-operations?tenant_slug="+encodedTenantSlug, &revenueOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.FinanceOverviewResponse{
		Service:           handler.ServiceName,
		TenantSlug:        tenantSlug,
		GeneratedAt:       time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:  buildFinanceExecutiveSummary(financeControl),
		ServicePulse:      servicePulse,
		Tenant360:         tenant360,
		FinanceControl:    financeControl,
		RevenueOperations: revenueOperations,
	})
}

func buildFinanceExecutiveSummary(financeControl map[string]any) dto.FinanceExecutiveSummary {
	currentBalanceCents := readMapInt(financeControl, "treasury", "currentBalanceCents")
	failedPaymentAttempts := readMapInt(financeControl, "billing", "failedAttempts")
	status := "stable"

	if failedPaymentAttempts > 0 || currentBalanceCents < 100000 {
		status = "attention"
	}

	if currentBalanceCents < 0 || failedPaymentAttempts > 3 {
		status = "critical"
	}

	return dto.FinanceExecutiveSummary{
		Status:                       status,
		CurrentBalanceCents:          currentBalanceCents,
		MonthlyRecurringRevenueCents: readMapInt(financeControl, "billing", "monthlyRecurringRevenueCents"),
		ReceivablesPaidCents:         readMapInt(financeControl, "receivables", "paidAmountCents"),
		PayablesPaidCents:            readMapInt(financeControl, "payables", "paidAmountCents"),
		FailedPaymentAttempts:        failedPaymentAttempts,
		ActiveSubscriptions:          readMapInt(financeControl, "billing", "activeSubscriptions"),
		PeriodClosures:               readMapInt(financeControl, "governance", "periodClosures"),
		NetOperationalMarginCents:    readMapInt(financeControl, "profitability", "netOperationalMarginCents"),
	}
}
