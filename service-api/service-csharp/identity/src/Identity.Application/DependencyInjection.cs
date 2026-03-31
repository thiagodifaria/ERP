// Este modulo registra dependencias da camada de aplicacao.
// Casos de uso e orquestracao entram aqui conforme o servico crescer.
namespace Identity.Application;

public static class DependencyInjection
{
  public static IServiceCollection AddIdentityApplication(this IServiceCollection services)
  {
    return services;
  }
}
