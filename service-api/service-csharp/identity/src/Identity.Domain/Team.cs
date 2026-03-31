// Team representa a segmentacao operacional de usuarios dentro do tenant.
namespace Identity.Domain;

public sealed class Team
{
  public Team(long id, long tenantId, long? companyId, Guid publicId, string name, string status)
  {
    Id = id;
    TenantId = tenantId;
    CompanyId = companyId;
    PublicId = publicId;
    Name = name;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long? CompanyId { get; }

  public Guid PublicId { get; }

  public string Name { get; }

  public string Status { get; }
}
