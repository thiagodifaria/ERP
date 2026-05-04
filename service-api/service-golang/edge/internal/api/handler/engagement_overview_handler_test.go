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

func TestEngagementOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/engagement-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewEngagementOverviewHandler("edge", "http://analytics.local", engagementOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.EngagementOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.Templates != 2 {
		t.Fatalf("expected templates 2, got %d", response.ExecutiveSummary.Templates)
	}

	if response.ExecutiveSummary.BusinessLinked != 4 {
		t.Fatalf("expected business linked touchpoints 4, got %d", response.ExecutiveSummary.BusinessLinked)
	}

	if response.ExecutiveSummary.ProviderLinkedEvents != 3 {
		t.Fatalf("expected provider linked events 3, got %d", response.ExecutiveSummary.ProviderLinkedEvents)
	}

	if response.ExecutiveSummary.DeliveryRate != 0.75 {
		t.Fatalf("expected delivery rate 0.75, got %f", response.ExecutiveSummary.DeliveryRate)
	}
}

func TestEngagementOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/engagement-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewEngagementOverviewHandler("edge", "http://analytics.local", engagementOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestEngagementOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/engagement-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewEngagementOverviewHandler("edge", "http://analytics.local", engagementOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type engagementOverviewReader struct {
	fail bool
}

func (reader engagementOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"engagement": map[string]any{
					"campaignsTotal":       2,
					"templatesTotal":       2,
					"deliveriesTotal":      4,
					"failedDeliveries":     1,
					"convertedTouchpoints": 1,
				},
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"engagement": map[string]any{
				"campaigns":  2,
				"templates":  2,
				"deliveries": 4,
			},
		}
	case strings.Contains(requestURL, "/engagement-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"campaigns": map[string]any{
				"total":  2,
				"active": 1,
			},
			"templates": map[string]any{
				"total": 2,
			},
			"touchpoints": map[string]any{
				"converted":      1,
				"businessLinked": 4,
			},
			"providers": map[string]any{
				"businessLinkedEvents": 3,
			},
			"deliveries": map[string]any{
				"total":        4,
				"delivered":    3,
				"failed":       1,
				"deliveryRate": 0.75,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
