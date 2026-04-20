// SecurityAuditEventResponse descreve os eventos publicos de auditoria de seguranca.
namespace Identity.Contracts;

public sealed record SecurityAuditEventResponse(
  Guid PublicId,
  Guid? ActorUserPublicId,
  Guid? SubjectUserPublicId,
  string EventCode,
  string Severity,
  string Summary,
  DateTimeOffset CreatedAt);
