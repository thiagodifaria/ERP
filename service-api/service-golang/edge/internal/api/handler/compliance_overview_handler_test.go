package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

func TestComplianceOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/compliance-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewComplianceOverviewHandler("edge", "http://analytics.local", complianceOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.ComplianceOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.FiscalDocuments != 4 || response.ExecutiveSummary.GrantedConsents != 3 {
		t.Fatalf("unexpected executive summary payload: %+v", response.ExecutiveSummary)
	}
}

type complianceOverviewReader struct{}

func (reader complianceOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}
	switch {
	case strings.Contains(requestURL, "/compliance-control"):
		payload = map[string]any{
			"readiness": map[string]any{"status": "attention"},
			"fiscal": map[string]any{"documents": 4, "cancelled": 1},
			"documentOperations": map[string]any{"events": 6},
			"privacy": map[string]any{"pending": 1},
			"consents": map[string]any{"granted": 3},
			"retention": map[string]any{"restrictedDocuments": 2},
			"audit": map[string]any{"events": 10},
		}
	default:
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
