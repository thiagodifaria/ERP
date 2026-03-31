// Este modulo registra dependencias da camada de aplicacao.
// Casos de uso e orquestracao entram aqui conforme o servico crescer.
using Microsoft.Extensions.DependencyInjection;

namespace Identity.Application;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityApplication(this IServiceCollection services)
  {
    services.AddScoped<CreateBootstrapTenant>();
    services.AddScoped<ListBootstrapTenants>();
    services.AddScoped<ListBootstrapRoles>();

    return services;
  }
}
