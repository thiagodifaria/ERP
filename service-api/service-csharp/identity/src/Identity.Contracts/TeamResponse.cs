// TeamResponse descreve a saida publica minima para leitura de times.
namespace Identity.Contracts;

public sealed record TeamResponse(
  long Id,
  Guid PublicId,
  long TenantId,
  long? CompanyId,
  string Name,
  string Status);
