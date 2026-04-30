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

func TestDocumentsOverviewReturnsExecutiveCockpit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/documents-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewDocumentsOverviewHandler("edge", "http://analytics.local", documentsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.DocumentsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.ExecutiveSummary.Status != "attention" {
		t.Fatalf("expected status attention, got %s", response.ExecutiveSummary.Status)
	}

	if response.ExecutiveSummary.AttachmentsTotal != 6 {
		t.Fatalf("expected attachments total 6, got %d", response.ExecutiveSummary.AttachmentsTotal)
	}

	if response.ExecutiveSummary.ExternalStorageAssets != 2 {
		t.Fatalf("expected external storage assets 2, got %d", response.ExecutiveSummary.ExternalStorageAssets)
	}
}

func TestDocumentsOverviewRequiresTenantSlug(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/documents-overview", nil)
	recorder := httptest.NewRecorder()
	handler := NewDocumentsOverviewHandler("edge", "http://analytics.local", documentsOverviewReader{})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestDocumentsOverviewReturnsBadGatewayWhenAnalyticsFails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/documents-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()
	handler := NewDocumentsOverviewHandler("edge", "http://analytics.local", documentsOverviewReader{fail: true})

	handler.Overview(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, recorder.Code)
	}
}

type documentsOverviewReader struct {
	fail bool
}

func (reader documentsOverviewReader) GetJSON(_ context.Context, requestURL string, target any) error {
	if reader.fail {
		return fmt.Errorf("boom")
	}

	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"documents": map[string]any{
					"attachmentsTotal":       6,
					"completedUploadSessions": 1,
				},
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"documents": map[string]any{
				"attachments":          6,
				"restrictedAttachments": 2,
			},
		}
	case strings.Contains(requestURL, "/document-governance"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"inventory": map[string]any{
				"attachmentsTotal":    6,
				"activeAttachments":   5,
				"archivedAttachments": 1,
				"retentionLongTerm":   4,
			},
			"visibility": map[string]any{
				"restricted": 2,
			},
			"storage": map[string]any{
				"external": 2,
			},
			"uploads": map[string]any{
				"pending":   0,
				"completed": 1,
			},
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}
