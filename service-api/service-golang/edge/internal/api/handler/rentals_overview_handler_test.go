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

func TestRentalsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/rentals-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewRentalsOverviewHandler("edge", "http://analytics.local", rentalsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.RentalsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.Contracts != 2 {
		t.Fatalf("expected contracts 2, got %d", response.ExecutiveSummary.Contracts)
	}

	if response.ExecutiveSummary.Attachments != 1 {
		t.Fatalf("expected attachments 1, got %d", response.ExecutiveSummary.Attachments)
	}
}

func TestRentalsOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/rentals-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewRentalsOverviewHandler("edge", "http://analytics.local", rentalsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestRentalsOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/rentals-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewRentalsOverviewHandler("edge", "http://analytics.local", rentalsOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type rentalsOverviewReader struct {
	fail bool
}

func (reader rentalsOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"rentals": map[string]any{
					"contractsTotal": 2,
				},
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"rentals": map[string]any{
				"attachments": 1,
			},
		}
	case strings.Contains(requestURL, "/rental-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"contracts": map[string]any{
				"total":  2,
				"active": 1,
			},
			"charges": map[string]any{
				"scheduled":              1,
				"paid":                   1,
				"cancelled":              1,
				"outstandingAmountCents": 165000,
				"overdue":                1,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
