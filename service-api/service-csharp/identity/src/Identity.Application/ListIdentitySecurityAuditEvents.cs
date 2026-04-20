using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class ListIdentitySecurityAuditEvents
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly IIdentitySecurityStore _securityStore;

  public ListIdentitySecurityAuditEvents(
    ITenantCatalog tenantCatalog,
    IIdentitySecurityStore securityStore)
  {
    _tenantCatalog = tenantCatalog;
    _securityStore = securityStore;
  }

  public OperationResult<IReadOnlyCollection<SecurityAuditEventResponse>> Execute(string tenantSlug)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<IReadOnlyCollection<SecurityAuditEventResponse>>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var events = _securityStore.ListSecurityAuditByTenantId(tenant.Id, 100)
      .Select(auditEvent => auditEvent.ToResponse())
      .ToArray();

    return OperationResult<IReadOnlyCollection<SecurityAuditEventResponse>>.Success(events);
  }
}
