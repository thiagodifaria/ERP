package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

func TestTenantOverviewReturnsAnalyticsCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewTenantOverviewHandler("edge", "http://analytics.local", tenantOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.TenantOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}

	if response.PipelineSummary["tenantSlug"] != "bootstrap-ops" {
		t.Fatalf("expected pipeline summary tenant bootstrap-ops")
	}
}

func TestTenantOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewTenantOverviewHandler("edge", "http://analytics.local", tenantOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestTenantOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewTenantOverviewHandler("edge", "http://analytics.local", tenantOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type tenantOverviewReader struct {
	fail bool
}

func (reader tenantOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{
		"requestUrl": requestURL,
		"tenantSlug": "bootstrap-ops",
		"status":     "ready",
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
