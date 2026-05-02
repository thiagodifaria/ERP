# infra

This directory contains runtime infrastructure for local development, containers, observability and deployment.

Planned responsibilities:

- local compose and supporting manifests
- gateway, Keycloak, Kafka, Redis and PostgreSQL runtime definitions
- observability stack wiring
- environment-specific deployment assets

Current local platform bootstrap:

- Keycloak realm import for local identity integration
- OpenFGA local authorization plane bootstrap for relationship-based access experiments
- Kafka single-node local broker for event backbone bootstrap
- Prometheus, Grafana and Blackbox Exporter for baseline observability

Operational conventions:

- infra manifests here are runtime assets, not service business logic
- new local dependencies should expose health or probe coverage when viable
- compose additions should remain container-first and compatible with `scripts/test.sh`
- provider adapters that depend on paid services should degrade explicitly instead of silently mocking critical behavior
