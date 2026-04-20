// Estes testes protegem os contratos publicos mais importantes da API de identidade.
using System.Net;
using System.Net.Http.Json;
using Identity.Contracts;
using Microsoft.AspNetCore.Mvc.Testing;
using Xunit;

namespace Identity.ContractTests;

public sealed class IdentityContractsTests : IClassFixture<WebApplicationFactory<Program>>
{
  private readonly HttpClient _client;

  public IdentityContractsTests(WebApplicationFactory<Program> factory)
  {
    _client = factory.CreateClient();
  }

  [Fact]
  public async Task TenantListContractShouldExposePublicFields()
  {
    var response = await _client.GetAsync("/api/identity/tenants");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.All(payload, tenant =>
    {
      Assert.NotEqual(0, tenant.Id);
      Assert.NotEqual(Guid.Empty, tenant.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(tenant.Slug));
      Assert.False(string.IsNullOrWhiteSpace(tenant.DisplayName));
      Assert.False(string.IsNullOrWhiteSpace(tenant.Status));
    });
  }

  [Fact]
  public async Task CreateTenantContractShouldReturnCreatedResource()
  {
    var request = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Tenant");

    var response = await _client.PostAsJsonAsync("/api/identity/tenants", request);

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(payload);
    Assert.Equal(request.Slug, payload!.Slug);
    Assert.Equal(request.DisplayName, payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task SnapshotContractShouldExposeConsistentNestedStructure()
  {
    var response = await _client.GetAsync("/api/identity/tenants/bootstrap-ops/snapshot");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TenantAccessSnapshotResponse>();

    Assert.NotNull(payload);
    Assert.Equal("bootstrap-ops", payload!.Tenant.Slug);
    Assert.True(payload.Counts.Companies >= 1);
    Assert.True(payload.Counts.Users >= 1);
    Assert.True(payload.Counts.Teams >= 1);
    Assert.True(payload.Counts.Roles >= 1);
    Assert.All(payload.Users, user =>
    {
      Assert.NotEqual(Guid.Empty, user.User.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(user.User.Email));
      Assert.All(user.Roles, role => Assert.False(string.IsNullOrWhiteSpace(role.RoleCode)));
    });
    Assert.All(payload.Teams, team =>
    {
      Assert.NotEqual(Guid.Empty, team.Team.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(team.Team.Name));
      Assert.All(team.Members, member => Assert.NotEqual(Guid.Empty, member.UserPublicId));
    });
  }

  [Fact]
  public async Task UpdateCompanyContractShouldReturnUpdatedResourceShape()
  {
    var companies = await _client.GetFromJsonAsync<CompanyResponse[]>("/api/identity/tenants/bootstrap-ops/companies");

    Assert.NotNull(companies);
    Assert.NotEmpty(companies);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/bootstrap-ops/companies/{companies![0].PublicId}",
      new UpdateCompanyRequest("Bootstrap Ops Contract", "Bootstrap Ops Contract LTDA", "12345678901234"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal("Bootstrap Ops Contract", payload.DisplayName);
    Assert.Equal("Bootstrap Ops Contract LTDA", payload.LegalName);
    Assert.Equal("12345678901234", payload.TaxId);
  }

  [Fact]
  public async Task CompanyDetailContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Company Detail");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var companiesResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/companies");
    var companiesPayload = await companiesResponse.Content.ReadFromJsonAsync<CompanyResponse[]>();

    Assert.NotNull(companiesPayload);
    Assert.NotEmpty(companiesPayload);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/companies/{companiesPayload![0].PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<CompanyResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal(tenantPayload.DisplayName, payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task UpdateUserContractShouldReturnUpdatedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract User Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.user@{tenantPayload.Slug}.local",
        "Contract User",
        "Contract",
        "User"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}",
      new UpdateUserRequest(
        $"contract.user.prime@{tenantPayload.Slug}.local",
        "Contract User Prime",
        "Contract",
        "Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal($"contract.user.prime@{tenantPayload.Slug}.local", payload.Email);
    Assert.Equal("Contract User Prime", payload.DisplayName);
    Assert.Equal("Prime", payload.FamilyName);
  }

  [Fact]
  public async Task UserDetailContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract User Detail");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.user.detail@{tenantPayload.Slug}.local",
        "Contract User Detail",
        "Contract",
        "Detail"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal($"contract.user.detail@{tenantPayload.Slug}.local", payload.Email);
    Assert.Equal("Contract User Detail", payload.DisplayName);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task RevokeUserRoleContractShouldReturnRevokedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Role Revoke");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/users",
      new CreateUserRequest(
        $"contract.revoke@{tenantPayload.Slug}.local",
        "Contract Revoke User",
        "Contract",
        "Revoke"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser!.PublicId}/roles",
      new AssignUserRoleRequest("viewer"));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{createdUser.PublicId}/roles/viewer");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserRoleResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.UserPublicId);
    Assert.NotEqual(Guid.Empty, payload.RolePublicId);
    Assert.Equal("viewer", payload.RoleCode);
    Assert.False(string.IsNullOrWhiteSpace(payload.RoleDisplayName));
  }

  [Fact]
  public async Task UpdateTeamContractShouldReturnUpdatedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Team Update");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teams = await _client.GetFromJsonAsync<TeamResponse[]>(
      $"/api/identity/tenants/{tenantPayload!.Slug}/teams");

