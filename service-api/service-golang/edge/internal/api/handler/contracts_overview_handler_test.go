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

func TestContractsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/contracts-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewContractsOverviewHandler("edge", "http://analytics.local", contractsReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.ContractsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.HTTPSpecs != 9 {
		t.Fatalf("expected nine HTTP specs, got %d", response.ExecutiveSummary.HTTPSpecs)
	}

	if !response.ExecutiveSummary.SchemaRegistryReady {
		t.Fatalf("expected schema registry readiness")
	}
}

type contractsReader struct{}

func (reader contractsReader) GetJSON(_ context.Context, requestURL string, target any) error {
	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/integration-readiness"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	case strings.Contains(requestURL, "/contract-governance"):
		payload = map[string]any{
			"catalog": map[string]any{
				"httpSpecs":    9,
				"eventSchemas": 6,
				"adrs":         1,
			},
			"readiness": map[string]any{
				"status":            "stable",
				"navigableApiReady": true,
				"registryReady":     true,
			},
		}
	case strings.Contains(requestURL, "/hardening-review"):
		payload = map[string]any{"tenantSlug": "bootstrap-ops"}
	default:
		return fmt.Errorf("unexpected url: %s", requestURL)
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
