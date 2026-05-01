// Este arquivo concentra a superficie HTTP do finance.
// O servico separa projeção analitica do ciclo financeiro operacional.
using System.Data;
using System.Globalization;
using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Npgsql;
using NpgsqlTypes;

namespace Finance.Api;

public static class Server
{
  public static IEndpointRouteBuilder MapFinanceRoutes(this IEndpointRouteBuilder app)
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

    app.MapPost("/api/finance/projections/ingest", async Task<IResult> (ProjectionIngestRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var tenantSlug = ResolveTenantSlug(request?.TenantSlug, configuration);
      var limit = request?.Limit is > 0 and <= 500 ? request.Limit.Value : 100;

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.Ok(new ProjectionIngestResponse(tenantSlug, 0, 0, 0, 0));
      }

      var pendingEvents = await LoadPendingOutboxEvents(connection, transaction, tenantId.Value, limit);
      var created = 0;
      var updated = 0;
      var processed = 0;

      foreach (var outboxEvent in pendingEvents)
      {
        var outcome = await ApplyOutboxEvent(connection, transaction, tenantId.Value, outboxEvent);
        created += outcome.Created;
        updated += outcome.Updated;
        processed += outcome.Processed;
      }

      await transaction.CommitAsync();
      return TypedResults.Ok(new ProjectionIngestResponse(tenantSlug, pendingEvents.Count, processed, created, updated));
    });

    app.MapGet("/api/finance/projections", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var projections = await ListProjections(connection, resolvedTenantSlug, status);
      return TypedResults.Ok<IReadOnlyList<ProjectionResponse>>(projections);
    });

    app.MapGet("/api/finance/projections/summary", async Task<IResult> (string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var summary = await BuildProjectionSummary(connection, resolvedTenantSlug);
      return TypedResults.Ok(summary);
    });

    app.MapPost("/api/finance/operations/sync", async Task<IResult> (OperationsSyncRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var tenantSlug = ResolveTenantSlug(request?.TenantSlug, configuration);

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.Ok(new OperationsSyncResponse(tenantSlug, 0, 0, 0, 0));
      }

      var salesReceivables = await SyncReceivableEntries(connection, transaction, tenantId.Value);
      var rentalReceivables = await SyncRentalReceivableEntries(connection, transaction, tenantId.Value);
      var rentalSettlements = await SyncRentalReceivableSettlements(connection, transaction, tenantId.Value);
      var commissions = await SyncCommissionEntries(connection, transaction, tenantId.Value);
      await InsertActivityEvent(
        connection,
        transaction,
        tenantId.Value,
        "receivable_synced",
        "finance.operations",
        null,
        $"Finance operations sync consolidated sales and rentals with {salesReceivables.Created + rentalReceivables.Created} receivables created, {salesReceivables.Updated + rentalReceivables.Updated} updated and {rentalSettlements.Created} rental settlements materialized.",
        "finance-sync",
        JsonSerializer.Serialize(new
        {
          sales = new { created = salesReceivables.Created, updated = salesReceivables.Updated },
          rentals = new { created = rentalReceivables.Created, updated = rentalReceivables.Updated, settlements = rentalSettlements.Created },
          commissions = new { created = commissions.Created, updated = commissions.Updated }
        }));

      await transaction.CommitAsync();
      return TypedResults.Ok(new OperationsSyncResponse(
        tenantSlug,
        salesReceivables.Created + rentalReceivables.Created,
        salesReceivables.Updated + rentalReceivables.Updated,
        commissions.Created,
        commissions.Updated));
    });

    app.MapGet("/api/finance/receivables", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var receivables = await ListReceivables(connection, resolvedTenantSlug, status);
      return TypedResults.Ok<IReadOnlyList<ReceivableResponse>>(receivables);
    });

    app.MapPost("/api/finance/receivables/{publicId}/settlements", async Task<IResult> (string publicId, ReceivableSettlementRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.SettlementReference))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_settlement_reference", "Settlement reference is required."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var existingByReference = await FindReceivableSettlementByReference(connection, transaction, tenantId.Value, request.SettlementReference.Trim());
      if (existingByReference is not null)
      {
        if (!string.Equals(existingByReference.ReceivablePublicId, publicId, StringComparison.OrdinalIgnoreCase))
        {
          await transaction.RollbackAsync();
          return TypedResults.Conflict(new ErrorResponse("settlement_reference_conflict", "Settlement reference already belongs to another receivable."));
        }

        await transaction.CommitAsync();
        return TypedResults.Ok(existingByReference with { Idempotent = true });
      }

      var receivable = await FindReceivableForSettlement(connection, transaction, tenantId.Value, publicId);
      if (receivable is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("receivable_not_found", "Receivable was not found."));
      }

      if (receivable.Status == "cancelled")
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("receivable_not_settleable", "Cancelled receivables cannot be settled."));
      }

      if (receivable.Status == "paid")
      {
        await transaction.RollbackAsync();
        return TypedResults.Conflict(new ErrorResponse("receivable_already_paid", "Receivable is already paid."));
      }

      var targetAmount = request.AmountCents ?? receivable.AmountCents;
      if (targetAmount != receivable.AmountCents)
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("invalid_settlement_amount", "Settlement amount must match the receivable amount."));
      }

      var settledAt = ParseUtcOrNow(request.SettledAt);
      if (settledAt is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("invalid_settlement_timestamp", "Settlement timestamp must be a valid UTC instant."));
      }

      var settlement = await CreateReceivableSettlement(
        connection,
        transaction,
        tenantId.Value,
        receivable,
        request.SettlementReference.Trim(),
        targetAmount,
        settledAt.Value);

      await transaction.CommitAsync();
      return TypedResults.Ok(settlement);
    });

    app.MapGet("/api/finance/commissions", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var commissions = await ListCommissions(connection, resolvedTenantSlug, status);
      return TypedResults.Ok<IReadOnlyList<CommissionResponse>>(commissions);
    });

    app.MapGet("/api/finance/commissions/summary", async Task<IResult> (string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var summary = await BuildCommissionSummary(connection, resolvedTenantSlug);
      return TypedResults.Ok(summary);
    });

    app.MapPost("/api/finance/commissions/{publicId}/block", async Task<IResult> (string publicId, CommissionLifecycleRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var commission = await FindCommissionForLifecycle(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (commission is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("commission_not_found", "Commission was not found."));
      }

      if (commission.Status == "blocked")
      {
        await transaction.CommitAsync();
        return TypedResults.Ok(await GetCommissionByPublicId(connection, request.TenantSlug.Trim(), publicId)
          ?? throw new InvalidOperationException("commission_lookup_failed"));
      }

      if (commission.Status != "pending")
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("commission_transition_invalid", "Only pending commissions can be blocked."));
      }

      var response = await UpdateCommissionStatus(connection, transaction, commission, "blocked", request.Actor, request.Reason);
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapPost("/api/finance/commissions/{publicId}/release", async Task<IResult> (string publicId, CommissionLifecycleRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var commission = await FindCommissionForLifecycle(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (commission is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("commission_not_found", "Commission was not found."));
      }

      if (commission.Status == "released")
      {
        await transaction.CommitAsync();
        return TypedResults.Ok(await GetCommissionByPublicId(connection, request.TenantSlug.Trim(), publicId)
          ?? throw new InvalidOperationException("commission_lookup_failed"));
      }

      if (commission.Status is not ("pending" or "blocked"))
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("commission_transition_invalid", "Only pending or blocked commissions can be released."));
      }

      var response = await UpdateCommissionStatus(connection, transaction, commission, "released", request.Actor, request.Reason);
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapGet("/api/finance/payables", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var payables = await ListPayables(connection, resolvedTenantSlug, status);
      return TypedResults.Ok<IReadOnlyList<PayableResponse>>(payables);
    });

    app.MapPost("/api/finance/payables", async Task<IResult> (CreatePayableRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null
        || string.IsNullOrWhiteSpace(request.Category)
        || string.IsNullOrWhiteSpace(request.VendorName)
        || string.IsNullOrWhiteSpace(request.Description)
        || request.AmountCents <= 0
        || !DateOnly.TryParseExact(request.DueDate, "yyyy-MM-dd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var dueDate))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_payable", "Payable payload is invalid."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var amountCents = request.AmountCents.GetValueOrDefault();
      var payable = await CreatePayable(
        connection,
        transaction,
        tenantId.Value,
        tenantSlug,
        request.Category.Trim(),
        request.VendorName.Trim(),
        request.Description.Trim(),
        amountCents,
        dueDate);

      await transaction.CommitAsync();
      return TypedResults.Ok(payable);
    });

    app.MapPatch("/api/finance/payables/{publicId}/status", async Task<IResult> (string publicId, UpdatePayableStatusRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.Status))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_payable_status", "Target payable status is required."));
      }

      var normalizedStatus = NormalizePayableStatus(request.Status);
      if (normalizedStatus.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_payable_status", "Target payable status is invalid."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var payable = await FindPayableForUpdate(connection, transaction, tenantId.Value, publicId);
      if (payable is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("payable_not_found", "Payable was not found."));
      }

      if (payable.Status == normalizedStatus)
      {
        await transaction.CommitAsync();
        return TypedResults.Ok(payable);
      }

      if (!CanTransitionPayableStatus(payable.Status, normalizedStatus))
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("invalid_payable_status_transition", "Payable status transition is invalid."));
      }

      var paymentReference = string.IsNullOrWhiteSpace(request.PaymentReference) ? null : request.PaymentReference.Trim();
      if (normalizedStatus == "paid" && paymentReference is not null)
      {
        var existingReference = await FindPayableByPaymentReference(connection, transaction, tenantId.Value, paymentReference);
        if (existingReference is not null && !string.Equals(existingReference.PublicId, payable.PublicId, StringComparison.OrdinalIgnoreCase))
        {
          await transaction.RollbackAsync();
          return TypedResults.Conflict(new ErrorResponse("payment_reference_conflict", "Payment reference already belongs to another payable."));
        }
      }

      var paidAt = normalizedStatus == "paid" ? ParseUtcOrNow(request.PaidAt) : null;
      if (normalizedStatus == "paid" && paidAt is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.BadRequest(new ErrorResponse("invalid_payable_paid_at", "Payable paid timestamp must be a valid UTC instant."));
      }

      var updated = await UpdatePayableStatus(
        connection,
        transaction,
        tenantId.Value,
        publicId,
        normalizedStatus,
        paymentReference,
        paidAt);

      await transaction.CommitAsync();
      return TypedResults.Ok(updated);
    });

    app.MapGet("/api/finance/costs", async Task<IResult> (string? tenantSlug, string? category, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var costs = await ListCosts(connection, resolvedTenantSlug, category);
      return TypedResults.Ok<IReadOnlyList<CostEntryResponse>>(costs);
    });

    app.MapPost("/api/finance/costs", async Task<IResult> (CreateCostEntryRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null
        || string.IsNullOrWhiteSpace(request.Category)
        || string.IsNullOrWhiteSpace(request.Summary)
        || request.AmountCents <= 0
        || !DateOnly.TryParseExact(request.IncurredOn, "yyyy-MM-dd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var incurredOn))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cost_entry", "Cost entry payload is invalid."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var amountCents = request.AmountCents.GetValueOrDefault();
      var cost = await CreateCostEntry(
        connection,
        transaction,
        tenantId.Value,
        tenantSlug,
        request.Category.Trim(),
        request.Summary.Trim(),
        amountCents,
        incurredOn,
        string.IsNullOrWhiteSpace(request.SalePublicId) ? null : request.SalePublicId.Trim());

      await transaction.CommitAsync();
      return TypedResults.Ok(cost);
    });

    app.MapGet("/api/finance/cash-accounts", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      var normalizedStatus = NormalizeCashAccountStatus(status);
      if (status is not null && normalizedStatus.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cash_account_status", "Cash account status is invalid."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      var accounts = await ListCashAccounts(connection, resolvedTenantSlug, normalizedStatus);
      return TypedResults.Ok<IReadOnlyList<CashAccountResponse>>(accounts);
    });

    app.MapPost("/api/finance/cash-accounts", async Task<IResult> (CreateCashAccountRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null
        || string.IsNullOrWhiteSpace(request.Code)
        || string.IsNullOrWhiteSpace(request.DisplayName)
        || string.IsNullOrWhiteSpace(request.CurrencyCode)
        || string.IsNullOrWhiteSpace(request.Provider)
        || request.OpeningBalanceCents < 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cash_account", "Cash account payload is invalid."));
      }

      var currencyCode = request.CurrencyCode.Trim().ToUpperInvariant();
      if (currencyCode.Length != 3)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cash_account_currency", "Currency code must have three uppercase characters."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var account = await CreateCashAccount(
        connection,
        transaction,
        tenantId.Value,
        tenantSlug,
        request.Code.Trim(),
        request.DisplayName.Trim(),
        currencyCode,
        request.Provider.Trim(),
        request.OpeningBalanceCents.GetValueOrDefault());

      await transaction.CommitAsync();
      return TypedResults.Ok(account);
    });

    app.MapPost("/api/finance/treasury/sync", async Task<IResult> (TreasurySyncRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.CashAccountPublicId))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_treasury_sync", "Cash account public id is required."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var account = await FindCashAccount(connection, transaction, tenantId.Value, request.CashAccountPublicId.Trim());
      if (account is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("cash_account_not_found", "Cash account was not found."));
      }

      var sync = await SyncTreasuryMovements(connection, transaction, tenantId.Value, account.Id, tenantSlug, account.PublicId);
      await transaction.CommitAsync();
      return TypedResults.Ok(sync);
    });

    app.MapGet("/api/finance/cash-movements", async Task<IResult> (string? tenantSlug, string? cashAccountPublicId, string? direction, string? movementType, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      var normalizedDirection = NormalizeCashMovementDirection(direction);
      if (direction is not null && normalizedDirection.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cash_movement_direction", "Cash movement direction is invalid."));
      }

      var normalizedMovementType = NormalizeCashMovementType(movementType);
      if (movementType is not null && normalizedMovementType.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_cash_movement_type", "Cash movement type is invalid."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      var movements = await ListCashMovements(connection, resolvedTenantSlug, string.IsNullOrWhiteSpace(cashAccountPublicId) ? null : cashAccountPublicId.Trim(), normalizedDirection, normalizedMovementType);
      return TypedResults.Ok<IReadOnlyList<CashMovementResponse>>(movements);
    });

    app.MapGet("/api/finance/cash-movements/summary", async Task<IResult> (string? tenantSlug, string? cashAccountPublicId, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var summary = await BuildCashMovementSummary(connection, resolvedTenantSlug, string.IsNullOrWhiteSpace(cashAccountPublicId) ? null : cashAccountPublicId.Trim());
      return TypedResults.Ok(summary);
    });

    app.MapGet("/api/finance/reports/treasury", async Task<IResult> (string? tenantSlug, string? cashAccountPublicId, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var report = await BuildTreasuryReport(connection, resolvedTenantSlug, string.IsNullOrWhiteSpace(cashAccountPublicId) ? null : cashAccountPublicId.Trim());
      return TypedResults.Ok(report);
    });

    app.MapPost("/api/finance/period-closures", async Task<IResult> (CreatePeriodClosureRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.PeriodKey))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_period_closure", "Period key is required."));
      }

      var periodKey = request.PeriodKey.Trim();
      if (!TryResolveCurrentPeriodKey(periodKey, out var currentPeriodKey) || !string.Equals(periodKey, currentPeriodKey, StringComparison.Ordinal))
      {
        return TypedResults.BadRequest(new ErrorResponse("unsupported_period_closure", "Only the current UTC period can be closed at this stage."));
      }

      var tenantSlug = ResolveTenantSlug(request.TenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        await transaction.RollbackAsync();
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found."));
      }

      var existingClosure = await FindPeriodClosure(connection, transaction, tenantId.Value, periodKey);
      if (existingClosure is not null)
      {
        await transaction.CommitAsync();
        return TypedResults.Ok(existingClosure with { AlreadyClosed = true });
      }

      var snapshot = await BuildOperationalReport(connection, tenantSlug);
      var created = await CreatePeriodClosure(connection, transaction, tenantId.Value, tenantSlug, periodKey, snapshot);
      await transaction.CommitAsync();
      return TypedResults.Ok(created);
    });

    app.MapGet("/api/finance/period-closures", async Task<IResult> (string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var closures = await ListPeriodClosures(connection, resolvedTenantSlug);
      return TypedResults.Ok<IReadOnlyList<PeriodClosureResponse>>(closures);
    });

    app.MapGet("/api/finance/period-closures/{periodKey}", async Task<IResult> (string periodKey, string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var closure = await GetPeriodClosureDetail(connection, resolvedTenantSlug, periodKey);
      return closure is null
        ? TypedResults.NotFound(new ErrorResponse("period_closure_not_found", "Period closure was not found."))
        : TypedResults.Ok(closure);
    });

    app.MapGet("/api/finance/reports/operations", async Task<IResult> (string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var report = await BuildOperationalReport(connection, resolvedTenantSlug);
      return TypedResults.Ok(report);
    });

    app.MapGet("/api/finance/activity", async Task<IResult> (string? tenantSlug, string? entityType, string? entityPublicId, string? activityType, NpgsqlDataSource dataSource, IConfiguration configuration) =>
    {
      var resolvedTenantSlug = ResolveTenantSlug(tenantSlug, configuration);
      await using var connection = await dataSource.OpenConnectionAsync();
      var events = await ListActivityEvents(connection, resolvedTenantSlug, entityType, entityPublicId, activityType);
      return TypedResults.Ok<IReadOnlyList<FinanceActivityResponse>>(events);
    });

    return app;
  }

  private static async Task<bool> CanReachDatabase(NpgsqlDataSource dataSource)
  {
    try
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var command = new NpgsqlCommand("SELECT 1", connection);
      await command.ExecuteScalarAsync();
      return true;
    }
    catch
    {
      return false;
    }
  }

  private static string ResolveTenantSlug(string? tenantSlug, IConfiguration configuration)
    => string.IsNullOrWhiteSpace(tenantSlug)
      ? configuration["FINANCE_BOOTSTRAP_TENANT_SLUG"] ?? "bootstrap-ops"
      : tenantSlug.Trim();

  private static async Task<long?> LookupTenantId(NpgsqlConnection connection, NpgsqlTransaction? transaction, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantSlug);

    var scalar = await command.ExecuteScalarAsync();
    return scalar switch
    {
      long value => value,
      int value => value,
      _ => null
    };
  }

  private static async Task<List<SalesOutboxEvent>> LoadPendingOutboxEvents(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, int limit)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT public_id::text, aggregate_type, aggregate_public_id::text, event_type, payload::text
      FROM sales.outbox_events
      WHERE tenant_id = $1
        AND status = 'pending'
      ORDER BY created_at, id
      LIMIT $2
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(limit);

    var response = new List<SalesOutboxEvent>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new SalesOutboxEvent(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4)));
    }

    return response;
  }

  private static async Task<ProjectionMutation> ApplyOutboxEvent(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, SalesOutboxEvent outboxEvent)
  {
    using var payload = JsonDocument.Parse(outboxEvent.Payload);
    var root = payload.RootElement;

    return outboxEvent.EventType switch
    {
      "sale.created" => await UpsertProjection(
        connection,
        transaction,
        tenantId,
        outboxEvent,
        projectionKind: "sale-booking",
        salePublicId: root.GetProperty("salePublicId").GetString() ?? outboxEvent.AggregatePublicId,
        invoicePublicId: null,
        status: "forecast",
        amountCents: root.GetProperty("amountCents").GetInt64(),
        dueDate: null),
      "invoice.created" => await UpsertProjection(
        connection,
        transaction,
        tenantId,
        outboxEvent,
        projectionKind: "invoice",
        salePublicId: root.GetProperty("salePublicId").GetString() ?? string.Empty,
        invoicePublicId: root.GetProperty("invoicePublicId").GetString() ?? outboxEvent.AggregatePublicId,
        status: "open",
        amountCents: root.GetProperty("amountCents").GetInt64(),
        dueDate: root.TryGetProperty("dueDate", out var dueDateElement) ? dueDateElement.GetString() : null),
      "invoice.status_changed" => await UpdateProjectionStatus(
        connection,
        transaction,
        tenantId,
        outboxEvent,
        outboxEvent.AggregatePublicId,
        NormalizeProjectionStatus(root.GetProperty("status").GetString())),
      "sale.renegotiated" => await UpdateSaleProjectionAmount(
        connection,
        transaction,
        tenantId,
        outboxEvent,
        root.GetProperty("salePublicId").GetString() ?? outboxEvent.AggregatePublicId,
        root.GetProperty("amountCents").GetInt64()),
      "sale.status_changed" => await UpdateSaleProjectionStatus(
        connection,
        transaction,
        tenantId,
        outboxEvent,
        root.GetProperty("salePublicId").GetString() ?? outboxEvent.AggregatePublicId,
        NormalizeProjectionStatus(root.GetProperty("status").GetString())),
      _ => await MarkOutboxProcessed(connection, transaction, tenantId, outboxEvent.PublicId)
    };
  }

  private static async Task<ProjectionMutation> UpsertProjection(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    SalesOutboxEvent outboxEvent,
    string projectionKind,
    string salePublicId,
    string? invoicePublicId,
    string status,
    long amountCents,
    string? dueDate)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.receivable_projections (
        tenant_id,
        public_id,
        source_event_public_id,
        projection_kind,
        sale_public_id,
        invoice_public_id,
        status,
        amount_cents,
        due_date,
        snapshot_payload
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2::uuid,
        $3,
        $4::uuid,
        NULLIF($5, '')::uuid,
        $6,
        $7,
        NULLIF($8, '')::date,
        $9::jsonb
      )
      ON CONFLICT (source_event_public_id)
      DO UPDATE SET
        projection_kind = EXCLUDED.projection_kind,
        sale_public_id = EXCLUDED.sale_public_id,
        invoice_public_id = EXCLUDED.invoice_public_id,
        status = EXCLUDED.status,
        amount_cents = EXCLUDED.amount_cents,
        due_date = EXCLUDED.due_date,
        snapshot_payload = EXCLUDED.snapshot_payload
      RETURNING xmax = 0
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(outboxEvent.PublicId);
    command.Parameters.AddWithValue(projectionKind);
    command.Parameters.AddWithValue(salePublicId);
    command.Parameters.AddWithValue(invoicePublicId ?? string.Empty);
    command.Parameters.AddWithValue(status);
    command.Parameters.AddWithValue(amountCents);
    command.Parameters.AddWithValue(dueDate ?? string.Empty);
    command.Parameters.AddWithValue(outboxEvent.Payload);

    var created = (bool)(await command.ExecuteScalarAsync() ?? false);
    var markProcessed = await MarkOutboxProcessed(connection, transaction, tenantId, outboxEvent.PublicId);
    return new ProjectionMutation(created ? 1 : 0, created ? 0 : 1, markProcessed.Processed);
  }

  private static async Task<ProjectionMutation> UpdateProjectionStatus(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    SalesOutboxEvent outboxEvent,
    string invoicePublicId,
    string status)
  {
    await using var command = new NpgsqlCommand(
      """
      UPDATE finance.receivable_projections
      SET
        status = $3,
        snapshot_payload = $4::jsonb
      WHERE tenant_id = $1
        AND invoice_public_id = $2::uuid
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(invoicePublicId);
    command.Parameters.AddWithValue(status);
    command.Parameters.AddWithValue(outboxEvent.Payload);

    var updated = await command.ExecuteNonQueryAsync();
    var markProcessed = await MarkOutboxProcessed(connection, transaction, tenantId, outboxEvent.PublicId);
    return new ProjectionMutation(0, updated > 0 ? 1 : 0, markProcessed.Processed);
  }

  private static async Task<ProjectionMutation> UpdateSaleProjectionStatus(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    SalesOutboxEvent outboxEvent,
    string salePublicId,
    string status)
  {
    await using var command = new NpgsqlCommand(
      """
      UPDATE finance.receivable_projections
      SET
        status = $3,
        snapshot_payload = $4::jsonb
      WHERE tenant_id = $1
        AND sale_public_id = $2::uuid
        AND projection_kind = 'sale-booking'
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(salePublicId);
    command.Parameters.AddWithValue(status);
    command.Parameters.AddWithValue(outboxEvent.Payload);

    var updated = await command.ExecuteNonQueryAsync();
    var markProcessed = await MarkOutboxProcessed(connection, transaction, tenantId, outboxEvent.PublicId);
    return new ProjectionMutation(0, updated > 0 ? 1 : 0, markProcessed.Processed);
  }

  private static async Task<ProjectionMutation> UpdateSaleProjectionAmount(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    SalesOutboxEvent outboxEvent,
    string salePublicId,
    long amountCents)
  {
    await using var command = new NpgsqlCommand(
      """
      UPDATE finance.receivable_projections
      SET
        amount_cents = $3,
        snapshot_payload = $4::jsonb
      WHERE tenant_id = $1
        AND sale_public_id = $2::uuid
        AND projection_kind = 'sale-booking'
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(salePublicId);
    command.Parameters.AddWithValue(amountCents);
    command.Parameters.AddWithValue(outboxEvent.Payload);

    var updated = await command.ExecuteNonQueryAsync();
    var markProcessed = await MarkOutboxProcessed(connection, transaction, tenantId, outboxEvent.PublicId);
    return new ProjectionMutation(0, updated > 0 ? 1 : 0, markProcessed.Processed);
  }

  private static async Task<ProjectionMutation> MarkOutboxProcessed(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string outboxPublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      UPDATE sales.outbox_events
      SET status = 'processed',
          processed_at = timezone('utc', now())
      WHERE tenant_id = $1
        AND public_id = $2::uuid
        AND status = 'pending'
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(outboxPublicId);

    var processed = await command.ExecuteNonQueryAsync();
    return new ProjectionMutation(0, 0, processed);
  }

  private static async Task<SyncMutation> SyncReceivableEntries(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      WITH source AS (
        SELECT
          invoice.public_id AS source_invoice_public_id,
          sale.public_id AS sale_public_id,
          sale.customer_public_id,
          NULL::uuid AS contract_public_id,
          'sales_invoice'::varchar AS source_kind,
          CASE
            WHEN invoice.status = 'paid' THEN 'paid'
            WHEN invoice.status = 'cancelled' THEN 'cancelled'
            ELSE 'open'
          END AS status,
          invoice.amount_cents,
          invoice.due_date,
          invoice.paid_at,
          jsonb_build_object(
            'invoicePublicId', invoice.public_id,
            'salePublicId', sale.public_id,
            'customerPublicId', sale.customer_public_id,
            'sourceKind', 'sales_invoice',
            'invoiceStatus', invoice.status,
            'amountCents', invoice.amount_cents,
            'dueDate', to_char(invoice.due_date, 'YYYY-MM-DD'),
            'paidAt', CASE WHEN invoice.paid_at IS NULL THEN NULL ELSE to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') END
          ) AS snapshot_payload
        FROM sales.invoices AS invoice
        INNER JOIN sales.sales AS sale
          ON sale.id = invoice.sale_id
        WHERE invoice.tenant_id = $1
      ),
      upserted AS (
        INSERT INTO finance.receivable_entries (
          tenant_id,
          public_id,
          source_invoice_public_id,
          sale_public_id,
          customer_public_id,
          contract_public_id,
          source_kind,
          status,
          amount_cents,
          due_date,
          paid_at,
          last_synced_at,
          snapshot_payload
        )
        SELECT
          $1,
          gen_random_uuid(),
          source.source_invoice_public_id,
          source.sale_public_id,
          source.customer_public_id,
          source.contract_public_id,
          source.source_kind,
          source.status,
          source.amount_cents,
          source.due_date,
          source.paid_at,
          timezone('utc', now()),
          source.snapshot_payload
        FROM source
        ON CONFLICT (source_invoice_public_id)
        DO UPDATE SET
          sale_public_id = EXCLUDED.sale_public_id,
          customer_public_id = EXCLUDED.customer_public_id,
          contract_public_id = EXCLUDED.contract_public_id,
          source_kind = EXCLUDED.source_kind,
          status = EXCLUDED.status,
          amount_cents = EXCLUDED.amount_cents,
          due_date = EXCLUDED.due_date,
          paid_at = EXCLUDED.paid_at,
          last_synced_at = EXCLUDED.last_synced_at,
          snapshot_payload = EXCLUDED.snapshot_payload
        RETURNING xmax = 0 AS inserted
      )
      SELECT
        COALESCE(sum(CASE WHEN inserted THEN 1 ELSE 0 END), 0) AS created_count,
        COALESCE(sum(CASE WHEN inserted THEN 0 ELSE 1 END), 0) AS updated_count
      FROM upserted
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return new SyncMutation(0, 0);
    }

    return new SyncMutation(ConvertToInt(reader.GetValue(0)), ConvertToInt(reader.GetValue(1)));
  }

  private static async Task<SyncMutation> SyncRentalReceivableEntries(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      WITH source AS (
        SELECT
          charge.public_id AS source_invoice_public_id,
          NULL::uuid AS sale_public_id,
          contract.customer_public_id,
          contract.public_id AS contract_public_id,
          'rental_charge'::varchar AS source_kind,
          CASE
            WHEN charge.status = 'paid' THEN 'paid'
            WHEN charge.status = 'cancelled' THEN 'cancelled'
            ELSE 'open'
          END AS status,
          charge.amount_cents,
          charge.due_date,
          charge.paid_at,
          jsonb_build_object(
            'chargePublicId', charge.public_id,
            'contractPublicId', contract.public_id,
            'customerPublicId', contract.customer_public_id,
            'sourceKind', 'rental_charge',
            'chargeStatus', charge.status,
            'amountCents', charge.amount_cents,
            'dueDate', to_char(charge.due_date, 'YYYY-MM-DD'),
            'paidAt', CASE WHEN charge.paid_at IS NULL THEN NULL ELSE to_char(charge.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"') END,
            'paymentReference', COALESCE(charge.payment_reference, '')
          ) AS snapshot_payload
        FROM rentals.contract_charges AS charge
        INNER JOIN rentals.contracts AS contract
          ON contract.id = charge.contract_id
        WHERE charge.tenant_id = $1
      ),
      upserted AS (
        INSERT INTO finance.receivable_entries (
          tenant_id,
          public_id,
          source_invoice_public_id,
          sale_public_id,
          customer_public_id,
          contract_public_id,
          source_kind,
          status,
          amount_cents,
          due_date,
          paid_at,
          last_synced_at,
          snapshot_payload
        )
        SELECT
          $1,
          gen_random_uuid(),
          source.source_invoice_public_id,
          source.sale_public_id,
          source.customer_public_id,
          source.contract_public_id,
          source.source_kind,
          source.status,
          source.amount_cents,
          source.due_date,
          source.paid_at,
          timezone('utc', now()),
          source.snapshot_payload
        FROM source
        ON CONFLICT (source_invoice_public_id)
        DO UPDATE SET
          sale_public_id = EXCLUDED.sale_public_id,
          customer_public_id = EXCLUDED.customer_public_id,
          contract_public_id = EXCLUDED.contract_public_id,
          source_kind = EXCLUDED.source_kind,
          status = EXCLUDED.status,
          amount_cents = EXCLUDED.amount_cents,
          due_date = EXCLUDED.due_date,
          paid_at = EXCLUDED.paid_at,
          last_synced_at = EXCLUDED.last_synced_at,
          snapshot_payload = EXCLUDED.snapshot_payload
        RETURNING xmax = 0 AS inserted
      )
      SELECT
        COALESCE(sum(CASE WHEN inserted THEN 1 ELSE 0 END), 0) AS created_count,
        COALESCE(sum(CASE WHEN inserted THEN 0 ELSE 1 END), 0) AS updated_count
      FROM upserted
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return new SyncMutation(0, 0);
    }

    return new SyncMutation(ConvertToInt(reader.GetValue(0)), ConvertToInt(reader.GetValue(1)));
  }

  private static async Task<SyncMutation> SyncRentalReceivableSettlements(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      WITH candidates AS (
        SELECT
          receivable.id AS receivable_entry_id,
          receivable.public_id::text AS receivable_public_id,
          charge.payment_reference,
          charge.amount_cents,
          charge.paid_at
        FROM finance.receivable_entries AS receivable
        INNER JOIN rentals.contract_charges AS charge
          ON charge.public_id = receivable.source_invoice_public_id
        WHERE receivable.tenant_id = $1
          AND receivable.source_kind = 'rental_charge'
          AND receivable.status = 'paid'
          AND charge.status = 'paid'
          AND charge.payment_reference IS NOT NULL
          AND charge.payment_reference <> ''
          AND charge.paid_at IS NOT NULL
          AND NOT EXISTS (
            SELECT 1
            FROM finance.receivable_settlements AS settlement
            WHERE settlement.receivable_entry_id = receivable.id
          )
      ),
      inserted AS (
        INSERT INTO finance.receivable_settlements (
          tenant_id,
          receivable_entry_id,
          public_id,
          settlement_reference,
          amount_cents,
          settled_at,
          snapshot_payload
        )
        SELECT
          $1,
          candidate.receivable_entry_id,
          gen_random_uuid(),
          candidate.payment_reference,
          candidate.amount_cents,
          candidate.paid_at,
          jsonb_build_object(
            'receivablePublicId', candidate.receivable_public_id,
            'settlementReference', candidate.payment_reference,
            'amountCents', candidate.amount_cents,
            'settledAt', to_char(candidate.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
            'sourceKind', 'rental_charge'
          )
        FROM candidates AS candidate
        RETURNING 1
      )
      SELECT COUNT(*) FROM inserted
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    var created = ConvertToInt(await command.ExecuteScalarAsync() ?? 0);
    return new SyncMutation(created, 0);
  }

  private static async Task<SyncMutation> SyncCommissionEntries(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      WITH source AS (
        SELECT
          commission.public_id AS source_commission_public_id,
          sale.public_id AS sale_public_id,
          commission.recipient_user_public_id,
          commission.role_code,
          commission.rate_bps,
          commission.amount_cents,
          commission.status,
          jsonb_build_object(
            'commissionPublicId', commission.public_id,
            'salePublicId', sale.public_id,
            'recipientUserPublicId', commission.recipient_user_public_id,
            'roleCode', commission.role_code,
            'rateBps', commission.rate_bps,
            'amountCents', commission.amount_cents,
            'status', commission.status
          ) AS snapshot_payload
        FROM sales.commissions AS commission
        INNER JOIN sales.sales AS sale
          ON sale.id = commission.sale_id
        WHERE commission.tenant_id = $1
      ),
      upserted AS (
        INSERT INTO finance.commission_entries (
          tenant_id,
          public_id,
          source_commission_public_id,
          sale_public_id,
          recipient_user_public_id,
          role_code,
          rate_bps,
          amount_cents,
          status,
          snapshot_payload
        )
        SELECT
          $1,
          gen_random_uuid(),
          source.source_commission_public_id,
          source.sale_public_id,
          source.recipient_user_public_id,
          source.role_code,
          source.rate_bps,
          source.amount_cents,
          source.status,
          source.snapshot_payload
        FROM source
        ON CONFLICT (source_commission_public_id)
        DO UPDATE SET
          sale_public_id = EXCLUDED.sale_public_id,
          recipient_user_public_id = EXCLUDED.recipient_user_public_id,
          role_code = EXCLUDED.role_code,
          rate_bps = EXCLUDED.rate_bps,
          amount_cents = EXCLUDED.amount_cents,
          status = EXCLUDED.status,
          snapshot_payload = EXCLUDED.snapshot_payload
        RETURNING xmax = 0 AS inserted
      )
      SELECT
        COALESCE(sum(CASE WHEN inserted THEN 1 ELSE 0 END), 0) AS created_count,
        COALESCE(sum(CASE WHEN inserted THEN 0 ELSE 1 END), 0) AS updated_count
      FROM upserted
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return new SyncMutation(0, 0);
    }

    return new SyncMutation(ConvertToInt(reader.GetValue(0)), ConvertToInt(reader.GetValue(1)));
  }

  private static async Task<IReadOnlyList<ProjectionResponse>> ListProjections(NpgsqlConnection connection, string tenantSlug, string? status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        projection.public_id::text,
        tenant.slug,
        projection.projection_kind,
        projection.sale_public_id::text,
        COALESCE(projection.invoice_public_id::text, ''),
        projection.status,
        projection.amount_cents,
        COALESCE(to_char(projection.due_date, 'YYYY-MM-DD'), ''),
        to_char(projection.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(projection.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.receivable_projections AS projection
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = projection.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR projection.status = $2)
      ORDER BY projection.created_at, projection.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(string.IsNullOrWhiteSpace(status) ? string.Empty : NormalizeProjectionStatus(status));

    var response = new List<ProjectionResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new ProjectionResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetInt64(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9)));
    }

    return response;
  }

  private static async Task<ProjectionSummaryResponse> BuildProjectionSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(projection.amount_cents) FILTER (WHERE projection.status IN ('forecast', 'open')), 0) AS pipeline_amount_cents,
        COALESCE(SUM(projection.amount_cents) FILTER (WHERE projection.status = 'paid'), 0) AS paid_amount_cents,
        COUNT(*) FILTER (WHERE projection.status = 'forecast') AS forecast_count,
        COUNT(*) FILTER (WHERE projection.status = 'open') AS open_count,
        COUNT(*) FILTER (WHERE projection.status = 'paid') AS paid_count,
        COUNT(*) FILTER (WHERE projection.status = 'cancelled') AS cancelled_count
      FROM finance.receivable_projections AS projection
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = projection.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return new ProjectionSummaryResponse(tenantSlug, 0, 0, 0, new ProjectionStatusCounts(0, 0, 0, 0));
    }

    return new ProjectionSummaryResponse(
      tenantSlug,
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      ConvertToInt64(reader.GetValue(2)),
      new ProjectionStatusCounts(
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4)),
        ConvertToInt(reader.GetValue(5)),
        ConvertToInt(reader.GetValue(6))));
  }

  private static async Task<IReadOnlyList<ReceivableResponse>> ListReceivables(NpgsqlConnection connection, string tenantSlug, string? status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        receivable.public_id::text,
        tenant.slug,
        receivable.source_kind,
        receivable.source_invoice_public_id::text,
        COALESCE(receivable.sale_public_id::text, ''),
        COALESCE(receivable.contract_public_id::text, ''),
        receivable.customer_public_id::text,
        receivable.status,
        receivable.amount_cents,
        to_char(receivable.due_date, 'YYYY-MM-DD'),
        COALESCE(to_char(receivable.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(settlement.settlement_reference, ''),
        to_char(receivable.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(receivable.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.receivable_entries AS receivable
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = receivable.tenant_id
      LEFT JOIN finance.receivable_settlements AS settlement
        ON settlement.receivable_entry_id = receivable.id
      WHERE tenant.slug = $1
        AND ($2 = '' OR receivable.status = $2)
      ORDER BY receivable.due_date, receivable.created_at, receivable.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(NormalizeOperationalStatus(status));

    var response = new List<ReceivableResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new ReceivableResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetInt64(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11),
        reader.GetString(12),
        reader.GetString(13)));
    }

    return response;
  }

  private static async Task<ReceivableLookup?> FindReceivableForSettlement(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string receivablePublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT id, public_id::text, status, amount_cents
      FROM finance.receivable_entries
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(receivablePublicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new ReceivableLookup(
      ConvertToInt64(reader.GetValue(0)),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt64(3));
  }

  private static async Task<ReceivableSettlementResponse?> FindReceivableSettlementByReference(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string settlementReference)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        settlement.public_id::text,
        receivable.public_id::text,
        tenant.slug,
        settlement.settlement_reference,
        settlement.amount_cents,
        to_char(settlement.settled_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.receivable_settlements AS settlement
      INNER JOIN finance.receivable_entries AS receivable
        ON receivable.id = settlement.receivable_entry_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = settlement.tenant_id
      WHERE settlement.tenant_id = $1
        AND settlement.settlement_reference = $2
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(settlementReference);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new ReceivableSettlementResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      false);
  }

  private static async Task<ReceivableSettlementResponse> CreateReceivableSettlement(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    ReceivableLookup receivable,
    string settlementReference,
    long amountCents,
    DateTime settledAt)
  {
    var payload = JsonSerializer.Serialize(new
    {
      receivablePublicId = receivable.PublicId,
      settlementReference,
      amountCents,
      settledAt = settledAt.ToUniversalTime().ToString("O", CultureInfo.InvariantCulture)
    });

    await using (var insertSettlement = new NpgsqlCommand(
      """
      INSERT INTO finance.receivable_settlements (
        tenant_id,
        receivable_entry_id,
        public_id,
        settlement_reference,
        amount_cents,
        settled_at,
        snapshot_payload
      )
      VALUES (
        $1,
        $2,
        gen_random_uuid(),
        $3,
        $4,
        $5,
        $6::jsonb
      )
      RETURNING public_id::text
      """,
      connection,
      transaction))
    {
      insertSettlement.Parameters.AddWithValue(tenantId);
      insertSettlement.Parameters.AddWithValue(receivable.Id);
      insertSettlement.Parameters.AddWithValue(settlementReference);
      insertSettlement.Parameters.AddWithValue(amountCents);
      insertSettlement.Parameters.AddWithValue(settledAt);
      insertSettlement.Parameters.AddWithValue(payload);

      var settlementPublicId = (string)(await insertSettlement.ExecuteScalarAsync() ?? string.Empty);

      await using var updateReceivable = new NpgsqlCommand(
        """
        UPDATE finance.receivable_entries
        SET status = 'paid',
            paid_at = $3,
            snapshot_payload = jsonb_set(snapshot_payload, '{settlementReference}', to_jsonb($4::text), true)
        WHERE tenant_id = $1
          AND id = $2
        """,
        connection,
        transaction);
      updateReceivable.Parameters.AddWithValue(tenantId);
      updateReceivable.Parameters.AddWithValue(receivable.Id);
      updateReceivable.Parameters.AddWithValue(settledAt);
      updateReceivable.Parameters.AddWithValue(settlementReference);
      await updateReceivable.ExecuteNonQueryAsync();

      var tenantSlug = await LookupTenantSlug(connection, transaction, tenantId) ?? string.Empty;
      var response = new ReceivableSettlementResponse(
        settlementPublicId,
        receivable.PublicId,
        tenantSlug,
        settlementReference,
        amountCents,
        settledAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture),
        false);

      await InsertActivityEvent(
        connection,
        transaction,
        tenantId,
        "receivable_settled",
        "finance.receivable",
        receivable.PublicId,
        $"Receivable settled with reference {settlementReference}.",
        "finance-receivables",
        payload);
      return response;
    }
  }

  private static async Task<IReadOnlyList<CommissionResponse>> ListCommissions(NpgsqlConnection connection, string tenantSlug, string? status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        commission.public_id::text,
        tenant.slug,
        commission.source_commission_public_id::text,
        commission.sale_public_id::text,
        commission.recipient_user_public_id::text,
        commission.role_code,
        commission.rate_bps,
        commission.amount_cents,
        commission.status,
        to_char(commission.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(commission.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.commission_entries AS commission
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = commission.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR commission.status = $2)
      ORDER BY commission.created_at, commission.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(NormalizeCommissionStatus(status));

    var response = new List<CommissionResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new CommissionResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetInt32(6),
        reader.GetInt64(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10)));
    }

    return response;
  }

  private static async Task<CommissionSummaryResponse> BuildCommissionSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(commission.amount_cents), 0) AS total_amount_cents,
        COUNT(*) FILTER (WHERE commission.status = 'pending') AS pending_count,
        COUNT(*) FILTER (WHERE commission.status = 'blocked') AS blocked_count,
        COUNT(*) FILTER (WHERE commission.status = 'released') AS released_count,
        COALESCE(SUM(commission.amount_cents) FILTER (WHERE commission.status = 'pending'), 0) AS pending_amount_cents,
        COALESCE(SUM(commission.amount_cents) FILTER (WHERE commission.status = 'blocked'), 0) AS blocked_amount_cents,
        COALESCE(SUM(commission.amount_cents) FILTER (WHERE commission.status = 'released'), 0) AS released_amount_cents
      FROM finance.commission_entries AS commission
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = commission.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return new CommissionSummaryResponse(tenantSlug, 0, 0, new CommissionStatusCounts(0, 0, 0), new CommissionAmountBuckets(0, 0, 0));
    }

    return new CommissionSummaryResponse(
      tenantSlug,
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      new CommissionStatusCounts(
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4))),
      new CommissionAmountBuckets(
        ConvertToInt64(reader.GetValue(5)),
        ConvertToInt64(reader.GetValue(6)),
        ConvertToInt64(reader.GetValue(7))));
  }

  private static async Task<InternalCommission?> FindCommissionForLifecycle(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        commission.id,
        commission.tenant_id,
        commission.public_id::text,
        tenant.slug,
        commission.sale_public_id::text,
        commission.status,
        commission.amount_cents
      FROM finance.commission_entries AS commission
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = commission.tenant_id
      WHERE tenant.slug = $1
        AND commission.public_id = $2::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new InternalCommission(
      ConvertToInt64(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetInt64(6));
  }

  private static async Task<CommissionResponse?> GetCommissionByPublicId(NpgsqlConnection connection, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        commission.public_id::text,
        tenant.slug,
        commission.source_commission_public_id::text,
        commission.sale_public_id::text,
        commission.recipient_user_public_id::text,
        commission.role_code,
        commission.rate_bps,
        commission.amount_cents,
        commission.status,
        to_char(commission.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(commission.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.commission_entries AS commission
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = commission.tenant_id
      WHERE tenant.slug = $1
        AND commission.public_id = $2::uuid
      LIMIT 1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new CommissionResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetInt32(6),
      reader.GetInt64(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10));
  }

  private static async Task<CommissionResponse> UpdateCommissionStatus(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalCommission commission,
    string nextStatus,
    string? actor,
    string? reason)
  {
    await using (var command = new NpgsqlCommand(
      """
      UPDATE finance.commission_entries
      SET status = $3
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(commission.TenantId);
      command.Parameters.AddWithValue(commission.PublicId);
      command.Parameters.AddWithValue(nextStatus);
      await command.ExecuteNonQueryAsync();
    }

    var payload = JsonSerializer.Serialize(new
    {
      commissionPublicId = commission.PublicId,
      salePublicId = commission.SalePublicId,
      previousStatus = commission.Status,
      status = nextStatus,
      reason = string.IsNullOrWhiteSpace(reason) ? string.Empty : reason.Trim()
    });

    await InsertActivityEvent(
      connection,
      transaction,
      commission.TenantId,
      nextStatus == "blocked" ? "commission_blocked" : "commission_released",
      "finance.commission",
      commission.PublicId,
      $"Commission status changed from {commission.Status} to {nextStatus}.",
      string.IsNullOrWhiteSpace(actor) ? "finance-commissions" : actor.Trim(),
      payload);

    return await GetCommissionByPublicId(connection, commission.TenantSlug, commission.PublicId)
      ?? throw new InvalidOperationException("commission_lookup_failed");
  }

  private static async Task<IReadOnlyList<PayableResponse>> ListPayables(NpgsqlConnection connection, string tenantSlug, string? status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        payable.public_id::text,
        tenant.slug,
        payable.category,
        payable.vendor_name,
        payable.description,
        payable.amount_cents,
        to_char(payable.due_date, 'YYYY-MM-DD'),
        payable.status,
        COALESCE(payable.payment_reference, ''),
        COALESCE(to_char(payable.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(payable.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(payable.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.payables AS payable
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = payable.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR payable.status = $2)
      ORDER BY payable.due_date, payable.created_at, payable.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(NormalizePayableStatus(status));

    var response = new List<PayableResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new PayableResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetInt64(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11)));
    }

    return response;
  }

  private static async Task<PayableResponse> CreatePayable(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string tenantSlug,
    string category,
    string vendorName,
    string description,
    long amountCents,
    DateOnly dueDate)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.payables (
        tenant_id,
        public_id,
        category,
        vendor_name,
        description,
        amount_cents,
        due_date
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        $3,
        $4,
        $5,
        $6
      )
      RETURNING
        public_id::text,
        category,
        vendor_name,
        description,
        amount_cents,
        to_char(due_date, 'YYYY-MM-DD'),
        status,
        COALESCE(payment_reference, ''),
        COALESCE(to_char(paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(category);
    command.Parameters.AddWithValue(vendorName);
    command.Parameters.AddWithValue(description);
    command.Parameters.AddWithValue(amountCents);
    command.Parameters.AddWithValue(dueDate);

    string responsePublicId;
    string responseCategory;
    string responseVendorName;
    string responseDescription;
    long responseAmountCents;
    string responseDueDate;
    string responseStatus;
    string responsePaymentReference;
    string responsePaidAt;
    string responseCreatedAt;
    string responseUpdatedAt;

    await using (var reader = await command.ExecuteReaderAsync())
    {
      await reader.ReadAsync();
      responsePublicId = reader.GetString(0);
      responseCategory = reader.GetString(1);
      responseVendorName = reader.GetString(2);
      responseDescription = reader.GetString(3);
      responseAmountCents = reader.GetInt64(4);
      responseDueDate = reader.GetString(5);
      responseStatus = reader.GetString(6);
      responsePaymentReference = reader.GetString(7);
      responsePaidAt = reader.GetString(8);
      responseCreatedAt = reader.GetString(9);
      responseUpdatedAt = reader.GetString(10);
    }

    var response = new PayableResponse(
      responsePublicId,
      tenantSlug,
      responseCategory,
      responseVendorName,
      responseDescription,
      responseAmountCents,
      responseDueDate,
      responseStatus,
      responsePaymentReference,
      responsePaidAt,
      responseCreatedAt,
      responseUpdatedAt);

    await InsertActivityEvent(
      connection,
      transaction,
      tenantId,
      "payable_created",
      "finance.payable",
      response.PublicId,
      $"Payable created for vendor {vendorName} in category {category}.",
      "finance-payables",
      JsonSerializer.Serialize(new { response.PublicId, category, vendorName, amountCents, dueDate = dueDate.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture) }));
    return response;
  }

  private static async Task<PayableResponse?> FindPayableForUpdate(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        payable.public_id::text,
        tenant.slug,
        payable.category,
        payable.vendor_name,
        payable.description,
        payable.amount_cents,
        to_char(payable.due_date, 'YYYY-MM-DD'),
        payable.status,
        COALESCE(payable.payment_reference, ''),
        COALESCE(to_char(payable.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(payable.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(payable.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.payables AS payable
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = payable.tenant_id
      WHERE payable.tenant_id = $1
        AND payable.public_id = $2::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new PayableResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetInt64(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10),
      reader.GetString(11));
  }

  private static async Task<PayableReferenceLookup?> FindPayableByPaymentReference(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string paymentReference)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT public_id::text, status
      FROM finance.payables
      WHERE tenant_id = $1
        AND payment_reference = $2
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(paymentReference);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new PayableReferenceLookup(reader.GetString(0), reader.GetString(1));
  }

  private static async Task<PayableResponse> UpdatePayableStatus(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string publicId,
    string status,
    string? paymentReference,
    DateTime? paidAt)
  {
    await using var command = new NpgsqlCommand(
      """
      UPDATE finance.payables
      SET
        status = $3,
        payment_reference = CASE WHEN $4 = '' THEN payment_reference ELSE $4 END,
        paid_at = CASE
          WHEN $3 = 'paid' THEN $5
          WHEN $3 = 'cancelled' THEN NULL
          ELSE paid_at
        END
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING
        public_id::text,
        category,
        vendor_name,
        description,
        amount_cents,
        to_char(due_date, 'YYYY-MM-DD'),
        status,
        COALESCE(payment_reference, ''),
        COALESCE(to_char(paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(publicId);
    command.Parameters.AddWithValue(status);
    command.Parameters.AddWithValue(paymentReference ?? string.Empty);
    command.Parameters.Add(new NpgsqlParameter
    {
      NpgsqlDbType = NpgsqlDbType.TimestampTz,
      Value = paidAt is null ? DBNull.Value : paidAt.Value
    });

    string responsePublicId;
    string category;
    string vendorName;
    string description;
    long amountCents;
    string dueDate;
    string responseStatus;
    string responsePaymentReference;
    string responsePaidAt;
    string createdAt;
    string updatedAt;

    await using (var reader = await command.ExecuteReaderAsync())
    {
      await reader.ReadAsync();
      responsePublicId = reader.GetString(0);
      category = reader.GetString(1);
      vendorName = reader.GetString(2);
      description = reader.GetString(3);
      amountCents = reader.GetInt64(4);
      dueDate = reader.GetString(5);
      responseStatus = reader.GetString(6);
      responsePaymentReference = reader.GetString(7);
      responsePaidAt = reader.GetString(8);
      createdAt = reader.GetString(9);
      updatedAt = reader.GetString(10);
    }

    var tenantSlug = await LookupTenantSlug(connection, transaction, tenantId) ?? string.Empty;
    var response = new PayableResponse(
      responsePublicId,
      tenantSlug,
      category,
      vendorName,
      description,
      amountCents,
      dueDate,
      responseStatus,
      responsePaymentReference,
      responsePaidAt,
      createdAt,
      updatedAt);

    await InsertActivityEvent(
      connection,
      transaction,
      tenantId,
      "payable_status_changed",
      "finance.payable",
      responsePublicId,
      $"Payable status changed to {responseStatus}.",
      "finance-payables",
      JsonSerializer.Serialize(new { responsePublicId, status = responseStatus, paymentReference = responsePaymentReference, paidAt = responsePaidAt }));
    return response;
  }

  private static async Task<IReadOnlyList<CostEntryResponse>> ListCosts(NpgsqlConnection connection, string tenantSlug, string? category)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        cost.public_id::text,
        tenant.slug,
        cost.category,
        cost.summary,
        cost.amount_cents,
        to_char(cost.incurred_on, 'YYYY-MM-DD'),
        COALESCE(cost.sale_public_id::text, ''),
        to_char(cost.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(cost.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.cost_entries AS cost
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = cost.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR cost.category = $2)
      ORDER BY cost.incurred_on DESC, cost.created_at DESC, cost.id DESC
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(string.IsNullOrWhiteSpace(category) ? string.Empty : category.Trim().ToLowerInvariant());

    var response = new List<CostEntryResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new CostEntryResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetInt64(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8)));
    }

    return response;
  }

  private static async Task<CostEntryResponse> CreateCostEntry(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string tenantSlug,
    string category,
    string summary,
    long amountCents,
    DateOnly incurredOn,
    string? salePublicId)
  {
    var payload = JsonSerializer.Serialize(new
    {
      category,
      summary,
      amountCents,
      incurredOn = incurredOn.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
      salePublicId
    });

    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.cost_entries (
        tenant_id,
        public_id,
        category,
        summary,
        amount_cents,
        incurred_on,
        sale_public_id,
        snapshot_payload
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        $3,
        $4,
        $5,
        NULLIF($6, '')::uuid,
        $7::jsonb
      )
      RETURNING
        public_id::text,
        category,
        summary,
        amount_cents,
        to_char(incurred_on, 'YYYY-MM-DD'),
        COALESCE(sale_public_id::text, ''),
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(category.ToLowerInvariant());
    command.Parameters.AddWithValue(summary);
    command.Parameters.AddWithValue(amountCents);
    command.Parameters.AddWithValue(incurredOn);
    command.Parameters.AddWithValue(salePublicId ?? string.Empty);
    command.Parameters.AddWithValue(payload);

    string responsePublicId;
    string responseCategory;
    string responseSummary;
    long responseAmountCents;
    string responseIncurredOn;
    string responseSalePublicId;
    string responseCreatedAt;
    string responseUpdatedAt;

    await using (var reader = await command.ExecuteReaderAsync())
    {
      await reader.ReadAsync();
      responsePublicId = reader.GetString(0);
      responseCategory = reader.GetString(1);
      responseSummary = reader.GetString(2);
      responseAmountCents = reader.GetInt64(3);
      responseIncurredOn = reader.GetString(4);
      responseSalePublicId = reader.GetString(5);
      responseCreatedAt = reader.GetString(6);
      responseUpdatedAt = reader.GetString(7);
    }

    var response = new CostEntryResponse(
      responsePublicId,
      tenantSlug,
      responseCategory,
      responseSummary,
      responseAmountCents,
      responseIncurredOn,
      responseSalePublicId,
      responseCreatedAt,
      responseUpdatedAt);

    await InsertActivityEvent(
      connection,
      transaction,
      tenantId,
      "cost_created",
      "finance.cost",
      response.PublicId,
      $"Cost entry created in category {category}.",
      "finance-costs",
      payload);
    return response;
  }

  private static async Task<IReadOnlyList<CashAccountResponse>> ListCashAccounts(NpgsqlConnection connection, string tenantSlug, string? status)
  {
    var sql =
      """
      SELECT
        account.public_id::text,
        tenant.slug,
        account.code,
        account.display_name,
        account.currency_code,
        account.provider,
        account.status,
        account.opening_balance_cents,
        to_char(account.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(account.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.cash_accounts AS account
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = account.tenant_id
      WHERE tenant.slug = $1
      """;

    if (!string.IsNullOrWhiteSpace(status))
    {
      sql += " AND account.status = $2";
    }

    sql += " ORDER BY account.created_at, account.id";

    await using var command = new NpgsqlCommand(sql, connection);
    command.Parameters.AddWithValue(tenantSlug);
    if (!string.IsNullOrWhiteSpace(status))
    {
      command.Parameters.AddWithValue(status);
    }

    var response = new List<CashAccountResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new CashAccountResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetInt64(7),
        reader.GetString(8),
        reader.GetString(9)));
    }

    return response;
  }

  private static async Task<CashAccountResponse> CreateCashAccount(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string tenantSlug,
    string code,
    string displayName,
    string currencyCode,
    string provider,
    long openingBalanceCents)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.cash_accounts (
        tenant_id,
        public_id,
        code,
        display_name,
        currency_code,
        provider,
        status,
        opening_balance_cents
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        $3,
        $4,
        $5,
        'active',
        $6
      )
      RETURNING
        public_id::text,
        code,
        display_name,
        currency_code,
        provider,
        status,
        opening_balance_cents,
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(code);
    command.Parameters.AddWithValue(displayName);
    command.Parameters.AddWithValue(currencyCode);
    command.Parameters.AddWithValue(provider);
    command.Parameters.AddWithValue(openingBalanceCents);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new CashAccountResponse(
      reader.GetString(0),
      tenantSlug,
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetInt64(6),
      reader.GetString(7),
      reader.GetString(8));
  }

  private static async Task<CashAccountLookup?> FindCashAccount(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        id,
        public_id::text,
        code,
        opening_balance_cents
      FROM finance.cash_accounts
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new CashAccountLookup(
      reader.GetInt64(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt64(3));
  }

  private static async Task<SyncTreasuryResponse> SyncTreasuryMovements(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    long cashAccountId,
    string tenantSlug,
    string cashAccountPublicId)
  {
    var receivableSync = await SyncTreasuryReceivableSettlements(connection, transaction, tenantId, cashAccountId);
    var payableSync = await SyncTreasuryPayables(connection, transaction, tenantId, cashAccountId);
    var costSync = await SyncTreasuryCosts(connection, transaction, tenantId, cashAccountId);
    var response = new SyncTreasuryResponse(
      tenantSlug,
      cashAccountPublicId,
      receivableSync.Created + payableSync.Created + costSync.Created,
      receivableSync.Skipped + payableSync.Skipped + costSync.Skipped);

    await InsertActivityEvent(
      connection,
      transaction,
      tenantId,
      "treasury_synced",
      "finance.cash_account",
      cashAccountPublicId,
      $"Treasury sync processed {response.CreatedMovements} movements and skipped {response.SkippedMovements}.",
      "finance-treasury",
      JsonSerializer.Serialize(new
      {
        response.CashAccountPublicId,
        response.CreatedMovements,
        response.SkippedMovements,
        receivable = new { receivableSync.Created, receivableSync.Skipped },
        payable = new { payableSync.Created, payableSync.Skipped },
        cost = new { costSync.Created, costSync.Skipped }
      }));
    return response;
  }

  private static async Task<TreasurySyncMutation> SyncTreasuryReceivableSettlements(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, long cashAccountId)
  {
    var candidates = await LoadReceivableSettlementCandidates(connection, transaction, tenantId);
    var created = 0;
    var skipped = 0;

    foreach (var candidate in candidates)
    {
      if (await TryInsertCashMovement(connection, transaction, tenantId, cashAccountId, "receivable_settlement", "inflow", candidate))
      {
        created += 1;
      }
      else
      {
        skipped += 1;
      }
    }

    return new TreasurySyncMutation(created, skipped);
  }

  private static async Task<TreasurySyncMutation> SyncTreasuryPayables(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, long cashAccountId)
  {
    var candidates = await LoadPaidPayableCandidates(connection, transaction, tenantId);
    var created = 0;
    var skipped = 0;

    foreach (var candidate in candidates)
    {
      if (await TryInsertCashMovement(connection, transaction, tenantId, cashAccountId, "payable_payment", "outflow", candidate))
      {
        created += 1;
      }
      else
      {
        skipped += 1;
      }
    }

    return new TreasurySyncMutation(created, skipped);
  }

  private static async Task<TreasurySyncMutation> SyncTreasuryCosts(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, long cashAccountId)
  {
    var candidates = await LoadCostEntryCandidates(connection, transaction, tenantId);
    var created = 0;
    var skipped = 0;

    foreach (var candidate in candidates)
    {
      if (await TryInsertCashMovement(connection, transaction, tenantId, cashAccountId, "cost_entry", "outflow", candidate))
      {
        created += 1;
      }
      else
      {
        skipped += 1;
      }
    }

    return new TreasurySyncMutation(created, skipped);
  }

  private static async Task<IReadOnlyList<TreasurySourceCandidate>> LoadReceivableSettlementCandidates(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        settlement.public_id::text,
        settlement.settlement_reference,
        settlement.amount_cents,
        settlement.settled_at,
        COALESCE(customer.name, customer.email, receivable.customer_public_id::text),
        'Recebimento liquidado a partir do financeiro operacional.'
      FROM finance.receivable_settlements AS settlement
      INNER JOIN finance.receivable_entries AS receivable
        ON receivable.id = settlement.receivable_entry_id
      LEFT JOIN crm.customers AS customer
        ON customer.public_id = receivable.customer_public_id
      WHERE settlement.tenant_id = $1
      ORDER BY settlement.settled_at, settlement.id
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    var response = new List<TreasurySourceCandidate>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new TreasurySourceCandidate(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetInt64(2),
        reader.GetFieldValue<DateTime>(3),
        reader.GetString(4),
        reader.GetString(5)));
    }

    return response;
  }

  private static async Task<IReadOnlyList<TreasurySourceCandidate>> LoadPaidPayableCandidates(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        payable.public_id::text,
        COALESCE(payable.payment_reference, payable.public_id::text),
        payable.amount_cents,
        COALESCE(payable.paid_at, payable.updated_at),
        payable.vendor_name,
        payable.description
      FROM finance.payables AS payable
      WHERE payable.tenant_id = $1
        AND payable.status = 'paid'
      ORDER BY COALESCE(payable.paid_at, payable.updated_at), payable.id
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    var response = new List<TreasurySourceCandidate>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new TreasurySourceCandidate(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetInt64(2),
        reader.GetFieldValue<DateTime>(3),
        reader.GetString(4),
        reader.GetString(5)));
    }

    return response;
  }

  private static async Task<IReadOnlyList<TreasurySourceCandidate>> LoadCostEntryCandidates(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        cost.public_id::text,
        cost.public_id::text,
        cost.amount_cents,
        timezone('utc', cost.incurred_on::timestamp),
        cost.category,
        cost.summary
      FROM finance.cost_entries AS cost
      WHERE cost.tenant_id = $1
      ORDER BY cost.incurred_on, cost.id
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);

    var response = new List<TreasurySourceCandidate>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new TreasurySourceCandidate(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetInt64(2),
        reader.GetFieldValue<DateTime>(3),
        reader.GetString(4),
        reader.GetString(5)));
    }

    return response;
  }

  private static async Task<bool> TryInsertCashMovement(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    long cashAccountId,
    string movementType,
    string direction,
    TreasurySourceCandidate candidate)
  {
    var payload = JsonSerializer.Serialize(new
    {
      movementType,
      direction,
      referenceCode = candidate.ReferenceCode,
      counterpartyName = candidate.CounterpartyName,
      description = candidate.Description
    });

    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.cash_movements (
        tenant_id,
        cash_account_id,
        public_id,
        movement_type,
        direction,
        source_public_id,
        reference_code,
        amount_cents,
        counterparty_name,
        description,
        effective_at,
        snapshot_payload
      )
      VALUES (
        $1,
        $2,
        gen_random_uuid(),
        $3,
        $4,
        $5::uuid,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11::jsonb
      )
      ON CONFLICT DO NOTHING
      RETURNING public_id::text
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(cashAccountId);
    command.Parameters.AddWithValue(movementType);
    command.Parameters.AddWithValue(direction);
    command.Parameters.AddWithValue(candidate.SourcePublicId);
    command.Parameters.AddWithValue(candidate.ReferenceCode);
    command.Parameters.AddWithValue(candidate.AmountCents);
    command.Parameters.AddWithValue(candidate.CounterpartyName);
    command.Parameters.AddWithValue(candidate.Description);
    command.Parameters.AddWithValue(candidate.EffectiveAt);
    command.Parameters.AddWithValue(payload);

    var inserted = await command.ExecuteScalarAsync();
    return inserted is not null;
  }

  private static async Task<IReadOnlyList<CashMovementResponse>> ListCashMovements(
    NpgsqlConnection connection,
    string tenantSlug,
    string? cashAccountPublicId,
    string? direction,
    string? movementType)
  {
    var sql =
      """
      SELECT
        movement.public_id::text,
        tenant.slug,
        account.public_id::text,
        account.code,
        movement.movement_type,
        movement.direction,
        COALESCE(movement.source_public_id::text, ''),
        movement.reference_code,
        movement.amount_cents,
        movement.counterparty_name,
        movement.description,
        to_char(movement.effective_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(movement.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(movement.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.cash_movements AS movement
      INNER JOIN finance.cash_accounts AS account
        ON account.id = movement.cash_account_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = movement.tenant_id
      WHERE tenant.slug = $1
      """;

    var parameters = new List<object?> { tenantSlug };
    if (!string.IsNullOrWhiteSpace(cashAccountPublicId))
    {
      parameters.Add(cashAccountPublicId);
      sql += $" AND account.public_id = ${parameters.Count}::uuid";
    }

    if (!string.IsNullOrWhiteSpace(direction))
    {
      parameters.Add(direction);
      sql += $" AND movement.direction = ${parameters.Count}";
    }

    if (!string.IsNullOrWhiteSpace(movementType))
    {
      parameters.Add(movementType);
      sql += $" AND movement.movement_type = ${parameters.Count}";
    }

    sql += " ORDER BY movement.effective_at DESC, movement.id DESC";

    await using var command = new NpgsqlCommand(sql, connection);
    foreach (var parameter in parameters)
    {
      command.Parameters.AddWithValue(parameter ?? string.Empty);
    }

    var response = new List<CashMovementResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new CashMovementResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetInt64(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11),
        reader.GetString(12),
        reader.GetString(13)));
    }

    return response;
  }

  private static async Task<CashMovementSummaryResponse> BuildCashMovementSummary(NpgsqlConnection connection, string tenantSlug, string? cashAccountPublicId)
  {
    var accountFilterSql = string.IsNullOrWhiteSpace(cashAccountPublicId) ? string.Empty : " AND account.public_id = $2::uuid";
    var movementFilterSql = string.IsNullOrWhiteSpace(cashAccountPublicId) ? string.Empty : " AND account.public_id = $2::uuid";

    var sql =
      $"""
      SELECT
        (
          SELECT COALESCE(sum(account.opening_balance_cents), 0)
          FROM finance.cash_accounts AS account
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = account.tenant_id
          WHERE tenant.slug = $1
          {accountFilterSql}
        ) AS opening_balance_cents,
        (
          SELECT count(*)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
          {movementFilterSql}
        ) AS total,
        (
          SELECT COALESCE(sum(movement.amount_cents), 0)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.direction = 'inflow'
          {movementFilterSql}
        ) AS inflow_cents,
        (
          SELECT COALESCE(sum(movement.amount_cents), 0)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.direction = 'outflow'
          {movementFilterSql}
        ) AS outflow_cents,
        (
          SELECT count(*)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.movement_type = 'receivable_settlement'
          {movementFilterSql}
        ) AS receivable_settlements,
        (
          SELECT count(*)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.movement_type = 'payable_payment'
          {movementFilterSql}
        ) AS payable_payments,
        (
          SELECT count(*)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.movement_type = 'cost_entry'
          {movementFilterSql}
        ) AS cost_entries,
        (
          SELECT count(*)
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
            AND movement.movement_type = 'manual_adjustment'
          {movementFilterSql}
        ) AS manual_adjustments,
        (
          SELECT COALESCE(to_char(max(movement.effective_at) AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '')
          FROM finance.cash_movements AS movement
          INNER JOIN finance.cash_accounts AS account
            ON account.id = movement.cash_account_id
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = movement.tenant_id
          WHERE tenant.slug = $1
          {movementFilterSql}
        ) AS latest_effective_at
      """;

    var parameters = new List<object?> { tenantSlug };
    if (!string.IsNullOrWhiteSpace(cashAccountPublicId))
    {
      parameters.Add(cashAccountPublicId);
    }

    await using var command = new NpgsqlCommand(sql, connection);
    foreach (var parameter in parameters)
    {
      command.Parameters.AddWithValue(parameter ?? string.Empty);
    }

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    var openingBalanceCents = ConvertToInt64(reader.GetValue(0));
    var inflowCents = ConvertToInt64(reader.GetValue(2));
    var outflowCents = ConvertToInt64(reader.GetValue(3));

    return new CashMovementSummaryResponse(
      tenantSlug,
      cashAccountPublicId ?? string.Empty,
      ConvertToInt(reader.GetValue(1)),
      inflowCents,
      outflowCents,
      openingBalanceCents,
      openingBalanceCents + inflowCents - outflowCents,
      new TreasuryMovementTypeCounts(
        ConvertToInt(reader.GetValue(4)),
        ConvertToInt(reader.GetValue(5)),
        ConvertToInt(reader.GetValue(6)),
        ConvertToInt(reader.GetValue(7))),
      reader.GetString(8));
  }

  private static async Task<TreasuryReportResponse> BuildTreasuryReport(NpgsqlConnection connection, string tenantSlug, string? cashAccountPublicId)
  {
    var accounts = await ListCashAccounts(connection, tenantSlug, null);
    var summary = await BuildCashMovementSummary(connection, tenantSlug, cashAccountPublicId);

    await using var command = new NpgsqlCommand(
      """
      SELECT
        (
          SELECT COALESCE(sum(receivable.amount_cents), 0)
          FROM finance.receivable_entries AS receivable
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = receivable.tenant_id
          WHERE tenant.slug = $1
            AND receivable.status = 'open'
        ) AS pending_receivables_cents,
        (
          SELECT COALESCE(sum(payable.amount_cents), 0)
          FROM finance.payables AS payable
          INNER JOIN identity.tenants AS tenant
            ON tenant.id = payable.tenant_id
          WHERE tenant.slug = $1
            AND payable.status = 'open'
        ) AS pending_payables_cents
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    long pendingReceivablesCents = 0;
    long pendingPayablesCents = 0;
    if (await reader.ReadAsync())
    {
      pendingReceivablesCents = ConvertToInt64(reader.GetValue(0));
      pendingPayablesCents = ConvertToInt64(reader.GetValue(1));
    }

    return new TreasuryReportResponse(
      tenantSlug,
      DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture),
      accounts,
      summary,
      new TreasuryLiquidityResponse(
        summary.CurrentBalanceCents,
        pendingReceivablesCents,
        pendingPayablesCents,
        summary.CurrentBalanceCents + pendingReceivablesCents - pendingPayablesCents));
  }

  private static async Task<IReadOnlyList<FinanceActivityResponse>> ListActivityEvents(
    NpgsqlConnection connection,
    string tenantSlug,
    string? entityType,
    string? entityPublicId,
    string? activityType)
  {
    var sql =
      """
      SELECT
        event.public_id::text,
        tenant.slug,
        event.activity_type,
        event.entity_type,
        COALESCE(event.entity_public_id::text, ''),
        event.summary,
        event.actor,
        event.payload::text,
        to_char(event.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.activity_events AS event
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = event.tenant_id
      WHERE tenant.slug = $1
      """;

    var parameters = new List<object?> { tenantSlug };
    if (!string.IsNullOrWhiteSpace(entityType))
    {
      sql += $" AND event.entity_type = ${parameters.Count + 1}";
      parameters.Add(entityType.Trim().ToLowerInvariant());
    }
    if (!string.IsNullOrWhiteSpace(entityPublicId))
    {
      sql += $" AND event.entity_public_id = ${parameters.Count + 1}::uuid";
      parameters.Add(entityPublicId.Trim());
    }
    if (!string.IsNullOrWhiteSpace(activityType))
    {
      sql += $" AND event.activity_type = ${parameters.Count + 1}";
      parameters.Add(activityType.Trim().ToLowerInvariant());
    }
    sql += " ORDER BY event.created_at DESC, event.id DESC";

    await using var command = new NpgsqlCommand(sql, connection);
    for (var index = 0; index < parameters.Count; index++)
    {
      command.Parameters.AddWithValue(parameters[index] ?? string.Empty);
    }

    var response = new List<FinanceActivityResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      using var payloadDocument = JsonDocument.Parse(reader.GetString(7));
      response.Add(new FinanceActivityResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        payloadDocument.RootElement.Clone(),
        reader.GetString(8)));
    }

    return response;
  }

  private static async Task InsertActivityEvent(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string activityType,
    string entityType,
    string? entityPublicId,
    string summary,
    string actor,
    string payload)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.activity_events (
        tenant_id,
        public_id,
        activity_type,
        entity_type,
        entity_public_id,
        summary,
        actor,
        payload
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        $3,
        NULLIF($4, '')::uuid,
        $5,
        $6,
        $7::jsonb
      )
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(activityType.Trim().ToLowerInvariant());
    command.Parameters.AddWithValue(entityType.Trim().ToLowerInvariant());
    command.Parameters.AddWithValue(entityPublicId ?? string.Empty);
    command.Parameters.AddWithValue(summary);
    command.Parameters.AddWithValue(string.IsNullOrWhiteSpace(actor) ? "system" : actor.Trim());
    command.Parameters.AddWithValue(string.IsNullOrWhiteSpace(payload) ? "{}" : payload);
    await command.ExecuteNonQueryAsync();
  }

  private static async Task<OperationalReportResponse> BuildOperationalReport(NpgsqlConnection connection, string tenantSlug)
  {
    var receivables = await BuildReceivableOperationsSummary(connection, tenantSlug);
    var payables = await BuildPayableOperationsSummary(connection, tenantSlug);
    var costs = await BuildCostOperationsSummary(connection, tenantSlug);
    var commissions = await BuildCommissionSummary(connection, tenantSlug);
    var closures = await BuildClosureOperationsSummary(connection, tenantSlug);

    return new OperationalReportResponse(
      tenantSlug,
      DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture),
      receivables,
      payables,
      costs,
      commissions,
      closures);
  }

  private static async Task<ReceivableOperationsSummary> BuildReceivableOperationsSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(receivable.amount_cents), 0) AS total_amount_cents,
        COUNT(*) FILTER (WHERE receivable.status = 'open') AS open_count,
        COUNT(*) FILTER (WHERE receivable.status = 'paid') AS paid_count,
        COUNT(*) FILTER (WHERE receivable.status = 'cancelled') AS cancelled_count,
        COALESCE(SUM(receivable.amount_cents) FILTER (WHERE receivable.status = 'open'), 0) AS open_amount_cents,
        COALESCE(SUM(receivable.amount_cents) FILTER (WHERE receivable.status = 'paid'), 0) AS paid_amount_cents,
        COALESCE(SUM(receivable.amount_cents) FILTER (WHERE receivable.status = 'cancelled'), 0) AS cancelled_amount_cents,
        COUNT(*) FILTER (
          WHERE receivable.status = 'open'
            AND receivable.due_date < timezone('utc', now())::date
        ) AS overdue_count,
        COALESCE(SUM(receivable.amount_cents) FILTER (
          WHERE receivable.status = 'open'
            AND receivable.due_date < timezone('utc', now())::date
        ), 0) AS overdue_amount_cents
      FROM finance.receivable_entries AS receivable
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = receivable.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    return new ReceivableOperationsSummary(
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      new AmountStatusCounts(
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4))),
      new AmountStatusBuckets(
        ConvertToInt64(reader.GetValue(5)),
        ConvertToInt64(reader.GetValue(6)),
        ConvertToInt64(reader.GetValue(7))),
      ConvertToInt(reader.GetValue(8)),
      ConvertToInt64(reader.GetValue(9)));
  }

  private static async Task<PayableOperationsSummary> BuildPayableOperationsSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(payable.amount_cents), 0) AS total_amount_cents,
        COUNT(*) FILTER (WHERE payable.status = 'open') AS open_count,
        COUNT(*) FILTER (WHERE payable.status = 'paid') AS paid_count,
        COUNT(*) FILTER (WHERE payable.status = 'cancelled') AS cancelled_count,
        COALESCE(SUM(payable.amount_cents) FILTER (WHERE payable.status = 'open'), 0) AS open_amount_cents,
        COALESCE(SUM(payable.amount_cents) FILTER (WHERE payable.status = 'paid'), 0) AS paid_amount_cents,
        COALESCE(SUM(payable.amount_cents) FILTER (WHERE payable.status = 'cancelled'), 0) AS cancelled_amount_cents
      FROM finance.payables AS payable
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = payable.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    return new PayableOperationsSummary(
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      new AmountStatusCounts(
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4))),
      new AmountStatusBuckets(
        ConvertToInt64(reader.GetValue(5)),
        ConvertToInt64(reader.GetValue(6)),
        ConvertToInt64(reader.GetValue(7))));
  }

  private static async Task<CostOperationsSummary> BuildCostOperationsSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(cost.amount_cents), 0) AS total_amount_cents,
        COALESCE(to_char(MAX(cost.incurred_on), 'YYYY-MM-DD'), '') AS latest_incurred_on
      FROM finance.cost_entries AS cost
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = cost.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    return new CostOperationsSummary(
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      reader.GetString(2));
  }

  private static async Task<ClosureOperationsSummary> BuildClosureOperationsSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(MAX(period_key), '') AS latest_period_key
      FROM finance.period_closures AS closure
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = closure.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    return new ClosureOperationsSummary(
      ConvertToInt(reader.GetValue(0)),
      reader.GetString(1));
  }

  private static async Task<PeriodClosureResponse?> FindPeriodClosure(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string periodKey)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        closure.public_id::text,
        tenant.slug,
        closure.period_key,
        to_char(closure.closed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.period_closures AS closure
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = closure.tenant_id
      WHERE closure.tenant_id = $1
        AND closure.period_key = $2
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(periodKey);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new PeriodClosureResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      false);
  }

  private static async Task<PeriodClosureResponse> CreatePeriodClosure(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string tenantSlug,
    string periodKey,
    OperationalReportResponse snapshot)
  {
    var snapshotJson = JsonSerializer.Serialize(snapshot);
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO finance.period_closures (
        tenant_id,
        public_id,
        period_key,
        closed_at,
        snapshot_payload
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        timezone('utc', now()),
        $3::jsonb
      )
      RETURNING
        public_id::text,
        period_key,
        to_char(closed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(periodKey);
    command.Parameters.AddWithValue(snapshotJson);

    string responsePublicId;
    string responsePeriodKey;
    string responseClosedAt;

    await using (var reader = await command.ExecuteReaderAsync())
    {
      await reader.ReadAsync();
      responsePublicId = reader.GetString(0);
      responsePeriodKey = reader.GetString(1);
      responseClosedAt = reader.GetString(2);
    }

    var response = new PeriodClosureResponse(
      responsePublicId,
      tenantSlug,
      responsePeriodKey,
      responseClosedAt,
      false);

    await InsertActivityEvent(
      connection,
      transaction,
      tenantId,
      "period_closed",
      "finance.period_closure",
      response.PublicId,
      $"Financial period {periodKey} was closed.",
      "finance-closure",
      snapshotJson);
    return response;
  }

  private static async Task<IReadOnlyList<PeriodClosureResponse>> ListPeriodClosures(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        closure.public_id::text,
        tenant.slug,
        closure.period_key,
        to_char(closure.closed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM finance.period_closures AS closure
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = closure.tenant_id
      WHERE tenant.slug = $1
      ORDER BY closure.period_key DESC, closure.created_at DESC
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    var response = new List<PeriodClosureResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new PeriodClosureResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        false));
    }

    return response;
  }

  private static async Task<PeriodClosureDetailResponse?> GetPeriodClosureDetail(NpgsqlConnection connection, string tenantSlug, string periodKey)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        closure.public_id::text,
        tenant.slug,
        closure.period_key,
        to_char(closure.closed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        closure.snapshot_payload::text
      FROM finance.period_closures AS closure
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = closure.tenant_id
      WHERE tenant.slug = $1
        AND closure.period_key = $2
      LIMIT 1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(periodKey);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    using var document = JsonDocument.Parse(reader.GetString(4));
    return new PeriodClosureDetailResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      document.RootElement.Clone());
  }

  private static async Task<string?> LookupTenantSlug(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT slug
      FROM identity.tenants
      WHERE id = $1
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    return await command.ExecuteScalarAsync() as string;
  }

  private static string NormalizeProjectionStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "active" => "forecast",
      "invoiced" => "open",
      "" => "open",
      _ => normalized
    };
  }

  private static string NormalizeOperationalStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "sent" => "open",
      "draft" => "open",
      "overdue" => "open",
      _ => normalized
    };
  }

  private static string NormalizeCommissionStatus(string? status)
    => (status ?? string.Empty).Trim().ToLowerInvariant();

  private static string NormalizePayableStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "open" or "paid" or "cancelled" => normalized,
      _ => string.Empty
    };
  }

  private static string NormalizeCashAccountStatus(string? status)
  {
    var normalized = (status ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "active" or "inactive" => normalized,
      _ => string.Empty
    };
  }

  private static string NormalizeCashMovementDirection(string? direction)
  {
    var normalized = (direction ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "inflow" or "outflow" => normalized,
      _ => string.Empty
    };
  }

  private static string NormalizeCashMovementType(string? movementType)
  {
    var normalized = (movementType ?? string.Empty).Trim().ToLowerInvariant();
    return normalized switch
    {
      "" => string.Empty,
      "receivable_settlement" or "payable_payment" or "cost_entry" or "manual_adjustment" => normalized,
      _ => string.Empty
    };
  }

  private static bool CanTransitionPayableStatus(string currentStatus, string targetStatus)
  {
    return currentStatus switch
    {
      "open" => targetStatus is "paid" or "cancelled",
      "paid" => false,
      "cancelled" => false,
      _ => false
    };
  }

  private static DateTime? ParseUtcOrNow(string? value)
  {
    if (string.IsNullOrWhiteSpace(value))
    {
      return DateTime.UtcNow;
    }

    if (!DateTimeOffset.TryParse(value, CultureInfo.InvariantCulture, DateTimeStyles.AssumeUniversal, out var parsed))
    {
      return null;
    }

    return parsed.UtcDateTime;
  }

  private static bool TryResolveCurrentPeriodKey(string periodKey, out string currentPeriodKey)
  {
    currentPeriodKey = DateTime.UtcNow.ToString("yyyy-MM", CultureInfo.InvariantCulture);
    return periodKey.Length == 7
      && periodKey[4] == '-'
      && int.TryParse(periodKey[..4], out _)
      && int.TryParse(periodKey[5..], out var month)
      && month is >= 1 and <= 12;
  }

  private static int ConvertToInt(object value)
    => value switch
    {
      int intValue => intValue,
      long longValue => (int)longValue,
      decimal decimalValue => (int)decimalValue,
      _ => Convert.ToInt32(value, CultureInfo.InvariantCulture)
    };

  private static long ConvertToInt64(object value)
    => value switch
    {
      long longValue => longValue,
      int intValue => intValue,
      decimal decimalValue => (long)decimalValue,
      _ => Convert.ToInt64(value, CultureInfo.InvariantCulture)
    };
}

