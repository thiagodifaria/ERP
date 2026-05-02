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

func TestHardeningOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/hardening-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewHardeningOverviewHandler("edge", "http://analytics.local", hardeningReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.HardeningOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected attention status, got %s", response.ExecutiveSummary.Status)
	}

	if !response.ExecutiveSummary.BackupRestoreValidated {
		t.Fatalf("expected backup restore validation to be true")
	}
}

func TestHardeningOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/hardening-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewHardeningOverviewHandler("edge", "http://analytics.local", hardeningReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

type hardeningReader struct{}

func (reader hardeningReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/platform-reliability"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"stability": map[string]any{
				"deadLetterEvents":      1,
				"failedPaymentAttempts": 2,
			},
		}
	case strings.Contains(requestURL, "/hardening-review"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"summary": map[string]any{
				"status":          "attention",
				"stableChecks":    6,
				"attentionChecks": 3,
				"criticalChecks":  0,
			},
			"reviews": map[string]any{
				"security":      map[string]any{"auditEvents": 14},
				"performance":   map[string]any{"latestBenchmarkStatus": "stable"},
				"backupRestore": map[string]any{"validated": true},
			},
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
