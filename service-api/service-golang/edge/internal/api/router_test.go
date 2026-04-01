package api

import (
  "context"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func TestRouterShouldExposeHealthDetails(t *testing.T) {
  router := NewRouter(
    telemetry.New("edge-test"),
    handler.NewHealthHandler(
      "edge",
      stubHealthChecker{},
      []integration.ServiceEndpoint{
        {Name: "identity", BaseURL: "http://identity.local"},
        {Name: "crm", BaseURL: "http://crm.local"},
      },
    ),
    handler.NewOpsHandler(
      "edge",
      stubHealthChecker{},
      []integration.ServiceEndpoint{
        {Name: "identity", BaseURL: "http://identity.local"},
        {Name: "crm", BaseURL: "http://crm.local"},
      },
    ),
    handler.NewTenantOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
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

  if response.Service != "edge" {
    t.Fatalf("expected service edge, got %s", response.Service)
  }

  if len(response.Dependencies) != 3 {
    t.Fatalf("expected 3 dependencies, got %d", len(response.Dependencies))
  }
}

func TestRouterShouldExposeOpsHealth(t *testing.T) {
  router := NewRouter(
    telemetry.New("edge-test"),
    handler.NewHealthHandler("edge", stubHealthChecker{}, nil),
    handler.NewOpsHandler(
      "edge",
      stubHealthChecker{},
      []integration.ServiceEndpoint{
        {Name: "identity", BaseURL: "http://identity.local"},
        {Name: "analytics", BaseURL: "http://analytics.local"},
      },
    ),
    handler.NewTenantOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
  )

  request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/health", nil)
  recorder := httptest.NewRecorder()

  router.ServeHTTP(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
  }

  var response dto.OpsHealthResponse
  if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
    t.Fatalf("unexpected decode error: %v", err)
  }

  if response.Service != "edge" {
    t.Fatalf("expected service edge, got %s", response.Service)
  }

  if response.Summary.Total != 2 {
    t.Fatalf("expected summary total 2, got %d", response.Summary.Total)
  }
}

func TestRouterShouldExposeTenantOverview(t *testing.T) {
  router := NewRouter(
    telemetry.New("edge-test"),
    handler.NewHealthHandler("edge", stubHealthChecker{}, nil),
    handler.NewOpsHandler("edge", stubHealthChecker{}, nil),
    handler.NewTenantOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
  )

  request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops", nil)
  recorder := httptest.NewRecorder()

  router.ServeHTTP(recorder, request)

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
}

type stubHealthChecker struct{}

func (stubHealthChecker) Check(_ context.Context, endpoint integration.ServiceEndpoint) dto.DependencyResponse {
  return dto.DependencyResponse{
    Name:   endpoint.Name,
    Status: "ready",
  }
}

func (stubHealthChecker) Details(_ context.Context, endpoint integration.ServiceEndpoint) dto.ServiceHealthSnapshot {
  return dto.ServiceHealthSnapshot{
    Name:   endpoint.Name,
    Status: "ready",
    Dependencies: []dto.DependencyResponse{
      {Name: "router", Status: "ready"},
    },
  }
}

func (stubHealthChecker) GetJSON(_ context.Context, requestURL string, target any) error {
  payload := map[string]any{
    "sourceUrl": requestURL,
    "status":    "ready",
  }
  bytes, err := json.Marshal(payload)
  if err != nil {
    return err
  }

  return json.Unmarshal(bytes, target)
}
