package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

func TestRouterShouldExposeHealthDetails(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(http.MethodGet, "/health/details", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("X-Correlation-Id") != "pending-correlation" {
		t.Fatalf("expected fallback correlation id header")
	}

	var response dto.ReadinessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Service != "crm" {
		t.Fatalf("expected service crm, got %s", response.Service)
	}

	if len(response.Dependencies) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(response.Dependencies))
	}
}

func TestRouterShouldExposeLeadList(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 lead, got %d", len(response))
	}
}

func TestRouterShouldExposeLeadSummary(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/summary", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadSummaryResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Total != 1 {
		t.Fatalf("expected total 1, got %d", response.Total)
	}
}

func TestRouterShouldExposeLeadByPublicID(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/"+persistence.BootstrapLeadPublicID, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.PublicID != persistence.BootstrapLeadPublicID {
		t.Fatalf("expected public id %s, got %s", persistence.BootstrapLeadPublicID, response.PublicID)
	}
}

func TestRouterShouldUpdateLeadOwner(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(
		http.MethodPatch,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/owner",
		bytes.NewBufferString(`{"ownerUserId":"0195e7a0-7a9c-7c1f-8a44-4a6e70000024"}`),
	)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.OwnerUserID != "0195e7a0-7a9c-7c1f-8a44-4a6e70000024" {
		t.Fatalf("expected owner 0195e7a0-7a9c-7c1f-8a44-4a6e70000024, got %s", response.OwnerUserID)
	}
}

func TestRouterShouldUpdateLeadProfile(t *testing.T) {
	router := NewRouter(
		telemetry.New("crm-test"),
		persistence.NewInMemoryLeadRepository(),
	)
	request := httptest.NewRequest(
		http.MethodPatch,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID,
		bytes.NewBufferString(`{"name":"Bootstrap Prime","email":"bootstrap.prime@example.com","source":"instagram"}`),
	)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Email != "bootstrap.prime@example.com" {
		t.Fatalf("expected updated email bootstrap.prime@example.com, got %s", response.Email)
	}
}
