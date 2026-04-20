// RefreshSessionRequest descreve a renovacao publica de sessao.
namespace Identity.Contracts;

public sealed record RefreshSessionRequest(
  string RefreshToken);
