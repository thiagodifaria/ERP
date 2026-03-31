// CreateCompanyRequest descreve a entrada publica minima para criacao de empresa.
namespace Identity.Contracts;

public sealed record CreateCompanyRequest(
  string DisplayName,
  string? LegalName,
  string? TaxId);
