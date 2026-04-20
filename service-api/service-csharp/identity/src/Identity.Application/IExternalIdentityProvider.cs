// IExternalIdentityProvider define a integracao de autenticacao externa.
// O contrato abstrai Keycloak sem vazar transporte para a aplicacao.
namespace Identity.Application;

public interface IExternalIdentityProvider
{
  ExternalIdentityUser EnsureUser(ExternalIdentityUpsertRequest request);

  IdentityProviderTokenResult PasswordGrant(string email, string password);

  IdentityProviderTokenResult RefreshGrant(string refreshToken);
}

public sealed record ExternalIdentityUpsertRequest(
  string Email,
  string? GivenName,
  string? FamilyName,
  string? DisplayName,
  bool Enabled,
  string? SubjectId,
  string? Password);

public sealed record ExternalIdentityUser(
  string SubjectId,
  string Email,
  bool Enabled);

public sealed record IdentityProviderTokenResult(
  string SubjectId,
  string AccessToken,
  string RefreshToken,
  DateTimeOffset AccessExpiresAt,
  DateTimeOffset RefreshExpiresAt);

public sealed class ExternalIdentityAuthenticationException : Exception
{
  public ExternalIdentityAuthenticationException(string message)
    : base(message)
  {
  }
}

public sealed class ExternalIdentityProviderException : Exception
{
  public ExternalIdentityProviderException(string message)
    : base(message)
  {
  }
}
