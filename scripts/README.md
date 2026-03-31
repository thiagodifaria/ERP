# scripts

This directory must stay small and direct.

Target scripts:

- `build.sh`
- `up.sh`
- `down.sh`
- `logs.sh`
- `db.sh`
- `test.sh`
- `perf.sh`
- `simulate.sh`
- `deploy.sh`

Rules:

- prefer flags and subcommands over script sprawl
- keep scripts readable
- do not hide critical business logic here
