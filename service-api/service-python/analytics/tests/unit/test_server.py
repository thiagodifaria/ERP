"""Cobertura inicial do bootstrap e do primeiro relatorio publico."""

from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_live_route_returns_service_health() -> None:
    response = client.get("/health/live")

    assert response.status_code == 200
    assert response.json() == {"service": "analytics", "status": "live"}


def test_details_route_exposes_analytics_dependencies() -> None:
    response = client.get("/health/details")
    payload = response.json()

    assert response.status_code == 200
    assert payload["status"] == "ready"
    assert any(dependency["name"] == "report-engine" for dependency in payload["dependencies"])
    assert any(dependency["name"] == "forecast-model" for dependency in payload["dependencies"])
    assert any(dependency["name"] == "simulation-catalog" for dependency in payload["dependencies"])


def test_pipeline_summary_returns_operational_payload() -> None:
    response = client.get("/api/analytics/reports/pipeline-summary?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["metrics"]["leadsCaptured"] == 128
    assert payload["bySource"]["whatsapp"] == 46
    assert payload["backlog"]["runningAutomations"] == 7


def test_service_pulse_returns_cross_service_payload() -> None:
    response = client.get("/api/analytics/reports/service-pulse?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["services"]["crm"]["totalLeads"] == 128
    assert payload["services"]["sales"]["salesTotal"] == 12
    assert payload["services"]["finance"]["cashAccounts"] == 2
    assert payload["services"]["billing"]["activeSubscriptions"] == 11
    assert payload["services"]["documents"]["attachmentsTotal"] == 24
    assert payload["services"]["documents"]["completedUploadSessions"] == 6
    assert payload["services"]["engagement"]["templatesTotal"] == 3
    assert payload["services"]["rentals"]["contractsTotal"] == 12
    assert payload["services"]["workflowControl"]["activeDefinitions"] == 6
    assert payload["services"]["webhookHub"]["forwarded"] == 87


def test_sales_journey_returns_commercial_vertical_payload() -> None:
    response = client.get("/api/analytics/reports/sales-journey?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["funnel"]["salesWon"] == 12
    assert payload["opportunities"]["total"] == 34
    assert payload["sales"]["bookedRevenueCents"] == 1775000
    assert payload["automation"]["runtimeCompleted"] == 11


def test_tenant_360_returns_tenant_operational_snapshot() -> None:
    response = client.get("/api/analytics/reports/tenant-360?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["identity"]["companies"] == 3
    assert payload["commercial"]["assignedLeads"] == 96
    assert payload["commercial"]["sales"] == 12
    assert payload["engagement"]["deliveries"] == 17
    assert payload["documents"]["attachments"] == 24
    assert payload["documents"]["restrictedAttachments"] == 10
    assert payload["rentals"]["contracts"] == 12
    assert payload["finance"]["cashAccounts"] == 2
    assert payload["billing"]["activeSubscriptions"] == 11
    assert payload["automation"]["workflowRuns"] == 41


def test_document_governance_returns_document_operational_payload() -> None:
    response = client.get("/api/analytics/reports/document-governance?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["inventory"]["attachmentsTotal"] == 24
    assert payload["inventory"]["retentionLongTerm"] == 11
    assert payload["visibility"]["restricted"] == 10
    assert payload["storage"]["drivers"]["r2"] == 9
    assert payload["uploads"]["completed"] == 6
    assert payload["ownership"]["crm.customer"] == 9


def test_engagement_operations_returns_engagement_operational_payload() -> None:
    response = client.get("/api/analytics/reports/engagement-operations?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["campaigns"]["total"] == 2
    assert payload["templates"]["active"] == 2
    assert payload["deliveries"]["byProvider"]["whatsapp_cloud"] == 8
    assert payload["governance"]["activeProviders"] == 4


def test_integration_readiness_returns_external_operations_payload() -> None:
    response = client.get("/api/analytics/reports/integration-readiness?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["providers"]["configured"] == 4
    assert payload["flows"]["inboundLeads"] == 3
    assert payload["webhookHub"]["deadLetterEvents"] == 1
    assert payload["readiness"]["callbackTraceabilityReady"] is True


def test_automation_board_returns_delivery_and_runtime_board() -> None:
    response = client.get("/api/analytics/reports/automation-board?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["catalog"]["definitionsActive"] == 6
    assert payload["runtime"]["completedExecutions"] == 28
    assert payload["runtime"]["byWorkflow"][0]["workflowDefinitionKey"] == "lead-follow-up"
    assert payload["runtime"]["byWorkflow"][0]["retriesTotal"] == 3
    assert payload["control"]["byWorkflow"][0]["workflowDefinitionKey"] == "lead-follow-up"
    assert payload["delivery"]["forwarded"] == 87


def test_workflow_definition_health_returns_definition_level_health() -> None:
    response = client.get("/api/analytics/reports/workflow-definition-health?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["summary"]["definitionsTotal"] == 2
    assert payload["summary"]["stable"] == 1
    assert payload["summary"]["attention"] == 1
    assert payload["definitions"][0]["workflowDefinitionKey"] == "lead-follow-up"
    assert payload["definitions"][0]["health"] == "stable"
    assert payload["definitions"][1]["workflowDefinitionKey"] == "proposal-reminder"
    assert payload["definitions"][1]["attentionReasons"] == ["definition-not-active"]


def test_delivery_reliability_returns_webhook_operational_footprint() -> None:
    response = client.get("/api/analytics/reports/delivery-reliability?provider=stripe")
    payload = response.json()

    assert response.status_code == 200
    assert payload["provider"] == "stripe"
    assert payload["lifecycle"]["totalEvents"] == 93
    assert payload["statusFootprint"]["forwarded"] == 87
    assert payload["providerLeaderboard"][0]["provider"] == "stripe"


def test_revenue_operations_returns_financial_operational_payload() -> None:
    response = client.get("/api/analytics/reports/revenue-operations?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["sales"]["bookedRevenueCents"] == 1775000
    assert payload["invoices"]["paidAmountCents"] == 845000
    assert payload["collections"]["invoiceCoverageRate"] == 0.75
    assert payload["risk"]["overdueInvoices"] == 1


def test_finance_control_returns_treasury_and_billing_payload() -> None:
    response = client.get("/api/analytics/reports/finance-control?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["treasury"]["accountsTotal"] == 3
    assert payload["receivables"]["paidAmountCents"] == 845000
    assert payload["billing"]["activeSubscriptions"] == 11
    assert payload["billing"]["recoveryCasesCritical"] == 2
    assert payload["profitability"]["netOperationalMarginCents"] == 453000
    assert payload["governance"]["failedPaymentAttempts"] == 3
    assert payload["governance"]["recoveryActions"] == 18


def test_collections_control_returns_recovery_operational_payload() -> None:
    response = client.get("/api/analytics/reports/collections-control?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["portfolio"]["casesTotal"] == 12
    assert payload["portfolio"]["criticalCases"] == 4
    assert payload["promises"]["activePromises"] == 2
    assert payload["throughput"]["touchpoints"] == 18
    assert payload["governance"]["pendingActions"] == 3


def test_rental_operations_returns_rental_operational_payload() -> None:
    response = client.get("/api/analytics/reports/rental-operations?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["contracts"]["active"] == 10
    assert payload["charges"]["scheduled"] == 18
    assert payload["charges"]["collectedAmountCents"] == 1665000
    assert payload["governance"]["attachments"] == 7


def test_cost_estimator_returns_scenario_based_cost_payload() -> None:
    response = client.get("/api/analytics/reports/cost-estimator?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["costs"]["estimatedMonthlyCostCents"] == 58820
    assert payload["projection"]["risk"] == "attention"
    assert payload["recommendations"]["needsTeamExpansion"] is True


def test_load_benchmark_returns_recent_performance_payload() -> None:
    response = client.get("/api/analytics/reports/load-benchmark?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["summary"]["totalRuns"] == 2
    assert payload["latest"]["status"] == "attention"
    assert payload["latest"]["p95LatencyMs"] == 332


def test_hardening_review_returns_operational_review_payload() -> None:
    response = client.get("/api/analytics/reports/hardening-review?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["summary"]["status"] == "attention"
    assert payload["reviews"]["security"]["mfaEnabledUsers"] == 2
    assert payload["reviews"]["backupRestore"]["validated"] is True
    assert payload["reviews"]["performance"]["latestBenchmarkStatus"] == "attention"
