# service-csharp

.NET is reserved for enterprise-sensitive services with stronger security and financial concerns.

Implemented or active services:

- identity
- finance
- billing

Reserved next services:

- billing
- shared

Standard layout:

- `src/<Service>.Api`
- `src/<Service>.Application`
- `src/<Service>.Domain`
- `src/<Service>.Infrastructure`
- `src/<Service>.Contracts`
- `tests/<Service>.UnitTests`
- `tests/<Service>.IntegrationTests`
- `tests/<Service>.ContractTests`
