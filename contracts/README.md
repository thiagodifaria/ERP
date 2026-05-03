# contracts

This directory stores versioned integration artifacts for the ERP platform.

Current structure:

- `http/`: OpenAPI documents for public HTTP surfaces
- `events/`: JSON Schemas for cross-service event payloads

Current baseline:

- HTTP specs for `identity`, `crm`, `sales`, `workflow-control`, `catalog`, `platform-control` and shared operational surfaces
- Event schemas for workflow, platform lifecycle and catalog item publication contracts

Contract rules:

- every artifact is versioned in git with the service code
- public route changes should update the matching OpenAPI file
- cross-service event changes should update the matching schema
- idempotent or async patterns must be documented alongside the HTTP surface
- artifacts here are intended to back future aggregated API UI and schema registry workflows
