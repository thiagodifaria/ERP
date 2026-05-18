"""Bootstrap HTTP do servico analytics."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.security import install_security_middleware
from app.reports.adapter_catalog import build_adapter_catalog
from app.infrastructure.postgres import postgres_ready
from app.reports.automation_board import build_automation_board
from app.reports.compliance_control import build_compliance_control
from app.reports.collections_control import build_collections_control
from app.reports.core_operations import build_core_operations
from app.reports.cost_estimator import build_cost_estimator
from app.reports.document_governance import build_document_governance
from app.reports.delivery_reliability import build_delivery_reliability
from app.reports.engagement_operations import build_engagement_operations
from app.reports.enterprise_runtime import (
    build_enterprise_runtime_fabric_readiness,
    build_financial_close_readiness,
    build_financial_close_snapshot,
    build_lakehouse_lineage,
    build_lakehouse_readiness,
    build_master_data_quality_score,
    build_reconciliation_run,
    close_financial_period,
    create_financial_close_period,
    create_master_data_merge_proposal,
    get_lakehouse_dataset,
    list_data_quality_findings,
    list_data_quality_rules,
    list_financial_close_periods,
    list_lakehouse_datasets,
    list_lakehouse_export_policies,
    list_master_data_duplicates,
    list_master_data_entities,
    list_reconciliation_findings,
)
from app.reports.external_intelligence import (
    build_brazil_registry_enrichment,
    build_document_intelligence_readiness,
    build_external_intelligence_readiness,
    build_external_risk_feed,
    build_fiscal_brazil_readiness,
    build_market_macro_risk,
)
from app.reports.finance_control import build_finance_control
from app.reports.go_live_control import build_go_live_control
from app.reports.hardening_review import build_hardening_review
from app.reports.integration_readiness import build_integration_readiness
from app.reports.load_benchmark import build_load_benchmark
from app.reports.pipeline_summary import build_pipeline_summary
from app.reports.platform_reliability import build_platform_reliability
from app.reports.production_readiness import build_production_readiness
from app.reports.relationship_intelligence import build_relationship_intelligence
from app.reports.rental_operations import build_rental_operations
from app.reports.revenue_operations import build_revenue_operations
from app.reports.risk_compliance import build_compliance_posture, build_domain_scores, build_risk_recommendations, build_service_scores, build_tenant_risk_score
from app.reports.saas_control import build_saas_control
from app.reports.sales_journey import build_sales_journey
from app.reports.semantic_metrics import build_metric_lineage, get_metric_definition, list_data_quality_checks, list_dataset_freshness, list_metric_definitions, list_metric_snapshots
from app.reports.service_pulse import build_service_pulse
from app.reports.tenant_360 import build_tenant_360
from app.reports.workflow_definition_health import build_workflow_definition_health
from app.reports.contract_governance import build_contract_governance
app = FastAPI(title=settings.service_name)
install_security_middleware(app, settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [
        {"name": "report-engine", "status": "ready"},
        {"name": "forecast-model", "status": "ready"},
        {"name": "simulation-catalog", "status": "ready"},
        {"name": "semantic-metrics", "status": "ready"},
        {"name": "risk-compliance-scoring", "status": "ready"},
        {"name": "reconciliation-center", "status": "ready"},
        {"name": "financial-close-center", "status": "ready"},
        {"name": "master-data-quality", "status": "ready"},
        {"name": "lakehouse-manifest", "status": "ready"},
        {"name": "external-intelligence-verification", "status": "ready"},
    ]

    if settings.repository_driver == "postgres":
        dependencies.insert(
            1,
            {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"},
        )
    else:
        dependencies.insert(1, {"name": "warehouse", "status": "pending-runtime-wiring"})

    return {
        "service": settings.service_name,
        "status": "ready",
        "dependencies": dependencies,
    }


@app.get("/api/analytics/reports/pipeline-summary")
def pipeline_summary(tenant_slug: str | None = None) -> dict:
    return build_pipeline_summary(tenant_slug)


@app.get("/api/analytics/reports/service-pulse")
def service_pulse(tenant_slug: str | None = None) -> dict:
    return build_service_pulse(tenant_slug)


@app.get("/api/analytics/reports/sales-journey")
def sales_journey(tenant_slug: str | None = None) -> dict:
    return build_sales_journey(tenant_slug)


@app.get("/api/analytics/reports/tenant-360")
def tenant_360(tenant_slug: str | None = None) -> dict:
    return build_tenant_360(tenant_slug)


@app.get("/api/analytics/reports/automation-board")
def automation_board(tenant_slug: str | None = None) -> dict:
    return build_automation_board(tenant_slug)


@app.get("/api/analytics/reports/workflow-definition-health")
def workflow_definition_health(tenant_slug: str | None = None) -> dict:
    return build_workflow_definition_health(tenant_slug)


@app.get("/api/analytics/reports/delivery-reliability")
def delivery_reliability(provider: str | None = None) -> dict:
    return build_delivery_reliability(provider)


@app.get("/api/analytics/reports/engagement-operations")
def engagement_operations(tenant_slug: str | None = None) -> dict:
    return build_engagement_operations(tenant_slug)


@app.get("/api/analytics/reports/integration-readiness")
def integration_readiness(tenant_slug: str | None = None) -> dict:
    return build_integration_readiness(tenant_slug)


@app.get("/api/analytics/reports/adapter-catalog")
def adapter_catalog() -> dict:
    return build_adapter_catalog()


@app.get("/api/analytics/reports/document-governance")
def document_governance(tenant_slug: str | None = None) -> dict:
    return build_document_governance(tenant_slug)


@app.get("/api/analytics/reports/revenue-operations")
def revenue_operations(tenant_slug: str | None = None) -> dict:
    return build_revenue_operations(tenant_slug)


@app.get("/api/analytics/reports/finance-control")
def finance_control(tenant_slug: str | None = None) -> dict:
    return build_finance_control(tenant_slug)


@app.get("/api/analytics/reports/collections-control")
def collections_control(tenant_slug: str | None = None) -> dict:
    return build_collections_control(tenant_slug)


@app.get("/api/analytics/reports/rental-operations")
def rental_operations(tenant_slug: str | None = None) -> dict:
    return build_rental_operations(tenant_slug)


@app.get("/api/analytics/reports/cost-estimator")
def cost_estimator(tenant_slug: str | None = None) -> dict:
    return build_cost_estimator(tenant_slug)


@app.get("/api/analytics/reports/load-benchmark")
def load_benchmark(tenant_slug: str | None = None) -> dict:
    return build_load_benchmark(tenant_slug)


@app.get("/api/analytics/reports/platform-reliability")
def platform_reliability(tenant_slug: str | None = None) -> dict:
    return build_platform_reliability(tenant_slug)


@app.get("/api/analytics/reports/hardening-review")
def hardening_review(tenant_slug: str | None = None) -> dict:
    return build_hardening_review(tenant_slug)


@app.get("/api/analytics/reports/saas-control")
def saas_control(tenant_slug: str | None = None) -> dict:
    return build_saas_control(tenant_slug)


@app.get("/api/analytics/reports/contract-governance")
def contract_governance() -> dict:
    return build_contract_governance()


@app.get("/api/analytics/reports/core-operations")
def core_operations(tenant_slug: str | None = None) -> dict:
    return build_core_operations(tenant_slug)


@app.get("/api/analytics/reports/relationship-intelligence")
def relationship_intelligence(tenant_slug: str | None = None) -> dict:
    return build_relationship_intelligence(tenant_slug)


@app.get("/api/analytics/reports/compliance-control")
def compliance_control(tenant_slug: str | None = None) -> dict:
    return build_compliance_control(tenant_slug)


@app.get("/api/analytics/reports/go-live-control")
def go_live_control(tenant_slug: str | None = None) -> dict:
    return build_go_live_control(tenant_slug)


@app.get("/api/analytics/reports/production-readiness")
def production_readiness(tenant_slug: str | None = None) -> dict:
    return build_production_readiness(tenant_slug)


@app.get("/api/analytics/metrics")
def metrics(domain: str | None = None) -> dict:
    return list_metric_definitions(domain)


@app.get("/api/analytics/metrics/{code}")
def metric_detail(code: str) -> dict:
    metric = get_metric_definition(code)
    if metric is None:
        raise HTTPException(status_code=404, detail={"code": "metric_not_found", "message": "Metric was not found."})
    return metric


@app.get("/api/analytics/metrics/{code}/snapshots")
def metric_snapshots(code: str) -> dict:
    try:
        return list_metric_snapshots(code)
    except ValueError as error:
        raise HTTPException(status_code=404, detail={"code": str(error), "message": "Metric was not found."}) from error


@app.get("/api/analytics/datasets/freshness")
def dataset_freshness() -> dict:
    return list_dataset_freshness()


@app.get("/api/analytics/data-quality")
def data_quality() -> dict:
    return list_data_quality_checks()


@app.get("/api/analytics/lineage")
def lineage() -> dict:
    return build_metric_lineage()


@app.get("/api/analytics/risk/tenant-score")
def risk_tenant_score(tenant_slug: str | None = None) -> dict:
    return build_tenant_risk_score(tenant_slug)


@app.get("/api/analytics/risk/domain-scores")
def risk_domain_scores(tenant_slug: str | None = None) -> dict:
    return build_domain_scores(tenant_slug)


@app.get("/api/analytics/risk/service-scores")
def risk_service_scores(tenant_slug: str | None = None) -> dict:
    return build_service_scores(tenant_slug)


@app.get("/api/analytics/risk/compliance-posture")
def risk_compliance_posture(tenant_slug: str | None = None) -> dict:
    return build_compliance_posture(tenant_slug)


@app.get("/api/analytics/risk/recommendations")
def risk_recommendations(tenant_slug: str | None = None) -> dict:
    return build_risk_recommendations(tenant_slug)


@app.get("/api/analytics/reconciliation/findings")
def reconciliation_findings(tenant_slug: str | None = None, severity: str | None = None) -> dict:
    return list_reconciliation_findings(tenant_slug, severity)


@app.post("/api/analytics/reconciliation/run")
def reconciliation_run(payload: dict | None = None) -> dict:
    payload = payload or {}
    return build_reconciliation_run(payload.get("tenantSlug"))


@app.get("/api/analytics/financial-close/periods")
def financial_close_periods(tenant_slug: str | None = None) -> dict:
    return list_financial_close_periods(tenant_slug)


@app.post("/api/analytics/financial-close/periods")
def post_financial_close_period(payload: dict) -> dict:
    try:
        return create_financial_close_period(payload.get("tenantSlug"), payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Financial close period payload is invalid."}) from error


@app.post("/api/analytics/financial-close/periods/{public_id}/close")
def post_financial_close(public_id: str, payload: dict | None = None) -> dict:
    payload = payload or {}
    return close_financial_period(public_id, payload.get("tenantSlug"))


@app.get("/api/analytics/financial-close/periods/{public_id}/snapshot")
def financial_close_snapshot(public_id: str, tenant_slug: str | None = None) -> dict:
    return build_financial_close_snapshot(public_id, tenant_slug)


@app.get("/api/analytics/financial-close/readiness")
def financial_close_readiness(tenant_slug: str | None = None) -> dict:
    return build_financial_close_readiness(tenant_slug)


@app.get("/api/analytics/master-data/entities")
def master_data_entities(tenant_slug: str | None = None) -> dict:
    return list_master_data_entities(tenant_slug)


@app.get("/api/analytics/master-data/quality-score")
def master_data_quality_score(tenant_slug: str | None = None) -> dict:
    return build_master_data_quality_score(tenant_slug)


@app.get("/api/analytics/master-data/duplicates")
def master_data_duplicates(tenant_slug: str | None = None) -> dict:
    return list_master_data_duplicates(tenant_slug)


@app.post("/api/analytics/master-data/merge-proposals")
def post_master_data_merge_proposal(payload: dict) -> dict:
    try:
        return create_master_data_merge_proposal(payload.get("tenantSlug"), payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Merge proposal payload is invalid."}) from error


@app.get("/api/analytics/data-quality/rules")
def data_quality_rules() -> dict:
    return list_data_quality_rules()


@app.get("/api/analytics/data-quality/findings")
def data_quality_findings(tenant_slug: str | None = None) -> dict:
    return list_data_quality_findings(tenant_slug)


@app.get("/api/analytics/lakehouse/datasets")
def lakehouse_datasets() -> dict:
    return list_lakehouse_datasets()


@app.get("/api/analytics/lakehouse/datasets/{datasetKey}")
def lakehouse_dataset(datasetKey: str) -> dict:
    dataset = get_lakehouse_dataset(datasetKey)
    if dataset is None:
        raise HTTPException(status_code=404, detail={"code": "dataset_not_found", "message": "Lakehouse dataset was not found."})
    return dataset


@app.get("/api/analytics/lakehouse/lineage")
def lakehouse_lineage() -> dict:
    return build_lakehouse_lineage()


@app.get("/api/analytics/lakehouse/export-policies")
def lakehouse_export_policies() -> dict:
    return list_lakehouse_export_policies()


@app.get("/api/analytics/lakehouse/readiness")
def lakehouse_readiness() -> dict:
    return build_lakehouse_readiness()


@app.get("/api/analytics/enterprise-runtime/readiness")
def enterprise_runtime_fabric_readiness() -> dict:
    return build_enterprise_runtime_fabric_readiness()


@app.get("/api/analytics/external-intelligence/readiness")
def external_intelligence_readiness(tenant_slug: str | None = None) -> dict:
    return build_external_intelligence_readiness(tenant_slug)


@app.get("/api/analytics/document-intelligence/readiness")
def document_intelligence_readiness(tenant_slug: str | None = None) -> dict:
    return build_document_intelligence_readiness(tenant_slug)


@app.get("/api/analytics/fiscal-brazil/readiness")
def fiscal_brazil_readiness(tenant_slug: str | None = None) -> dict:
    return build_fiscal_brazil_readiness(tenant_slug)


@app.get("/api/analytics/registry-enrichment/brazil")
def registry_enrichment_brazil(tenant_slug: str | None = None) -> dict:
    return build_brazil_registry_enrichment(tenant_slug)


@app.get("/api/analytics/market-macro-risk")
def market_macro_risk(tenant_slug: str | None = None) -> dict:
    return build_market_macro_risk(tenant_slug)


@app.get("/api/analytics/external-risk-feed")
def external_risk_feed(tenant_slug: str | None = None) -> dict:
    return build_external_risk_feed(tenant_slug)
