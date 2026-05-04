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

func TestIntegrationsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/integrations-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewIntegrationsOverviewHandler("edge", "http://analytics.local", integrationsReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.IntegrationsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected attention status, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.ConfiguredProviders != 4 {
		t.Fatalf("expected four configured providers, got %d", response.ExecutiveSummary.ConfiguredProviders)
	}

	if response.ExecutiveSummary.BusinessLinkedEvents != 3 {
		t.Fatalf("expected business linked events 3, got %d", response.ExecutiveSummary.BusinessLinkedEvents)
	}

	if response.ExecutiveSummary.CriticalProviderGaps != 1 {
		t.Fatalf("expected critical provider gaps 1, got %d", response.ExecutiveSummary.CriticalProviderGaps)
	}
}

func TestIntegrationsOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/integrations-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewIntegrationsOverviewHandler("edge", "http://analytics.local", integrationsReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

type integrationsReader struct{}

func (reader integrationsReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/engagement-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"providers":  map[string]any{"responsesTracked": 1},
		}
	case strings.Contains(requestURL, "/integration-readiness"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"providers": map[string]any{
				"configured":              4,
				"activeInboundProviders":  2,
				"activeOutboundProviders": 2,
			},
			"flows": map[string]any{
				"inboundLeads":               1,
				"workflowDispatches":         1,
				"businessEntityLinkedEvents": 3,
				"failedProviderEvents":       1,
			},
			"webhookHub": map[string]any{"deadLetterEvents": 1},
			"capabilityRegistry": map[string]any{
				"summary": map[string]any{
					"criticalUnconfiguredCapabilities": 1,
					"contractArtifacts":                7,
				},
			},
			"readiness": map[string]any{
				"status":            "attention",
				"openProviderRisks": 1,
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
