// Este adapter fornece usuarios basicos por tenant durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryUserRepository : IUserRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<long, List<User>> _usersByTenantId = [];

  public InMemoryUserRepository()
  {
    SeedDefaults(new Tenant(
      1,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000001"),
      "bootstrap-ops",
      "Bootstrap Ops",
      "active"));

    SeedDefaults(new Tenant(
      2,
      Guid.Parse("0195e7a0-7a9c-7c1f-8a44-4a6e50000002"),
      "northwind-group",
      "Northwind Group",
      "active"));
  }

  public IReadOnlyCollection<User> ListByTenantId(long tenantId)
  {
    lock (_sync)
    {
      if (_usersByTenantId.TryGetValue(tenantId, out var users))
      {
        return users.ToArray();
      }

      return [];
    }
  }

  public User? FindByTenantIdAndId(long tenantId, long userId)
  {
    lock (_sync)
    {
      if (!_usersByTenantId.TryGetValue(tenantId, out var users))
      {
        return null;
      }

      return users.FirstOrDefault(user => user.Id == userId);
    }
  }

  public User? FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    lock (_sync)
    {
      if (!_usersByTenantId.TryGetValue(tenantId, out var users))
      {
        return null;
      }

      return users.FirstOrDefault(user => user.PublicId == publicId);
    }
  }

  public User? FindByTenantIdAndEmail(long tenantId, string email)
  {
    lock (_sync)
    {
      if (!_usersByTenantId.TryGetValue(tenantId, out var users))
      {
        return null;
      }

      return users.FirstOrDefault(
        user => user.Email.Equals(email, StringComparison.OrdinalIgnoreCase));
    }
  }

  public User Add(User user)
  {
    lock (_sync)
    {
      if (_usersByTenantId.TryGetValue(user.TenantId, out var existingUsers)
        && existingUsers.Any(existing =>
          existing.Email.Equals(user.Email, StringComparison.OrdinalIgnoreCase)))
      {
        throw new InvalidOperationException("User email already exists for tenant.");
      }

      if (!_usersByTenantId.TryGetValue(user.TenantId, out var users))
      {
        users = [];
        _usersByTenantId[user.TenantId] = users;
      }

      users.Add(user);

      return user;
    }
  }

  public User Update(User user)
  {
    lock (_sync)
    {
      if (!_usersByTenantId.TryGetValue(user.TenantId, out var users))
      {
        users = [];
        _usersByTenantId[user.TenantId] = users;
      }

      for (var index = 0; index < users.Count; index++)
      {
        if (users[index].PublicId != user.PublicId)
        {
          continue;
        }

        users[index] = user;
        return user;
      }

      users.Add(user);
      return user;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return NextIdUnsafe();
    }
  }

  public IReadOnlyCollection<User> SeedDefaults(Tenant tenant)
  {
    lock (_sync)
    {
      if (_usersByTenantId.TryGetValue(tenant.Id, out var existingUsers) && existingUsers.Count > 0)
      {
        return existingUsers.ToArray();
      }

      var seededUsers =
        new[]
        {
          new User(
            NextIdUnsafe(),
            tenant.Id,
            null,
            PublicIds.NewUuidV7(),
            $"owner@{tenant.Slug}.local",
            $"{tenant.DisplayName} Owner",
            null,
            null,
            "active")
        }.ToList();

      _usersByTenantId[tenant.Id] = seededUsers;

      return seededUsers.ToArray();
    }
  }

  private long NextIdUnsafe()
  {
    return _usersByTenantId
      .SelectMany(entry => entry.Value)
      .Select(user => user.Id)
      .DefaultIfEmpty(0)
      .Max() + 1;
  }
}
