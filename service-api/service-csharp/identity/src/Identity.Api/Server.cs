// Este arquivo concentra as rotas minimas e o bootstrap HTTP do servico.
// Crescimento de endpoints deve manter a API fina.
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.HttpResults;
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
    app.MapGet("/health/details", () => TypedResults.Ok(BuildReadiness()));
    app.MapGet(
      "/api/identity/tenants",
      (ListBootstrapTenants useCase) => TypedResults.Ok(useCase.Execute()));
    app.MapGet(
      "/api/identity/tenants/{slug}/roles",
      Results<Ok<IReadOnlyCollection<RoleResponse>>, NotFound> (string slug, ListBootstrapRoles useCase) =>
      {
        var roles = useCase.Execute(slug);

        return roles.Count == 0
          ? TypedResults.NotFound()
          : TypedResults.Ok(roles);
      });

    return app;
  }

  private static ReadinessResponse BuildReadiness()
  {
    return new ReadinessResponse(
      "identity",
      "ready",
      [
        new DependencyHealthResponse("tenant-catalog", "ready"),
        new DependencyHealthResponse("bootstrap-api", "ready"),
        new DependencyHealthResponse("postgresql", "pending-runtime-wiring")
      ]);
  }
}
