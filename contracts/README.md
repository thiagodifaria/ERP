# contracts

This directory stores versioned integration artifacts for the ERP platform.

Current structure:

- `http/`: OpenAPI documents for public HTTP surfaces
- `events/`: JSON Schemas for cross-service event payloads

Contract rules:

- every artifact is versioned in git with the service code
- public route changes should update the matching OpenAPI file
- cross-service event changes should update the matching schema
- idempotent or async patterns must be documented alongside the HTTP surface
- artifacts here are intended to back future aggregated API UI and schema registry workflows
