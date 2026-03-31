// Contratos publicos da API ficam separados para manter fronteiras claras.
namespace Identity.Contracts;

public sealed record HealthResponse(string Service, string Status);
