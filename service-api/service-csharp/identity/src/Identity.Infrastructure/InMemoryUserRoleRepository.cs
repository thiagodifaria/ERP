// Este adapter fornece atribuicoes basicas de papeis por usuario durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryUserRoleRepository : IUserRoleRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<(long TenantId, long UserId), List<UserRole>> _userRolesByUser = [];

  public InMemoryUserRoleRepository(
    InMemoryTenantRepository tenantRepository,
    InMemoryUserRepository userRepository,
    InMemoryRoleRepository roleRepository)
  {
    foreach (var tenant in tenantRepository.List())
    {
      var user = userRepository.ListByTenantId(tenant.Id).FirstOrDefault();
      var role = roleRepository.FindByTenantIdAndCode(tenant.Id, "owner");

      if (user is not null && role is not null)
      {
        Add(new UserRole(
          NextId(),
          tenant.Id,
          user.Id,
          role.Id,
          DateTimeOffset.UtcNow));
      }
    }
  }

  public IReadOnlyCollection<UserRole> ListByTenantIdAndUserId(long tenantId, long userId)
  {
    lock (_sync)
    {
      if (_userRolesByUser.TryGetValue((tenantId, userId), out var userRoles))
      {
        return userRoles.ToArray();
      }

      return [];
    }
  }

  public UserRole? FindByTenantIdAndUserIdAndRoleId(long tenantId, long userId, long roleId)
  {
    lock (_sync)
    {
      if (!_userRolesByUser.TryGetValue((tenantId, userId), out var userRoles))
      {
        return null;
      }

      return userRoles.FirstOrDefault(userRole => userRole.RoleId == roleId);
    }
  }

  public UserRole Add(UserRole userRole)
  {
    lock (_sync)
    {
      if (!_userRolesByUser.TryGetValue((userRole.TenantId, userRole.UserId), out var userRoles))
      {
        userRoles = [];
        _userRolesByUser[(userRole.TenantId, userRole.UserId)] = userRoles;
      }

      userRoles.Add(userRole);

      return userRole;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return _userRolesByUser
        .SelectMany(entry => entry.Value)
        .Select(userRole => userRole.Id)
        .DefaultIfEmpty(0)
        .Max() + 1;
    }
  }
}
