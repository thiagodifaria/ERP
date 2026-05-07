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

func TestCoreOperationsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/core-operations?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewCoreOperationsOverviewHandler("edge", "http://analytics.local", coreOperationsReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.CoreOperationsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.CatalogItems != 12 || response.ExecutiveSummary.CriticalNotifications != 2 {
		t.Fatalf("unexpected executive summary payload: %+v", response.ExecutiveSummary)
	}
}

type coreOperationsReader struct{}

func (reader coreOperationsReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}
	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/core-operations"):
		payload = map[string]any{
			"readiness": map[string]any{"status": "attention"},
			"summary": map[string]any{"catalogItems": 12, "suppliers": 6, "supportCases": 8},
			"support": map[string]any{"overdue": 1},
			"notification": map[string]any{"unread": 5, "critical": 2},
		}
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
