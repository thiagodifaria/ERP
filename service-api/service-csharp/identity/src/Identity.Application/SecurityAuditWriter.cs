// SecurityAuditWriter concentra a gravacao de trilha de seguranca.
using Identity.Domain;

namespace Identity.Application;

public sealed class SecurityAuditWriter
{
  private readonly IIdentitySecurityStore _securityStore;

  public SecurityAuditWriter(IIdentitySecurityStore securityStore)
  {
    _securityStore = securityStore;
  }

  public void Record(
    long tenantId,
    Guid? actorUserPublicId,
    Guid? subjectUserPublicId,
    string eventCode,
    string severity,
    string summary)
  {
    _securityStore.AddSecurityAuditEvent(new SecurityAuditEvent(
      _securityStore.NextSecurityAuditEventId(),
      tenantId,
      PublicIds.NewUuidV7(),
      actorUserPublicId,
      subjectUserPublicId,
      eventCode,
      severity,
      summary,
      DateTimeOffset.UtcNow));
  }
}
