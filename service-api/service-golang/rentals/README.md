# rentals

The rentals service owns recurring contracts, contractual adjustments, terminations and the first recurring schedule linked to customer operations.

Current scope:

- health and readiness endpoints
- tenant-aware contract registry with explicit customer linkage
- recurring charge schedule generated from start date, end date and billing day
- charge settlement lifecycle with explicit `scheduled`, `paid` and `cancelled` states
- contractual adjustments with future charge recalculation
- terminations with future charge cancellation
- auditable contractual history
- initial attachment linkage through `documents`
- PostgreSQL-backed persistence and runtime smoke validation in container

Public routes:

- `GET /health/live`
- `GET /health/ready`
- `GET /health/details`
- `GET /api/rentals/contracts`
- `GET /api/rentals/contracts/summary`
- `POST /api/rentals/contracts`
- `GET /api/rentals/contracts/{publicId}`
- `GET /api/rentals/contracts/{publicId}/charges`
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status`
- `GET /api/rentals/contracts/{publicId}/history`
- `GET /api/rentals/contracts/{publicId}/adjustments`
- `POST /api/rentals/contracts/{publicId}/adjustments`
- `POST /api/rentals/contracts/{publicId}/terminate`
- `GET /api/rentals/contracts/{publicId}/attachments`
- `POST /api/rentals/contracts/{publicId}/attachments`

Query and payload conventions:

- `GET /api/rentals/contracts` accepts `tenantSlug`, `status` and `customerPublicId`
- `GET /api/rentals/contracts/{publicId}/charges` accepts `tenantSlug` and optional `status`
- `POST /api/rentals/contracts` creates an active contract with generated recurring charges
- `PATCH /api/rentals/contracts/{publicId}/charges/{chargePublicId}/status` settles or cancels a generated charge and records contractual history plus outbox side effects
- `POST /api/rentals/contracts/{publicId}/adjustments` updates only future scheduled charges from `effectiveAt`
- `POST /api/rentals/contracts/{publicId}/terminate` terminates the contract and cancels future charges after the effective date

Container-first validation:

- `docker run --rm -v ${PWD}:/workspace -w /workspace/service-api/service-golang/rentals golang:1.24-alpine go test ./...`
- `docker build -t erp-rentals ./service-api/service-golang/rentals`
- `bash scripts/test.sh smoke`
