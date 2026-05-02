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

func TestPlatformReliabilityOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/platform-reliability?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", platformReliabilityReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.PlatformReliabilityOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected attention status, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.DeadLetterEvents != 1 {
		t.Fatalf("expected one dead-letter event, got %d", response.ExecutiveSummary.DeadLetterEvents)
	}

	if response.ExecutiveSummary.WebhookForwardingRateBps != 9500 {
		t.Fatalf("expected webhook forwarding rate 9500 bps, got %d", response.ExecutiveSummary.WebhookForwardingRateBps)
	}
}

func TestPlatformReliabilityOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/platform-reliability", nil)
	recorder := httptest.NewRecorder()
	handler := NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", platformReliabilityReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestPlatformReliabilityOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/platform-reliability?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", platformReliabilityReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type platformReliabilityReader struct {
	fail bool
}

func (reader platformReliabilityReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"webhookHub": map[string]any{
					"totalEvents": 20,
					"deadLetter":  1,
				},
			},
		}
	case strings.Contains(requestURL, "/delivery-reliability"):
		payload = map[string]any{
			"lifecycle": map[string]any{
				"totalEvents":   20,
				"handledEvents": 19,
			},
		}
	case strings.Contains(requestURL, "/platform-reliability"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"stability": map[string]any{
				"status":                   "attention",
				"pendingWebhookEvents":     1,
				"deadLetterEvents":         1,
				"failedWorkflowExecutions": 1,
				"criticalRecoveryCases":    1,
				"failedPaymentAttempts":    2,
			},
			"serviceLevelObjectives": map[string]any{
				"webhookForwardingRate": 0.95,
				"workflowSuccessRate":   0.9,
				"billingRecoveryRate":   0.8,
			},
			"safeguards": map[string]any{
				"openCriticalRisks": 1,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
