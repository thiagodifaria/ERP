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
- `./scripts/db.sh seed identity`
- `./scripts/db.sh seed crm`
- `./scripts/db.sh summary identity smoke-identity-bootstrap`
- `./scripts/db.sh summary crm smoke-identity-bootstrap`  `# total, status e ownership`
- `./scripts/test.sh unit`
- `./scripts/test.sh integration`
- `./scripts/test.sh contract`
- `./scripts/test.sh smoke`
- `./scripts/test.sh all`

Rules:

- prefer flags and subcommands over script sprawl
- keep scripts readable
- do not hide critical business logic here
