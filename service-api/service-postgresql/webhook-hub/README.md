# webhook-hub

This domain owns inbound webhook intake logs, lifecycle transitions and future normalization buffers.

Current structure:

- migrations
- seeds
- inbound webhook event storage
- webhook event transition ledger
- unique guard for provider and external id
- provider, status and chronology indexes