public sealed record HealthResponse(string Service, string Status);
public sealed record DependencyResponse(string Name, string Status);
public sealed record ReadinessResponse(string Service, string Status, IReadOnlyList<DependencyResponse> Dependencies);
public sealed record ErrorResponse(string Code, string Message);
public sealed record ProjectionIngestRequest(string? TenantSlug, int? Limit);
public sealed record ProjectionIngestResponse(string TenantSlug, int LoadedEvents, int ProcessedEvents, int CreatedProjections, int UpdatedProjections);
public sealed record ProjectionResponse(
  string PublicId,
  string TenantSlug,
  string ProjectionKind,
  string SalePublicId,
  string InvoicePublicId,
  string Status,
  long AmountCents,
  string DueDate,
  string CreatedAt,
  string UpdatedAt);
public sealed record ProjectionStatusCounts(int Forecast, int Open, int Paid, int Cancelled);
public sealed record ProjectionSummaryResponse(string TenantSlug, int Total, long PipelineAmountCents, long PaidAmountCents, ProjectionStatusCounts Status);
public sealed record OperationsSyncRequest(string? TenantSlug);
public sealed record OperationsSyncResponse(string TenantSlug, int CreatedReceivables, int UpdatedReceivables, int CreatedCommissions, int UpdatedCommissions);
public sealed record ReceivableResponse(
  string PublicId,
  string TenantSlug,
  string SourceKind,
  string SourceInvoicePublicId,
  string SalePublicId,
  string ContractPublicId,
  string CustomerPublicId,
  string Status,
  long AmountCents,
  string DueDate,
  string PaidAt,
  string SettlementReference,
  string CreatedAt,
  string UpdatedAt);
