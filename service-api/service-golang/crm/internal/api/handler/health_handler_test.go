package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
)

func TestLiveReturnsCrmServiceStatus(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	recorder := httptest.NewRecorder()

	Live(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Service != "crm" {
		t.Fatalf("expected service crm, got %s", response.Service)
	}

	if response.Status != "live" {
		t.Fatalf("expected status live, got %s", response.Status)
	}
}

func TestReadyReturnsCrmServiceStatus(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	recorder := httptest.NewRecorder()

	Ready(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.HealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Service != "crm" {
		t.Fatalf("expected service crm, got %s", response.Service)
	}

	if response.Status != "ready" {
		t.Fatalf("expected status ready, got %s", response.Status)
	}
}

func TestDetailsReturnsReadinessPayload(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/details", nil)
	recorder := httptest.NewRecorder()

	Details(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.ReadinessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Service != "crm" {
		t.Fatalf("expected service crm, got %s", response.Service)
	}

	if response.Status != "ready" {
		t.Fatalf("expected status ready, got %s", response.Status)
	}

	if len(response.Dependencies) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(response.Dependencies))
	}
}
