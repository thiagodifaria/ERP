package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func buildTestRouter() http.Handler {
	return NewRouter(
		telemetry.New("edge-test"),
		handler.NewHealthHandler("edge", stubHealthChecker{}, nil),
		handler.NewOpsHandler("edge", stubHealthChecker{}, nil),
		handler.NewTenantOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewAutomationOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewEngagementOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewIntegrationsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewDocumentsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewCollectionsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewHardeningOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewSalesOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRevenueOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewFinanceOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRentalsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		"http://identity.local",
		stubAccessResolver{},
	)
}

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
		handler.NewAutomationOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewEngagementOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewIntegrationsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewDocumentsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewCollectionsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewHardeningOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewSalesOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRevenueOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewFinanceOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRentalsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		"http://identity.local",
		stubAccessResolver{},
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
		handler.NewAutomationOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewEngagementOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewIntegrationsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewDocumentsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewCollectionsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewPlatformReliabilityOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewHardeningOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewSalesOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRevenueOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewFinanceOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		handler.NewRentalsOverviewHandler("edge", "http://analytics.local", stubHealthChecker{}),
		"http://identity.local",
		stubAccessResolver{},
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
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
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

func TestRouterShouldExposeAutomationOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/automation-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.AutomationOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeSalesOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/sales-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.SalesOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeEngagementOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/engagement-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.EngagementOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeIntegrationsOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/integrations-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.IntegrationsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeDocumentsOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/documents-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.DocumentsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeCollectionsOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/collections-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.CollectionsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeHardeningOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/hardening-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.HardeningOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeRevenueOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/revenue-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.RevenueOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeFinanceOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/finance-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.FinanceOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldExposeRentalsOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/rentals-overview?tenantSlug=bootstrap-ops", nil)
	request.Header.Set("Authorization", "Bearer session-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.RentalsOverviewResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected tenant slug bootstrap-ops, got %s", response.TenantSlug)
	}
}

func TestRouterShouldRequireSessionForProtectedTenantOverview(t *testing.T) {
	router := buildTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
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
	payload := map[string]any{}

	switch {
	case strings.Contains(requestURL, "/service-pulse"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"services": map[string]any{
				"webhookHub": map[string]any{
					"forwarded": 1,
				},
				"rentals": map[string]any{
					"contractsTotal": 1,
				},
			},
		}
	case strings.Contains(requestURL, "/automation-board"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"catalog": map[string]any{
				"definitionsActive": 1,
			},
			"control": map[string]any{
				"runningRuns": 2,
			},
			"runtime": map[string]any{
				"completedExecutions": 2,
			},
		}
	case strings.Contains(requestURL, "/workflow-definition-health"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"summary": map[string]any{
				"stable":    1,
				"attention": 1,
				"critical":  0,
			},
		}
	case strings.Contains(requestURL, "/sales-journey"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"opportunities": map[string]any{
				"total": 2,
			},
			"proposals": map[string]any{
				"total": 2,
			},
			"sales": map[string]any{
				"total":              1,
				"bookedRevenueCents": 125000,
			},
			"engagement": map[string]any{
				"campaignsTotal":       2,
				"templatesTotal":       2,
				"deliveriesTotal":      4,
				"failedDeliveries":     1,
				"convertedTouchpoints": 1,
			},
			"automation": map[string]any{
				"runtimeCompleted": 2,
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
				"converted": 1,
			},
			"deliveries": map[string]any{
				"total":        4,
				"delivered":    3,
				"failed":       1,
				"deliveryRate": 0.75,
			},
		}
	case strings.Contains(requestURL, "/revenue-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"invoices": map[string]any{
				"total":           2,
				"openAmountCents": 125000,
				"paidAmountCents": 99000,
				"byStatus": map[string]any{
					"paid": 1,
				},
			},
			"collections": map[string]any{
				"collectionRate": 0.442,
			},
			"risk": map[string]any{
				"overdueInvoices": 1,
			},
		}
	case strings.Contains(requestURL, "/finance-control"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"billing": map[string]any{
				"failedAttempts": 2,
			},
		}
	case strings.Contains(requestURL, "/collections-control"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"portfolio": map[string]any{
				"casesTotal":           1,
				"criticalCases":        1,
				"openAmountCents":      0,
				"recoveredAmountCents": 4900,
			},
			"invoices": map[string]any{
				"invoicesInRecovery": 0,
			},
			"promises": map[string]any{
				"activePromises": 0,
			},
			"throughput": map[string]any{
				"recoveryRate": 1.0,
			},
			"governance": map[string]any{
				"nextActionsDue": 0,
			},
		}
	case strings.Contains(requestURL, "/rental-operations"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"contracts": map[string]any{
				"total":  1,
				"active": 0,
			},
			"charges": map[string]any{
				"scheduled":              1,
				"paid":                   1,
				"cancelled":              1,
				"outstandingAmountCents": 165000,
				"overdue":                0,
			},
		}
	case strings.Contains(requestURL, "/tenant-360"):
		payload = map[string]any{
			"tenantSlug": "bootstrap-ops",
			"identity": map[string]any{
				"companies": 1,
			},
			"commercial": map[string]any{
				"leads":          1,
				"qualifiedLeads": 1,
				"assignedLeads":  1,
				"leadNotes":      2,
				"opportunities":  1,
				"proposals":      1,
				"sales":          1,
			},
			"engagement": map[string]any{
				"campaigns":  2,
				"templates":  2,
				"deliveries": 4,
			},
			"rentals": map[string]any{
				"attachments": 1,
			},
			"automation": map[string]any{
				"workflowRuns":      1,
				"workflowRunEvents": 2,
			},
		}
	default:
		payload = map[string]any{
			"sourceUrl": requestURL,
			"status":    "ready",
		}
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

type stubAccessResolver struct{}

func (stubAccessResolver) ResolveTenantAccess(_ context.Context, _ string, tenantSlug string, _ string) (integration.AccessResolution, error) {
	return integration.AccessResolution{
		TenantSlug:   tenantSlug,
		UserPublicID: "01960d76-3c95-7c85-9f4e-b5b794f7a001",
		RoleCodes:    []string{"owner"},
		Authorized:   true,
		Status:       "active",
	}, nil
}
