"""Governanca operacional dos contratos HTTP, eventos e API portal."""

from datetime import datetime, timezone
from pathlib import Path


def build_contract_governance() -> dict:
    root = Path(__file__).resolve().parents[5]
    http_specs = sorted((root / "contracts" / "http").glob("*.yaml"))
    event_schemas = sorted((root / "contracts" / "events").glob("*.json"))
    adrs = sorted((root / "docs").glob("ADR-*.md"))
    api_portal = root / "docs" / "API_PORTAL.md"
    registry = root / "contracts" / "registry.json"

    return {
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "catalog": {
            "httpSpecs": len(http_specs),
            "eventSchemas": len(event_schemas),
            "adrs": len(adrs),
            "apiPortalReady": api_portal.exists(),
            "schemaRegistryReady": registry.exists(),
        },
        "patterns": {
            "idempotencyKeyReady": True,
            "accepted202Ready": True,
            "cursorPaginationReady": True,
            "bulkContractsReady": True,
        },
        "artifacts": {
            "http": [item.name for item in http_specs],
            "events": [item.name for item in event_schemas],
            "adrs": [item.name for item in adrs],
        },
        "readiness": {
            "status": "stable" if api_portal.exists() and registry.exists() else "attention",
            "navigableApiReady": api_portal.exists(),
            "registryReady": registry.exists(),
        },
    }
