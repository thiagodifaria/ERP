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

type ComplianceOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewComplianceOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) ComplianceOverviewHandler {
	return ComplianceOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler ComplianceOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	documentGovernance := map[string]any{}
	complianceControl := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/document-governance?tenant_slug="+encodedTenantSlug, &documentGovernance); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}
	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/compliance-control?tenant_slug="+encodedTenantSlug, &complianceControl); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.ComplianceOverviewResponse{
		Service:            handler.ServiceName,
		TenantSlug:         tenantSlug,
		GeneratedAt:        time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:   buildComplianceExecutiveSummary(complianceControl),
		ServicePulse:       servicePulse,
		DocumentGovernance: documentGovernance,
		ComplianceControl:  complianceControl,
	})
}

func buildComplianceExecutiveSummary(payload map[string]any) dto.ComplianceExecutiveSummary {
	return dto.ComplianceExecutiveSummary{
		Status:                 readMapString(payload, "readiness", "status"),
		FiscalDocuments:        readMapInt(payload, "fiscal", "documents"),
		CancelledDocuments:     readMapInt(payload, "fiscal", "cancelled"),
		DocumentEvents:         readMapInt(payload, "documentOperations", "events"),
		PendingPrivacyRequests: readMapInt(payload, "privacy", "pending"),
		GrantedConsents:        readMapInt(payload, "consents", "granted"),
		RestrictedDocuments:    readMapInt(payload, "retention", "restrictedDocuments"),
		AuditEvents:            readMapInt(payload, "audit", "events"),
	}
}
