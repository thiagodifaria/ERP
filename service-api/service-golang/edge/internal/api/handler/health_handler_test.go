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

func TestLiveReturnsEdgeServiceStatus(t *testing.T) {
  request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
  recorder := httptest.NewRecorder()
  handler := NewHealthHandler("edge", handlerTestChecker{}, nil)

  handler.Live(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
  }

  var response dto.HealthResponse
  if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
    t.Fatalf("unexpected decode error: %v", err)
  }

  if response.Service != "edge" {
    t.Fatalf("expected service edge, got %s", response.Service)
  }

  if response.Status != "live" {
    t.Fatalf("expected status live, got %s", response.Status)
  }
}

func TestReadyReturnsEdgeServiceStatus(t *testing.T) {
  request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
  recorder := httptest.NewRecorder()
  handler := NewHealthHandler("edge", handlerTestChecker{}, nil)

  handler.Ready(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
  }

  var response dto.HealthResponse
  if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
    t.Fatalf("unexpected decode error: %v", err)
  }

  if response.Service != "edge" {
    t.Fatalf("expected service edge, got %s", response.Service)
  }

  if response.Status != "ready" {
    t.Fatalf("expected status ready, got %s", response.Status)
  }
}

func TestDetailsReturnsReadinessPayload(t *testing.T) {
  request := httptest.NewRequest(http.MethodGet, "/health/details", nil)
  recorder := httptest.NewRecorder()
  handler := NewHealthHandler(
    "edge",
    handlerTestChecker{},
    []integration.ServiceEndpoint{
      {Name: "identity", BaseURL: "http://identity.local"},
      {Name: "analytics", BaseURL: "http://analytics.local"},
    },
  )

  handler.Details(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
  }

  var response dto.ReadinessResponse
  if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
    t.Fatalf("unexpected decode error: %v", err)
  }

  if response.Service != "edge" {
    t.Fatalf("expected service edge, got %s", response.Service)
  }

  if response.Status != "ready" {
    t.Fatalf("expected status ready, got %s", response.Status)
  }

  if len(response.Dependencies) != 3 {
    t.Fatalf("expected 3 dependencies, got %d", len(response.Dependencies))
  }
}

func TestDetailsReturnsDegradedWhenDependencyIsDown(t *testing.T) {
  request := httptest.NewRequest(http.MethodGet, "/health/details", nil)
  recorder := httptest.NewRecorder()
  handler := NewHealthHandler(
    "edge",
    handlerTestChecker{degradedName: "analytics"},
    []integration.ServiceEndpoint{
      {Name: "identity", BaseURL: "http://identity.local"},
      {Name: "analytics", BaseURL: "http://analytics.local"},
    },
  )

  handler.Details(recorder, request)

  var response dto.ReadinessResponse
  if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
    t.Fatalf("unexpected decode error: %v", err)
  }

  if response.Status != "degraded" {
    t.Fatalf("expected status degraded, got %s", response.Status)
  }
}

type handlerTestChecker struct {
  degradedName string
}

func (checker handlerTestChecker) Check(_ context.Context, endpoint integration.ServiceEndpoint) dto.DependencyResponse {
  status := "ready"
  if checker.degradedName == endpoint.Name {
    status = "not_ready"
  }

  return dto.DependencyResponse{
    Name:   endpoint.Name,
    Status: status,
  }
}
