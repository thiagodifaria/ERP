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

  public string BootstrapPassword { get; init; } = "change-me-unsafe-local-only-bootstrap";

  public string OpenFgaBaseUrl { get; init; } = "http://localhost:8080";

  public string OpenFgaStoreName { get; init; } = "erp-local-store";

  public string MfaIssuer { get; init; } = "ERP";

  public static IdentityInfrastructureOptions Load()
  {
    var options = new IdentityInfrastructureOptions
    {
      RepositoryDriver = GetEnvironment("IDENTITY_REPOSITORY_DRIVER", "memory"),
      IdentityProviderDriver = GetEnvironment("IDENTITY_PROVIDER_DRIVER", "memory"),
      AuthorizationDriver = GetEnvironment("IDENTITY_AUTHORIZATION_DRIVER", "memory"),
      PostgresConnectionString = ResolvePostgresConnectionString(),
      KeycloakBaseUrl = GetEnvironment("KEYCLOAK_BASE_URL", "http://localhost:8080"),
      KeycloakRealm = GetEnvironment("KEYCLOAK_REALM", "erp-local"),
      KeycloakClientId = GetEnvironment("KEYCLOAK_IDENTITY_CLIENT_ID", "erp-identity"),
      KeycloakAdminUsername = GetEnvironment("KEYCLOAK_ADMIN_USER", "admin"),
      KeycloakAdminPassword = GetEnvironment("KEYCLOAK_ADMIN_PASSWORD", "change-me-unsafe-local-only-keycloak"),
      BootstrapPassword = GetEnvironment("IDENTITY_BOOTSTRAP_PASSWORD", "change-me-unsafe-local-only-bootstrap"),
      OpenFgaBaseUrl = GetEnvironment("OPENFGA_BASE_URL", "http://localhost:8080"),
      OpenFgaStoreName = GetEnvironment("OPENFGA_STORE_NAME", "erp-local-store"),
      MfaIssuer = GetEnvironment("IDENTITY_MFA_ISSUER", "ERP")
    };

    ValidateNonLocalSecret("KEYCLOAK_ADMIN_PASSWORD", options.KeycloakAdminPassword);
    ValidateNonLocalSecret("IDENTITY_BOOTSTRAP_PASSWORD", options.BootstrapPassword);

    return options;
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
      Password = GetEnvironment("IDENTITY_POSTGRES_PASSWORD", GetEnvironment("ERP_POSTGRES_PASSWORD", "change-me-unsafe-local-only-postgres")),
      SslMode = Enum.Parse<SslMode>(GetEnvironment("IDENTITY_POSTGRES_SSL_MODE", "Disable"), ignoreCase: true)
    };

    ValidateNonLocalSecret("IDENTITY_POSTGRES_PASSWORD/ERP_POSTGRES_PASSWORD", builder.Password);

    return builder.ConnectionString;
  }

  private static void ValidateNonLocalSecret(string key, string value)
  {
    if (IsLocalRuntime())
    {
      return;
    }

    if (string.IsNullOrWhiteSpace(value)
      || value is "erp" or "admin" or "Change.Me123!" or "local-jwt-secret" or "local-service-token" or "documents-local-secret"
      || value.StartsWith("change-me-unsafe-local-only", StringComparison.Ordinal)
      || value.Length < 32)
    {
      throw new InvalidOperationException($"{key} must be provided by a secret manager outside local/test environments.");
    }
  }

  private static bool IsLocalRuntime()
  {
    var environment = GetEnvironment("ERP_ENV", "local").ToLowerInvariant();
    return environment is "local" or "dev" or "development" or "test";
  }

  private static string GetEnvironment(string key, string fallback)
  {
    var value = Environment.GetEnvironmentVariable(key);

    return string.IsNullOrWhiteSpace(value)
      ? fallback
      : value.Trim();
  }
}
