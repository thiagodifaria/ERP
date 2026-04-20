# infra

This directory contains runtime infrastructure for local development, containers, observability and deployment.

Planned responsibilities:

- local compose and supporting manifests
- gateway, Keycloak, Kafka, Redis and PostgreSQL runtime definitions
- observability stack wiring
- environment-specific deployment assets

Current local platform bootstrap:

- Keycloak realm import for local identity integration
- Kafka single-node local broker for event backbone bootstrap
- Prometheus, Grafana and Blackbox Exporter for baseline observability
