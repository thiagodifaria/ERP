# contracts

This directory stores versioned integration artifacts for the ERP platform.

Current structure:

- `http/`: OpenAPI documents for public HTTP surfaces
- `events/`: JSON Schemas for cross-service event payloads
- `registry.json`: simple machine-readable artifact index

Current baseline:

- HTTP specs for `analytics`, `edge`, `crm`, `documents`, `engagement`, `catalog`, `platform-control` and shared operational surfaces
- Event schemas for workflow, platform lifecycle, quotas, CNPJ enrichment, document signing and catalog item publication contracts

Contract rules:

- every artifact is versioned in git with the service code
- public route changes should update the matching OpenAPI file
- cross-service event changes should update the matching schema
- idempotent or async patterns must be documented alongside the HTTP surface
- artifacts here are intended to back future aggregated API UI and schema registry workflows
- `docs/API_PORTAL.md` documents the next step for a central navigable API portal
