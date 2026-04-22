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

      var receivables = await SyncReceivableEntries(connection, transaction, tenantId.Value);
      var commissions = await SyncCommissionEntries(connection, transaction, tenantId.Value);

      await transaction.CommitAsync();
      return TypedResults.Ok(new OperationsSyncResponse(
        tenantSlug,
        receivables.Created,
        receivables.Updated,
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
        receivable.source_invoice_public_id::text,
        receivable.sale_public_id::text,
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
        reader.GetInt64(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11)));
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
      return new ReceivableSettlementResponse(
        settlementPublicId,
        receivable.PublicId,
        tenantSlug,
        settlementReference,
        amountCents,
        settledAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture),
        false);
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

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new PayableResponse(
      reader.GetString(0),
      tenantSlug,
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10));
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
    return new PayableResponse(
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

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new CostEntryResponse(
      reader.GetString(0),
      tenantSlug,
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt64(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7));
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

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();

    return new PeriodClosureResponse(
      reader.GetString(0),
      tenantSlug,
      reader.GetString(1),
      reader.GetString(2),
      false);
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
  string SourceInvoicePublicId,
  string SalePublicId,
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
public sealed record CreatePeriodClosureRequest(string? TenantSlug, string? PeriodKey);
public sealed record PeriodClosureResponse(string PublicId, string TenantSlug, string PeriodKey, string ClosedAt, bool AlreadyClosed);
public sealed record PeriodClosureDetailResponse(string PublicId, string TenantSlug, string PeriodKey, string ClosedAt, JsonElement Snapshot);
internal sealed record SalesOutboxEvent(string PublicId, string AggregateType, string AggregatePublicId, string EventType, string Payload);
internal sealed record ProjectionMutation(int Created, int Updated, int Processed);
internal sealed record SyncMutation(int Created, int Updated);
internal sealed record ReceivableLookup(long Id, string PublicId, string Status, long AmountCents);
internal sealed record PayableReferenceLookup(string PublicId, string Status);
