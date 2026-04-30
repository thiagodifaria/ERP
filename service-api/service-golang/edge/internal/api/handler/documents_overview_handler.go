// Handler que consolida a governanca documental em um unico payload.
// O edge expoe inventario, visibilidade e lifecycle de upload com leitura executiva pronta.
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

type DocumentsOverviewHandler struct {
	ServiceName      string
	AnalyticsBaseURL string
	Reader           integration.JSONReader
}

func NewDocumentsOverviewHandler(serviceName string, analyticsBaseURL string, reader integration.JSONReader) DocumentsOverviewHandler {
	return DocumentsOverviewHandler{
		ServiceName:      serviceName,
		AnalyticsBaseURL: strings.TrimRight(analyticsBaseURL, "/"),
		Reader:           reader,
	}
}

func (handler DocumentsOverviewHandler) Overview(writer http.ResponseWriter, request *http.Request) {
	tenantSlug, ok := requiredTenantSlug(writer, request)
	if !ok {
		return
	}

	encodedTenantSlug := url.QueryEscape(tenantSlug)
	servicePulse := map[string]any{}
	tenant360 := map[string]any{}
	documentGovernance := map[string]any{}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/service-pulse?tenant_slug="+encodedTenantSlug, &servicePulse); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/tenant-360?tenant_slug="+encodedTenantSlug, &tenant360); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	if err := handler.Reader.GetJSON(request.Context(), handler.AnalyticsBaseURL+"/api/analytics/reports/document-governance?tenant_slug="+encodedTenantSlug, &documentGovernance); err != nil {
		respondAnalyticsDependencyError(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.DocumentsOverviewResponse{
		Service:            handler.ServiceName,
		TenantSlug:         tenantSlug,
		GeneratedAt:        time.Now().UTC().Format(time.RFC3339),
		ExecutiveSummary:   buildDocumentsExecutiveSummary(tenant360, documentGovernance),
		ServicePulse:       servicePulse,
		Tenant360:          tenant360,
		DocumentGovernance: documentGovernance,
	})
}

func buildDocumentsExecutiveSummary(tenant360 map[string]any, documentGovernance map[string]any) dto.DocumentsExecutiveSummary {
	pendingUploads := readMapInt(documentGovernance, "uploads", "pending")
	restrictedAttachments := readMapInt(documentGovernance, "visibility", "restricted")
	status := "stable"

	if pendingUploads > 0 || restrictedAttachments > 0 {
		status = "attention"
	}
	if pendingUploads > 2 || restrictedAttachments > 5 {
		status = "critical"
	}

	return dto.DocumentsExecutiveSummary{
		Status:                status,
		AttachmentsTotal:      readMapInt(documentGovernance, "inventory", "attachmentsTotal"),
		ActiveAttachments:     readMapInt(documentGovernance, "inventory", "activeAttachments"),
		ArchivedAttachments:   readMapInt(documentGovernance, "inventory", "archivedAttachments"),
		RestrictedAttachments: restrictedAttachments,
		PendingUploads:        pendingUploads,
		CompletedUploads:      readMapInt(documentGovernance, "uploads", "completed"),
		ExternalStorageAssets: readMapInt(documentGovernance, "storage", "external"),
		LongTermRetention:     readMapInt(documentGovernance, "inventory", "retentionLongTerm"),
	}
}
