using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class UpdateIdentityUserAccess
{
  private static readonly HashSet<string> AllowedStatuses = ["active", "inactive", "suspended"];

  private readonly ITenantCatalog _tenantCatalog;
  private readonly IUserRepository _userRepository;
  private readonly IIdentitySecurityStore _securityStore;
  private readonly IExternalIdentityProvider _identityProvider;
  private readonly SecurityAuditWriter _auditWriter;
  private readonly TenantAccessCoordinator _tenantAccessCoordinator;

  public UpdateIdentityUserAccess(
    ITenantCatalog tenantCatalog,
    IUserRepository userRepository,
    IIdentitySecurityStore securityStore,
    IExternalIdentityProvider identityProvider,
    SecurityAuditWriter auditWriter,
    TenantAccessCoordinator tenantAccessCoordinator)
  {
    _tenantCatalog = tenantCatalog;
    _userRepository = userRepository;
    _securityStore = securityStore;
    _identityProvider = identityProvider;
    _auditWriter = auditWriter;
    _tenantAccessCoordinator = tenantAccessCoordinator;
  }

  public OperationResult<UserResponse> Execute(string tenantSlug, Guid userPublicId, UpdateUserAccessRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);
    if (tenant is null)
    {
      return OperationResult<UserResponse>.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var user = _userRepository.FindByTenantIdAndPublicId(tenant.Id, userPublicId);
    if (user is null)
    {
      return OperationResult<UserResponse>.NotFound(
        new ErrorResponse("user_not_found", "User was not found."));
    }

    var status = request.Status.Trim().ToLowerInvariant();
    if (!AllowedStatuses.Contains(status))
    {
      return OperationResult<UserResponse>.BadRequest(
        new ErrorResponse("invalid_status", "Status must be active, inactive or suspended."));
    }

    var updatedUser = _userRepository.Update(user.ReviseStatus(status));
    var profile = _securityStore.GetOrCreateProfile(updatedUser.Id);

    if (!string.IsNullOrWhiteSpace(profile.IdentityProviderSubject))
    {
      _identityProvider.EnsureUser(new ExternalIdentityUpsertRequest(
        updatedUser.Email,
        updatedUser.GivenName,
        updatedUser.FamilyName,
        updatedUser.DisplayName,
        status == "active",
        profile.IdentityProviderSubject,
        null));
    }

    if (status != "active")
    {
      _securityStore.RevokeSessionsByUserId(tenant.Id, updatedUser.Id, DateTimeOffset.UtcNow);
    }

    _tenantAccessCoordinator.SyncAndListRoleCodes(tenant, updatedUser);
    _auditWriter.Record(
      tenant.Id,
      null,
      updatedUser.PublicId,
      status == "active" ? "access_restored" : "access_blocked",
      status == "active" ? "info" : "warning",
      $"Access status changed to {status} for {updatedUser.Email}.");

    return OperationResult<UserResponse>.Success(updatedUser.ToResponse());
  }
}
