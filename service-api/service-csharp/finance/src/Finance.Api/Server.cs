// Este arquivo concentra a superficie HTTP inicial do finance.
// O servico nasce pequeno, mas ja com idempotencia e leitura operacional.
using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Routing;
using Npgsql;

namespace Finance.Api;

public static class Server
{
  public static IEndpointRouteBuilder MapFinanceRoutes(this IEndpointRouteBuilder app)
  {
    app.MapGet("/health/live", () => TypedResults.Ok(new HealthResponse("finance", "live")));
    app.MapGet("/health/ready", async Task<Results<Ok<HealthResponse>, StatusCodeHttpResult>> (NpgsqlDataSource dataSource) =>
    {
      return await CanReachDatabase(dataSource)
        ? TypedResults.Ok(new HealthResponse("finance", "ready"))
        : TypedResults.StatusCode(StatusCodes.Status503ServiceUnavailable);
    });
    app.MapGet("/health/details", async Task<Ok<ReadinessResponse>> (NpgsqlDataSource dataSource) =>
    {
      return TypedResults.Ok(new ReadinessResponse(
        "finance",
        await CanReachDatabase(dataSource) ? "ready" : "degraded",
        [new DependencyResponse("postgresql", await CanReachDatabase(dataSource) ? "ready" : "degraded")]));
    });
    app.MapPost(
      "/api/finance/projections/ingest",
      async Task<Ok<ProjectionIngestResponse>> (ProjectionIngestRequest? request, NpgsqlDataSource dataSource, IConfiguration configuration) =>
      {
        var tenantSlug = string.IsNullOrWhiteSpace(request?.TenantSlug)
          ? configuration["FINANCE_BOOTSTRAP_TENANT_SLUG"] ?? "bootstrap-ops"
          : request!.TenantSlug.Trim();
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
    app.MapGet(
      "/api/finance/projections",
      async Task<Ok<IReadOnlyList<ProjectionResponse>>> (string? tenantSlug, string? status, NpgsqlDataSource dataSource, IConfiguration configuration) =>
      {
        var resolvedTenantSlug = string.IsNullOrWhiteSpace(tenantSlug)
          ? configuration["FINANCE_BOOTSTRAP_TENANT_SLUG"] ?? "bootstrap-ops"
          : tenantSlug.Trim();

        await using var connection = await dataSource.OpenConnectionAsync();
        var projections = await ListProjections(connection, resolvedTenantSlug, status);
        return TypedResults.Ok<IReadOnlyList<ProjectionResponse>>(projections);
      });
    app.MapGet(
      "/api/finance/projections/summary",
      async Task<Ok<ProjectionSummaryResponse>> (string? tenantSlug, NpgsqlDataSource dataSource, IConfiguration configuration) =>
      {
        var resolvedTenantSlug = string.IsNullOrWhiteSpace(tenantSlug)
          ? configuration["FINANCE_BOOTSTRAP_TENANT_SLUG"] ?? "bootstrap-ops"
          : tenantSlug.Trim();

        await using var connection = await dataSource.OpenConnectionAsync();
        var summary = await BuildProjectionSummary(connection, resolvedTenantSlug);
        return TypedResults.Ok(summary);
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

  private static async Task<long?> LookupTenantId(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug)
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
    return scalar is long value ? value : scalar as long?;
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
      reader.GetInt32(0),
      reader.GetInt64(1),
      reader.GetInt64(2),
      new ProjectionStatusCounts(
        reader.GetInt32(3),
        reader.GetInt32(4),
        reader.GetInt32(5),
        reader.GetInt32(6)));
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
}

public sealed record HealthResponse(string Service, string Status);
public sealed record DependencyResponse(string Name, string Status);
public sealed record ReadinessResponse(string Service, string Status, IReadOnlyList<DependencyResponse> Dependencies);
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
internal sealed record SalesOutboxEvent(string PublicId, string AggregateType, string AggregatePublicId, string EventType, string Payload);
internal sealed record ProjectionMutation(int Created, int Updated, int Processed);
