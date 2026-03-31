// Team representa a segmentacao operacional de usuarios dentro do tenant.
namespace Identity.Domain;

public sealed class Team
{
  public Team(long id, long tenantId, Guid publicId, string name, string status)
  {
    Id = id;
    TenantId = tenantId;
    PublicId = publicId;
    Name = name;
    Status = status;
  }

  public long Id { get; }

  public long TenantId { get; }

  public Guid PublicId { get; }

  public string Name { get; }

  public string Status { get; }
}
