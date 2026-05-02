# tests

This directory will host cross-service scenarios for the platform.

Current contract:

- `e2e`
- `performance`
- `simulation`
- `resilience`
- `hardening`

Scope rules:

- this directory is reserved for cross-service and platform-level suites
- service-specific unit, integration and contract tests stay inside each service subtree
- new platform suites should be invokable from `scripts/test.sh`
- runtime assertions that span multiple services should prefer this directory or the central smoke flow instead of leaking into one service only

Service-specific tests remain inside each service repository subtree.
