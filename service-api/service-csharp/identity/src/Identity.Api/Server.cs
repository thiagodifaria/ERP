// Este arquivo concentra as rotas minimas e o bootstrap HTTP do servico.
// Crescimento de endpoints deve manter a API fina.
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Routing;
using Identity.Application;
using Identity.Contracts;
using Identity.Infrastructure;

namespace Identity.Api;

public static class Server
{
  public static IEndpointRouteBuilder MapIdentityRoutes(this IEndpointRouteBuilder app)
  {
    app.MapGet("/health/live", () => TypedResults.Ok(new HealthResponse("identity", "live")));
    app.MapGet("/health/ready", () => TypedResults.Ok(new HealthResponse("identity", "ready")));
    app.MapGet("/health/details", () => TypedResults.Ok(BuildReadiness()));
    app.MapPost(
      "/api/identity/tenants",
      Results<Created<TenantResponse>, BadRequest<ErrorResponse>, Conflict<ErrorResponse>>
      (CreateTenantRequest request, CreateBootstrapTenant useCase) =>
      {
        var result = useCase.Execute(request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{result.Tenant!.Slug}",
          result.Tenant);
      });
    app.MapPost(
      "/api/identity/tenants/{slug}/companies",
      Results<Created<CompanyResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, CreateCompanyRequest request, CreateBootstrapCompany useCase) =>
      {
        var result = useCase.Execute(slug, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{slug}/companies/{result.Company!.PublicId}",
          result.Company);
      });
    app.MapPatch(
      "/api/identity/tenants/{slug}/companies/{companyPublicId:guid}",
      Results<Ok<CompanyResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, Guid companyPublicId, UpdateCompanyRequest request, UpdateBootstrapCompany useCase) =>
      {
        var result = useCase.Execute(slug, companyPublicId, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Ok(result.Company!);
      });
    app.MapPost(
      "/api/identity/tenants/{slug}/users",
      Results<Created<UserResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, CreateUserRequest request, CreateBootstrapUser useCase) =>
      {
        var result = useCase.Execute(slug, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{slug}/users/{result.User!.PublicId}",
          result.User);
      });
    app.MapPatch(
      "/api/identity/tenants/{slug}/users/{userPublicId:guid}",
      Results<Ok<UserResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, Guid userPublicId, UpdateUserRequest request, UpdateBootstrapUser useCase) =>
      {
        var result = useCase.Execute(slug, userPublicId, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Ok(result.User!);
      });
    app.MapPost(
      "/api/identity/tenants/{slug}/teams",
      Results<Created<TeamResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, CreateTeamRequest request, CreateBootstrapTeam useCase) =>
      {
        var result = useCase.Execute(slug, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{slug}/teams/{result.Team!.PublicId}",
          result.Team);
      });
    app.MapPatch(
      "/api/identity/tenants/{slug}/teams/{teamPublicId:guid}",
      Results<Ok<TeamResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, Guid teamPublicId, UpdateTeamRequest request, UpdateBootstrapTeam useCase) =>
      {
        var result = useCase.Execute(slug, teamPublicId, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Ok(result.Team!);
      });
    app.MapPost(
      "/api/identity/tenants/{slug}/teams/{teamPublicId:guid}/members",
      Results<Created<TeamMembershipResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, Guid teamPublicId, AddTeamMemberRequest request, AddBootstrapTeamMember useCase) =>
      {
        var result = useCase.Execute(slug, teamPublicId, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{slug}/teams/{teamPublicId}/members/{result.Membership!.UserPublicId}",
          result.Membership);
      });
    app.MapPost(
      "/api/identity/tenants/{slug}/users/{userPublicId:guid}/roles",
      Results<Created<UserRoleResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>, Conflict<ErrorResponse>>
      (string slug, Guid userPublicId, AssignUserRoleRequest request, AssignBootstrapUserRole useCase) =>
      {
        var result = useCase.Execute(slug, userPublicId, request);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        if (result.IsConflict)
        {
          return TypedResults.Conflict(result.Error!);
        }

        return TypedResults.Created(
          $"/api/identity/tenants/{slug}/users/{userPublicId}/roles/{result.UserRole!.RoleCode}",
          result.UserRole);
      });
    app.MapDelete(
      "/api/identity/tenants/{slug}/users/{userPublicId:guid}/roles/{roleCode}",
      Results<Ok<UserRoleResponse>, BadRequest<ErrorResponse>, NotFound<ErrorResponse>>
      (string slug, Guid userPublicId, string roleCode, RevokeBootstrapUserRole useCase) =>
      {
        var result = useCase.Execute(slug, userPublicId, roleCode);

        if (result.IsBadRequest)
        {
          return TypedResults.BadRequest(result.Error!);
        }

        if (result.IsNotFound)
        {
          return TypedResults.NotFound(result.Error!);
        }

        return TypedResults.Ok(result.UserRole!);
      });
    app.MapGet(
      "/api/identity/tenants",
      (ListBootstrapTenants useCase) => TypedResults.Ok(useCase.Execute()));
    app.MapGet(
      "/api/identity/tenants/{slug}",
      Results<Ok<TenantResponse>, NotFound> (string slug, GetBootstrapTenantBySlug useCase) =>
      {
        var tenant = useCase.Execute(slug);

        return tenant is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(tenant);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/snapshot",
      Results<Ok<TenantAccessSnapshotResponse>, NotFound> (
        string slug,
        GetBootstrapTenantAccessSnapshot useCase) =>
      {
        var snapshot = useCase.Execute(slug);

        return snapshot is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(snapshot);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/companies",
      Results<Ok<IReadOnlyCollection<CompanyResponse>>, NotFound> (string slug, ListBootstrapCompanies useCase) =>
      {
        var companies = useCase.Execute(slug);

        return companies is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(companies);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/users",
      Results<Ok<IReadOnlyCollection<UserResponse>>, NotFound> (string slug, ListBootstrapUsers useCase) =>
      {
        var users = useCase.Execute(slug);

        return users is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(users);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/users/{userPublicId:guid}/roles",
      Results<Ok<IReadOnlyCollection<UserRoleResponse>>, NotFound> (
        string slug,
        Guid userPublicId,
        ListBootstrapUserRoles useCase) =>
      {
        var userRoles = useCase.Execute(slug, userPublicId);

        return userRoles is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(userRoles);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/teams",
      Results<Ok<IReadOnlyCollection<TeamResponse>>, NotFound> (string slug, ListBootstrapTeams useCase) =>
      {
        var teams = useCase.Execute(slug);

        return teams is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(teams);
      });
    app.MapGet(
      "/api/identity/tenants/{slug}/teams/{teamPublicId:guid}/members",
      Results<Ok<IReadOnlyCollection<TeamMembershipResponse>>, NotFound> (
        string slug,
        Guid teamPublicId,
        ListBootstrapTeamMembers useCase) =>
      {
        var members = useCase.Execute(slug, teamPublicId);

        return members is null
          ? TypedResults.NotFound()
          : TypedResults.Ok(members);
      });
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
    var repositoryDriver = IdentityInfrastructureOptions.Load().RepositoryDriver;
    var postgresqlStatus = repositoryDriver.Equals("postgres", StringComparison.OrdinalIgnoreCase)
      ? "ready"
      : "pending-runtime-wiring";

    return new ReadinessResponse(
      "identity",
      "ready",
      [
        new DependencyHealthResponse("tenant-catalog", "ready"),
        new DependencyHealthResponse("bootstrap-api", "ready"),
        new DependencyHealthResponse("postgresql", postgresqlStatus)
      ]);
  }
}
