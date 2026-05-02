using System.Globalization;
using System.Text.Json;
using Npgsql;
using NpgsqlTypes;

namespace Billing.Api;

public static class Server
{
  public static void MapBillingRoutes(this WebApplication app)
  {
    app.MapGet("/health/live", () =>
      TypedResults.Ok(new HealthResponse("billing", "live")));

    app.MapGet("/health/ready", async Task<IResult> (NpgsqlDataSource dataSource) =>
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      await using var command = new NpgsqlCommand("SELECT 1", connection);
      await command.ExecuteScalarAsync();
      return TypedResults.Ok(new HealthResponse("billing", "ready"));
    });

    app.MapGet("/health/details", async Task<IResult> (NpgsqlDataSource dataSource) =>
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      var dependencies = await BuildDependencies(connection);
      return TypedResults.Ok(new ReadinessResponse("billing", dependencies.All(dep => dep.Status == "ready") ? "ready" : "degraded", dependencies));
    });

    app.MapGet("/api/billing/plans", async Task<IResult> (string? active, NpgsqlDataSource dataSource) =>
    {
      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListPlans(connection, ParseActiveFilter(active)));
    });

    app.MapPost("/api/billing/plans", async Task<IResult> (CreatePlanRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_payload", "Billing plan payload is required."));
      }

      var code = request.Code?.Trim().ToLowerInvariant() ?? string.Empty;
      var name = request.Name?.Trim() ?? string.Empty;
      var description = request.Description?.Trim() ?? string.Empty;
      var intervalUnit = NormalizeIntervalUnit(request.IntervalUnit);

      if (code.Length < 3)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_code", "Billing plan code must contain at least 3 characters."));
      }

      if (name.Length < 3)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_name", "Billing plan name must contain at least 3 characters."));
      }

      if (request.AmountCents is null || request.AmountCents <= 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_amount", "Billing plan amount must be greater than zero."));
      }

      if (intervalUnit.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_interval", "Billing plan interval must be monthly or yearly."));
      }

      var intervalCount = request.IntervalCount.GetValueOrDefault(1);
      var gracePeriodDays = request.GracePeriodDays.GetValueOrDefault(0);
      var maxRetries = request.MaxRetries.GetValueOrDefault(0);

      if (intervalCount <= 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_interval_count", "Billing plan interval count must be greater than zero."));
      }

      if (gracePeriodDays < 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_grace_period", "Billing plan grace period days cannot be negative."));
      }

      if (maxRetries < 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_plan_max_retries", "Billing plan max retries cannot be negative."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      if (await PlanCodeExists(connection, transaction, code))
      {
        return TypedResults.Conflict(new ErrorResponse("billing_plan_conflict", "Billing plan code already exists."));
      }

      var response = await CreatePlan(connection, transaction, code, name, description, request.AmountCents.Value, intervalUnit, intervalCount, gracePeriodDays, maxRetries);
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapGet("/api/billing/subscriptions", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListSubscriptions(connection, tenantSlug.Trim(), NormalizeSubscriptionStatus(status)));
    });

    app.MapGet("/api/billing/subscriptions/{publicId}", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      var subscription = await FindSubscription(connection, tenantSlug.Trim(), publicId);
      return subscription is null
        ? TypedResults.NotFound(new ErrorResponse("billing_subscription_not_found", "Billing subscription was not found."))
        : TypedResults.Ok(subscription);
    });

    app.MapGet("/api/billing/subscriptions/{publicId}/events", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListSubscriptionEvents(connection, tenantSlug.Trim(), publicId));
    });

    app.MapGet("/api/billing/events", async Task<IResult> (string? tenantSlug, string? subscriptionPublicId, string? invoicePublicId, string? eventCode, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListBillingEvents(connection, tenantSlug.Trim(), subscriptionPublicId, invoicePublicId, eventCode));
    });

    app.MapPost("/api/billing/subscriptions", async Task<IResult> (CreateSubscriptionRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_subscription_payload", "Billing subscription payload is required."));
      }

      var tenantSlug = request.TenantSlug?.Trim() ?? string.Empty;
      var planPublicId = request.PlanPublicId?.Trim() ?? string.Empty;

      if (tenantSlug.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (planPublicId.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("plan_public_id_required", "Plan public id is required."));
      }

      if (!TryParseDate(request.StartedOn, out var startedOn))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_subscription_start", "Subscription start date must be a valid YYYY-MM-DD value."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var tenantId = await LookupTenantId(connection, transaction, tenantSlug);
      if (tenantId is null)
      {
        return TypedResults.NotFound(new ErrorResponse("tenant_not_found", "Tenant was not found for billing subscription creation."));
      }

      var plan = await FindPlanInternal(connection, transaction, planPublicId);
      if (plan is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_plan_not_found", "Billing plan was not found."));
      }

      var response = await CreateSubscription(
        connection,
        transaction,
        tenantId.Value,
        tenantSlug,
        plan,
        request.ExternalReference?.Trim() ?? string.Empty,
        startedOn ?? DateOnly.FromDateTime(DateTime.UtcNow));
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapPost("/api/billing/subscriptions/{publicId}/suspend", async Task<IResult> (string publicId, SubscriptionStatusRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (!TryParseUtcTimestamp(request.EffectiveAt, out var effectiveAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_effective_at", "Effective timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var subscription = await FindSubscriptionInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (subscription is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_subscription_not_found", "Billing subscription was not found."));
      }

      var response = await UpdateSubscriptionStatus(connection, transaction, subscription, "suspended", effectiveAt ?? DateTime.UtcNow, request.Reason?.Trim() ?? "Subscription suspended by operator.", "billing");
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapPost("/api/billing/subscriptions/{publicId}/reactivate", async Task<IResult> (string publicId, SubscriptionStatusRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (!TryParseUtcTimestamp(request.EffectiveAt, out var effectiveAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_effective_at", "Effective timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var subscription = await FindSubscriptionInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (subscription is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_subscription_not_found", "Billing subscription was not found."));
      }

      var response = await UpdateSubscriptionStatus(connection, transaction, subscription, "active", effectiveAt ?? DateTime.UtcNow, request.Reason?.Trim() ?? "Subscription reactivated.", "billing");
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapGet("/api/billing/invoices", async Task<IResult> (string? tenantSlug, string? status, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListInvoices(connection, tenantSlug.Trim(), NormalizeInvoiceStatus(status)));
    });

    app.MapGet("/api/billing/invoices/{publicId}", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      var invoice = await FindInvoice(connection, tenantSlug.Trim(), publicId);
      return invoice is null
        ? TypedResults.NotFound(new ErrorResponse("billing_invoice_not_found", "Billing invoice was not found."))
        : TypedResults.Ok(invoice);
    });

    app.MapPost("/api/billing/subscriptions/{publicId}/invoices", async Task<IResult> (string publicId, CreateInvoiceRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (!TryParseDate(request.DueDate, out var dueDate))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_invoice_due_date", "Billing invoice due date must be a valid YYYY-MM-DD value."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var subscription = await FindSubscriptionInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (subscription is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_subscription_not_found", "Billing subscription was not found."));
      }

      var amountCents = request.AmountCents.GetValueOrDefault(subscription.PlanAmountCents);
      if (amountCents <= 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_invoice_amount", "Billing invoice amount must be greater than zero."));
      }

      var response = await CreateInvoice(connection, transaction, subscription, amountCents, dueDate ?? subscription.CurrentPeriodEnd, request.Number?.Trim() ?? string.Empty);
      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapGet("/api/billing/invoices/{publicId}/attempts", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListPaymentAttempts(connection, tenantSlug.Trim(), publicId));
    });

    app.MapPost("/api/billing/invoices/{publicId}/attempts", async Task<IResult> (string publicId, CreatePaymentAttemptRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var status = NormalizeAttemptStatus(request.Status);
      if (status.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_attempt_status", "Billing payment attempt status must be succeeded or failed."));
      }

      var provider = request.Provider?.Trim().ToLowerInvariant() ?? string.Empty;
      var idempotencyKey = request.IdempotencyKey?.Trim() ?? string.Empty;

      if (provider.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("provider_required", "Billing payment attempt provider is required."));
      }

      if (idempotencyKey.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("idempotency_key_required", "Billing payment attempt idempotency key is required."));
      }

      if (!TryParseUtcTimestamp(request.AttemptedAt, out var attemptedAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_attempted_at", "Billing payment attempt timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var invoice = await FindInvoiceInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (invoice is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_invoice_not_found", "Billing invoice was not found."));
      }

      var outcome = await ApplyPaymentAttempt(
        connection,
        transaction,
        invoice,
        provider,
        status,
        idempotencyKey,
        request.ExternalReference?.Trim() ?? string.Empty,
        request.FailureReason?.Trim() ?? string.Empty,
        attemptedAt ?? DateTime.UtcNow,
        "billing");

      await transaction.CommitAsync();
      return TypedResults.Ok(outcome);
    });

    app.MapPost("/api/billing/webhook-events/process", async Task<IResult> (ProcessWebhookEventRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var webhookEventPublicId = request.WebhookEventPublicId?.Trim() ?? string.Empty;
      if (webhookEventPublicId.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("webhook_event_required", "Webhook event public id is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var webhookEvent = await FindWebhookEvent(connection, transaction, webhookEventPublicId);
      if (webhookEvent is null)
      {
        return TypedResults.NotFound(new ErrorResponse("webhook_event_not_found", "Webhook event was not found."));
      }

      if (webhookEvent.Status == "rejected")
      {
        return TypedResults.Conflict(new ErrorResponse("webhook_event_rejected", "Rejected webhook events cannot be processed by billing."));
      }

      var outcome = await ProcessWebhookEventAgainstBilling(connection, transaction, request.TenantSlug.Trim(), webhookEvent);
      if (outcome.ErrorCode == "billing_invoice_not_found")
      {
        return TypedResults.NotFound(new ErrorResponse("billing_invoice_not_found", "No billing invoice matched the webhook event external id."));
      }

      if (outcome.ErrorCode == "unsupported_webhook_event")
      {
        return TypedResults.BadRequest(new ErrorResponse("unsupported_webhook_event", "Webhook event type is not supported by billing."));
      }

      await transaction.CommitAsync();
      return TypedResults.Ok(new WebhookProcessResponse(webhookEvent.PublicId, webhookEvent.EventType, webhookEvent.ExternalId, outcome.OutcomeStatus, outcome.AttemptOutcome.Idempotent, outcome.AttemptOutcome.Attempt, outcome.AttemptOutcome.Invoice, outcome.AttemptOutcome.Subscription));
    });

    app.MapGet("/api/billing/webhook-events/pending", async Task<IResult> (string? tenantSlug, string? provider, int? limit, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var normalizedProvider = provider?.Trim().ToLowerInvariant() ?? string.Empty;
      var boundedLimit = limit.GetValueOrDefault(25);
      if (boundedLimit <= 0)
      {
        boundedLimit = 25;
      }

      if (boundedLimit > 100)
      {
        boundedLimit = 100;
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListPendingWebhookEvents(connection, tenantSlug.Trim(), normalizedProvider, boundedLimit));
    });

    app.MapPost("/api/billing/webhook-events/process-batch", async Task<IResult> (ProcessWebhookBatchRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var normalizedProvider = request.Provider?.Trim().ToLowerInvariant() ?? string.Empty;
      var boundedLimit = request.Limit.GetValueOrDefault(25);
      if (boundedLimit <= 0)
      {
        boundedLimit = 25;
      }

      if (boundedLimit > 100)
      {
        boundedLimit = 100;
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var candidates = await ListPendingWebhookEvents(connection, request.TenantSlug.Trim(), normalizedProvider, boundedLimit, transaction);
      var processed = new List<WebhookProcessResponse>();
      var skipped = new List<WebhookBatchSkipResponse>();

      foreach (var candidate in candidates)
      {
        var webhookEvent = await FindWebhookEvent(connection, transaction, candidate.PublicId);
        if (webhookEvent is null)
        {
          skipped.Add(new WebhookBatchSkipResponse(candidate.PublicId, "webhook_event_not_found", "Webhook event was not found during batch processing."));
          continue;
        }

        if (webhookEvent.Status == "rejected")
        {
          skipped.Add(new WebhookBatchSkipResponse(candidate.PublicId, "webhook_event_rejected", "Rejected webhook events cannot be processed by billing."));
          continue;
        }

        var outcome = await ProcessWebhookEventAgainstBilling(connection, transaction, request.TenantSlug.Trim(), webhookEvent);
        if (outcome.ErrorCode.Length > 0)
        {
          skipped.Add(new WebhookBatchSkipResponse(candidate.PublicId, outcome.ErrorCode, outcome.ErrorMessage));
          continue;
        }

        processed.Add(new WebhookProcessResponse(
          webhookEvent.PublicId,
          webhookEvent.EventType,
          webhookEvent.ExternalId,
          outcome.OutcomeStatus,
          outcome.AttemptOutcome.Idempotent,
          outcome.AttemptOutcome.Attempt,
          outcome.AttemptOutcome.Invoice,
          outcome.AttemptOutcome.Subscription));
      }

      await transaction.CommitAsync();
      return TypedResults.Ok(new WebhookBatchProcessResponse(request.TenantSlug.Trim(), processed.Count, skipped.Count, processed, skipped));
    });

    app.MapGet("/api/billing/reports/operations", async Task<IResult> (string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await BuildOperationsReport(connection, tenantSlug.Trim()));
    });

    app.MapGet("/api/billing/recovery/cases", async Task<IResult> (string? tenantSlug, string? status, string? severity, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListRecoveryCases(connection, tenantSlug.Trim(), NormalizeRecoveryStatus(status), NormalizeRecoverySeverity(severity)));
    });

    app.MapGet("/api/billing/recovery/cases/{publicId}", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      var recoveryCase = await FindRecoveryCase(connection, tenantSlug.Trim(), publicId);
      return recoveryCase is null
        ? TypedResults.NotFound(new ErrorResponse("billing_recovery_case_not_found", "Billing recovery case was not found."))
        : TypedResults.Ok(recoveryCase);
    });

    app.MapGet("/api/billing/recovery/cases/{publicId}/actions", async Task<IResult> (string publicId, string? tenantSlug, NpgsqlDataSource dataSource) =>
    {
      if (string.IsNullOrWhiteSpace(tenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      return TypedResults.Ok(await ListRecoveryActions(connection, tenantSlug.Trim(), publicId));
    });

    app.MapPost("/api/billing/invoices/{publicId}/recovery/open", async Task<IResult> (string publicId, OpenRecoveryCaseRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (!TryParseUtcTimestamp(request.NextActionAt, out var nextActionAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_next_action_at", "Next action timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var invoice = await FindInvoiceInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (invoice is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_invoice_not_found", "Billing invoice was not found."));
      }

      var recoveryCase = await OpenOrRefreshRecoveryCase(
        connection,
        transaction,
        invoice,
        invoice.RetryCount > 0 ? invoice.RetryCount : 1,
        invoice.RetryCount >= invoice.MaxRetries && invoice.MaxRetries > 0 ? "critical" : "attention",
        nextActionAt,
        request.WorkflowDefinitionKey?.Trim() ?? string.Empty,
        request.Notes?.Trim() ?? string.Empty,
        request.Actor?.Trim() ?? "billing");

      await transaction.CommitAsync();
      return TypedResults.Ok(recoveryCase);
    });

    app.MapPost("/api/billing/recovery/cases/{publicId}/touchpoints", async Task<IResult> (string publicId, RegisterRecoveryTouchpointRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var channel = NormalizeRecoveryChannel(request.Channel);
      if (channel.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_recovery_channel", "Recovery contact channel must be email, whatsapp, phone or manual."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var recoveryCase = await FindRecoveryCaseInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (recoveryCase is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_recovery_case_not_found", "Billing recovery case was not found."));
      }

      var response = await RegisterRecoveryTouchpoint(
        connection,
        transaction,
        recoveryCase,
        request.Actor?.Trim() ?? "billing",
        channel,
        request.Provider?.Trim().ToLowerInvariant() ?? string.Empty,
        request.TouchpointPublicId?.Trim() ?? string.Empty,
        request.DeliveryPublicId?.Trim() ?? string.Empty,
        request.Notes?.Trim() ?? string.Empty);

      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapPost("/api/billing/recovery/cases/{publicId}/promise", async Task<IResult> (string publicId, RegisterRecoveryPromiseRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      if (!TryParseDate(request.PromisedPaymentDate, out var promisedPaymentDate) || promisedPaymentDate is null)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_promised_payment_date", "Promised payment date must be a valid YYYY-MM-DD value."));
      }

      if (!TryParseUtcTimestamp(request.NextActionAt, out var nextActionAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_next_action_at", "Next action timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var recoveryCase = await FindRecoveryCaseInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (recoveryCase is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_recovery_case_not_found", "Billing recovery case was not found."));
      }

      var response = await RegisterRecoveryPromise(
        connection,
        transaction,
        recoveryCase,
        request.Actor?.Trim() ?? "billing",
        promisedPaymentDate.Value,
        nextActionAt,
        request.Notes?.Trim() ?? string.Empty);

      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });

    app.MapPost("/api/billing/recovery/cases/{publicId}/close", async Task<IResult> (string publicId, CloseRecoveryCaseRequest? request, NpgsqlDataSource dataSource) =>
    {
      if (request is null || string.IsNullOrWhiteSpace(request.TenantSlug))
      {
        return TypedResults.BadRequest(new ErrorResponse("tenant_slug_required", "Tenant slug is required."));
      }

      var status = NormalizeRecoveryTerminalStatus(request.Status);
      if (status.Length == 0)
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_recovery_status", "Recovery close status must be recovered, defaulted or closed."));
      }

      if (!TryParseUtcTimestamp(request.ResolvedAt, out var resolvedAt))
      {
        return TypedResults.BadRequest(new ErrorResponse("invalid_resolved_at", "Resolved timestamp must be a valid UTC instant."));
      }

      await using var connection = await dataSource.OpenConnectionAsync();
      await using var transaction = await connection.BeginTransactionAsync();

      var recoveryCase = await FindRecoveryCaseInternal(connection, transaction, request.TenantSlug.Trim(), publicId);
      if (recoveryCase is null)
      {
        return TypedResults.NotFound(new ErrorResponse("billing_recovery_case_not_found", "Billing recovery case was not found."));
      }

      var response = await CloseRecoveryCase(
        connection,
        transaction,
        recoveryCase,
        status,
        request.ResolutionCode?.Trim() ?? string.Empty,
        request.Actor?.Trim() ?? "billing",
        request.Notes?.Trim() ?? string.Empty,
        resolvedAt ?? DateTime.UtcNow);

      await transaction.CommitAsync();
      return TypedResults.Ok(response);
    });
  }

  private static async Task<IReadOnlyList<DependencyStatus>> BuildDependencies(NpgsqlConnection connection)
  {
    var dependencies = new List<DependencyStatus>();

    await using (var command = new NpgsqlCommand("SELECT 1", connection))
    {
      await command.ExecuteScalarAsync();
      dependencies.Add(new DependencyStatus("postgresql", "ready"));
    }

    await using (var command = new NpgsqlCommand("SELECT to_regclass('webhook_hub.webhook_events') IS NOT NULL", connection))
    {
      var available = Convert.ToBoolean(await command.ExecuteScalarAsync() ?? false, CultureInfo.InvariantCulture);
      dependencies.Add(new DependencyStatus("webhook-hub", available ? "ready" : "pending-runtime-wiring"));
    }

    return dependencies;
  }

  private static async Task<bool> PlanCodeExists(NpgsqlConnection connection, NpgsqlTransaction transaction, string code)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT EXISTS (
        SELECT 1
        FROM billing.plans
        WHERE code = $1
      )
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(code);
    return Convert.ToBoolean(await command.ExecuteScalarAsync() ?? false, CultureInfo.InvariantCulture);
  }

  private static async Task<IReadOnlyList<PlanResponse>> ListPlans(NpgsqlConnection connection, bool? active)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        public_id::text,
        code,
        name,
        description,
        amount_cents,
        currency_code,
        interval_unit,
        interval_count,
        grace_period_days,
        max_retries,
        active,
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.plans
      WHERE ($1::text = '' OR active = $2)
      ORDER BY code, created_at
      """,
      connection);
    command.Parameters.AddWithValue(active is null ? string.Empty : active.Value ? "true" : "false");
    command.Parameters.AddWithValue(active ?? false);

    var response = new List<PlanResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new PlanResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetInt64(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetInt32(7),
        reader.GetInt32(8),
        reader.GetInt32(9),
        reader.GetBoolean(10),
        reader.GetString(11),
        reader.GetString(12)));
    }

    return response;
  }

  private static async Task<PlanResponse> CreatePlan(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    string code,
    string name,
    string description,
    long amountCents,
    string intervalUnit,
    int intervalCount,
    int gracePeriodDays,
    int maxRetries)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO billing.plans (
        public_id,
        code,
        name,
        description,
        amount_cents,
        interval_unit,
        interval_count,
        grace_period_days,
        max_retries
      )
      VALUES (
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
      )
      RETURNING
        public_id::text,
        code,
        name,
        description,
        amount_cents,
        currency_code,
        interval_unit,
        interval_count,
        grace_period_days,
        max_retries,
        active,
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(code);
    command.Parameters.AddWithValue(name);
    command.Parameters.AddWithValue(description);
    command.Parameters.AddWithValue(amountCents);
    command.Parameters.AddWithValue(intervalUnit);
    command.Parameters.AddWithValue(intervalCount);
    command.Parameters.AddWithValue(gracePeriodDays);
    command.Parameters.AddWithValue(maxRetries);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new PlanResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetInt32(7),
      reader.GetInt32(8),
      reader.GetInt32(9),
      reader.GetBoolean(10),
      reader.GetString(11),
      reader.GetString(12));
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
    var value = await command.ExecuteScalarAsync();
    return value is null ? null : Convert.ToInt64(value, CultureInfo.InvariantCulture);
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

  private static async Task<InternalPlan?> FindPlanInternal(NpgsqlConnection connection, NpgsqlTransaction transaction, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        id,
        public_id::text,
        code,
        amount_cents,
        interval_unit,
        interval_count,
        grace_period_days,
        max_retries,
        active
      FROM billing.plans
      WHERE public_id = $1::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new InternalPlan(
      reader.GetInt64(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt64(3),
      reader.GetString(4),
      reader.GetInt32(5),
      reader.GetInt32(6),
      reader.GetInt32(7),
      reader.GetBoolean(8));
  }

  private static async Task<IReadOnlyList<SubscriptionResponse>> ListSubscriptions(NpgsqlConnection connection, string tenantSlug, string status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        subscription.public_id::text,
        tenant.slug,
        plan.public_id::text,
        plan.code,
        subscription.external_reference,
        subscription.status,
        to_char(subscription.current_period_start, 'YYYY-MM-DD'),
        to_char(subscription.current_period_end, 'YYYY-MM-DD'),
        COALESCE(to_char(subscription.grace_ends_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(subscription.suspended_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(subscription.cancelled_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(subscription.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(subscription.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscriptions AS subscription
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = subscription.tenant_id
      INNER JOIN billing.plans AS plan
        ON plan.id = subscription.plan_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR subscription.status = $2)
      ORDER BY subscription.created_at, subscription.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(status);

    var response = new List<SubscriptionResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new SubscriptionResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11),
        reader.GetString(12)));
    }

    return response;
  }

  private static async Task<SubscriptionResponse?> FindSubscription(NpgsqlConnection connection, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        subscription.public_id::text,
        tenant.slug,
        plan.public_id::text,
        plan.code,
        subscription.external_reference,
        subscription.status,
        to_char(subscription.current_period_start, 'YYYY-MM-DD'),
        to_char(subscription.current_period_end, 'YYYY-MM-DD'),
        COALESCE(to_char(subscription.grace_ends_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(subscription.suspended_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(subscription.cancelled_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        to_char(subscription.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(subscription.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscriptions AS subscription
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = subscription.tenant_id
      INNER JOIN billing.plans AS plan
        ON plan.id = subscription.plan_id
      WHERE tenant.slug = $1
        AND subscription.public_id = $2::uuid
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

    return new SubscriptionResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10),
      reader.GetString(11),
      reader.GetString(12));
  }

  private static async Task<InternalSubscription?> FindSubscriptionInternal(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        subscription.id,
        subscription.tenant_id,
        subscription.public_id::text,
        tenant.slug,
        subscription.plan_id,
        plan.public_id::text,
        plan.code,
        plan.amount_cents,
        plan.interval_unit,
        plan.interval_count,
        plan.grace_period_days,
        plan.max_retries,
        subscription.external_reference,
        subscription.status,
        subscription.current_period_start,
        subscription.current_period_end
      FROM billing.subscriptions AS subscription
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = subscription.tenant_id
      INNER JOIN billing.plans AS plan
        ON plan.id = subscription.plan_id
      WHERE tenant.slug = $1
        AND subscription.public_id = $2::uuid
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

    return new InternalSubscription(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetInt64(7),
      reader.GetString(8),
      reader.GetInt32(9),
      reader.GetInt32(10),
      reader.GetInt32(11),
      reader.GetString(12),
      reader.GetString(13),
      reader.GetFieldValue<DateOnly>(14),
      reader.GetFieldValue<DateOnly>(15));
  }

  private static async Task<SubscriptionResponse> CreateSubscription(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string tenantSlug,
    InternalPlan plan,
    string externalReference,
    DateOnly startedOn)
  {
    var currentPeriodEnd = ComputePeriodEnd(startedOn, plan.IntervalUnit, plan.IntervalCount);
    string subscriptionPublicId;
    string createdAt;
    string updatedAt;

    await using (var command = new NpgsqlCommand(
      """
      INSERT INTO billing.subscriptions (
        tenant_id,
        public_id,
        plan_id,
        external_reference,
        status,
        current_period_start,
        current_period_end
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2,
        $3,
        'active',
        $4,
        $5
      )
      RETURNING
        public_id::text,
        to_char(created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(tenantId);
      command.Parameters.AddWithValue(plan.Id);
      command.Parameters.AddWithValue(externalReference);
      command.Parameters.AddWithValue(startedOn);
      command.Parameters.AddWithValue(currentPeriodEnd);

      await using var reader = await command.ExecuteReaderAsync();
      await reader.ReadAsync();
      subscriptionPublicId = reader.GetString(0);
      createdAt = reader.GetString(1);
      updatedAt = reader.GetString(2);
    }

    await InsertSubscriptionEvent(
      connection,
      transaction,
      tenantId,
      subscriptionPublicId,
      null,
      "subscription_created",
      "billing",
      "Billing subscription created.",
      JsonSerializer.Serialize(new
      {
        subscriptionPublicId,
        tenantSlug,
        planPublicId = plan.PublicId,
        planCode = plan.Code,
        currentPeriodStart = startedOn.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
        currentPeriodEnd = currentPeriodEnd.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture)
      }));

    return new SubscriptionResponse(
      subscriptionPublicId,
      tenantSlug,
      plan.PublicId,
      plan.Code,
      externalReference,
      "active",
      startedOn.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
      currentPeriodEnd.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
      string.Empty,
      string.Empty,
      string.Empty,
      createdAt,
      updatedAt);
  }

  private static async Task<SubscriptionResponse> UpdateSubscriptionStatus(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalSubscription subscription,
    string status,
    DateTime effectiveAt,
    string reason,
    string actor)
  {
    string updatedAt;
    await using (var command = new NpgsqlCommand(
      """
      UPDATE billing.subscriptions
      SET
        status = $3,
        grace_ends_at = CASE WHEN $3 = 'grace_period' THEN COALESCE(grace_ends_at, $4) ELSE NULL END,
        suspended_at = CASE WHEN $3 = 'suspended' THEN $4 ELSE NULL END,
        cancelled_at = CASE WHEN $3 = 'cancelled' THEN $4 ELSE cancelled_at END
      WHERE id = $1
        AND tenant_id = $2
      RETURNING to_char(updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(subscription.Id);
      command.Parameters.AddWithValue(subscription.TenantId);
      command.Parameters.AddWithValue(status);
      command.Parameters.AddWithValue(effectiveAt);
      updatedAt = Convert.ToString(await command.ExecuteScalarAsync(), CultureInfo.InvariantCulture) ?? string.Empty;
    }

    var eventCode = status switch
    {
      "suspended" => "subscription_suspended",
      "active" => "subscription_reactivated",
      "cancelled" => "subscription_cancelled",
      _ => "subscription_status_changed"
    };

    await InsertSubscriptionEvent(
      connection,
      transaction,
      subscription.TenantId,
      subscription.PublicId,
      null,
      eventCode,
      actor,
      reason,
      JsonSerializer.Serialize(new
      {
        subscriptionPublicId = subscription.PublicId,
        status,
        reason,
        effectiveAt = effectiveAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture)
      }));

    var response = await FindSubscription(connection, subscription.TenantSlug, subscription.PublicId);
    return response ?? new SubscriptionResponse(
      subscription.PublicId,
      subscription.TenantSlug,
      subscription.PlanPublicId,
      subscription.PlanCode,
      subscription.ExternalReference,
      status,
      subscription.CurrentPeriodStart.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
      subscription.CurrentPeriodEnd.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture),
      status == "grace_period" ? effectiveAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture) : string.Empty,
      status == "suspended" ? effectiveAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture) : string.Empty,
      status == "cancelled" ? effectiveAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture) : string.Empty,
      updatedAt,
      updatedAt);
  }

  private static async Task<IReadOnlyList<SubscriptionEventResponse>> ListSubscriptionEvents(NpgsqlConnection connection, string tenantSlug, string subscriptionPublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        event.public_id::text,
        tenant.slug,
        event.subscription_public_id::text,
        COALESCE(event.invoice_public_id::text, ''),
        event.event_code,
        event.actor,
        event.summary,
        event.payload::text,
        to_char(event.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscription_events AS event
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = event.tenant_id
      WHERE tenant.slug = $1
        AND event.subscription_public_id = $2::uuid
      ORDER BY event.created_at, event.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(subscriptionPublicId);

    var response = new List<SubscriptionEventResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new SubscriptionEventResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8)));
    }

    return response;
  }

  private static async Task<IReadOnlyList<SubscriptionEventResponse>> ListBillingEvents(
    NpgsqlConnection connection,
    string tenantSlug,
    string? subscriptionPublicId,
    string? invoicePublicId,
    string? eventCode)
  {
    var sql =
      """
      SELECT
        event.public_id::text,
        tenant.slug,
        event.subscription_public_id::text,
        COALESCE(event.invoice_public_id::text, ''),
        event.event_code,
        event.actor,
        event.summary,
        event.payload::text,
        to_char(event.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscription_events AS event
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = event.tenant_id
      WHERE tenant.slug = $1
      """;

    var parameters = new List<object?> { tenantSlug };
    if (!string.IsNullOrWhiteSpace(subscriptionPublicId))
    {
      sql += $" AND event.subscription_public_id = ${parameters.Count + 1}::uuid";
      parameters.Add(subscriptionPublicId.Trim());
    }
    if (!string.IsNullOrWhiteSpace(invoicePublicId))
    {
      sql += $" AND event.invoice_public_id = ${parameters.Count + 1}::uuid";
      parameters.Add(invoicePublicId.Trim());
    }
    if (!string.IsNullOrWhiteSpace(eventCode))
    {
      sql += $" AND event.event_code = ${parameters.Count + 1}";
      parameters.Add(eventCode.Trim());
    }
    sql += " ORDER BY event.created_at DESC, event.id DESC";

    await using var command = new NpgsqlCommand(sql, connection);
    for (var index = 0; index < parameters.Count; index++)
    {
      command.Parameters.AddWithValue(parameters[index] ?? string.Empty);
    }

    var response = new List<SubscriptionEventResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new SubscriptionEventResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8)));
    }

    return response;
  }

  private static async Task<InvoiceResponse> CreateInvoice(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalSubscription subscription,
    long amountCents,
    DateOnly dueDate,
    string explicitNumber)
  {
    var invoiceNumber = explicitNumber.Length > 0
      ? explicitNumber
      : await GenerateInvoiceNumber(connection, transaction, subscription.TenantId, subscription.TenantSlug);

    string invoicePublicId;

    await using (var command = new NpgsqlCommand(
      """
      INSERT INTO billing.subscription_invoices (
        tenant_id,
        subscription_id,
        public_id,
        number,
        status,
        amount_cents,
        due_date
      )
      VALUES (
        $1,
        $2,
        gen_random_uuid(),
        $3,
        'open',
        $4,
        $5
      )
      RETURNING public_id::text
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(subscription.TenantId);
      command.Parameters.AddWithValue(subscription.Id);
      command.Parameters.AddWithValue(invoiceNumber);
      command.Parameters.AddWithValue(amountCents);
      command.Parameters.AddWithValue(dueDate);
      invoicePublicId = Convert.ToString(await command.ExecuteScalarAsync(), CultureInfo.InvariantCulture) ?? string.Empty;
    }

    await InsertSubscriptionEvent(
      connection,
      transaction,
      subscription.TenantId,
      subscription.PublicId,
      invoicePublicId,
      "invoice_created",
      "billing",
      "Billing invoice created for subscription.",
      JsonSerializer.Serialize(new
      {
        subscriptionPublicId = subscription.PublicId,
        invoicePublicId,
        invoiceNumber,
        amountCents,
        dueDate = dueDate.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture)
      }));

    var response = await FindInvoice(connection, subscription.TenantSlug, invoicePublicId);
    return response ?? new InvoiceResponse(invoicePublicId, subscription.TenantSlug, subscription.PublicId, invoiceNumber, "open", amountCents, dueDate.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture), 0, string.Empty, string.Empty, string.Empty, string.Empty);
  }

  private static async Task<string> GenerateInvoiceNumber(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT COUNT(*) + 1
      FROM billing.subscription_invoices
      WHERE tenant_id = $1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    var sequence = Convert.ToInt32(await command.ExecuteScalarAsync() ?? 1, CultureInfo.InvariantCulture);
    var normalizedTenant = new string(tenantSlug.Where(char.IsLetterOrDigit).ToArray()).ToUpperInvariant();
    return $"{normalizedTenant}-BILL-{sequence:0000}";
  }

  private static async Task<IReadOnlyList<InvoiceResponse>> ListInvoices(NpgsqlConnection connection, string tenantSlug, string status)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        invoice.public_id::text,
        tenant.slug,
        subscription.public_id::text,
        invoice.number,
        invoice.status,
        invoice.amount_cents,
        to_char(invoice.due_date, 'YYYY-MM-DD'),
        invoice.retry_count,
        COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(invoice.gateway_reference, ''),
        to_char(invoice.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(invoice.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscription_invoices AS invoice
      INNER JOIN billing.subscriptions AS subscription
        ON subscription.id = invoice.subscription_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR invoice.status = $2)
      ORDER BY invoice.created_at, invoice.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(status);

    var response = new List<InvoiceResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new InvoiceResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetInt64(5),
        reader.GetString(6),
        reader.GetInt32(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11)));
    }

    return response;
  }

  private static async Task<InvoiceResponse?> FindInvoice(NpgsqlConnection connection, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        invoice.public_id::text,
        tenant.slug,
        subscription.public_id::text,
        invoice.number,
        invoice.status,
        invoice.amount_cents,
        to_char(invoice.due_date, 'YYYY-MM-DD'),
        invoice.retry_count,
        COALESCE(to_char(invoice.paid_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(invoice.gateway_reference, ''),
        to_char(invoice.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(invoice.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.subscription_invoices AS invoice
      INNER JOIN billing.subscriptions AS subscription
        ON subscription.id = invoice.subscription_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
        AND invoice.public_id = $2::uuid
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

    return new InvoiceResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetInt64(5),
      reader.GetString(6),
      reader.GetInt32(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10),
      reader.GetString(11));
  }

  private static async Task<InternalInvoice?> FindInvoiceInternal(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        invoice.id,
        invoice.tenant_id,
        invoice.public_id::text,
        tenant.slug,
        subscription.id,
        subscription.public_id::text,
        subscription.status,
        plan.grace_period_days,
        plan.max_retries,
        invoice.retry_count,
        invoice.status,
        invoice.number
      FROM billing.subscription_invoices AS invoice
      INNER JOIN billing.subscriptions AS subscription
        ON subscription.id = invoice.subscription_id
      INNER JOIN billing.plans AS plan
        ON plan.id = subscription.plan_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
        AND invoice.public_id = $2::uuid
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

    return new InternalInvoice(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetInt32(7),
      reader.GetInt32(8),
      reader.GetInt32(9),
      reader.GetString(10),
      reader.GetString(11));
  }

  private static async Task<InternalInvoice?> FindInvoiceByExternalReference(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug, string externalId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        invoice.id,
        invoice.tenant_id,
        invoice.public_id::text,
        tenant.slug,
        subscription.id,
        subscription.public_id::text,
        subscription.status,
        plan.grace_period_days,
        plan.max_retries,
        invoice.retry_count,
        invoice.status,
        invoice.number
      FROM billing.subscription_invoices AS invoice
      INNER JOIN billing.subscriptions AS subscription
        ON subscription.id = invoice.subscription_id
      INNER JOIN billing.plans AS plan
        ON plan.id = subscription.plan_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
        AND (
          invoice.public_id::text = $2
          OR invoice.number = $2
          OR invoice.gateway_reference = $2
        )
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(externalId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new InternalInvoice(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetInt64(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetInt32(7),
      reader.GetInt32(8),
      reader.GetInt32(9),
      reader.GetString(10),
      reader.GetString(11));
  }

  private static async Task<IReadOnlyList<PaymentAttemptResponse>> ListPaymentAttempts(NpgsqlConnection connection, string tenantSlug, string invoicePublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        attempt.public_id::text,
        invoice.public_id::text,
        tenant.slug,
        attempt.attempt_number,
        attempt.provider,
        attempt.status,
        attempt.idempotency_key,
        COALESCE(attempt.external_reference, ''),
        COALESCE(attempt.failure_reason, ''),
        to_char(attempt.attempted_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.payment_attempts AS attempt
      INNER JOIN billing.subscription_invoices AS invoice
        ON invoice.id = attempt.invoice_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = attempt.tenant_id
      WHERE tenant.slug = $1
        AND invoice.public_id = $2::uuid
      ORDER BY attempt.attempt_number, attempt.created_at
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(invoicePublicId);

    var response = new List<PaymentAttemptResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new PaymentAttemptResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetInt32(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9)));
    }

    return response;
  }

  private static async Task<AttemptOutcomeResponse> ApplyPaymentAttempt(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalInvoice invoice,
    string provider,
    string status,
    string idempotencyKey,
    string externalReference,
    string failureReason,
    DateTime attemptedAt,
    string actor)
  {
    var existingAttempt = await FindPaymentAttemptByIdempotencyKey(connection, transaction, invoice.TenantId, idempotencyKey);
    if (existingAttempt is not null)
    {
      var existingInvoice = await FindInvoice(connection, invoice.TenantSlug, invoice.PublicId)
        ?? new InvoiceResponse(invoice.PublicId, invoice.TenantSlug, invoice.SubscriptionPublicId, invoice.Number, invoice.Status, 0, string.Empty, invoice.RetryCount, string.Empty, string.Empty, string.Empty, string.Empty);
      var existingSubscription = await FindSubscription(connection, invoice.TenantSlug, invoice.SubscriptionPublicId)
        ?? new SubscriptionResponse(invoice.SubscriptionPublicId, invoice.TenantSlug, string.Empty, string.Empty, string.Empty, invoice.SubscriptionStatus, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty);
      return new AttemptOutcomeResponse(existingAttempt, existingInvoice, existingSubscription, true);
    }

    var attemptNumber = invoice.RetryCount + 1;
    string attemptPublicId;

    await using (var insertAttempt = new NpgsqlCommand(
      """
      INSERT INTO billing.payment_attempts (
        tenant_id,
        invoice_id,
        public_id,
        attempt_number,
        provider,
        status,
        idempotency_key,
        external_reference,
        failure_reason,
        attempted_at
      )
      VALUES (
        $1,
        $2,
        gen_random_uuid(),
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9
      )
      RETURNING public_id::text
      """,
      connection,
      transaction))
    {
      insertAttempt.Parameters.AddWithValue(invoice.TenantId);
      insertAttempt.Parameters.AddWithValue(invoice.Id);
      insertAttempt.Parameters.AddWithValue(attemptNumber);
      insertAttempt.Parameters.AddWithValue(provider);
      insertAttempt.Parameters.AddWithValue(status);
      insertAttempt.Parameters.AddWithValue(idempotencyKey);
      insertAttempt.Parameters.AddWithValue(externalReference);
      insertAttempt.Parameters.AddWithValue(failureReason);
      insertAttempt.Parameters.AddWithValue(attemptedAt);
      attemptPublicId = Convert.ToString(await insertAttempt.ExecuteScalarAsync(), CultureInfo.InvariantCulture) ?? string.Empty;
    }

    string nextSubscriptionStatus = invoice.SubscriptionStatus;
    var existingRecoveryCase = await FindRecoveryCaseByInvoice(connection, transaction, invoice.TenantId, invoice.PublicId);

    if (status == "succeeded")
    {
      await using var updateInvoice = new NpgsqlCommand(
        """
        UPDATE billing.subscription_invoices
        SET
          status = 'paid',
          retry_count = $3,
          paid_at = $4,
          last_attempt_at = $4,
          gateway_reference = CASE WHEN $5 = '' THEN gateway_reference ELSE $5 END
        WHERE id = $1
          AND tenant_id = $2
        """,
        connection,
        transaction);
      updateInvoice.Parameters.AddWithValue(invoice.Id);
      updateInvoice.Parameters.AddWithValue(invoice.TenantId);
      updateInvoice.Parameters.AddWithValue(attemptNumber);
      updateInvoice.Parameters.AddWithValue(attemptedAt);
      updateInvoice.Parameters.AddWithValue(externalReference);
      await updateInvoice.ExecuteNonQueryAsync();

      await using var updateSubscription = new NpgsqlCommand(
        """
        UPDATE billing.subscriptions
        SET
          status = 'active',
          grace_ends_at = NULL,
          suspended_at = NULL
        WHERE id = $1
          AND tenant_id = $2
        """,
        connection,
        transaction);
      updateSubscription.Parameters.AddWithValue(invoice.SubscriptionId);
      updateSubscription.Parameters.AddWithValue(invoice.TenantId);
      await updateSubscription.ExecuteNonQueryAsync();

      nextSubscriptionStatus = "active";

      if (existingRecoveryCase is not null && existingRecoveryCase.Status is "open" or "contacted" or "promised")
      {
        await CloseRecoveryCase(
          connection,
          transaction,
          existingRecoveryCase,
          "recovered",
          "invoice_paid",
          actor,
          "Invoice recovered after successful payment.",
          attemptedAt);
      }

      await InsertSubscriptionEvent(
        connection,
        transaction,
        invoice.TenantId,
        invoice.SubscriptionPublicId,
        invoice.PublicId,
        "invoice_paid",
        actor,
        "Billing invoice paid successfully.",
        JsonSerializer.Serialize(new
        {
          invoicePublicId = invoice.PublicId,
          attemptPublicId,
          idempotencyKey,
          provider,
          externalReference,
          attemptedAt = attemptedAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture)
        }));
    }
    else
    {
      await using var updateInvoice = new NpgsqlCommand(
        """
        UPDATE billing.subscription_invoices
        SET
          status = 'failed',
          retry_count = $3,
          last_attempt_at = $4,
          gateway_reference = CASE WHEN $5 = '' THEN gateway_reference ELSE $5 END
        WHERE id = $1
          AND tenant_id = $2
        """,
        connection,
        transaction);
      updateInvoice.Parameters.AddWithValue(invoice.Id);
      updateInvoice.Parameters.AddWithValue(invoice.TenantId);
      updateInvoice.Parameters.AddWithValue(attemptNumber);
      updateInvoice.Parameters.AddWithValue(attemptedAt);
      updateInvoice.Parameters.AddWithValue(externalReference);
      await updateInvoice.ExecuteNonQueryAsync();

      if (attemptNumber >= invoice.MaxRetries)
      {
        if (invoice.GracePeriodDays > 0)
        {
          var graceEndsAt = attemptedAt.AddDays(invoice.GracePeriodDays);
          await using var updateSubscription = new NpgsqlCommand(
            """
            UPDATE billing.subscriptions
            SET
              status = 'grace_period',
              grace_ends_at = $3,
              suspended_at = NULL
            WHERE id = $1
              AND tenant_id = $2
            """,
            connection,
            transaction);
          updateSubscription.Parameters.AddWithValue(invoice.SubscriptionId);
          updateSubscription.Parameters.AddWithValue(invoice.TenantId);
          updateSubscription.Parameters.AddWithValue(graceEndsAt);
          await updateSubscription.ExecuteNonQueryAsync();
          nextSubscriptionStatus = "grace_period";

          await InsertSubscriptionEvent(
            connection,
            transaction,
            invoice.TenantId,
            invoice.SubscriptionPublicId,
            invoice.PublicId,
            "subscription_grace_period_started",
            actor,
            "Billing subscription entered grace period after failed payment retries.",
            JsonSerializer.Serialize(new
            {
              invoicePublicId = invoice.PublicId,
              attemptPublicId,
              attemptNumber,
              maxRetries = invoice.MaxRetries,
              graceEndsAt = graceEndsAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture)
            }));
        }
        else
        {
          await using var updateSubscription = new NpgsqlCommand(
            """
            UPDATE billing.subscriptions
            SET
              status = 'suspended',
              suspended_at = $3,
              grace_ends_at = NULL
            WHERE id = $1
              AND tenant_id = $2
            """,
            connection,
            transaction);
          updateSubscription.Parameters.AddWithValue(invoice.SubscriptionId);
          updateSubscription.Parameters.AddWithValue(invoice.TenantId);
          updateSubscription.Parameters.AddWithValue(attemptedAt);
          await updateSubscription.ExecuteNonQueryAsync();
          nextSubscriptionStatus = "suspended";

          await InsertSubscriptionEvent(
            connection,
            transaction,
            invoice.TenantId,
            invoice.SubscriptionPublicId,
            invoice.PublicId,
            "subscription_suspended",
            actor,
            "Billing subscription suspended after failed payment retries.",
            JsonSerializer.Serialize(new
            {
              invoicePublicId = invoice.PublicId,
              attemptPublicId,
              attemptNumber,
              maxRetries = invoice.MaxRetries,
              suspendedAt = attemptedAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture)
            }));
        }
      }

      await OpenOrRefreshRecoveryCase(
        connection,
        transaction,
        invoice,
        attemptNumber,
        attemptNumber >= invoice.MaxRetries && invoice.MaxRetries > 0 ? "critical" : "attention",
        attemptedAt.AddHours(attemptNumber >= invoice.MaxRetries && invoice.MaxRetries > 0 ? 6 : 24),
        existingRecoveryCase?.WorkflowDefinitionKey ?? string.Empty,
        failureReason.Length > 0 ? failureReason : "Billing recovery case refreshed after failed payment attempt.",
        actor);

      await InsertSubscriptionEvent(
        connection,
        transaction,
        invoice.TenantId,
        invoice.SubscriptionPublicId,
        invoice.PublicId,
        "invoice_payment_failed",
        actor,
        failureReason.Length > 0 ? failureReason : "Billing payment attempt failed.",
        JsonSerializer.Serialize(new
        {
          invoicePublicId = invoice.PublicId,
          attemptPublicId,
          idempotencyKey,
          provider,
          externalReference,
          failureReason,
          attemptedAt = attemptedAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture)
        }));
    }

    var invoiceResponse = await FindInvoice(connection, invoice.TenantSlug, invoice.PublicId)
      ?? new InvoiceResponse(invoice.PublicId, invoice.TenantSlug, invoice.SubscriptionPublicId, invoice.Number, status == "succeeded" ? "paid" : "failed", 0, string.Empty, attemptNumber, status == "succeeded" ? attemptedAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture) : string.Empty, externalReference, string.Empty, string.Empty);

    var subscriptionResponse = await FindSubscription(connection, invoice.TenantSlug, invoice.SubscriptionPublicId)
      ?? new SubscriptionResponse(invoice.SubscriptionPublicId, invoice.TenantSlug, string.Empty, string.Empty, string.Empty, nextSubscriptionStatus, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty);

    var attemptResponse = await FindPaymentAttemptByPublicId(connection, invoice.TenantSlug, attemptPublicId)
      ?? new PaymentAttemptResponse(attemptPublicId, invoice.PublicId, invoice.TenantSlug, attemptNumber, provider, status, idempotencyKey, externalReference, failureReason, attemptedAt.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture));

    return new AttemptOutcomeResponse(attemptResponse, invoiceResponse, subscriptionResponse, false);
  }

  private static async Task<PaymentAttemptResponse?> FindPaymentAttemptByIdempotencyKey(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string idempotencyKey)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        attempt.public_id::text,
        invoice.public_id::text,
        tenant.slug,
        attempt.attempt_number,
        attempt.provider,
        attempt.status,
        attempt.idempotency_key,
        COALESCE(attempt.external_reference, ''),
        COALESCE(attempt.failure_reason, ''),
        to_char(attempt.attempted_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.payment_attempts AS attempt
      INNER JOIN billing.subscription_invoices AS invoice
        ON invoice.id = attempt.invoice_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = attempt.tenant_id
      WHERE attempt.tenant_id = $1
        AND attempt.idempotency_key = $2
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(idempotencyKey);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new PaymentAttemptResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt32(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9));
  }

  private static async Task<PaymentAttemptResponse?> FindPaymentAttemptByPublicId(NpgsqlConnection connection, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        attempt.public_id::text,
        invoice.public_id::text,
        tenant.slug,
        attempt.attempt_number,
        attempt.provider,
        attempt.status,
        attempt.idempotency_key,
        COALESCE(attempt.external_reference, ''),
        COALESCE(attempt.failure_reason, ''),
        to_char(attempt.attempted_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.payment_attempts AS attempt
      INNER JOIN billing.subscription_invoices AS invoice
        ON invoice.id = attempt.invoice_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = attempt.tenant_id
      WHERE tenant.slug = $1
        AND attempt.public_id = $2::uuid
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

    return new PaymentAttemptResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetInt32(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9));
  }

  private static async Task<WebhookEvent?> FindWebhookEvent(NpgsqlConnection connection, NpgsqlTransaction transaction, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        public_id::text,
        provider,
        event_type,
        external_id,
        COALESCE(payload_summary, ''),
        status,
        received_at
      FROM webhook_hub.webhook_events
      WHERE public_id = $1::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(publicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new WebhookEvent(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetFieldValue<DateTime>(6));
  }

  private static async Task<IReadOnlyList<PendingWebhookEventResponse>> ListPendingWebhookEvents(
    NpgsqlConnection connection,
    string tenantSlug,
    string provider,
    int limit,
    NpgsqlTransaction? transaction = null)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        event.public_id::text,
        event.provider,
        event.event_type,
        event.external_id,
        COALESCE(event.payload_summary, ''),
        event.status,
        to_char(event.received_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        invoice.public_id::text,
        invoice.status
      FROM webhook_hub.webhook_events AS event
      INNER JOIN billing.subscription_invoices AS invoice
        ON invoice.external_reference = event.external_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR event.provider = $2)
        AND event.status IN ('validated', 'queued', 'processing', 'forwarded', 'failed')
      ORDER BY event.received_at, event.id
      LIMIT $3
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(provider);
    command.Parameters.AddWithValue(limit);

    var responses = new List<PendingWebhookEventResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      responses.Add(new PendingWebhookEventResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8)));
    }

    return responses;
  }

  private static async Task<WebhookBillingProcessResult> ProcessWebhookEventAgainstBilling(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    string tenantSlug,
    WebhookEvent webhookEvent)
  {
    var invoice = await FindInvoiceByExternalReference(connection, transaction, tenantSlug, webhookEvent.ExternalId);
    if (invoice is null)
    {
      return WebhookBillingProcessResult.Error("billing_invoice_not_found", "No billing invoice matched the webhook event external id.");
    }

    var outcomeStatus = webhookEvent.EventType switch
    {
      "payment.succeeded" or "billing.invoice.paid" or "invoice.paid" => "succeeded",
      "payment.failed" or "billing.invoice.failed" or "invoice.failed" => "failed",
      _ => string.Empty
    };

    if (outcomeStatus.Length == 0)
    {
      return WebhookBillingProcessResult.Error("unsupported_webhook_event", "Webhook event type is not supported by billing.");
    }

    var attemptOutcome = await ApplyPaymentAttempt(
      connection,
      transaction,
      invoice,
      webhookEvent.Provider,
      outcomeStatus,
      $"webhook:{webhookEvent.PublicId}:{webhookEvent.EventType}",
      webhookEvent.ExternalId,
      outcomeStatus == "failed" ? webhookEvent.PayloadSummary : string.Empty,
      webhookEvent.ReceivedAt,
      "webhook-hub");

    return WebhookBillingProcessResult.Success(outcomeStatus, attemptOutcome);
  }

  private static async Task InsertSubscriptionEvent(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    long tenantId,
    string subscriptionPublicId,
    string? invoicePublicId,
    string eventCode,
    string actor,
    string summary,
    string payload)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO billing.subscription_events (
        tenant_id,
        public_id,
        subscription_public_id,
        invoice_public_id,
        event_code,
        actor,
        summary,
        payload
      )
      VALUES (
        $1,
        gen_random_uuid(),
        $2::uuid,
        CASE WHEN $3 = '' THEN NULL ELSE $3::uuid END,
        $4,
        $5,
        $6,
        $7::jsonb
      )
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(subscriptionPublicId);
    command.Parameters.AddWithValue(invoicePublicId ?? string.Empty);
    command.Parameters.AddWithValue(eventCode);
    command.Parameters.AddWithValue(actor);
    command.Parameters.AddWithValue(summary);
    command.Parameters.AddWithValue(payload);
    await command.ExecuteNonQueryAsync();
  }

  private static async Task<BillingOperationsReportResponse> BuildOperationsReport(NpgsqlConnection connection, string tenantSlug)
  {
    var now = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ", CultureInfo.InvariantCulture);
    var subscriptions = await BuildSubscriptionSummary(connection, tenantSlug);
    var invoices = await BuildInvoiceSummary(connection, tenantSlug);
    var attempts = await BuildAttemptSummary(connection, tenantSlug);
    var recovery = await BuildRecoverySummary(connection, tenantSlug);
    var plans = await BuildPlanSummary(connection);

    return new BillingOperationsReportResponse(tenantSlug, now, plans, subscriptions, invoices, attempts, recovery);
  }

  private static async Task<PlanSummary> BuildPlanSummary(NpgsqlConnection connection)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COUNT(*) FILTER (WHERE active) AS active
      FROM billing.plans
      """,
      connection);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new PlanSummary(ConvertToInt(reader.GetValue(0)), ConvertToInt(reader.GetValue(1)));
  }

  private static async Task<SubscriptionSummary> BuildSubscriptionSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COUNT(*) FILTER (WHERE subscription.status = 'active') AS active_count,
        COUNT(*) FILTER (WHERE subscription.status = 'grace_period') AS grace_count,
        COUNT(*) FILTER (WHERE subscription.status = 'suspended') AS suspended_count,
        COUNT(*) FILTER (WHERE subscription.status = 'cancelled') AS cancelled_count
      FROM billing.subscriptions AS subscription
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = subscription.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new SubscriptionSummary(
      ConvertToInt(reader.GetValue(0)),
      new SubscriptionStatusCounts(
        ConvertToInt(reader.GetValue(1)),
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4))));
  }

  private static async Task<InvoiceSummary> BuildInvoiceSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COALESCE(SUM(invoice.amount_cents), 0) AS total_amount_cents,
        COUNT(*) FILTER (WHERE invoice.status = 'open') AS open_count,
        COUNT(*) FILTER (WHERE invoice.status = 'paid') AS paid_count,
        COUNT(*) FILTER (WHERE invoice.status = 'failed') AS failed_count,
        COUNT(*) FILTER (WHERE invoice.status = 'void') AS void_count,
        COALESCE(SUM(invoice.amount_cents) FILTER (WHERE invoice.status = 'open'), 0) AS open_amount_cents,
        COALESCE(SUM(invoice.amount_cents) FILTER (WHERE invoice.status = 'paid'), 0) AS paid_amount_cents,
        COALESCE(SUM(invoice.amount_cents) FILTER (WHERE invoice.status = 'failed'), 0) AS failed_amount_cents
      FROM billing.subscription_invoices AS invoice
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = invoice.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new InvoiceSummary(
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt64(reader.GetValue(1)),
      new InvoiceStatusCounts(
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4)),
        ConvertToInt(reader.GetValue(5))),
      new InvoiceAmountBuckets(
        ConvertToInt64(reader.GetValue(6)),
        ConvertToInt64(reader.GetValue(7)),
        ConvertToInt64(reader.GetValue(8))));
  }

  private static async Task<AttemptSummary> BuildAttemptSummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COUNT(*) FILTER (WHERE attempt.status = 'succeeded') AS succeeded_count,
        COUNT(*) FILTER (WHERE attempt.status = 'failed') AS failed_count
      FROM billing.payment_attempts AS attempt
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = attempt.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new AttemptSummary(
      ConvertToInt(reader.GetValue(0)),
      ConvertToInt(reader.GetValue(1)),
      ConvertToInt(reader.GetValue(2)));
  }

  private static async Task<RecoverySummary> BuildRecoverySummary(NpgsqlConnection connection, string tenantSlug)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        COUNT(*) AS total,
        COUNT(*) FILTER (WHERE recovery.status = 'open') AS open_count,
        COUNT(*) FILTER (WHERE recovery.status = 'contacted') AS contacted_count,
        COUNT(*) FILTER (WHERE recovery.status = 'promised') AS promised_count,
        COUNT(*) FILTER (WHERE recovery.status = 'recovered') AS recovered_count,
        COUNT(*) FILTER (WHERE recovery.status = 'defaulted') AS defaulted_count,
        COUNT(*) FILTER (WHERE recovery.status = 'closed') AS closed_count,
        COUNT(*) FILTER (WHERE recovery.severity = 'attention') AS attention_count,
        COUNT(*) FILTER (WHERE recovery.severity = 'critical') AS critical_count,
        COUNT(*) FILTER (
          WHERE recovery.status IN ('open', 'contacted', 'promised')
            AND recovery.next_action_at IS NOT NULL
            AND recovery.next_action_at <= timezone('utc', now())
        ) AS due_actions_count,
        COUNT(*) FILTER (
          WHERE recovery.status = 'promised'
            AND recovery.promised_payment_date IS NOT NULL
            AND recovery.promised_payment_date <= timezone('utc', now())::date
        ) AS promises_due_count
      FROM billing.recovery_cases AS recovery
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = recovery.tenant_id
      WHERE tenant.slug = $1
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);

    await using var reader = await command.ExecuteReaderAsync();
    await reader.ReadAsync();
    return new RecoverySummary(
      ConvertToInt(reader.GetValue(0)),
      new RecoveryStatusCounts(
        ConvertToInt(reader.GetValue(1)),
        ConvertToInt(reader.GetValue(2)),
        ConvertToInt(reader.GetValue(3)),
        ConvertToInt(reader.GetValue(4)),
        ConvertToInt(reader.GetValue(5)),
        ConvertToInt(reader.GetValue(6))),
      new RecoverySeverityCounts(
        ConvertToInt(reader.GetValue(7)),
        ConvertToInt(reader.GetValue(8))),
      ConvertToInt(reader.GetValue(9)),
      ConvertToInt(reader.GetValue(10)));
  }

  private static async Task<IReadOnlyList<RecoveryCaseResponse>> ListRecoveryCases(NpgsqlConnection connection, string tenantSlug, string status, string severity)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        recovery.public_id::text,
        tenant.slug,
        recovery.subscription_public_id::text,
        recovery.invoice_public_id::text,
        COALESCE(recovery.source_attempt_public_id::text, ''),
        recovery.workflow_definition_key,
        recovery.status,
        recovery.severity,
        recovery.contact_channel,
        recovery.provider_hint,
        recovery.last_failed_attempt_number,
        COALESCE(to_char(recovery.next_action_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(recovery.promised_payment_date, 'YYYY-MM-DD'), ''),
        COALESCE(to_char(recovery.resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        recovery.resolution_code,
        recovery.notes,
        to_char(recovery.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(recovery.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.recovery_cases AS recovery
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = recovery.tenant_id
      WHERE tenant.slug = $1
        AND ($2 = '' OR recovery.status = $2)
        AND ($3 = '' OR recovery.severity = $3)
      ORDER BY recovery.created_at, recovery.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(status);
    command.Parameters.AddWithValue(severity);

    var response = new List<RecoveryCaseResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new RecoveryCaseResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetInt32(10),
        reader.GetString(11),
        reader.GetString(12),
        reader.GetString(13),
        reader.GetString(14),
        reader.GetString(15),
        reader.GetString(16),
        reader.GetString(17)));
    }

    return response;
  }

  private static async Task<RecoveryCaseResponse?> FindRecoveryCase(NpgsqlConnection connection, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        recovery.public_id::text,
        tenant.slug,
        recovery.subscription_public_id::text,
        recovery.invoice_public_id::text,
        COALESCE(recovery.source_attempt_public_id::text, ''),
        recovery.workflow_definition_key,
        recovery.status,
        recovery.severity,
        recovery.contact_channel,
        recovery.provider_hint,
        recovery.last_failed_attempt_number,
        COALESCE(to_char(recovery.next_action_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        COALESCE(to_char(recovery.promised_payment_date, 'YYYY-MM-DD'), ''),
        COALESCE(to_char(recovery.resolved_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), ''),
        recovery.resolution_code,
        recovery.notes,
        to_char(recovery.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
        to_char(recovery.updated_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.recovery_cases AS recovery
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = recovery.tenant_id
      WHERE tenant.slug = $1
        AND recovery.public_id = $2::uuid
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

    return new RecoveryCaseResponse(
      reader.GetString(0),
      reader.GetString(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetInt32(10),
      reader.GetString(11),
      reader.GetString(12),
      reader.GetString(13),
      reader.GetString(14),
      reader.GetString(15),
      reader.GetString(16),
      reader.GetString(17));
  }

  private static async Task<InternalRecoveryCase?> FindRecoveryCaseInternal(NpgsqlConnection connection, NpgsqlTransaction transaction, string tenantSlug, string publicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        recovery.id,
        recovery.tenant_id,
        recovery.public_id::text,
        tenant.slug,
        recovery.subscription_public_id::text,
        recovery.invoice_public_id::text,
        COALESCE(recovery.source_attempt_public_id::text, ''),
        recovery.workflow_definition_key,
        recovery.status,
        recovery.severity,
        recovery.contact_channel,
        recovery.provider_hint,
        recovery.last_failed_attempt_number,
        recovery.next_action_at,
        recovery.promised_payment_date,
        recovery.resolved_at,
        recovery.resolution_code,
        recovery.notes
      FROM billing.recovery_cases AS recovery
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = recovery.tenant_id
      WHERE tenant.slug = $1
        AND recovery.public_id = $2::uuid
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

    return new InternalRecoveryCase(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10),
      reader.GetString(11),
      reader.GetInt32(12),
      reader.IsDBNull(13) ? null : reader.GetFieldValue<DateTime>(13),
      reader.IsDBNull(14) ? null : reader.GetFieldValue<DateOnly>(14),
      reader.IsDBNull(15) ? null : reader.GetFieldValue<DateTime>(15),
      reader.GetString(16),
      reader.GetString(17));
  }

  private static async Task<InternalRecoveryCase?> FindRecoveryCaseByInvoice(NpgsqlConnection connection, NpgsqlTransaction transaction, long tenantId, string invoicePublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        recovery.id,
        recovery.tenant_id,
        recovery.public_id::text,
        tenant.slug,
        recovery.subscription_public_id::text,
        recovery.invoice_public_id::text,
        COALESCE(recovery.source_attempt_public_id::text, ''),
        recovery.workflow_definition_key,
        recovery.status,
        recovery.severity,
        recovery.contact_channel,
        recovery.provider_hint,
        recovery.last_failed_attempt_number,
        recovery.next_action_at,
        recovery.promised_payment_date,
        recovery.resolved_at,
        recovery.resolution_code,
        recovery.notes
      FROM billing.recovery_cases AS recovery
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = recovery.tenant_id
      WHERE recovery.tenant_id = $1
        AND recovery.invoice_public_id = $2::uuid
      LIMIT 1
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(tenantId);
    command.Parameters.AddWithValue(invoicePublicId);

    await using var reader = await command.ExecuteReaderAsync();
    if (!await reader.ReadAsync())
    {
      return null;
    }

    return new InternalRecoveryCase(
      reader.GetInt64(0),
      reader.GetInt64(1),
      reader.GetString(2),
      reader.GetString(3),
      reader.GetString(4),
      reader.GetString(5),
      reader.GetString(6),
      reader.GetString(7),
      reader.GetString(8),
      reader.GetString(9),
      reader.GetString(10),
      reader.GetString(11),
      reader.GetInt32(12),
      reader.IsDBNull(13) ? null : reader.GetFieldValue<DateTime>(13),
      reader.IsDBNull(14) ? null : reader.GetFieldValue<DateOnly>(14),
      reader.IsDBNull(15) ? null : reader.GetFieldValue<DateTime>(15),
      reader.GetString(16),
      reader.GetString(17));
  }

  private static async Task<IReadOnlyList<RecoveryActionResponse>> ListRecoveryActions(NpgsqlConnection connection, string tenantSlug, string recoveryCasePublicId)
  {
    await using var command = new NpgsqlCommand(
      """
      SELECT
        action.public_id::text,
        tenant.slug,
        recovery.public_id::text,
        action.action_code,
        action.actor,
        action.channel,
        action.provider,
        COALESCE(action.touchpoint_public_id::text, ''),
        COALESCE(action.delivery_public_id::text, ''),
        COALESCE(to_char(action.promised_payment_date, 'YYYY-MM-DD'), ''),
        action.notes,
        to_char(action.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
      FROM billing.recovery_actions AS action
      INNER JOIN billing.recovery_cases AS recovery
        ON recovery.id = action.recovery_case_id
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = action.tenant_id
      WHERE tenant.slug = $1
        AND recovery.public_id = $2::uuid
      ORDER BY action.created_at, action.id
      """,
      connection);
    command.Parameters.AddWithValue(tenantSlug);
    command.Parameters.AddWithValue(recoveryCasePublicId);

    var response = new List<RecoveryActionResponse>();
    await using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
      response.Add(new RecoveryActionResponse(
        reader.GetString(0),
        reader.GetString(1),
        reader.GetString(2),
        reader.GetString(3),
        reader.GetString(4),
        reader.GetString(5),
        reader.GetString(6),
        reader.GetString(7),
        reader.GetString(8),
        reader.GetString(9),
        reader.GetString(10),
        reader.GetString(11)));
    }

    return response;
  }

  private static async Task<RecoveryCaseResponse> OpenOrRefreshRecoveryCase(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalInvoice invoice,
    int failedAttemptNumber,
    string severity,
    DateTime? nextActionAt,
    string workflowDefinitionKey,
    string notes,
    string actor)
  {
    var existing = await FindRecoveryCaseByInvoice(connection, transaction, invoice.TenantId, invoice.PublicId);
    var normalizedNextActionAt = nextActionAt?.ToUniversalTime() ?? DateTime.UtcNow.AddHours(severity == "critical" ? 6 : 24);
    var summary = notes.Length > 0 ? notes : "Billing recovery case is tracking invoice collection risk.";

    if (existing is null)
    {
      string publicId;
      await using (var command = new NpgsqlCommand(
        """
        INSERT INTO billing.recovery_cases (
          tenant_id,
          public_id,
          subscription_public_id,
          invoice_public_id,
          workflow_definition_key,
          status,
          severity,
          last_failed_attempt_number,
          next_action_at,
          notes
        )
        VALUES (
          $1,
          gen_random_uuid(),
          $2::uuid,
          $3::uuid,
          $4,
          'open',
          $5,
          $6,
          $7,
          $8
        )
        RETURNING public_id::text
        """,
        connection,
        transaction))
      {
        command.Parameters.AddWithValue(invoice.TenantId);
        command.Parameters.AddWithValue(invoice.SubscriptionPublicId);
        command.Parameters.AddWithValue(invoice.PublicId);
        command.Parameters.AddWithValue(workflowDefinitionKey);
        command.Parameters.AddWithValue(severity);
        command.Parameters.AddWithValue(failedAttemptNumber);
        command.Parameters.AddWithValue(normalizedNextActionAt);
        command.Parameters.AddWithValue(summary);
        publicId = Convert.ToString(await command.ExecuteScalarAsync(), CultureInfo.InvariantCulture) ?? string.Empty;
      }

      var created = await FindRecoveryCaseInternal(connection, transaction, invoice.TenantSlug, publicId)
        ?? throw new InvalidOperationException("billing_recovery_case_creation_failed");
      await InsertRecoveryAction(connection, transaction, created, "case_opened", actor, string.Empty, string.Empty, string.Empty, string.Empty, null, summary);
      return await FindRecoveryCase(connection, invoice.TenantSlug, publicId)
        ?? throw new InvalidOperationException("billing_recovery_case_lookup_failed");
    }

    await using (var command = new NpgsqlCommand(
      """
      UPDATE billing.recovery_cases
      SET
        workflow_definition_key = CASE WHEN $3 = '' THEN workflow_definition_key ELSE $3 END,
        status = CASE WHEN status IN ('recovered', 'defaulted', 'closed') THEN 'open' ELSE status END,
        severity = $4,
        last_failed_attempt_number = GREATEST(last_failed_attempt_number, $5),
        next_action_at = $6,
        resolved_at = NULL,
        resolution_code = '',
        notes = CASE WHEN $7 = '' THEN notes ELSE $7 END
      WHERE id = $1
        AND tenant_id = $2
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(existing.Id);
      command.Parameters.AddWithValue(existing.TenantId);
      command.Parameters.AddWithValue(workflowDefinitionKey);
      command.Parameters.AddWithValue(severity);
      command.Parameters.AddWithValue(failedAttemptNumber);
      command.Parameters.AddWithValue(normalizedNextActionAt);
      command.Parameters.AddWithValue(summary);
      await command.ExecuteNonQueryAsync();
    }

    var refreshed = await FindRecoveryCaseInternal(connection, transaction, invoice.TenantSlug, existing.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_refresh_failed");
    await InsertRecoveryAction(connection, transaction, refreshed, "case_refreshed", actor, string.Empty, string.Empty, string.Empty, string.Empty, null, summary);
    return await FindRecoveryCase(connection, invoice.TenantSlug, existing.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_lookup_failed");
  }

  private static async Task<RecoveryCaseResponse> RegisterRecoveryTouchpoint(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalRecoveryCase recoveryCase,
    string actor,
    string channel,
    string provider,
    string touchpointPublicId,
    string deliveryPublicId,
    string notes)
  {
    await using (var command = new NpgsqlCommand(
      """
      UPDATE billing.recovery_cases
      SET
        status = CASE WHEN status IN ('open', 'promised') THEN 'contacted' ELSE status END,
        contact_channel = $3,
        provider_hint = $4,
        next_action_at = NULL,
        notes = CASE WHEN $5 = '' THEN notes ELSE $5 END
      WHERE id = $1
        AND tenant_id = $2
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(recoveryCase.Id);
      command.Parameters.AddWithValue(recoveryCase.TenantId);
      command.Parameters.AddWithValue(channel);
      command.Parameters.AddWithValue(provider);
      command.Parameters.AddWithValue(notes);
      await command.ExecuteNonQueryAsync();
    }

    var refreshed = await FindRecoveryCaseInternal(connection, transaction, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_touchpoint_failed");
    await InsertRecoveryAction(connection, transaction, refreshed, "touchpoint_registered", actor, channel, provider, touchpointPublicId, deliveryPublicId, null, notes);
    return await FindRecoveryCase(connection, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_lookup_failed");
  }

  private static async Task<RecoveryCaseResponse> RegisterRecoveryPromise(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalRecoveryCase recoveryCase,
    string actor,
    DateOnly promisedPaymentDate,
    DateTime? nextActionAt,
    string notes)
  {
    await using (var command = new NpgsqlCommand(
      """
      UPDATE billing.recovery_cases
      SET
        status = 'promised',
        promised_payment_date = $3,
        next_action_at = $4,
        notes = CASE WHEN $5 = '' THEN notes ELSE $5 END
      WHERE id = $1
        AND tenant_id = $2
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(recoveryCase.Id);
      command.Parameters.AddWithValue(recoveryCase.TenantId);
      command.Parameters.AddWithValue(promisedPaymentDate);
      command.Parameters.AddWithValue(nextActionAt?.ToUniversalTime() ?? DateTime.UtcNow.AddDays(1));
      command.Parameters.AddWithValue(notes);
      await command.ExecuteNonQueryAsync();
    }

    var refreshed = await FindRecoveryCaseInternal(connection, transaction, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_promise_failed");
    await InsertRecoveryAction(connection, transaction, refreshed, "promise_registered", actor, string.Empty, string.Empty, string.Empty, string.Empty, promisedPaymentDate, notes);
    return await FindRecoveryCase(connection, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_lookup_failed");
  }

  private static async Task<RecoveryCaseResponse> CloseRecoveryCase(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalRecoveryCase recoveryCase,
    string status,
    string resolutionCode,
    string actor,
    string notes,
    DateTime resolvedAt)
  {
    await using (var command = new NpgsqlCommand(
      """
      UPDATE billing.recovery_cases
      SET
        status = $3,
        resolved_at = $4,
        resolution_code = $5,
        next_action_at = NULL,
        notes = CASE WHEN $6 = '' THEN notes ELSE $6 END
      WHERE id = $1
        AND tenant_id = $2
      """,
      connection,
      transaction))
    {
      command.Parameters.AddWithValue(recoveryCase.Id);
      command.Parameters.AddWithValue(recoveryCase.TenantId);
      command.Parameters.AddWithValue(status);
      command.Parameters.AddWithValue(resolvedAt.ToUniversalTime());
      command.Parameters.AddWithValue(resolutionCode);
      command.Parameters.AddWithValue(notes);
      await command.ExecuteNonQueryAsync();
    }

    var refreshed = await FindRecoveryCaseInternal(connection, transaction, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_close_failed");
    await InsertRecoveryAction(connection, transaction, refreshed, status == "recovered" ? "case_recovered" : "case_closed", actor, string.Empty, string.Empty, string.Empty, string.Empty, null, notes.Length > 0 ? notes : resolutionCode);
    return await FindRecoveryCase(connection, recoveryCase.TenantSlug, recoveryCase.PublicId)
      ?? throw new InvalidOperationException("billing_recovery_case_lookup_failed");
  }

  private static async Task InsertRecoveryAction(
    NpgsqlConnection connection,
    NpgsqlTransaction transaction,
    InternalRecoveryCase recoveryCase,
    string actionCode,
    string actor,
    string channel,
    string provider,
    string touchpointPublicId,
    string deliveryPublicId,
    DateOnly? promisedPaymentDate,
    string notes)
  {
    await using var command = new NpgsqlCommand(
      """
      INSERT INTO billing.recovery_actions (
        tenant_id,
        recovery_case_id,
        public_id,
        action_code,
        actor,
        channel,
        provider,
        touchpoint_public_id,
        delivery_public_id,
        promised_payment_date,
        notes
      )
      VALUES (
        $1,
        $2,
        gen_random_uuid(),
        $3,
        $4,
        $5,
        $6,
        CASE WHEN $7 = '' THEN NULL ELSE $7::uuid END,
        CASE WHEN $8 = '' THEN NULL ELSE $8::uuid END,
        $9,
        $10
      )
      """,
      connection,
      transaction);
    command.Parameters.AddWithValue(recoveryCase.TenantId);
    command.Parameters.AddWithValue(recoveryCase.Id);
    command.Parameters.AddWithValue(actionCode);
    command.Parameters.AddWithValue(actor);
    command.Parameters.AddWithValue(channel);
    command.Parameters.AddWithValue(provider);
    command.Parameters.AddWithValue(touchpointPublicId);
    command.Parameters.AddWithValue(deliveryPublicId);
    command.Parameters.AddWithValue(promisedPaymentDate.HasValue ? promisedPaymentDate.Value : DBNull.Value);
    command.Parameters.AddWithValue(notes);
    await command.ExecuteNonQueryAsync();
  }

  private static bool? ParseActiveFilter(string? value)
  {
    if (string.IsNullOrWhiteSpace(value))
    {
      return null;
    }

    return value.Trim().ToLowerInvariant() switch
    {
      "true" or "1" or "yes" => true,
      "false" or "0" or "no" => false,
      _ => null
    };
  }

  private static string NormalizeIntervalUnit(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "monthly" => "monthly",
      "yearly" => "yearly",
      _ => string.Empty
    };

  private static string NormalizeSubscriptionStatus(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "active" => "active",
      "grace_period" => "grace_period",
      "suspended" => "suspended",
      "cancelled" => "cancelled",
      _ => string.Empty
    };

  private static string NormalizeInvoiceStatus(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "draft" => "draft",
      "open" => "open",
      "paid" => "paid",
      "failed" => "failed",
      "void" => "void",
      _ => string.Empty
    };

  private static string NormalizeAttemptStatus(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "succeeded" => "succeeded",
      "failed" => "failed",
      _ => string.Empty
    };

  private static string NormalizeRecoveryStatus(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "open" => "open",
      "contacted" => "contacted",
      "promised" => "promised",
      "recovered" => "recovered",
      "defaulted" => "defaulted",
      "closed" => "closed",
      _ => string.Empty
    };

  private static string NormalizeRecoverySeverity(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "attention" => "attention",
      "critical" => "critical",
      _ => string.Empty
    };

  private static string NormalizeRecoveryTerminalStatus(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "recovered" => "recovered",
      "defaulted" => "defaulted",
      "closed" => "closed",
      _ => string.Empty
    };

  private static string NormalizeRecoveryChannel(string? value)
    => value?.Trim().ToLowerInvariant() switch
    {
      "email" => "email",
      "whatsapp" => "whatsapp",
      "phone" => "phone",
      "manual" => "manual",
      _ => string.Empty
    };

  private static bool TryParseDate(string? value, out DateOnly? parsed)
  {
    if (string.IsNullOrWhiteSpace(value))
    {
      parsed = null;
      return true;
    }

    if (DateOnly.TryParseExact(value.Trim(), "yyyy-MM-dd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var result))
    {
      parsed = result;
      return true;
    }

    parsed = null;
    return false;
  }

  private static bool TryParseUtcTimestamp(string? value, out DateTime? parsed)
  {
    if (string.IsNullOrWhiteSpace(value))
    {
      parsed = null;
      return true;
    }

    if (DateTime.TryParse(
      value.Trim(),
      CultureInfo.InvariantCulture,
      DateTimeStyles.AdjustToUniversal | DateTimeStyles.AssumeUniversal,
      out var result))
    {
      parsed = result.ToUniversalTime();
      return true;
    }

    parsed = null;
    return false;
  }

  private static DateOnly ComputePeriodEnd(DateOnly start, string intervalUnit, int intervalCount)
    => intervalUnit switch
    {
      "yearly" => start.AddYears(intervalCount).AddDays(-1),
      _ => start.AddMonths(intervalCount).AddDays(-1)
    };

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
public sealed record DependencyStatus(string Name, string Status);
public sealed record ReadinessResponse(string Service, string Status, IReadOnlyList<DependencyStatus> Dependencies);
public sealed record ErrorResponse(string Code, string Message);
public sealed record PlanResponse(
  string PublicId,
  string Code,
  string Name,
  string Description,
  long AmountCents,
  string CurrencyCode,
  string IntervalUnit,
  int IntervalCount,
  int GracePeriodDays,
  int MaxRetries,
  bool Active,
  string CreatedAt,
  string UpdatedAt);
public sealed record SubscriptionResponse(
  string PublicId,
  string TenantSlug,
  string PlanPublicId,
  string PlanCode,
  string ExternalReference,
  string Status,
  string CurrentPeriodStart,
  string CurrentPeriodEnd,
  string GraceEndsAt,
  string SuspendedAt,
  string CancelledAt,
  string CreatedAt,
  string UpdatedAt);
public sealed record InvoiceResponse(
  string PublicId,
  string TenantSlug,
  string SubscriptionPublicId,
  string Number,
  string Status,
  long AmountCents,
  string DueDate,
  int RetryCount,
  string PaidAt,
  string GatewayReference,
  string CreatedAt,
  string UpdatedAt);
public sealed record PaymentAttemptResponse(
  string PublicId,
  string InvoicePublicId,
  string TenantSlug,
  int AttemptNumber,
  string Provider,
  string Status,
  string IdempotencyKey,
  string ExternalReference,
  string FailureReason,
  string AttemptedAt);
public sealed record AttemptOutcomeResponse(PaymentAttemptResponse Attempt, InvoiceResponse Invoice, SubscriptionResponse Subscription, bool Idempotent);
public sealed record SubscriptionEventResponse(
  string PublicId,
  string TenantSlug,
  string SubscriptionPublicId,
  string InvoicePublicId,
  string EventCode,
  string Actor,
  string Summary,
  string Payload,
  string CreatedAt);
public sealed record WebhookProcessResponse(
  string WebhookEventPublicId,
  string EventType,
  string ExternalId,
  string OutcomeStatus,
  bool Idempotent,
  PaymentAttemptResponse Attempt,
  InvoiceResponse Invoice,
  SubscriptionResponse Subscription);
public sealed record BillingOperationsReportResponse(
  string TenantSlug,
  string GeneratedAt,
  PlanSummary Plans,
  SubscriptionSummary Subscriptions,
  InvoiceSummary Invoices,
  AttemptSummary Attempts,
  RecoverySummary Recovery);
public sealed record PlanSummary(int Total, int Active);
public sealed record SubscriptionSummary(int Total, SubscriptionStatusCounts Status);
public sealed record SubscriptionStatusCounts(int Active, int GracePeriod, int Suspended, int Cancelled);
public sealed record InvoiceSummary(int Total, long TotalAmountCents, InvoiceStatusCounts Status, InvoiceAmountBuckets Amounts);
public sealed record InvoiceStatusCounts(int Open, int Paid, int Failed, int Void);
public sealed record InvoiceAmountBuckets(long OpenAmountCents, long PaidAmountCents, long FailedAmountCents);
public sealed record AttemptSummary(int Total, int Succeeded, int Failed);
public sealed record RecoverySummary(int Total, RecoveryStatusCounts Status, RecoverySeverityCounts Severity, int DueActions, int PromisesDue);
public sealed record RecoveryStatusCounts(int Open, int Contacted, int Promised, int Recovered, int Defaulted, int Closed);
public sealed record RecoverySeverityCounts(int Attention, int Critical);
public sealed record RecoveryCaseResponse(
  string PublicId,
  string TenantSlug,
  string SubscriptionPublicId,
  string InvoicePublicId,
  string SourceAttemptPublicId,
  string WorkflowDefinitionKey,
  string Status,
  string Severity,
  string ContactChannel,
  string ProviderHint,
  int LastFailedAttemptNumber,
  string NextActionAt,
  string PromisedPaymentDate,
  string ResolvedAt,
  string ResolutionCode,
  string Notes,
  string CreatedAt,
  string UpdatedAt);
public sealed record RecoveryActionResponse(
  string PublicId,
  string TenantSlug,
  string RecoveryCasePublicId,
  string ActionCode,
  string Actor,
  string Channel,
  string Provider,
  string TouchpointPublicId,
  string DeliveryPublicId,
  string PromisedPaymentDate,
  string Notes,
  string CreatedAt);

public sealed record CreatePlanRequest(
  string? Code,
  string? Name,
  string? Description,
  long? AmountCents,
  string? IntervalUnit,
  int? IntervalCount,
  int? GracePeriodDays,
  int? MaxRetries);
public sealed record CreateSubscriptionRequest(string? TenantSlug, string? PlanPublicId, string? ExternalReference, string? StartedOn);
public sealed record SubscriptionStatusRequest(string? TenantSlug, string? Reason, string? EffectiveAt);
public sealed record CreateInvoiceRequest(string? TenantSlug, string? Number, long? AmountCents, string? DueDate);
public sealed record CreatePaymentAttemptRequest(
  string? TenantSlug,
  string? Provider,
  string? Status,
  string? IdempotencyKey,
  string? ExternalReference,
  string? FailureReason,
  string? AttemptedAt);
public sealed record ProcessWebhookEventRequest(string? TenantSlug, string? WebhookEventPublicId);
public sealed record ProcessWebhookBatchRequest(string? TenantSlug, string? Provider, int? Limit);
public sealed record PendingWebhookEventResponse(
  string PublicId,
  string Provider,
  string EventType,
  string ExternalId,
  string PayloadSummary,
  string Status,
  string ReceivedAt,
  string InvoicePublicId,
  string InvoiceStatus);
public sealed record WebhookBatchSkipResponse(string WebhookEventPublicId, string Code, string Message);
public sealed record WebhookBatchProcessResponse(
  string TenantSlug,
  int ProcessedCount,
  int SkippedCount,
  IReadOnlyList<WebhookProcessResponse> Processed,
  IReadOnlyList<WebhookBatchSkipResponse> Skipped);
public sealed record OpenRecoveryCaseRequest(string? TenantSlug, string? Actor, string? WorkflowDefinitionKey, string? NextActionAt, string? Notes);
public sealed record RegisterRecoveryTouchpointRequest(string? TenantSlug, string? Actor, string? Channel, string? Provider, string? TouchpointPublicId, string? DeliveryPublicId, string? Notes);
public sealed record RegisterRecoveryPromiseRequest(string? TenantSlug, string? Actor, string? PromisedPaymentDate, string? NextActionAt, string? Notes);
public sealed record CloseRecoveryCaseRequest(string? TenantSlug, string? Actor, string? Status, string? ResolutionCode, string? ResolvedAt, string? Notes);

internal sealed record InternalPlan(
  long Id,
  string PublicId,
  string Code,
  long AmountCents,
  string IntervalUnit,
  int IntervalCount,
  int GracePeriodDays,
  int MaxRetries,
  bool Active);
internal sealed record InternalSubscription(
  long Id,
  long TenantId,
  string PublicId,
  string TenantSlug,
  long PlanId,
  string PlanPublicId,
  string PlanCode,
  long PlanAmountCents,
  string PlanIntervalUnit,
  int PlanIntervalCount,
  int GracePeriodDays,
  int MaxRetries,
  string ExternalReference,
  string Status,
  DateOnly CurrentPeriodStart,
  DateOnly CurrentPeriodEnd);
internal sealed record InternalInvoice(
  long Id,
  long TenantId,
  string PublicId,
  string TenantSlug,
  long SubscriptionId,
  string SubscriptionPublicId,
  string SubscriptionStatus,
  int GracePeriodDays,
  int MaxRetries,
  int RetryCount,
  string Status,
  string Number);
internal sealed record InternalRecoveryCase(
  long Id,
  long TenantId,
  string PublicId,
  string TenantSlug,
  string SubscriptionPublicId,
  string InvoicePublicId,
  string SourceAttemptPublicId,
  string WorkflowDefinitionKey,
  string Status,
  string Severity,
  string ContactChannel,
  string ProviderHint,
  int LastFailedAttemptNumber,
  DateTime? NextActionAt,
  DateOnly? PromisedPaymentDate,
  DateTime? ResolvedAt,
  string ResolutionCode,
  string Notes);
internal sealed record WebhookEvent(
  string PublicId,
  string Provider,
  string EventType,
  string ExternalId,
  string PayloadSummary,
  string Status,
  DateTime ReceivedAt);
internal sealed record WebhookBillingProcessResult(
  string OutcomeStatus,
  string ErrorCode,
  string ErrorMessage,
  AttemptOutcomeResponse AttemptOutcome)
{
  public static WebhookBillingProcessResult Success(string outcomeStatus, AttemptOutcomeResponse attemptOutcome)
  {
    return new WebhookBillingProcessResult(outcomeStatus, string.Empty, string.Empty, attemptOutcome);
  }

  public static WebhookBillingProcessResult Error(string errorCode, string errorMessage)
  {
    return new WebhookBillingProcessResult(
      string.Empty,
      errorCode,
      errorMessage,
      new AttemptOutcomeResponse(
        new PaymentAttemptResponse(string.Empty, string.Empty, string.Empty, 0, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty),
        new InvoiceResponse(string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, 0, string.Empty, 0, string.Empty, string.Empty, string.Empty, string.Empty),
        new SubscriptionResponse(string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty, string.Empty),
        true));
  }
}
