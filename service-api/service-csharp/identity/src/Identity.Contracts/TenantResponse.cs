// TenantResponse descreve a saida publica minima para leitura de tenants.
namespace Identity.Contracts;

public sealed record TenantResponse(
  long Id,
  Guid PublicId,
  string Slug,
  string DisplayName,
  string Status);
