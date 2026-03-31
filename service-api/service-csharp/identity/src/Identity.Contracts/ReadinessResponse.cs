// ReadinessResponse descreve o estado detalhado de readiness do servico.
namespace Identity.Contracts;

public sealed record ReadinessResponse(
  string Service,
  string Status,
  IReadOnlyCollection<DependencyHealthResponse> Dependencies);
