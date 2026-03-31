// Este adapter fornece times basicos por tenant durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryTeamRepository : ITeamRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<long, List<Team>> _teamsByTenantId = [];

  public InMemoryTeamRepository()
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

  public IReadOnlyCollection<Team> ListByTenantId(long tenantId)
  {
    lock (_sync)
    {
      if (_teamsByTenantId.TryGetValue(tenantId, out var teams))
      {
        return teams.ToArray();
      }

      return [];
    }
  }

  public Team? FindByTenantIdAndName(long tenantId, string name)
  {
    lock (_sync)
    {
      if (!_teamsByTenantId.TryGetValue(tenantId, out var teams))
      {
        return null;
      }

      return teams.FirstOrDefault(
        team => team.Name.Equals(name, StringComparison.OrdinalIgnoreCase));
    }
  }

  public Team? FindByTenantIdAndPublicId(long tenantId, Guid publicId)
  {
    lock (_sync)
    {
      if (!_teamsByTenantId.TryGetValue(tenantId, out var teams))
      {
        return null;
      }

      return teams.FirstOrDefault(team => team.PublicId == publicId);
    }
  }

  public Team Add(Team team)
  {
    lock (_sync)
    {
      if (!_teamsByTenantId.TryGetValue(team.TenantId, out var teams))
      {
        teams = [];
        _teamsByTenantId[team.TenantId] = teams;
      }

      teams.Add(team);

      return team;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return NextIdUnsafe();
    }
  }

  public IReadOnlyCollection<Team> SeedDefaults(Tenant tenant)
  {
    lock (_sync)
    {
      if (_teamsByTenantId.TryGetValue(tenant.Id, out var existingTeams) && existingTeams.Count > 0)
      {
        return existingTeams.ToArray();
      }

      var seededTeams =
        new[]
        {
          new Team(
            NextIdUnsafe(),
            tenant.Id,
            null,
            PublicIds.NewUuidV7(),
            "Core",
            "active")
        }.ToList();

      _teamsByTenantId[tenant.Id] = seededTeams;

      return seededTeams.ToArray();
    }
  }

  private long NextIdUnsafe()
  {
    return _teamsByTenantId
      .SelectMany(entry => entry.Value)
      .Select(team => team.Id)
      .DefaultIfEmpty(0)
      .Max() + 1;
  }
}
