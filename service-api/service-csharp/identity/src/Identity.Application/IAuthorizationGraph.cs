// IAuthorizationGraph define a autorizacao fina usada pela identidade e pela borda.
namespace Identity.Application;

public interface IAuthorizationGraph
{
  void SyncTenantAccess(string tenantSlug, Guid userPublicId, IReadOnlyCollection<string> roleCodes, bool active);

  bool CanAccessTenant(string tenantSlug, Guid userPublicId);
}
