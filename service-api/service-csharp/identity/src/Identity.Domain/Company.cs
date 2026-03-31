// Company representa uma empresa ou unidade vinculada a um tenant.
// Ownership de estrutura organizacional basica fica neste agregado.
namespace Identity.Domain;

public sealed class Company
{
  public Company(long id, long tenantId, Guid publicId, string displayName, string status)
  {
    Id = id;
    TenantId = tenantId;
    PublicId = publicId;
    DisplayName = displayName;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public Guid PublicId { get; }

  public string DisplayName { get; }

  public string Status { get; }
}
