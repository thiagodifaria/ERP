// UserSecurityProfile concentra o estado tecnico de acesso do usuario.
// MFA e sincronizacao com provider externo vivem fora do agregado funcional User.
namespace Identity.Domain;

public sealed class UserSecurityProfile
{
  public UserSecurityProfile(
    long userId,
    string? identityProviderSubject,
    bool mfaEnabled,
    string? mfaSecret,
    DateTimeOffset? lastLoginAt)
  {
    UserId = userId;
    IdentityProviderSubject = identityProviderSubject;
    MfaEnabled = mfaEnabled;
    MfaSecret = mfaSecret;
    LastLoginAt = lastLoginAt;
  }

  public long UserId { get; }

  public string? IdentityProviderSubject { get; }

  public bool MfaEnabled { get; }

  public string? MfaSecret { get; }

  public DateTimeOffset? LastLoginAt { get; }

  public UserSecurityProfile AttachIdentityProviderSubject(string subject)
  {
    return new UserSecurityProfile(UserId, subject, MfaEnabled, MfaSecret, LastLoginAt);
  }

  public UserSecurityProfile StartMfaEnrollment(string mfaSecret)
  {
    return new UserSecurityProfile(UserId, IdentityProviderSubject, false, mfaSecret, LastLoginAt);
  }

  public UserSecurityProfile EnableMfa()
  {
    return new UserSecurityProfile(UserId, IdentityProviderSubject, true, MfaSecret, LastLoginAt);
  }

  public UserSecurityProfile DisableMfa()
  {
    return new UserSecurityProfile(UserId, IdentityProviderSubject, false, null, LastLoginAt);
  }

  public UserSecurityProfile RecordLogin(DateTimeOffset lastLoginAt)
  {
    return new UserSecurityProfile(UserId, IdentityProviderSubject, MfaEnabled, MfaSecret, lastLoginAt);
  }
}
