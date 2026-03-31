// Este modulo registra adapters, persistencia e integracoes tecnicas do servico.
// Dependencias externas nao devem vazar para dominio ou aplicacao.
using Identity.Domain;
using Microsoft.Extensions.DependencyInjection;

namespace Identity.Infrastructure;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityInfrastructure(this IServiceCollection services)
  {
    services.AddSingleton<ITenantCatalog, InMemoryTenantCatalog>();
    services.AddSingleton<IRoleCatalog, InMemoryRoleCatalog>();

    return services;
  }
}
