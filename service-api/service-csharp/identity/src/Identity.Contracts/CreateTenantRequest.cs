// CreateTenantRequest descreve a entrada publica minima para criacao de tenant.
namespace Identity.Contracts;

public sealed record CreateTenantRequest(
  string Slug,
  string DisplayName);
