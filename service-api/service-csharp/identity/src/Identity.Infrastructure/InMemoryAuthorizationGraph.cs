// Este grafo em memoria valida o escopo basico de tenant sem depender de runtime externo.
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class InMemoryAuthorizationGraph : IAuthorizationGraph
{
  private readonly Dictionary<string, HashSet<Guid>> _tenantMembers = new(StringComparer.OrdinalIgnoreCase);

  public void SyncTenantAccess(string tenantSlug, Guid userPublicId, IReadOnlyCollection<string> roleCodes, bool active)
  {
    if (!_tenantMembers.TryGetValue(tenantSlug, out var members))
    {
      members = [];
      _tenantMembers[tenantSlug] = members;
    }

    if (active && roleCodes.Count > 0)
    {
      members.Add(userPublicId);
      return;
    }

    members.Remove(userPublicId);
  }

  public bool CanAccessTenant(string tenantSlug, Guid userPublicId)
  {
    return _tenantMembers.TryGetValue(tenantSlug, out var members) && members.Contains(userPublicId);
  }
}
