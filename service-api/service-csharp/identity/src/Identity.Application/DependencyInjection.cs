// Este modulo registra dependencias da camada de aplicacao.
// Casos de uso e orquestracao entram aqui conforme o servico crescer.
using Microsoft.Extensions.DependencyInjection;

namespace Identity.Application;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityApplication(this IServiceCollection services)
  {
    services.AddScoped<CreateBootstrapCompany>();
    services.AddScoped<CreateBootstrapTenant>();
    services.AddScoped<CreateBootstrapUser>();
    services.AddScoped<ListBootstrapCompanies>();
    services.AddScoped<GetBootstrapTenantBySlug>();
    services.AddScoped<ListBootstrapTenants>();
    services.AddScoped<ListBootstrapRoles>();
    services.AddScoped<ListBootstrapUsers>();

    return services;
  }
}
