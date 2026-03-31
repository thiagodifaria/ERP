// TenantStructureCountsResponse resume a composicao basica do tenant no bootstrap.
namespace Identity.Contracts;

public sealed record TenantStructureCountsResponse(
  int Companies,
  int Users,
  int Teams,
  int Roles,
  int TeamMemberships,
  int UserRoles);
