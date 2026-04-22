# documents

This context owns attachment metadata, archive governance, retention defaults and
owner associations for tenant-scoped aggregates.

Current relational scope:

- `documents.attachments`

Operational notes:

- attachment rows preserve owner reference, file metadata and storage driver selection.
- governance fields cover file size, checksum, visibility, retention and archive timestamps.
- archive is metadata-only for now, preserving operational history without deleting the record.

Validation:

- `bash scripts/db.sh migrate documents`
- `bash scripts/db.sh summary documents`
- `bash scripts/test.sh smoke`
