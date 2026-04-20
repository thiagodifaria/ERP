// Este arquivo inicia a API do contexto financeiro.
// O bootstrap deve ficar pequeno para facilitar evolucao do servico.
using Finance.Api;
using Npgsql;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddSingleton(_ => NpgsqlDataSource.Create(BuildConnectionString(builder.Configuration)));

var app = builder.Build();

app.MapFinanceRoutes();

app.Run();

static string BuildConnectionString(ConfigurationManager configuration)
{
  var host = configuration["FINANCE_POSTGRES_HOST"] ?? "service-postgresql";
  var port = configuration["FINANCE_POSTGRES_PORT"] ?? "5432";
  var database = configuration["FINANCE_POSTGRES_DB"] ?? "erp";
  var user = configuration["FINANCE_POSTGRES_USER"] ?? "erp";
  var password = configuration["FINANCE_POSTGRES_PASSWORD"] ?? "erp";
  var sslMode = configuration["FINANCE_POSTGRES_SSL_MODE"] ?? "Disable";

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
