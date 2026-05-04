# contracts

This directory stores versioned integration artifacts for the ERP platform.

Current structure:

- `http/`: OpenAPI documents for public HTTP surfaces
- `events/`: JSON Schemas for cross-service event payloads
- `registry.json`: simple machine-readable artifact index
- `schema-registry.json`: materialized event schema registry
- `portal/`: central navigable API portal baseline
- `../docs/VERSIONAMENTO_CONTRATOS.md`: compatibility and versioning baseline

Current baseline:

- HTTP specs for `analytics`, `billing`, `catalog`, `crm`, `documents`, `edge`, `engagement`, `finance`, `identity`, `platform-control`, `rentals`, `sales`, `simulation`, `webhook-hub`, `workflow-control` and `workflow-runtime`
- Event schemas for inbound and outbound webhooks, platform lifecycle, quotas, CNPJ enrichment, document signing and catalog item publication contracts

Contract rules:

- every artifact is versioned in git with the service code
- public route changes should update the matching OpenAPI file
- cross-service event changes should update the matching schema
- idempotent or async patterns must be documented alongside the HTTP surface
- cursor pagination and partial success should be explicit whenever the route follows those patterns
- artifacts here are intended to back future aggregated API UI and schema registry workflows
- `contracts/portal/index.html` is the versioned baseline for a central navigable API portal
