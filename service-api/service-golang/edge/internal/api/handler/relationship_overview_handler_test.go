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

func TestRelationshipOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/relationship-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewRelationshipOverviewHandler("edge", "http://analytics.local", relationshipOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.RelationshipOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.AverageLeadScore != 68 || response.ExecutiveSummary.PipelineConfigs != 1 {
		t.Fatalf("unexpected executive summary payload: %+v", response.ExecutiveSummary)
	}
}

type relationshipOverviewReader struct{}

func (reader relationshipOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}
	switch {
	case strings.Contains(requestURL, "/pipeline-summary"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/relationship-intelligence"):
		payload = map[string]any{
			"readiness": map[string]any{"status": "attention"},
			"scoring": map[string]any{"average": 68, "hot": 3},
			"pipeline": map[string]any{"configs": 1, "stages": 5},
			"support": map[string]any{"openCases": 3, "overdueCases": 1},
			"forecast": map[string]any{"weightedPipelineCents": 2314000, "bookedRevenueCents": 1775000},
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
