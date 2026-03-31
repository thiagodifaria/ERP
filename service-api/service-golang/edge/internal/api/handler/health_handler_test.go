package handler

import (
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

func TestLiveReturnsEdgeServiceStatus(t *testing.T) {
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

  Ready(recorder, request)

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
