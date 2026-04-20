// Este modulo registra adapters, persistencia e integracoes tecnicas do servico.
// Dependencias externas nao devem vazar para dominio ou aplicacao.
using Identity.Domain;
using Microsoft.Extensions.DependencyInjection;
using Identity.Application;

namespace Identity.Infrastructure;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityInfrastructure(this IServiceCollection services)
  {
    var options = IdentityInfrastructureOptions.Load();

    if (options.RepositoryDriver.Equals("postgres", StringComparison.OrdinalIgnoreCase))
    {
      services.AddSingleton(new PostgresIdentityRepositoryBundle(options.PostgresConnectionString));
      services.AddSingleton(new PostgresIdentitySecurityStore(options.PostgresConnectionString));
      services.AddSingleton<ITenantCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ITenantRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ICompanyCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ICompanyRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IUserCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IUserRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ITeamCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ITeamRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ITeamMembershipCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<ITeamMembershipRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IRoleCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IRoleRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IUserRoleCatalog>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IUserRoleRepository>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentityRepositoryBundle>());
      services.AddSingleton<IIdentitySecurityStore>(serviceProvider => serviceProvider.GetRequiredService<PostgresIdentitySecurityStore>());

      RegisterExternalDrivers(services, options);

      return services;
    }

    services.AddSingleton<InMemoryTenantRepository>();
    services.AddSingleton<InMemoryCompanyRepository>();
    services.AddSingleton<InMemoryUserRepository>();
    services.AddSingleton<InMemoryTeamRepository>();
    services.AddSingleton<InMemoryRoleRepository>();
    services.AddSingleton<InMemoryTeamMembershipRepository>();
    services.AddSingleton<InMemoryUserRoleRepository>();
    services.AddSingleton<ITenantCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTenantRepository>());
    services.AddSingleton<ITenantRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTenantRepository>());
    services.AddSingleton<ICompanyCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryCompanyRepository>());
    services.AddSingleton<ICompanyRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryCompanyRepository>());
    services.AddSingleton<IUserCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryUserRepository>());
    services.AddSingleton<IUserRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryUserRepository>());
    services.AddSingleton<ITeamCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTeamRepository>());
    services.AddSingleton<ITeamRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTeamRepository>());
    services.AddSingleton<ITeamMembershipCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTeamMembershipRepository>());
    services.AddSingleton<ITeamMembershipRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryTeamMembershipRepository>());
    services.AddSingleton<IRoleCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryRoleRepository>());
    services.AddSingleton<IRoleRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryRoleRepository>());
    services.AddSingleton<IUserRoleCatalog>(serviceProvider => serviceProvider.GetRequiredService<InMemoryUserRoleRepository>());
    services.AddSingleton<IUserRoleRepository>(serviceProvider => serviceProvider.GetRequiredService<InMemoryUserRoleRepository>());
    services.AddSingleton<IIdentitySecurityStore, InMemoryIdentitySecurityStore>();
    RegisterExternalDrivers(services, options);

    return services;
  }

  private static void RegisterExternalDrivers(IServiceCollection services, IdentityInfrastructureOptions options)
  {
    services.AddSingleton(options);
    services.AddSingleton<IIdentityAccessSettings>(options);
    services.AddSingleton<ITotpService, TotpService>();

    if (options.IdentityProviderDriver.Equals("keycloak", StringComparison.OrdinalIgnoreCase))
    {
      services.AddSingleton<IExternalIdentityProvider, KeycloakIdentityProvider>();
    }
    else
    {
      services.AddSingleton<IExternalIdentityProvider, InMemoryExternalIdentityProvider>();
    }

    if (options.AuthorizationDriver.Equals("openfga", StringComparison.OrdinalIgnoreCase))
    {
      services.AddSingleton<IAuthorizationGraph, OpenFgaAuthorizationGraph>();
    }
    else
    {
      services.AddSingleton<IAuthorizationGraph, InMemoryAuthorizationGraph>();
    }
  }
}
