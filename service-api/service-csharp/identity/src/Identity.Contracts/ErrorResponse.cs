// ErrorResponse padroniza erros publicos simples no bootstrap do servico.
namespace Identity.Contracts;

public sealed record ErrorResponse(
  string Code,
  string Message);
