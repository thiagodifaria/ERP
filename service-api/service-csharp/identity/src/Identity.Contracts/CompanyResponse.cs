// CompanyResponse descreve a saida publica minima para leitura de empresas.
namespace Identity.Contracts;

public sealed record CompanyResponse(
  long Id,
  Guid PublicId,
  long TenantId,
  string DisplayName,
  string? LegalName,
  string? TaxId,
  string Status);
