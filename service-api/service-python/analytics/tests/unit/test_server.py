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
    assert payload["identityClosure"]["acceptanceReady"] is True
    assert "mfa-enforcement" in payload["identityClosure"]["controls"]


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
    assert payload["documentClosure"]["acceptanceReady"] is True
    assert "version-history" in payload["documentClosure"]["controls"]


def test_engagement_operations_returns_engagement_operational_payload() -> None:
    response = client.get("/api/analytics/reports/engagement-operations?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["campaigns"]["total"] == 2
    assert payload["templates"]["active"] == 2
    assert payload["touchpoints"]["businessLinked"] == 18
    assert payload["deliveries"]["byProvider"]["whatsapp_cloud"] == 8
    assert payload["providers"]["businessLinkedEvents"] == 6
    assert payload["governance"]["activeProviders"] == 4
    assert payload["engagementClosure"]["acceptanceReady"] is True
    assert "provider-callback-idempotency" in payload["engagementClosure"]["controls"]


def test_integration_readiness_returns_external_operations_payload() -> None:
    response = client.get("/api/analytics/reports/integration-readiness?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["providers"]["configured"] == 4
    assert payload["flows"]["inboundLeads"] == 3
    assert payload["flows"]["businessEntityLinkedEvents"] == 6
    assert payload["webhookHub"]["deadLetterEvents"] == 1
    assert payload["capabilityRegistry"]["summary"]["contractArtifacts"] >= 11
    assert payload["readiness"]["callbackTraceabilityReady"] is True
    assert payload["readiness"]["businessEntityLinkageReady"] is True
    assert payload["providerClosure"]["acceptanceReady"] is True
    assert "capability-registry" in payload["providerClosure"]["controls"]


def test_adapter_catalog_returns_provider_and_contract_capabilities() -> None:
    response = client.get("/api/analytics/reports/adapter-catalog")
    payload = response.json()

    assert response.status_code == 200
    assert payload["summary"]["contractArtifacts"] >= 11
    assert payload["engagement"]["summary"]["fallback"] >= 1
    assert payload["documents"]["capabilities"][0]["provider"] == "local"
    assert payload["documentSigning"]["capabilities"][0]["provider"] == "local"
    assert payload["crmEnrichment"]["capabilities"][0]["provider"] == "local"
    assert payload["webhookHub"]["controls"]["dlqReady"] is True


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
    assert payload["workflowClosure"]["acceptanceReady"] is True
    assert "compensation-metadata" in payload["workflowClosure"]["controls"]


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


def test_sales_journey_returns_commercial_closure_payload() -> None:
    response = client.get("/api/analytics/reports/sales-journey?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["salesClosure"]["acceptanceReady"] is True
    assert "commission-governance" in payload["salesClosure"]["controls"]


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
    assert payload["financeClosure"]["acceptanceReady"] is True
    assert "period-closure" in payload["financeClosure"]["controls"]
    assert payload["billingClosure"]["acceptanceReady"] is True
    assert "payment-attempt-idempotency" in payload["billingClosure"]["controls"]


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
    assert payload["rentalClosure"]["acceptanceReady"] is True
    assert "documents-linkage" in payload["rentalClosure"]["controls"]


def test_cost_estimator_returns_scenario_based_cost_payload() -> None:
    response = client.get("/api/analytics/reports/cost-estimator?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["costs"]["estimatedMonthlyCostCents"] == 58820
    assert payload["projection"]["risk"] == "attention"
    assert payload["recommendations"]["needsTeamExpansion"] is True
    assert payload["intelligenceClosure"]["acceptanceReady"] is True
    assert "cost-estimation" in payload["intelligenceClosure"]["controls"]


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
    assert payload["reviews"]["providerCapabilities"]["fallbackCapabilities"] >= 1
    assert payload["reviews"]["contractGovernance"]["httpSpecs"] >= 7
    assert payload["reviews"]["performance"]["latestBenchmarkStatus"] == "attention"
    assert payload["reviews"]["operationalRunbooks"]["acceptanceReady"] is True
    assert "backup-restore" in payload["reviews"]["operationalRunbooks"]["testSuites"]
    assert "permissions" in payload["reviews"]["operationalRunbooks"]["coveredAreas"]
    assert payload["platformClosure"]["acceptanceReady"] is True
    assert "compose-stack" in payload["platformClosure"]["controls"]
    assert payload["hardeningClosure"]["acceptanceReady"] is True
    assert "permission-review" in payload["hardeningClosure"]["controls"]
    assert payload["productionMaturityClosure"]["acceptanceReady"] is True
    assert payload["productionMaturityClosure"]["readModelMode"] == "operational-near-realtime"
    assert payload["productionMaturityClosure"]["gatewayControls"]["cache"] is True
    assert "claims-policy-authorization" in payload["productionMaturityClosure"]["controls"]
    assert "expanded-domain-contracts" in payload["productionMaturityClosure"]["controls"]
    assert "workflow-runtime" in payload["productionMaturityClosure"]["expandedContracts"]


def test_saas_control_returns_usage_and_lifecycle_payload() -> None:
    response = client.get("/api/analytics/reports/saas-control?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["entitlements"]["enabled"] >= 5
    assert payload["readiness"]["onboardingReady"] is True
    assert payload["saasClosure"]["acceptanceReady"] is True
    assert "usage-metering" in payload["saasClosure"]["controls"]


def test_contract_governance_returns_registry_payload() -> None:
    response = client.get("/api/analytics/reports/contract-governance")
    payload = response.json()

    assert response.status_code == 200
    assert payload["catalog"]["httpSpecs"] >= 7
    assert payload["patterns"]["idempotencyKeyReady"] is True
    assert payload["foundationClosure"]["acceptanceReady"] is True
    assert "monorepo-layout" in payload["foundationClosure"]["controls"]
    assert payload["contractClosure"]["acceptanceReady"] is True
    assert "api-portal" in payload["contractClosure"]["controls"]


def test_core_operations_returns_new_product_context_payload() -> None:
    response = client.get("/api/analytics/reports/core-operations?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["summary"]["catalogItems"] == 12
    assert payload["support"]["overdue"] == 1
    assert payload["coreClosure"]["acceptanceReady"] is True
    assert "catalog-consumers" in payload["coreClosure"]["controls"]


def test_relationship_intelligence_returns_pipeline_and_forecast_payload() -> None:
    response = client.get("/api/analytics/reports/relationship-intelligence?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["pipeline"]["configs"] == 1
    assert payload["territories"]["rules"] == 1
    assert payload["approvals"]["policies"] == 1
    assert payload["conversations"]["threads"] == 4
    assert payload["bulkOperations"]["partialSuccessTracking"] is True
    assert payload["forecast"]["confidence"] == "attention"
    assert payload["forecast"]["scenarioCount"] == 2
    assert payload["readiness"]["bulkReady"] is True
    assert payload["crmClosure"]["acceptanceReady"] is True
    assert "deterministic-email-dedup" in payload["crmClosure"]["controls"]
    assert payload["relationshipClosure"]["acceptanceReady"] is True
    assert "forecast-scenarios" in payload["relationshipClosure"]["controls"]


def test_compliance_control_returns_fiscal_and_privacy_payload() -> None:
    response = client.get("/api/analytics/reports/compliance-control?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["fiscal"]["documents"] == 4
    assert "completed" in payload["privacy"]
    assert "executions" in payload["retention"]
    assert "sensitiveOperations" in payload["audit"]
    assert payload["readiness"]["fiscalReady"] is True
    assert payload["complianceClosure"]["acceptanceReady"] is True
    assert "privacy-request-execution" in payload["complianceClosure"]["controls"]


def test_go_live_control_returns_rollout_and_adoption_payload() -> None:
    response = client.get("/api/analytics/reports/go-live-control?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["rollouts"]["completed"] == 1
    assert "adoptionPct" in payload["adoption"]
    assert "recommended" in payload["adjustments"]
    assert "total" in payload["bottlenecks"]
    assert payload["readiness"]["rolloutReady"] is True
    assert payload["readiness"]["rollbackReady"] is True
    assert payload["releaseControls"]["acceptanceReady"] is True
    assert "rollback" in payload["releaseControls"]["controls"]
    assert "hardening" in payload["releaseControls"]["testSuites"]
    assert payload["goLiveClosure"]["acceptanceReady"] is True
