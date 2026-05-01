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

func TestCollectionsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/collections-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewCollectionsOverviewHandler("edge", "http://analytics.local", collectionsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.CollectionsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.CasesTotal != 1 {
		t.Fatalf("expected cases total 1, got %d", response.ExecutiveSummary.CasesTotal)
	}

	if response.ExecutiveSummary.RecoveryRateBps != 10000 {
		t.Fatalf("expected recovery rate bps 10000, got %d", response.ExecutiveSummary.RecoveryRateBps)
	}
}

func TestCollectionsOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/collections-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewCollectionsOverviewHandler("edge", "http://analytics.local", collectionsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestCollectionsOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/collections-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewCollectionsOverviewHandler("edge", "http://analytics.local", collectionsOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type collectionsOverviewReader struct {
	fail bool
}

func (reader collectionsOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"billing": map[string]any{
					"failedAttempts": 2,
				},
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"billing": map[string]any{
				"recoveryCases": 1,
			},
		}
	case strings.Contains(requestURL, "/finance-control"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"billing": map[string]any{
				"failedAttempts": 2,
			},
		}
	case strings.Contains(requestURL, "/collections-control"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"portfolio": map[string]any{
				"casesTotal":           1,
				"criticalCases":        1,
				"openAmountCents":      0,
				"recoveredAmountCents": 4900,
			},
			"invoices": map[string]any{
				"invoicesInRecovery": 0,
			},
			"promises": map[string]any{
				"activePromises": 0,
			},
			"throughput": map[string]any{
				"recoveryRate": 1.0,
			},
			"governance": map[string]any{
				"nextActionsDue": 0,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
