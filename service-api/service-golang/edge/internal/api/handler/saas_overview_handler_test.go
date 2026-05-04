package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

func TestSaaSOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/saas-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewSaaSOverviewHandler("edge", "http://analytics.local", saasReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.SaaSOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.EntitlementsTotal != 8 {
		t.Fatalf("expected eight entitlements, got %d", response.ExecutiveSummary.EntitlementsTotal)
	}

	if response.ExecutiveSummary.ActiveQuotas != 4 {
		t.Fatalf("expected four quotas, got %d", response.ExecutiveSummary.ActiveQuotas)
	}
}

type saasReader struct{}

func (reader saasReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/saas-control"):
		payload = map[string]any{
			"readiness":    map[string]any{"status": "stable"},
			"entitlements": map[string]any{"total": 8, "enabled": 7},
			"quotas":       map[string]any{"active": 4},
			"blocks":       map[string]any{"active": 1},
			"metering":     map[string]any{"trackedMetrics": 5},
			"lifecycle":    map[string]any{"queued": 1, "completed": 3},
		}
	default:
		return fmt.Errorf("unexpected url: %s", requestURL)
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
