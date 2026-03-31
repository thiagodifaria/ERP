// Este caso de uso cria times basicos vinculados a um tenant existente.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class CreateBootstrapTeam
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamRepository _teamRepository;

  public CreateBootstrapTeam(ITenantCatalog tenantCatalog, ITeamRepository teamRepository)
  {
    _tenantCatalog = tenantCatalog;
    _teamRepository = teamRepository;
  }

  public CreateBootstrapTeamResult Execute(string tenantSlug, CreateTeamRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return CreateBootstrapTeamResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var name = request.Name.Trim();

    if (string.IsNullOrWhiteSpace(name))
    {
      return CreateBootstrapTeamResult.BadRequest(
        new ErrorResponse("invalid_team_name", "Team name is required."));
    }

    if (_teamRepository.FindByTenantIdAndName(tenant.Id, name) is not null)
    {
      return CreateBootstrapTeamResult.Conflict(
        new ErrorResponse("team_name_conflict", "Team name already exists for tenant."));
    }

    var team = new Team(
      _teamRepository.NextId(),
      tenant.Id,
      null,
      PublicIds.NewUuidV7(),
      name,
      "active");

    var createdTeam = _teamRepository.Add(team);

    return CreateBootstrapTeamResult.Success(new TeamResponse(
      createdTeam.Id,
      createdTeam.PublicId,
      createdTeam.TenantId,
      createdTeam.CompanyId,
      createdTeam.Name,
      createdTeam.Status));
  }
}
