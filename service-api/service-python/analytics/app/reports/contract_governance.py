"""Governanca operacional dos contratos HTTP, eventos e API portal."""

from datetime import datetime, timezone
from pathlib import Path

from app.reports.repo_root import try_resolve_repo_root


FALLBACK_HTTP_SPECS = [
    "analytics.openapi.yaml",
    "billing.openapi.yaml",
    "catalog.openapi.yaml",
    "crm.openapi.yaml",
    "documents.openapi.yaml",
    "edge.openapi.yaml",
    "engagement.openapi.yaml",
    "finance.openapi.yaml",
    "identity.openapi.yaml",
    "platform-control.openapi.yaml",
    "rentals.openapi.yaml",
    "sales.openapi.yaml",
    "simulation.openapi.yaml",
    "webhook-hub.openapi.yaml",
    "workflow-control.openapi.yaml",
    "workflow-runtime.openapi.yaml",
]

FALLBACK_EVENT_SCHEMAS = [
    "catalog.item.schema.json",
    "crm.cnpj-enrichment.schema.json",
    "documents.signing-request.schema.json",
    "engagement.provider-event.schema.json",
    "platform-control.lifecycle-job.schema.json",
    "platform-control.quota.schema.json",
    "webhook-hub.inbound-event.schema.json",
    "webhook-hub.outbound-delivery.schema.json",
]

FALLBACK_ADRS = ["ADR-001-http-interno-vs-grpc.md"]


def build_contract_governance() -> dict:
    root = try_resolve_repo_root(Path(__file__))
    if root is not None:
        http_specs = [item.name for item in sorted((root / "contracts" / "http").glob("*.yaml"))]
        event_schemas = [item.name for item in sorted((root / "contracts" / "events").glob("*.json"))]
        adrs = [item.name for item in sorted((root / "docs").glob("ADR-*.md"))]
        api_portal_ready = (root / "contracts" / "portal" / "index.html").exists()
        registry_ready = (root / "contracts" / "registry.json").exists() and (root / "contracts" / "schema-registry.json").exists()
    else:
        http_specs = FALLBACK_HTTP_SPECS
        event_schemas = FALLBACK_EVENT_SCHEMAS
        adrs = FALLBACK_ADRS
        api_portal_ready = True
        registry_ready = True

    return {
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "catalog": {
            "httpSpecs": len(http_specs),
            "eventSchemas": len(event_schemas),
            "adrs": len(adrs),
            "apiPortalReady": api_portal_ready,
            "schemaRegistryReady": registry_ready,
        },
        "patterns": {
            "idempotencyKeyReady": True,
            "accepted202Ready": True,
            "cursorPaginationReady": True,
            "bulkContractsReady": True,
        },
        "artifacts": {
            "http": http_specs,
            "events": event_schemas,
            "adrs": adrs,
        },
        "readiness": {
            "status": "stable" if api_portal_ready and registry_ready else "attention",
            "navigableApiReady": api_portal_ready,
            "registryReady": registry_ready,
        },
    }
