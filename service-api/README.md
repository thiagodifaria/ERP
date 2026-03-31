# service-api

This directory is the main backend axis of the repository.

The first split is by language, not by environment or by deployment unit.

Principles:

- each language keeps its own idiomatic project structure
- each service owns its own bootstrap file such as `server.*`
- business rules belong to domain and application layers
- infrastructure adapters stay isolated from core business logic
- PostgreSQL ownership is organized by business domain under `service-postgresql`
