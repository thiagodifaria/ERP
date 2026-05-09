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

func TestGoLiveOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/go-live-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewGoLiveOverviewHandler("edge", "http://analytics.local", goLiveOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.GoLiveOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if !response.ExecutiveSummary.RolloutReady || !response.ExecutiveSummary.RollbackReady || response.ExecutiveSummary.CompletedRollouts != 1 {
		t.Fatalf("unexpected executive summary payload: %+v", response.ExecutiveSummary)
	}
	if response.ExecutiveSummary.AdoptionPct != 84 || response.ExecutiveSummary.PendingAdjustments != 1 || !response.ExecutiveSummary.AcceptanceReady {
		t.Fatalf("unexpected go-live adjustments payload: %+v", response.ExecutiveSummary)
	}
}

type goLiveOverviewReader struct{}

func (reader goLiveOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}
	switch {
	case strings.Contains(requestURL, "/go-live-control"):
		payload = map[string]any{
			"readiness":       map[string]any{"status": "stable", "rolloutReady": true, "rollbackReady": true, "metricsObserved": true},
			"rollouts":        map[string]any{"planned": 0, "running": 0, "completed": 1, "rolledBack": 0},
			"adoption":        map[string]any{"trackedMetrics": 4, "totalQuantity": 4096, "adoptionPct": 84},
			"adjustments":     map[string]any{"recommended": 1},
			"bottlenecks":     map[string]any{"critical": 0},
			"releaseControls": map[string]any{"acceptanceReady": true},
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