    Assert.NotNull(teams);
    Assert.NotEmpty(teams);

    var response = await _client.PatchAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teams![0].PublicId}",
      new UpdateTeamRequest("Core Contract Prime"));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal("Core Contract Prime", payload.Name);
    Assert.Equal("active", payload.Status);
  }

  [Fact]
  public async Task RemoveTeamMemberContractShouldReturnRemovedResourceShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Membership Removal");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var teamsResponse = await _client.GetAsync($"/api/identity/tenants/{tenantPayload!.Slug}/teams");
    var teamsPayload = await teamsResponse.Content.ReadFromJsonAsync<TeamResponse[]>();

    Assert.NotNull(teamsPayload);

    var createUserResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users",
      new CreateUserRequest(
        $"contract.membership.remove@{tenantPayload.Slug}.local",
        "Contract Membership Remove",
        "Contract",
        "Remove"));

    var createdUser = await createUserResponse.Content.ReadFromJsonAsync<UserResponse>();

    Assert.NotNull(createdUser);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload![0].PublicId}/members",
      new AddTeamMemberRequest(createdUser!.PublicId));

    var response = await _client.DeleteAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/teams/{teamsPayload[0].PublicId}/members/{createdUser.PublicId}");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<TeamMembershipResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.TeamPublicId);
    Assert.NotEqual(Guid.Empty, payload.UserPublicId);
    Assert.False(string.IsNullOrWhiteSpace(payload.UserEmail));
    Assert.False(string.IsNullOrWhiteSpace(payload.UserDisplayName));
  }

  [Fact]
  public async Task ErrorContractShouldExposeCodeAndMessage()
  {
    var response = await _client.PostAsJsonAsync(
      "/api/identity/tenants/bootstrap-ops/users",
      new CreateUserRequest("invalid-email", "Invalid User", null, null));

    Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<ErrorResponse>();

    Assert.NotNull(payload);
    Assert.False(string.IsNullOrWhiteSpace(payload!.Code));
    Assert.False(string.IsNullOrWhiteSpace(payload.Message));
  }

  [Fact]
  public async Task InviteContractShouldExposeAssignmentsAndToken()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Invite Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.invite@{tenantPayload.Slug}.local",
        "Contract Invite User",
        ["viewer"],
        null,
        7));

    Assert.Equal(HttpStatusCode.Created, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.NotEqual(Guid.Empty, payload.UserPublicId);
    Assert.False(string.IsNullOrWhiteSpace(payload.Email));
    Assert.False(string.IsNullOrWhiteSpace(payload.InviteToken));
    Assert.Contains("viewer", payload.RoleCodes);
    Assert.Equal("pending", payload.Status);
  }

  [Fact]
  public async Task SessionLoginContractShouldExposeOpaqueTokensAndRoles()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Session Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.session@{tenantPayload.Slug}.local",
        "Contract Session User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Contract Session User", "Contract", "Session", "PhaseTwo123"));

    var response = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        $"contract.session@{tenantPayload.Slug}.local",
        "PhaseTwo123",
        null));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.NotEqual(Guid.Empty, payload.UserPublicId);
    Assert.False(string.IsNullOrWhiteSpace(payload.SessionToken));
    Assert.False(string.IsNullOrWhiteSpace(payload.RefreshToken));
    Assert.True(payload.ExpiresAt > DateTimeOffset.UtcNow);
    Assert.True(payload.RefreshExpiresAt > payload.ExpiresAt);
    Assert.Contains("viewer", payload.RoleCodes);
  }

  [Fact]
  public async Task SecurityAuditContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Audit Flow");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.audit@{tenantPayload.Slug}.local",
        "Contract Audit User",
        ["viewer"],
        null,
        7));

    var response = await _client.GetAsync($"/api/identity/tenants/{tenantPayload.Slug}/security/audit");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<SecurityAuditEventResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.All(payload, auditEvent =>
    {
      Assert.NotEqual(Guid.Empty, auditEvent.PublicId);
      Assert.False(string.IsNullOrWhiteSpace(auditEvent.EventCode));
      Assert.False(string.IsNullOrWhiteSpace(auditEvent.Severity));
      Assert.False(string.IsNullOrWhiteSpace(auditEvent.Summary));
    });
  }

  [Fact]
  public async Task PasswordRecoveryContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Password Recovery");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.password@{tenantPayload.Slug}.local",
        "Contract Password User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Contract Password User", "Contract", "Password", "PhaseTwo123"));

    var response = await _client.PostAsJsonAsync(
      "/api/identity/password-recovery",
      new StartPasswordRecoveryRequest(
        tenantPayload.Slug,
        $"contract.password@{tenantPayload.Slug}.local",
        30));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<PasswordRecoveryResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.NotEqual(Guid.Empty, payload.UserPublicId);
    Assert.Equal("pending", payload.Status);
    Assert.False(string.IsNullOrWhiteSpace(payload.ResetToken));
    Assert.True(payload.ExpiresAt > DateTimeOffset.UtcNow);
  }

  [Fact]
  public async Task UserSessionListContractShouldExposePublicFields()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Session List");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.session.list@{tenantPayload.Slug}.local",
        "Contract Session List User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Contract Session List User", "Contract", "Session", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));

    var response = await _client.GetAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/users/{acceptedPayload.User.PublicId}/sessions");

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserSessionResponse[]>();

    Assert.NotNull(payload);
    Assert.NotEmpty(payload);
    Assert.All(payload!, session =>
    {
      Assert.NotEqual(Guid.Empty, session.PublicId);
      Assert.NotEqual(Guid.Empty, session.UserPublicId);
      Assert.False(string.IsNullOrWhiteSpace(session.Status));
      Assert.True(session.ExpiresAt > DateTimeOffset.UtcNow);
      Assert.True(session.RefreshExpiresAt >= session.ExpiresAt);
    });
  }

  [Fact]
  public async Task ResendInviteContractShouldExposeUpdatedInviteShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Invite Resend");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.invite.resend@{tenantPayload.Slug}.local",
        "Contract Invite Resend User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var response = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload.Slug}/invites/{invitePayload!.PublicId}/resend",
      new ResendInviteRequest(10));

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(payload);
    Assert.Equal(invitePayload.PublicId, payload!.PublicId);
    Assert.Equal("pending", payload.Status);
    Assert.Equal($"contract.invite.resend@{tenantPayload.Slug}.local", payload.Email);
    Assert.False(string.IsNullOrWhiteSpace(payload.InviteToken));
    Assert.NotEqual(invitePayload.InviteToken, payload.InviteToken);
    Assert.True(payload.ExpiresAt > DateTimeOffset.UtcNow);
  }

  [Fact]
  public async Task LogoutContractShouldExposeRevokedSessionShape()
  {
    var tenantRequest = new CreateTenantRequest(
      $"tenant-{Guid.NewGuid():N}"[..15],
      "Contract Logout");

    var tenantResponse = await _client.PostAsJsonAsync("/api/identity/tenants", tenantRequest);
    var tenantPayload = await tenantResponse.Content.ReadFromJsonAsync<TenantResponse>();

    Assert.NotNull(tenantPayload);

    var inviteResponse = await _client.PostAsJsonAsync(
      $"/api/identity/tenants/{tenantPayload!.Slug}/invites",
      new CreateInviteRequest(
        $"contract.logout@{tenantPayload.Slug}.local",
        "Contract Logout User",
        ["viewer"],
        null,
        7));
    var invitePayload = await inviteResponse.Content.ReadFromJsonAsync<InviteResponse>();

    Assert.NotNull(invitePayload);

    var acceptResponse = await _client.PostAsJsonAsync(
      $"/api/identity/invites/{invitePayload!.InviteToken}/accept",
      new AcceptInviteRequest("Contract Logout User", "Contract", "Logout", "PhaseTwo123"));
    var acceptedPayload = await acceptResponse.Content.ReadFromJsonAsync<AcceptInviteResponse>();

    Assert.NotNull(acceptedPayload);

    var loginResponse = await _client.PostAsJsonAsync(
      "/api/identity/sessions/login",
      new LoginSessionRequest(
        tenantPayload.Slug,
        acceptedPayload!.User.Email,
        "PhaseTwo123",
        null));
    var loginPayload = await loginResponse.Content.ReadFromJsonAsync<SessionResponse>();

    Assert.NotNull(loginPayload);

    using var request = new HttpRequestMessage(HttpMethod.Post, "/api/identity/sessions/logout");
    request.Headers.Authorization = new System.Net.Http.Headers.AuthenticationHeaderValue("Bearer", loginPayload!.SessionToken);

    var response = await _client.SendAsync(request);

    Assert.Equal(HttpStatusCode.OK, response.StatusCode);

    var payload = await response.Content.ReadFromJsonAsync<UserSessionResponse>();

    Assert.NotNull(payload);
    Assert.NotEqual(Guid.Empty, payload!.PublicId);
    Assert.Equal(acceptedPayload.User.PublicId, payload.UserPublicId);
    Assert.Equal("revoked", payload.Status);
    Assert.NotNull(payload.RevokedAt);
  }
}
