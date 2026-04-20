// Este provider simula o comportamento de um diretório externo durante testes.
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class InMemoryExternalIdentityProvider : IExternalIdentityProvider
{
  private readonly Dictionary<string, InMemoryIdentityAccount> _accountsByEmail = new(StringComparer.OrdinalIgnoreCase);
  private readonly Dictionary<string, InMemoryIdentityAccount> _accountsBySubject = new(StringComparer.Ordinal);
  private readonly Dictionary<string, string> _refreshTokenIndex = new(StringComparer.Ordinal);

  public ExternalIdentityUser EnsureUser(ExternalIdentityUpsertRequest request)
  {
    var email = request.Email.Trim().ToLowerInvariant();
    var account = ResolveAccount(request.SubjectId, email);

    if (account is null)
    {
      account = new InMemoryIdentityAccount(Guid.NewGuid().ToString("N"), email);
    }

    if (!account.Email.Equals(email, StringComparison.OrdinalIgnoreCase))
    {
      _accountsByEmail.Remove(account.Email);
      account.Email = email;
    }

    account.Enabled = request.Enabled;
    account.DisplayName = request.DisplayName;
    account.GivenName = request.GivenName;
    account.FamilyName = request.FamilyName;

    if (!string.IsNullOrWhiteSpace(request.Password))
    {
      account.Password = request.Password;
    }

    _accountsByEmail[email] = account;
    _accountsBySubject[account.SubjectId] = account;

    return new ExternalIdentityUser(account.SubjectId, account.Email, account.Enabled);
  }

  public IdentityProviderTokenResult PasswordGrant(string email, string password)
  {
    if (!_accountsByEmail.TryGetValue(email.Trim().ToLowerInvariant(), out var account)
      || !account.Enabled
      || string.IsNullOrWhiteSpace(account.Password)
      || !account.Password.Equals(password, StringComparison.Ordinal))
    {
      throw new ExternalIdentityAuthenticationException("Invalid credentials.");
    }

    return IssueTokens(account);
  }

  public IdentityProviderTokenResult RefreshGrant(string refreshToken)
  {
    if (!_refreshTokenIndex.TryGetValue(refreshToken, out var subjectId)
      || !_accountsBySubject.TryGetValue(subjectId, out var account)
      || !account.Enabled)
    {
      throw new ExternalIdentityAuthenticationException("Refresh token is invalid.");
    }

    return IssueTokens(account);
  }

  private IdentityProviderTokenResult IssueTokens(InMemoryIdentityAccount account)
  {
    var accessExpiresAt = DateTimeOffset.UtcNow.AddMinutes(15);
    var refreshExpiresAt = DateTimeOffset.UtcNow.AddHours(12);
    var accessToken = $"access-{Guid.NewGuid():N}";
    var refreshToken = $"refresh-{Guid.NewGuid():N}";

    _refreshTokenIndex[refreshToken] = account.SubjectId;

    return new IdentityProviderTokenResult(
      account.SubjectId,
      accessToken,
      refreshToken,
      accessExpiresAt,
      refreshExpiresAt);
  }

  private InMemoryIdentityAccount? ResolveAccount(string? subjectId, string email)
  {
    if (!string.IsNullOrWhiteSpace(subjectId)
      && _accountsBySubject.TryGetValue(subjectId, out var bySubject))
    {
      return bySubject;
    }

    return _accountsByEmail.TryGetValue(email, out var byEmail)
      ? byEmail
      : null;
  }

  private sealed class InMemoryIdentityAccount
  {
    public InMemoryIdentityAccount(string subjectId, string email)
    {
      SubjectId = subjectId;
      Email = email;
    }

    public string SubjectId { get; }

    public string Email { get; set; }

    public string? Password { get; set; }

    public bool Enabled { get; set; }

    public string? DisplayName { get; set; }

    public string? GivenName { get; set; }

    public string? FamilyName { get; set; }
  }
}
