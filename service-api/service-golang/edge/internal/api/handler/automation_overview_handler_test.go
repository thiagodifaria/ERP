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

func TestAutomationOverviewReturnsExecutiveAutomationCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/automation-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewAutomationOverviewHandler("edge", "http://analytics.local", automationOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.AutomationOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.ActiveDefinitions != 1 {
		t.Fatalf("expected active definitions 1, got %d", response.ExecutiveSummary.ActiveDefinitions)
	}

	if response.ExecutiveSummary.ForwardedWebhookEvents != 1 {
		t.Fatalf("expected forwarded webhook events 1, got %d", response.ExecutiveSummary.ForwardedWebhookEvents)
	}
}

func TestAutomationOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/automation-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewAutomationOverviewHandler("edge", "http://analytics.local", automationOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestAutomationOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/automation-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewAutomationOverviewHandler("edge", "http://analytics.local", automationOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type automationOverviewReader struct {
	fail bool
}

func (reader automationOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
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
					"forwarded": 1,
				},
			},
		}
	case strings.Contains(requestURL, "/automation-board"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"catalog": map[string]any{
				"definitionsActive": 1,
			},
			"control": map[string]any{
				"runningRuns": 2,
			},
			"runtime": map[string]any{
				"completedExecutions": 2,
			},
		}
	case strings.Contains(requestURL, "/workflow-definition-health"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"summary": map[string]any{
				"stable":    1,
				"attention": 1,
				"critical":  0,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