public sealed record ReceivableSettlementRequest(string? TenantSlug, string? SettlementReference, long? AmountCents, string? SettledAt);
public sealed record ReceivableSettlementResponse(
  string PublicId,
  string ReceivablePublicId,
  string TenantSlug,
  string SettlementReference,
  long AmountCents,
  string SettledAt,
  bool Idempotent);
public sealed record CommissionResponse(
  string PublicId,
  string TenantSlug,
  string SourceCommissionPublicId,
  string SalePublicId,
  string RecipientUserPublicId,
  string RoleCode,
  int RateBps,
  long AmountCents,
  string Status,
  string CreatedAt,
  string UpdatedAt);
public sealed record CommissionStatusCounts(int Pending, int Blocked, int Released);
public sealed record CommissionAmountBuckets(long PendingAmountCents, long BlockedAmountCents, long ReleasedAmountCents);
public sealed record CommissionSummaryResponse(string TenantSlug, int Total, long TotalAmountCents, CommissionStatusCounts Status, CommissionAmountBuckets Amounts);
public sealed record CommissionLifecycleRequest(string? TenantSlug, string? Actor, string? Reason);
public sealed record PayableResponse(
  string PublicId,
  string TenantSlug,
  string Category,
  string VendorName,
  string Description,
  long AmountCents,
  string DueDate,
  string Status,
  string PaymentReference,
  string PaidAt,
  string CreatedAt,
  string UpdatedAt);
