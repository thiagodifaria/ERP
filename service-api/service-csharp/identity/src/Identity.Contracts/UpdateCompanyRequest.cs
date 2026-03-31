// UpdateCompanyRequest descreve a entrada publica minima para atualizar empresa existente.
namespace Identity.Contracts;

public sealed record UpdateCompanyRequest(
  string DisplayName,
  string? LegalName,
  string? TaxId);
