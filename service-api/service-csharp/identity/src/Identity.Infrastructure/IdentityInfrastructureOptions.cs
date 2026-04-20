// Este arquivo centraliza a resolucao de configuracao da infraestrutura.
// O objetivo e alternar drivers sem espalhar variaveis de ambiente pelo servico.
using Npgsql;
using Identity.Application;

namespace Identity.Infrastructure;

public sealed class IdentityInfrastructureOptions : IIdentityAccessSettings
{
  public string RepositoryDriver { get; init; } = "memory";

  public string IdentityProviderDriver { get; init; } = "memory";

  public string AuthorizationDriver { get; init; } = "memory";

  public string PostgresConnectionString { get; init; } = string.Empty;

  public string KeycloakBaseUrl { get; init; } = "http://localhost:8080";

  public string KeycloakRealm { get; init; } = "erp-local";

  public string KeycloakClientId { get; init; } = "erp-identity";

  public string KeycloakAdminUsername { get; init; } = "admin";

  public string KeycloakAdminPassword { get; init; } = "admin";

  public string BootstrapPassword { get; init; } = "Change.Me123!";

  public string OpenFgaBaseUrl { get; init; } = "http://localhost:8080";

  public string OpenFgaStoreName { get; init; } = "erp-local-store";

  public string MfaIssuer { get; init; } = "ERP";

  public static IdentityInfrastructureOptions Load()
  {
    return new IdentityInfrastructureOptions
    {
      RepositoryDriver = GetEnvironment("IDENTITY_REPOSITORY_DRIVER", "memory"),
      IdentityProviderDriver = GetEnvironment("IDENTITY_PROVIDER_DRIVER", "memory"),
      AuthorizationDriver = GetEnvironment("IDENTITY_AUTHORIZATION_DRIVER", "memory"),
      PostgresConnectionString = ResolvePostgresConnectionString(),
      KeycloakBaseUrl = GetEnvironment("KEYCLOAK_BASE_URL", "http://localhost:8080"),
      KeycloakRealm = GetEnvironment("KEYCLOAK_REALM", "erp-local"),
      KeycloakClientId = GetEnvironment("KEYCLOAK_IDENTITY_CLIENT_ID", "erp-identity"),
      KeycloakAdminUsername = GetEnvironment("KEYCLOAK_ADMIN_USER", "admin"),
      KeycloakAdminPassword = GetEnvironment("KEYCLOAK_ADMIN_PASSWORD", "admin"),
      BootstrapPassword = GetEnvironment("IDENTITY_BOOTSTRAP_PASSWORD", "Change.Me123!"),
      OpenFgaBaseUrl = GetEnvironment("OPENFGA_BASE_URL", "http://localhost:8080"),
      OpenFgaStoreName = GetEnvironment("OPENFGA_STORE_NAME", "erp-local-store"),
      MfaIssuer = GetEnvironment("IDENTITY_MFA_ISSUER", "ERP")
    };
  }

  private static string ResolvePostgresConnectionString()
  {
    var explicitConnectionString = Environment.GetEnvironmentVariable("IDENTITY_POSTGRES_CONNECTION_STRING");
    if (!string.IsNullOrWhiteSpace(explicitConnectionString))
    {
      return explicitConnectionString.Trim();
    }

    var builder = new NpgsqlConnectionStringBuilder
    {
      Host = GetEnvironment("IDENTITY_POSTGRES_HOST", "localhost"),
      Port = int.Parse(GetEnvironment("IDENTITY_POSTGRES_PORT", "5432")),
      Database = GetEnvironment("IDENTITY_POSTGRES_DB", GetEnvironment("ERP_POSTGRES_DB", "erp")),
      Username = GetEnvironment("IDENTITY_POSTGRES_USER", GetEnvironment("ERP_POSTGRES_USER", "erp")),
      Password = GetEnvironment("IDENTITY_POSTGRES_PASSWORD", GetEnvironment("ERP_POSTGRES_PASSWORD", "erp")),
      SslMode = Enum.Parse<SslMode>(GetEnvironment("IDENTITY_POSTGRES_SSL_MODE", "Disable"), ignoreCase: true)
    };

    return builder.ConnectionString;
  }

  private static string GetEnvironment(string key, string fallback)
  {
    var value = Environment.GetEnvironmentVariable(key);

    return string.IsNullOrWhiteSpace(value)
      ? fallback
      : value.Trim();
  }
}
