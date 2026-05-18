from __future__ import annotations

from datetime import datetime, timezone


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


METRIC_DEFINITIONS = [
    {
        "code": "revenue.net_operational_margin",
        "domain": "finance",
        "owner": "finance-control",
        "name": "Net operational margin",
        "description": "Receita operacional menos custos, despesas e inadimplencia operacional.",
        "unit": "cents",
        "grain": "tenant/month",
        "formula": "booked_revenue - operational_costs - overdue_risk",
        "sources": ["finance", "billing", "sales"],
        "dimensions": ["tenant", "period", "cost_center"],
        "freshnessTargetMinutes": 60,
        "qualityPolicy": "finance-critical",
        "status": "active",
    },
    {
        "code": "crm.pipeline_conversion_rate",
        "domain": "crm",
        "owner": "relationship-intelligence",
        "name": "Pipeline conversion rate",
        "description": "Percentual de leads convertidos em oportunidades e vendas.",
        "unit": "basis_points",
        "grain": "tenant/day",
        "formula": "won_sales / captured_leads",
        "sources": ["crm", "sales"],
        "dimensions": ["tenant", "source", "territory"],
        "freshnessTargetMinutes": 30,
        "qualityPolicy": "commercial-operational",
        "status": "active",
    },
    {
        "code": "operations.incident_mttr",
        "domain": "platform-control",
        "owner": "incident-command",
        "name": "Mean time to resolution",
        "description": "Tempo medio de resolucao de incidentes por severidade.",
        "unit": "minutes",
        "grain": "tenant/service/severity",
        "formula": "sum(resolution_minutes) / resolved_incidents",
        "sources": ["platform-control", "analytics"],
        "dimensions": ["tenant", "service", "severity"],
        "freshnessTargetMinutes": 15,
        "qualityPolicy": "sre-operational",
        "status": "active",
    },
]

SNAPSHOTS = {
    "revenue.net_operational_margin": [
        {"period": "2026-05", "value": 453000, "quality": "passed", "capturedAt": "2026-05-15T17:00:00Z"},
        {"period": "2026-04", "value": 418000, "quality": "passed", "capturedAt": "2026-04-30T23:00:00Z"},
    ],
    "crm.pipeline_conversion_rate": [
        {"period": "2026-05-15", "value": 937, "quality": "warning", "capturedAt": "2026-05-15T17:00:00Z"},
    ],
    "operations.incident_mttr": [
        {"period": "2026-05", "value": 42, "quality": "passed", "capturedAt": "2026-05-15T17:00:00Z"},
    ],
}

DATASET_FRESHNESS = [
    {"dataset": "finance.receivables", "domain": "finance", "freshnessMinutes": 12, "targetMinutes": 60, "status": "fresh"},
    {"dataset": "crm.pipeline", "domain": "crm", "freshnessMinutes": 18, "targetMinutes": 30, "status": "fresh"},
    {"dataset": "documents.audit_events", "domain": "documents", "freshnessMinutes": 74, "targetMinutes": 60, "status": "attention"},
    {"dataset": "platform.incidents", "domain": "platform-control", "freshnessMinutes": 5, "targetMinutes": 15, "status": "fresh"},
]

DATA_QUALITY_CHECKS = [
    {"checkKey": "finance-ledger-balanced", "domain": "finance", "severity": "critical", "status": "passed", "failedRows": 0},
    {"checkKey": "crm-lead-owner-integrity", "domain": "crm", "severity": "warning", "status": "warning", "failedRows": 2},
    {"checkKey": "documents-audit-redaction", "domain": "documents", "severity": "critical", "status": "passed", "failedRows": 0},
    {"checkKey": "incident-timeline-append-only", "domain": "platform-control", "severity": "critical", "status": "passed", "failedRows": 0},
]


def list_metric_definitions(domain: str | None = None) -> dict:
    normalized_domain = (domain or "").strip()
    items = [item for item in METRIC_DEFINITIONS if not normalized_domain or item["domain"] == normalized_domain]
    return {"generatedAt": utc_now(), "summary": _summary(items), "items": items}


def get_metric_definition(code: str) -> dict | None:
    return next((item for item in METRIC_DEFINITIONS if item["code"] == code), None)


def list_metric_snapshots(code: str) -> dict:
    definition = get_metric_definition(code)
    if definition is None:
        raise ValueError("metric_not_found")
    snapshots = SNAPSHOTS.get(code, [])
    return {"metric": definition, "snapshots": snapshots, "latest": snapshots[0] if snapshots else None}


def list_dataset_freshness() -> dict:
    return {
        "generatedAt": utc_now(),
        "summary": {
            "datasets": len(DATASET_FRESHNESS),
            "fresh": sum(1 for item in DATASET_FRESHNESS if item["status"] == "fresh"),
            "attention": sum(1 for item in DATASET_FRESHNESS if item["status"] == "attention"),
        },
        "items": DATASET_FRESHNESS,
    }


def list_data_quality_checks() -> dict:
    return {
        "generatedAt": utc_now(),
        "summary": {
            "checks": len(DATA_QUALITY_CHECKS),
            "passed": sum(1 for item in DATA_QUALITY_CHECKS if item["status"] == "passed"),
            "warning": sum(1 for item in DATA_QUALITY_CHECKS if item["status"] == "warning"),
            "failed": sum(1 for item in DATA_QUALITY_CHECKS if item["status"] == "failed"),
        },
        "items": DATA_QUALITY_CHECKS,
    }


def build_metric_lineage() -> dict:
    edges = []
    for metric in METRIC_DEFINITIONS:
        for source in metric["sources"]:
            edges.append({"from": source, "to": metric["code"], "type": "source"})
        for dimension in metric["dimensions"]:
            edges.append({"from": metric["code"], "to": f"dimension.{dimension}", "type": "dimension"})
    return {"generatedAt": utc_now(), "nodes": len(METRIC_DEFINITIONS), "edges": edges}


def build_semantic_metrics_readiness() -> dict:
    freshness = list_dataset_freshness()["summary"]
    quality = list_data_quality_checks()["summary"]
    return {
        "acceptanceReady": freshness["attention"] <= 1 and quality["failed"] == 0,
        "controls": ["metric-definitions", "dataset-freshness", "data-quality", "lineage"],
        "metricCount": len(METRIC_DEFINITIONS),
        "freshnessAttention": freshness["attention"],
        "qualityWarnings": quality["warning"],
        "qualityFailures": quality["failed"],
    }


def _summary(items: list[dict]) -> dict:
    domains = sorted({item["domain"] for item in items})
    return {"metrics": len(items), "domains": domains, "active": sum(1 for item in items if item["status"] == "active")}

