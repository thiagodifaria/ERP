// Este arquivo inicia a API do contexto de billing.
// O bootstrap deve ficar pequeno para facilitar evolucao do servico.
using Billing.Api;
using Npgsql;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton(_ => NpgsqlDataSource.Create(BuildConnectionString(builder.Configuration)));

var app = builder.Build();

app.MapBillingRoutes();

app.Run();

static string BuildConnectionString(ConfigurationManager configuration)
{
  var host = configuration["BILLING_POSTGRES_HOST"] ?? "service-postgresql";
  var port = configuration["BILLING_POSTGRES_PORT"] ?? "5432";
  var database = configuration["BILLING_POSTGRES_DB"] ?? "erp";
  var user = configuration["BILLING_POSTGRES_USER"] ?? "erp";
  var password = configuration["BILLING_POSTGRES_PASSWORD"] ?? "erp";
  var sslMode = configuration["BILLING_POSTGRES_SSL_MODE"] ?? "Disable";

  var builder = new NpgsqlConnectionStringBuilder
  {
    Host = host,
    Port = int.Parse(port),
    Database = database,
    Username = user,
    Password = password,
    SslMode = Enum.TryParse<SslMode>(sslMode, ignoreCase: true, out var parsedSslMode)
      ? parsedSslMode
      : SslMode.Disable
  };

  return builder.ConnectionString;
}

public partial class Program;
