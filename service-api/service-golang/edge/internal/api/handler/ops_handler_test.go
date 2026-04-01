package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
)

func TestOpsHealthReturnsAggregatedSnapshot(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/health", nil)
	recorder := httptest.NewRecorder()
	handler := NewOpsHandler(
		"edge",
		opsTestChecker{},
		[]integration.ServiceEndpoint{
			{Name: "identity", BaseURL: "http://identity.local"},
			{Name: "analytics", BaseURL: "http://analytics.local"},
		},
	)

	handler.Health(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.OpsHealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Status != "ready" {
		t.Fatalf("expected status ready, got %s", response.Status)
	}

	if response.Summary.Ready != 2 {
		t.Fatalf("expected ready services 2, got %d", response.Summary.Ready)
	}
}

func TestOpsHealthReturnsDegradedSummary(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/health", nil)
	recorder := httptest.NewRecorder()
	handler := NewOpsHandler(
		"edge",
		opsTestChecker{degradedName: "analytics"},
		[]integration.ServiceEndpoint{
			{Name: "identity", BaseURL: "http://identity.local"},
			{Name: "analytics", BaseURL: "http://analytics.local"},
		},
	)

	handler.Health(recorder, request)

	var response dto.OpsHealthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Status != "degraded" {
		t.Fatalf("expected status degraded, got %s", response.Status)
	}

	if response.Summary.Degraded != 1 {
		t.Fatalf("expected degraded services 1, got %d", response.Summary.Degraded)
	}
}

type opsTestChecker struct {
	degradedName string
}

func (checker opsTestChecker) Check(_ context.Context, endpoint integration.ServiceEndpoint) dto.DependencyResponse {
	status := "ready"
	if checker.degradedName == endpoint.Name {
		status = "not_ready"
	}

	return dto.DependencyResponse{Name: endpoint.Name, Status: status}
}

func (checker opsTestChecker) Details(_ context.Context, endpoint integration.ServiceEndpoint) dto.ServiceHealthSnapshot {
	status := "ready"
	if checker.degradedName == endpoint.Name {
		status = "not_ready"
	}

	return dto.ServiceHealthSnapshot{
		Name:   endpoint.Name,
		Status: status,
		Dependencies: []dto.DependencyResponse{
			{Name: "router", Status: "ready"},
		},
	}
}
