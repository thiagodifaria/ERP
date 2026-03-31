// TenantAccessSnapshotResponse consolida a estrutura inicial de acesso e organizacao do tenant.
namespace Identity.Contracts;

public sealed record TenantAccessSnapshotResponse(
  TenantResponse Tenant,
  TenantStructureCountsResponse Counts,
  IReadOnlyCollection<CompanyResponse> Companies,
  IReadOnlyCollection<UserAccessSnapshotResponse> Users,
  IReadOnlyCollection<TeamAccessSnapshotResponse> Teams,
  IReadOnlyCollection<RoleResponse> Roles);
