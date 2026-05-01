// Handler que consolida collections e recovery em um cockpit executivo unico.
// O edge agrega cobranca critica, promessas e recuperacao sem reimplementar o plano analitico.
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

type CollectionsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewCollectionsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) CollectionsOverviewHandler {
	return CollectionsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler CollectionsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	financeControl := map[string]any{}
	collectionsControl := map[string]any{}

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

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/collections-control?tenant_slug="+encodedTenantSlug, &collectionsControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.CollectionsOverviewResponse{
		Service:            handler.ServiceName,
		TenantSlug:         tenantSlug,
		GeneratedAt:        time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:   buildCollectionsExecutiveSummary(financeControl, collectionsControl),
		ServicePulse:       servicePulse,
		Tenant360:          tenant360,
		FinanceControl:     financeControl,
		CollectionsControl: collectionsControl,
	})
}

func buildCollectionsExecutiveSummary(financeControl map[string]any, collectionsControl map[string]any) dto.CollectionsExecutiveSummary {
	criticalCases := readMapInt(collectionsControl, "portfolio", "criticalCases")
	invoicesInRecovery := readMapInt(collectionsControl, "invoices", "invoicesInRecovery")
	nextActionsDue := readMapInt(collectionsControl, "governance", "nextActionsDue")
	failedPaymentAttempts := readMapInt(financeControl, "billing", "failedAttempts")
	status := "stable"

	if criticalCases > 0 || invoicesInRecovery > 0 || failedPaymentAttempts > 0 {
		status = "attention"
	}

	if criticalCases > 2 || invoicesInRecovery > 2 || nextActionsDue > 2 || failedPaymentAttempts > 3 {
		status = "critical"
	}

	return dto.CollectionsExecutiveSummary{
		Status:                status,
		CasesTotal:            readMapInt(collectionsControl, "portfolio", "casesTotal"),
		CriticalCases:         criticalCases,
		InvoicesInRecovery:    invoicesInRecovery,
		OpenAmountCents:       readMapInt(collectionsControl, "portfolio", "openAmountCents"),
		RecoveredAmountCents:  readMapInt(collectionsControl, "portfolio", "recoveredAmountCents"),
		FailedPaymentAttempts: failedPaymentAttempts,
		ActivePromises:        readMapInt(collectionsControl, "promises", "activePromises"),
		NextActionsDue:        nextActionsDue,
		RecoveryRateBps:       int(readMapFloat(collectionsControl, "throughput", "recoveryRate") * 10000),
	}
}
