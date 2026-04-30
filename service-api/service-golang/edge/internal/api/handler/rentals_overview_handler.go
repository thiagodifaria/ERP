// Handler que consolida a operacao de locacoes em um unico payload.
// O edge expoe risco de carteira e governanca documental sem duplicar calculos analiticos.
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

type RentalsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewRentalsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) RentalsOverviewHandler {
	return RentalsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler RentalsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	rentalOperations := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/rental-operations?tenant_slug="+encodedTenantSlug, &rentalOperations); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.RentalsOverviewResponse{
		Service:          handler.ServiceName,
		TenantSlug:       tenantSlug,
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary: buildRentalsExecutiveSummary(tenant360, rentalOperations),
		ServicePulse:     servicePulse,
		Tenant360:        tenant360,
		RentalOperations: rentalOperations,
	})
}

func buildRentalsExecutiveSummary(tenant360 map[string]any, rentalOperations map[string]any) dto.RentalsExecutiveSummary {
	overdueCharges := readMapInt(rentalOperations, "charges", "overdue")
	outstandingAmountCents := readMapInt(rentalOperations, "charges", "outstandingAmountCents")
	status := "stable"

	if overdueCharges > 0 {
		status = "attention"
	}
	if overdueCharges > 2 || outstandingAmountCents > 500000 {
		status = "critical"
	}

	return dto.RentalsExecutiveSummary{
		Status:                 status,
		Contracts:              readMapInt(rentalOperations, "contracts", "total"),
		ActiveContracts:        readMapInt(rentalOperations, "contracts", "active"),
		ScheduledCharges:       readMapInt(rentalOperations, "charges", "scheduled"),
		PaidCharges:            readMapInt(rentalOperations, "charges", "paid"),
		CancelledCharges:       readMapInt(rentalOperations, "charges", "cancelled"),
		OutstandingAmountCents: outstandingAmountCents,
		OverdueCharges:         overdueCharges,
		Attachments:            readMapInt(tenant360, "rentals", "attachments"),
	}
}
