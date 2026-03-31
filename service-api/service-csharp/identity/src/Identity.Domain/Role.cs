// Role representa um papel de acesso com ownership do contexto de identidade.
namespace Identity.Domain;

public sealed class Role
{
  public Role(long id, long tenantId, Guid publicId, string code, string displayName, string status)
  {
    Id = id;
    TenantId = tenantId;
    PublicId = publicId;
    Code = code;
    DisplayName = displayName;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public Guid PublicId { get; }

  public string Code { get; }

  public string DisplayName { get; }

  public string Status { get; }
}
