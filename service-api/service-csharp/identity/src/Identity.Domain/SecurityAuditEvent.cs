// SecurityAuditEvent descreve eventos relevantes de acesso e seguranca.
// O registro permanece simples para permitir leitura publica sem expor payload sensivel.
namespace Identity.Domain;

public sealed class SecurityAuditEvent
{
  public SecurityAuditEvent(
    long id,
    long tenantId,
    Guid publicId,
    Guid? actorUserPublicId,
    Guid? subjectUserPublicId,
    string eventCode,
    string severity,
    string summary,
    DateTimeOffset createdAt)
  {
    Id = id;
    TenantId = tenantId;
    PublicId = publicId;
    ActorUserPublicId = actorUserPublicId;
    SubjectUserPublicId = subjectUserPublicId;
    EventCode = eventCode;
    Severity = severity;
    Summary = summary;
    CreatedAt = createdAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public Guid PublicId { get; }

  public Guid? ActorUserPublicId { get; }

  public Guid? SubjectUserPublicId { get; }

  public string EventCode { get; }

  public string Severity { get; }

  public string Summary { get; }

  public DateTimeOffset CreatedAt { get; }
}
