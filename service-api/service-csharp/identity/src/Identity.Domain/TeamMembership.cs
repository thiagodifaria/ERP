// TeamMembership liga usuarios a times de forma explicita e auditavel.
namespace Identity.Domain;

public sealed class TeamMembership
{
  public TeamMembership(long id, long tenantId, long teamId, long userId, DateTimeOffset createdAt)
  {
    Id = id;
    TenantId = tenantId;
    TeamId = teamId;
    UserId = userId;
    CreatedAt = createdAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long TeamId { get; }

  public long UserId { get; }

  public DateTimeOffset CreatedAt { get; }
}
