// UserRole representa a atribuicao direta de papeis a usuarios.
namespace Identity.Domain;

public sealed class UserRole
{
  public UserRole(long id, long tenantId, long userId, long roleId, DateTimeOffset createdAt)
  {
    Id = id;
    TenantId = tenantId;
    UserId = userId;
    RoleId = roleId;
    CreatedAt = createdAt;
  }

  public long Id { get; }

  public long TenantId { get; }

  public long UserId { get; }

  public long RoleId { get; }

  public DateTimeOffset CreatedAt { get; }
}
