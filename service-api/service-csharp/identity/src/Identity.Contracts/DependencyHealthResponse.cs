// DependencyHealthResponse descreve o estado observado de uma dependencia do servico.
namespace Identity.Contracts;

public sealed record DependencyHealthResponse(string Name, string Status);
