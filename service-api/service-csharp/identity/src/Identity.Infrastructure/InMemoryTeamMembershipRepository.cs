// Este adapter fornece memberships basicos de times durante o bootstrap do servico.
// A persistencia real pode substituir esta implementacao sem quebrar os contratos publicos.
using Identity.Domain;

namespace Identity.Infrastructure;

public sealed class InMemoryTeamMembershipRepository : ITeamMembershipRepository
{
  private readonly object _sync = new();
  private readonly Dictionary<(long TenantId, long TeamId), List<TeamMembership>> _membershipsByTeam = [];

  public InMemoryTeamMembershipRepository(
    InMemoryTenantRepository tenantRepository,
    InMemoryTeamRepository teamRepository,
    InMemoryUserRepository userRepository)
  {
    foreach (var tenant in tenantRepository.List())
    {
      var team = teamRepository.ListByTenantId(tenant.Id).FirstOrDefault();
      var user = userRepository.ListByTenantId(tenant.Id).FirstOrDefault();

      if (team is not null && user is not null)
      {
        Add(new TeamMembership(
          NextId(),
          tenant.Id,
          team.Id,
          user.Id,
          DateTimeOffset.UtcNow));
      }
    }
  }

  public IReadOnlyCollection<TeamMembership> ListByTenantIdAndTeamId(long tenantId, long teamId)
  {
    lock (_sync)
    {
      if (_membershipsByTeam.TryGetValue((tenantId, teamId), out var memberships))
      {
        return memberships.ToArray();
      }

      return [];
    }
  }

  public TeamMembership? FindByTenantIdAndTeamIdAndUserId(long tenantId, long teamId, long userId)
  {
    lock (_sync)
    {
      if (!_membershipsByTeam.TryGetValue((tenantId, teamId), out var memberships))
      {
        return null;
      }

      return memberships.FirstOrDefault(membership => membership.UserId == userId);
    }
  }

  public TeamMembership Add(TeamMembership membership)
  {
    lock (_sync)
    {
      if (!_membershipsByTeam.TryGetValue((membership.TenantId, membership.TeamId), out var memberships))
      {
        memberships = [];
        _membershipsByTeam[(membership.TenantId, membership.TeamId)] = memberships;
      }

      memberships.Add(membership);

      return membership;
    }
  }

  public bool RemoveByTenantIdAndTeamIdAndUserId(long tenantId, long teamId, long userId)
  {
    lock (_sync)
    {
      if (!_membershipsByTeam.TryGetValue((tenantId, teamId), out var memberships))
      {
        return false;
      }

      var removed = memberships.RemoveAll(membership => membership.UserId == userId);

      if (memberships.Count == 0)
      {
        _membershipsByTeam.Remove((tenantId, teamId));
      }

      return removed > 0;
    }
  }

  public long NextId()
  {
    lock (_sync)
    {
      return _membershipsByTeam
        .SelectMany(entry => entry.Value)
        .Select(membership => membership.Id)
        .DefaultIfEmpty(0)
        .Max() + 1;
    }
  }
}
