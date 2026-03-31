// Este arquivo concentra as rotas minimas e o bootstrap HTTP do servico.
// Crescimento de endpoints deve manter a API fina.
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Identity.Application;
using Identity.Contracts;

namespace Identity.Api;

public static class Server
{
  public static IEndpointRouteBuilder MapIdentityRoutes(this IEndpointRouteBuilder app)
  {
    app.MapGet("/health/live", () => TypedResults.Ok(new HealthResponse("identity", "live")));
    app.MapGet("/health/ready", () => TypedResults.Ok(new HealthResponse("identity", "ready")));
    app.MapGet(
      "/api/identity/tenants",
      (ListBootstrapTenants useCase) => TypedResults.Ok(useCase.Execute()));

    return app;
  }
}