public sealed record CreatePayableRequest(string? TenantSlug, string? Category, string? VendorName, string? Description, long? AmountCents, string? DueDate);
public sealed record UpdatePayableStatusRequest(string? TenantSlug, string? Status, string? PaymentReference, string? PaidAt);
public sealed record CostEntryResponse(
  string PublicId,
  string TenantSlug,
  string Category,
  string Summary,
  long AmountCents,
  string IncurredOn,
  string SalePublicId,
  string CreatedAt,
  string UpdatedAt);
public sealed record CreateCostEntryRequest(string? TenantSlug, string? Category, string? Summary, long? AmountCents, string? IncurredOn, string? SalePublicId);
public sealed record CashAccountResponse(
  string PublicId,
  string TenantSlug,
  string Code,
  string DisplayName,
  string CurrencyCode,
  string Provider,
  string Status,
  long OpeningBalanceCents,
  string CreatedAt,
  string UpdatedAt);
public sealed record CreateCashAccountRequest(string? TenantSlug, string? Code, string? DisplayName, string? CurrencyCode, string? Provider, long? OpeningBalanceCents);
public sealed record TreasurySyncRequest(string? TenantSlug, string? CashAccountPublicId);
public sealed record SyncTreasuryResponse(string TenantSlug, string CashAccountPublicId, int CreatedMovements, int SkippedMovements);
public sealed record CashMovementResponse(
  string PublicId,
  string TenantSlug,
  string CashAccountPublicId,
  string CashAccountCode,
  string MovementType,
  string Direction,
  string SourcePublicId,
  string ReferenceCode,
  long AmountCents,
  string CounterpartyName,
  string Description,
  string EffectiveAt,
  string CreatedAt,
  string UpdatedAt);
