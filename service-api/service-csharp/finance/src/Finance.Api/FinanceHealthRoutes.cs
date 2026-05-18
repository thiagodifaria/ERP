using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Npgsql;

namespace Finance.Api;

public static partial class Server
{
  private static void MapFinanceHealthRoutes(IEndpointRouteBuilder app)
  {
    app.MapGet("/health/live", () => TypedResults.Ok(new HealthResponse("finance", "live")));
    app.MapGet("/health/ready", async Task<IResult> (NpgsqlDataSource dataSource) =>
    {
      return await CanReachDatabase(dataSource)
        ? TypedResults.Ok(new HealthResponse("finance", "ready"))
        : TypedResults.StatusCode(StatusCodes.Status503ServiceUnavailable);
    });
    app.MapGet("/health/details", async Task<IResult> (NpgsqlDataSource dataSource) =>
    {
      var status = await CanReachDatabase(dataSource) ? "ready" : "degraded";
      return TypedResults.Ok(new ReadinessResponse("finance", status, [new DependencyResponse("postgresql", status)]));
    });
  }
}
