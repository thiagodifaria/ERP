// Este caso de uso atualiza times basicos ja existentes dentro de um tenant conhecido.
using Identity.Contracts;
using Identity.Domain;

namespace Identity.Application;

public sealed class UpdateBootstrapTeam
{
  private readonly ITenantCatalog _tenantCatalog;
  private readonly ITeamRepository _teamRepository;

  public UpdateBootstrapTeam(ITenantCatalog tenantCatalog, ITeamRepository teamRepository)
  {
    _tenantCatalog = tenantCatalog;
    _teamRepository = teamRepository;
  }

  public UpdateBootstrapTeamResult Execute(string tenantSlug, Guid teamPublicId, UpdateTeamRequest request)
  {
    var tenant = _tenantCatalog.FindBySlug(tenantSlug);

    if (tenant is null)
    {
      return UpdateBootstrapTeamResult.NotFound(
        new ErrorResponse("tenant_not_found", "Tenant was not found."));
    }

    var team = _teamRepository.FindByTenantIdAndPublicId(tenant.Id, teamPublicId);

    if (team is null)
    {
      return UpdateBootstrapTeamResult.NotFound(
        new ErrorResponse("team_not_found", "Team was not found."));
    }

    var name = request.Name.Trim();

    if (string.IsNullOrWhiteSpace(name))
    {
      return UpdateBootstrapTeamResult.BadRequest(
        new ErrorResponse("invalid_team_name", "Team name is required."));
    }

    var existingTeam = _teamRepository.FindByTenantIdAndName(tenant.Id, name);

    if (existingTeam is not null && existingTeam.PublicId != team.PublicId)
    {
      return UpdateBootstrapTeamResult.Conflict(
        new ErrorResponse("team_name_conflict", "Team name already exists for tenant."));
    }

    var updatedTeam = _teamRepository.Update(team.ReviseProfile(name));

    return UpdateBootstrapTeamResult.Success(new TeamResponse(
      updatedTeam.Id,
      updatedTeam.PublicId,
      updatedTeam.TenantId,
      updatedTeam.CompanyId,
      updatedTeam.Name,
      updatedTeam.Status));
  }
}