public sealed record TreasuryMovementTypeCounts(int ReceivableSettlements, int PayablePayments, int CostEntries, int ManualAdjustments);
public sealed record CashMovementSummaryResponse(
  string TenantSlug,
  string CashAccountPublicId,
  int Total,
  long InflowCents,
  long OutflowCents,
  long OpeningBalanceCents,
  long CurrentBalanceCents,
  TreasuryMovementTypeCounts ByType,
  string LatestEffectiveAt);
public sealed record TreasuryLiquidityResponse(long CurrentBalanceCents, long PendingReceivablesCents, long PendingPayablesCents, long ProjectedNetPositionCents);
public sealed record TreasuryReportResponse(
  string TenantSlug,
  string GeneratedAt,
  IReadOnlyList<CashAccountResponse> Accounts,
  CashMovementSummaryResponse Summary,
  TreasuryLiquidityResponse Liquidity);
public sealed record AmountStatusCounts(int Open, int Paid, int Cancelled);
public sealed record AmountStatusBuckets(long OpenAmountCents, long PaidAmountCents, long CancelledAmountCents);
public sealed record ReceivableOperationsSummary(int Total, long TotalAmountCents, AmountStatusCounts Status, AmountStatusBuckets Amounts, int OverdueCount, long OverdueAmountCents);
public sealed record PayableOperationsSummary(int Total, long TotalAmountCents, AmountStatusCounts Status, AmountStatusBuckets Amounts);
public sealed record CostOperationsSummary(int Total, long TotalAmountCents, string LatestIncurredOn);
public sealed record ClosureOperationsSummary(int Total, string LatestPeriodKey);
public sealed record OperationalReportResponse(
  string TenantSlug,
  string GeneratedAt,
  ReceivableOperationsSummary Receivables,
  PayableOperationsSummary Payables,
  CostOperationsSummary Costs,
  CommissionSummaryResponse Commissions,
  ClosureOperationsSummary Closures);
