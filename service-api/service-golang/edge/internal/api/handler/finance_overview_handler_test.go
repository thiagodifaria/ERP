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

func TestFinanceOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/finance-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewFinanceOverviewHandler("edge", "http://analytics.local", financeOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.FinanceOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.CurrentBalanceCents != 383000 {
		t.Fatalf("expected balance 383000, got %d", response.ExecutiveSummary.CurrentBalanceCents)
	}

	if response.ExecutiveSummary.ActiveSubscriptions != 1 {
		t.Fatalf("expected active subscriptions 1, got %d", response.ExecutiveSummary.ActiveSubscriptions)
	}
}

func TestFinanceOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/finance-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewFinanceOverviewHandler("edge", "http://analytics.local", financeOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestFinanceOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/finance-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewFinanceOverviewHandler("edge", "http://analytics.local", financeOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type financeOverviewReader struct {
	fail bool
}

func (reader financeOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"finance": map[string]any{
					"currentBalanceCents": 383000,
				},
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"finance": map[string]any{
				"cashAccounts": 1,
			},
		}
	case strings.Contains(requestURL, "/finance-control"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"treasury": map[string]any{
				"currentBalanceCents": 383000,
			},
			"billing": map[string]any{
				"activeSubscriptions":          1,
				"monthlyRecurringRevenueCents": 4900,
				"failedAttempts":               2,
			},
			"receivables": map[string]any{
				"paidAmountCents": 230000,
			},
			"payables": map[string]any{
				"paidAmountCents": 19000,
			},
			"profitability": map[string]any{
				"netOperationalMarginCents": 179600,
			},
			"governance": map[string]any{
				"periodClosures": 1,
			},
		}
	case strings.Contains(requestURL, "/revenue-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"invoices": map[string]any{
				"paidAmountCents": 105000,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
