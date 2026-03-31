// Este modulo registra adapters, persistencia e integracoes tecnicas do servico.
// Dependencias externas nao devem vazar para dominio ou aplicacao.
namespace Identity.Infrastructure;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityInfrastructure(this IServiceCollection services)
  {
    return services;
  }
}
