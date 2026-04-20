// IIdentityAccessSettings expõe as poucas configuracoes de acesso consumidas pela aplicacao.
namespace Identity.Application;

public interface IIdentityAccessSettings
{
  string BootstrapPassword { get; }

  string MfaIssuer { get; }
}
