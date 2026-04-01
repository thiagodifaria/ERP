# scripts

This directory must stay small and direct.

Current scripts:

- `build.sh`
- `up.sh`
- `down.sh`
- `logs.sh`
- `db.sh`
- `test.sh`

Useful commands:

- `./scripts/up.sh`
- `./scripts/down.sh`
- `./scripts/logs.sh identity`
- `./scripts/db.sh migrate all`
- `./scripts/db.sh migrate crm`
- `./scripts/db.sh migrate workflow-control`
- `./scripts/db.sh seed identity`
- `./scripts/db.sh seed crm`
- `./scripts/db.sh seed workflow-control`
- `./scripts/db.sh summary identity smoke-identity-bootstrap`
- `./scripts/db.sh summary crm smoke-identity-bootstrap`  `# total, status e ownership`
- `./scripts/db.sh summary workflow-control bootstrap-ops`  `# total, versoes, runs e distribuicao por status de execucao`
- `./scripts/test.sh unit`
- `./scripts/test.sh integration`
- `./scripts/test.sh contract`  `# contratos HTTP publicos de workflow-control, crm e identity`
- `./scripts/test.sh smoke`  `# reset relacional + bootstrap + smoke HTTP de workflow-control, crm e identity, incluindo versionamento, ledger de runs e filtros`
- `./scripts/test.sh all`

Rules:

- prefer flags and subcommands over script sprawl
- keep scripts readable
- do not hide critical business logic here

Current unit scope:

- Go: `edge`, `crm`
- TypeScript: `workflow-control`
- .NET: `identity`
- Rust: `webhook-hub`

Current contract scope:

- TypeScript: `workflow-control`
- Go: `crm`
- .NET: `identity`
