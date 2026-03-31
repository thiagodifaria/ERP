# service-golang

Go is reserved for high-concurrency transactional services and edge workloads.

Planned services:

- edge
- crm
- sales
- rentals
- documents
- shared

Standard layout:

- `cmd/<service>/main.go`
- `internal/api`
- `internal/application`
- `internal/domain`
- `internal/infrastructure`
- `internal/config`
- `internal/telemetry`
- `internal/bootstrap`
- `tests/unit`
- `tests/integration`
- `tests/contract`
