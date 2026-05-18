using Npgsql;

namespace Billing.Api;

public static partial class Server
{
  private static void MapBillingHealthRoutes(WebApplication app)
  {
    app.MapGet("/health/live", () =>
      TypedResults.Ok(new HealthResponse("billing", "live")));

    app.MapGet("/health/ready", async Task<IResult> (NpgsqlDataSource dataSource) =>
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var command = new NpgsqlCommand("SELECT 1", connection);
      await command.ExecuteScalarAsync();
      return TypedResults.Ok(new HealthResponse("billing", "ready"));
    });

    app.MapGet("/health/details", async Task<IResult> (NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      var dependencies = await BuildDependencies(connection, configuration);
      return TypedResults.Ok(new ReadinessResponse("billing", dependencies.All(dep => dep.Status == "ready") ? "ready" : "degraded", dependencies));
    });
  }
}
