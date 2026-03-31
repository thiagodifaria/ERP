// Este arquivo centraliza a resolucao de configuracao da infraestrutura.
// O objetivo e alternar drivers sem espalhar variaveis de ambiente pelo servico.
using Npgsql;

namespace Identity.Infrastructure;

public sealed class IdentityInfrastructureOptions
{
  public string RepositoryDriver { get; init; } = "memory";

  public string PostgresConnectionString { get; init; } = string.Empty;

  public static IdentityInfrastructureOptions Load()
  {
    return new IdentityInfrastructureOptions
    {
      RepositoryDriver = GetEnvironment("IDENTITY_REPOSITORY_DRIVER", "memory"),
      PostgresConnectionString = ResolvePostgresConnectionString()
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