public sealed record FinanceActivityResponse(
  string PublicId,
  string TenantSlug,
  string ActivityType,
  string EntityType,
  string EntityPublicId,
  string Summary,
  string Actor,
  JsonElement Payload,
  string CreatedAt);
public sealed record CreatePeriodClosureRequest(string? TenantSlug, string? PeriodKey);
public sealed record PeriodClosureResponse(string PublicId, string TenantSlug, string PeriodKey, string ClosedAt, bool AlreadyClosed);
public sealed record PeriodClosureDetailResponse(string PublicId, string TenantSlug, string PeriodKey, string ClosedAt, JsonElement Snapshot);
internal sealed record SalesOutboxEvent(string PublicId, string AggregateType, string AggregatePublicId, string EventType, string Payload);
internal sealed record ProjectionMutation(int Created, int Updated, int Processed);
internal sealed record SyncMutation(int Created, int Updated);
internal sealed record TreasurySyncMutation(int Created, int Skipped);
internal sealed record CashAccountLookup(long Id, string PublicId, string Code, long OpeningBalanceCents);
internal sealed record TreasurySourceCandidate(string SourcePublicId, string ReferenceCode, long AmountCents, DateTime EffectiveAt, string CounterpartyName, string Description);
internal sealed record ReceivableLookup(long Id, string PublicId, string Status, long AmountCents);
internal sealed record InternalCommission(long Id, long TenantId, string PublicId, string TenantSlug, string SalePublicId, string Status, long AmountCents);
internal sealed record PayableReferenceLookup(string PublicId, string Status);
