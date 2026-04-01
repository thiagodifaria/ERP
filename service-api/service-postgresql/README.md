# service-postgresql

PostgreSQL ownership is organized by business domain.

Each domain is expected to own:

- `migrations`
- `seeds`
- `views`
- `functions`
- `indexes`
- a short `README.md` describing data ownership

Domains already mapped in this repository:

- common
- identity
- crm
- sales
- rentals
- finance
- billing
- webhook-hub
- workflow-control
- workflow-runtime
- engagement
- documents
- analytics

Domains already carrying concrete relational work:

- `common`
- `identity`
- `crm`
- `webhook-hub`
