// Este modulo registra dependencias da camada de aplicacao.
// Casos de uso e orquestracao entram aqui conforme o servico crescer.
using Microsoft.Extensions.DependencyInjection;

namespace Identity.Application;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityApplication(this IServiceCollection services)
  {
    services.AddScoped<SecurityAuditWriter>();
    services.AddScoped<TenantAccessCoordinator>();
    services.AddScoped<AddBootstrapTeamMember>();
    services.AddScoped<AssignBootstrapUserRole>();
    services.AddScoped<AcceptIdentityInvite>();
    services.AddScoped<CompleteIdentityPasswordRecovery>();
    services.AddScoped<CreateIdentityInvite>();
    services.AddScoped<CreateBootstrapCompany>();
    services.AddScoped<CreateBootstrapTenant>();
    services.AddScoped<CreateBootstrapTeam>();
    services.AddScoped<CreateBootstrapUser>();
    services.AddScoped<DisableIdentityUserMfa>();
    services.AddScoped<ListIdentityInvites>();
    services.AddScoped<ListIdentityUserSessions>();
    services.AddScoped<ListIdentitySecurityAuditEvents>();
    services.AddScoped<LoginIdentitySession>();
    services.AddScoped<RefreshIdentitySession>();
    services.AddScoped<RemoveBootstrapTeamMember>();
    services.AddScoped<ResolveTenantAccess>();
    services.AddScoped<RevokeIdentitySession>();
    services.AddScoped<RevokeIdentityUserSessions>();
    services.AddScoped<RevokeBootstrapUserRole>();
    services.AddScoped<StartIdentityUserMfaEnrollment>();
    services.AddScoped<StartIdentityPasswordRecovery>();
    services.AddScoped<UpdateIdentityUserAccess>();
    services.AddScoped<VerifyIdentityUserMfa>();
    services.AddScoped<GetBootstrapCompanyByPublicId>();
    services.AddScoped<GetBootstrapUserByPublicId>();
    services.AddScoped<UpdateBootstrapTeam>();
    services.AddScoped<UpdateBootstrapCompany>();
    services.AddScoped<UpdateBootstrapUser>();
    services.AddScoped<GetBootstrapTenantAccessSnapshot>();
    services.AddScoped<ListBootstrapUserRoles>();
    services.AddScoped<ListBootstrapCompanies>();
    services.AddScoped<ListBootstrapTeamMembers>();
    services.AddScoped<ListBootstrapTeams>();
    services.AddScoped<GetBootstrapTenantBySlug>();
    services.AddScoped<ListBootstrapTenants>();
    services.AddScoped<ListBootstrapRoles>();
    services.AddScoped<ListBootstrapUsers>();

    return services;
  }
}
