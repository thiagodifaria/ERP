import { IncomingMessage, ServerResponse } from "node:http";
import { CreateCampaignRequest } from "./dto/create-campaign-request.js";
import { CreateTemplateRequest } from "./dto/create-template-request.js";
import { CreateTouchpointRequest } from "./dto/create-touchpoint-request.js";
import { CreateTouchpointDeliveryRequest } from "./dto/create-touchpoint-delivery-request.js";
import { HealthResponse, ReadinessResponse } from "./dto/health.js";
import { UpdateCampaignStatusRequest } from "./dto/update-campaign-status-request.js";
import { UpdateTemplateStatusRequest } from "./dto/update-template-status-request.js";
import { UpdateTouchpointDeliveryStatusRequest } from "./dto/update-touchpoint-delivery-status-request.js";
import { UpdateTouchpointStatusRequest } from "./dto/update-touchpoint-status-request.js";
import { runtime, services } from "../config/container.js";
import { CampaignChannel, CampaignStatus } from "../domain/campaign.js";
import { DeliveryStatus } from "../domain/delivery.js";
import { TemplateProvider, TemplateStatus } from "../domain/template.js";
import { TouchpointStatus } from "../domain/touchpoint.js";

function json(response: ServerResponse, statusCode: number, body: unknown): void {
  response.writeHead(statusCode, { "content-type": "application/json" });
  response.end(JSON.stringify(body));
}

function live(): HealthResponse {
  return {
    service: runtime.config.serviceName,
    status: "live"
  };
}

function ready(): HealthResponse {
  return {
    service: runtime.config.serviceName,
    status: "ready"
  };
}

async function details(): Promise<ReadinessResponse> {
  return {
    service: runtime.config.serviceName,
    status: "ready",
    dependencies: await runtime.readinessDependencies()
  };
}

function pathSegments(request: IncomingMessage): string[] {
  const pathname = request.url?.split("?")[0] ?? "/";
  return pathname.split("/").filter((segment) => segment.length > 0);
}

function searchParams(request: IncomingMessage): URLSearchParams {
  return new URL(request.url ?? "/", "http://engagement.local").searchParams;
}

async function readJson<T>(request: IncomingMessage): Promise<T> {
  const chunks: Buffer[] = [];

  for await (const chunk of request) {
    chunks.push(Buffer.from(chunk));
  }

  const rawBody = Buffer.concat(chunks).toString("utf8");

  if (rawBody.length === 0) {
    throw new Error("invalid_json");
  }

  return JSON.parse(rawBody) as T;
}

function handleDomainError(response: ServerResponse, error: unknown): void {
  const code = error instanceof Error ? error.message : "unexpected_error";

  if (
    code === "tenant_slug_required" ||
    code === "campaign_tenant_not_found" ||
    code === "template_tenant_not_found" ||
    code === "delivery_tenant_not_found" ||
    code === "campaign_key_invalid" ||
    code === "campaign_name_required" ||
    code === "campaign_description_required" ||
    code === "campaign_channel_invalid" ||
    code === "campaign_status_invalid" ||
    code === "campaign_touchpoint_goal_required" ||
    code === "campaign_budget_invalid" ||
    code === "campaign_public_id_invalid" ||
    code === "lead_public_id_invalid" ||
    code === "touchpoint_public_id_invalid" ||
    code === "touchpoint_contact_value_required" ||
    code === "touchpoint_source_required" ||
    code === "touchpoint_created_by_required" ||
    code === "touchpoint_status_invalid" ||
    code === "touchpoint_workflow_run_public_id_invalid" ||
    code === "template_public_id_invalid" ||
    code === "template_key_invalid" ||
    code === "template_name_required" ||
    code === "template_channel_invalid" ||
    code === "template_status_invalid" ||
    code === "template_provider_invalid" ||
    code === "template_body_required" ||
    code === "delivery_status_invalid" ||
    code === "delivery_provider_invalid" ||
    code === "delivery_channel_mismatch" ||
    code === "delivery_public_id_invalid" ||
    code === "delivery_sent_by_required" ||
    code === "invalid_json"
  ) {
    json(response, 400, { code, message: "Request payload is invalid." });
    return;
  }

  if (code === "campaign_key_conflict") {
    json(response, 409, { code, message: "Campaign key already exists for this tenant." });
    return;
  }

  if (code === "template_key_conflict") {
    json(response, 409, { code, message: "Template key already exists for this tenant." });
    return;
  }

  json(response, 500, {
    code: "unexpected_error",
    message: "Unexpected error."
  });
}

export async function route(request: IncomingMessage, response: ServerResponse): Promise<void> {
  if (request.url === "/health/live") {
    json(response, 200, live());
    return;
  }

  if (request.url === "/health/ready") {
    json(response, 200, ready());
    return;
  }

  if (request.url === "/health/details") {
    json(response, 200, await details());
    return;
  }

  if (request.method === "GET" && request.url?.startsWith("/api/engagement/campaigns")) {
    const segments = pathSegments(request);

    if (segments.length === 3) {
      const params = searchParams(request);
      json(
        response,
        200,
        await services.listCampaigns.execute({
          tenantSlug: params.get("tenantSlug") ?? undefined,
          status: (params.get("status") ?? undefined) as CampaignStatus | undefined,
          channel: (params.get("channel") ?? undefined) as CampaignChannel | undefined,
          q: params.get("q") ?? undefined
        })
      );
      return;
    }
  }

  if (request.method === "GET" && request.url?.startsWith("/api/engagement/templates")) {
    const segments = pathSegments(request);

    if (segments.length === 3) {
      const params = searchParams(request);
      json(
        response,
        200,
        await services.listTemplates.execute({
          tenantSlug: params.get("tenantSlug") ?? undefined,
          status: (params.get("status") ?? undefined) as TemplateStatus | undefined,
          channel: (params.get("channel") ?? undefined) as CampaignChannel | undefined,
          provider: (params.get("provider") ?? undefined) as TemplateProvider | undefined,
          q: params.get("q") ?? undefined
        })
      );
      return;
    }
  }

  if (request.method === "POST" && request.url === "/api/engagement/campaigns") {
    try {
      const payload = await readJson<CreateCampaignRequest>(request);
      json(response, 201, await services.createCampaign.execute(payload));
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (request.method === "POST" && request.url === "/api/engagement/templates") {
    try {
      const payload = await readJson<CreateTemplateRequest>(request);
      json(response, 201, await services.createTemplate.execute(payload));
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (request.method === "GET" && request.url?.startsWith("/api/engagement/touchpoints/summary")) {
    const params = searchParams(request);
    json(
      response,
      200,
      await services.getTouchpointSummary.execute({
        tenantSlug: params.get("tenantSlug") ?? undefined,
        campaignPublicId: params.get("campaignPublicId") ?? undefined,
        status: (params.get("status") ?? undefined) as TouchpointStatus | undefined,
        channel: (params.get("channel") ?? undefined) as CampaignChannel | undefined,
        leadPublicId: params.get("leadPublicId") ?? undefined
      })
    );
    return;
  }

  if (request.method === "GET" && request.url?.startsWith("/api/engagement/touchpoints")) {
    const segments = pathSegments(request);

    if (segments.length === 3) {
      const params = searchParams(request);
      json(
        response,
        200,
        await services.listTouchpoints.execute({
          tenantSlug: params.get("tenantSlug") ?? undefined,
          campaignPublicId: params.get("campaignPublicId") ?? undefined,
          status: (params.get("status") ?? undefined) as TouchpointStatus | undefined,
          channel: (params.get("channel") ?? undefined) as CampaignChannel | undefined,
          leadPublicId: params.get("leadPublicId") ?? undefined
        })
      );
      return;
    }
  }

  if (request.method === "GET" && request.url?.startsWith("/api/engagement/deliveries/summary")) {
    const params = searchParams(request);
    json(
      response,
      200,
      await services.getDeliverySummary.execute({
        tenantSlug: params.get("tenantSlug") ?? undefined,
        status: (params.get("status") ?? undefined) as DeliveryStatus | undefined,
        channel: (params.get("channel") ?? undefined) as CampaignChannel | undefined,
        provider: (params.get("provider") ?? undefined) as TemplateProvider | undefined
      })
    );
    return;
  }

  if (request.method === "POST" && request.url === "/api/engagement/touchpoints") {
    try {
      const payload = await readJson<CreateTouchpointRequest>(request);
      json(response, 201, await services.createTouchpoint.execute(payload));
      return;
    } catch (error) {
      if (error instanceof Error && error.message === "campaign_not_found") {
        json(response, 404, { code: error.message, message: "Campaign was not found." });
        return;
      }

      handleDomainError(response, error);
      return;
    }
  }

  const segments = pathSegments(request);

  if (
    request.method === "GET" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "templates"
  ) {
    const template = await services.getTemplateByPublicId.execute(segments[3]);

    if (template === null) {
      json(response, 404, { code: "template_not_found", message: "Template was not found." });
      return;
    }

    json(response, 200, template);
    return;
  }

  if (
    request.method === "PATCH" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "templates" &&
    segments[4] === "status"
  ) {
    try {
      const payload = await readJson<UpdateTemplateStatusRequest>(request);
      const template = await services.updateTemplateStatus.execute(segments[3], payload.status);

      if (template === null) {
        json(response, 404, { code: "template_not_found", message: "Template was not found." });
        return;
      }

      json(response, 200, template);
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "GET" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "campaigns"
  ) {
    const campaign = await services.getCampaignByPublicId.execute(segments[3]);

    if (campaign === null) {
      json(response, 404, { code: "campaign_not_found", message: "Campaign was not found." });
      return;
    }

    json(response, 200, campaign);
    return;
  }

  if (
    request.method === "GET" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "touchpoints" &&
    segments[4] === "deliveries"
  ) {
    try {
      const params = searchParams(request);
      json(response, 200, await services.listTouchpointDeliveries.execute(segments[3], params.get("tenantSlug") ?? undefined));
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "POST" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "touchpoints" &&
    segments[4] === "deliveries"
  ) {
    try {
      const payload = await readJson<CreateTouchpointDeliveryRequest>(request);
      json(response, 201, await services.createTouchpointDelivery.execute(segments[3], payload));
      return;
    } catch (error) {
      if (error instanceof Error && error.message === "touchpoint_not_found") {
        json(response, 404, { code: error.message, message: "Touchpoint was not found." });
        return;
      }
      if (error instanceof Error && error.message === "template_not_found") {
        json(response, 404, { code: error.message, message: "Template was not found." });
        return;
      }

      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "PATCH" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "campaigns" &&
    segments[4] === "status"
  ) {
    try {
      const payload = await readJson<UpdateCampaignStatusRequest>(request);
      const campaign = await services.updateCampaignStatus.execute(segments[3], payload.status);

      if (campaign === null) {
        json(response, 404, { code: "campaign_not_found", message: "Campaign was not found." });
        return;
      }

      json(response, 200, campaign);
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "PATCH" &&
    segments.length === 7 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "touchpoints" &&
    segments[4] === "deliveries" &&
    segments[6] === "status"
  ) {
    try {
      const payload = await readJson<UpdateTouchpointDeliveryStatusRequest>(request);
      const delivery = await services.updateTouchpointDeliveryStatus.execute(segments[5], payload);

      if (delivery === null) {
        json(response, 404, { code: "delivery_not_found", message: "Delivery was not found." });
        return;
      }

      json(response, 200, delivery);
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "GET" &&
    segments.length === 4 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "touchpoints"
  ) {
    try {
      const touchpoint = await services.getTouchpointByPublicId.execute(segments[3]);

      if (touchpoint === null) {
        json(response, 404, { code: "touchpoint_not_found", message: "Touchpoint was not found." });
        return;
      }

      json(response, 200, touchpoint);
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  if (
    request.method === "PATCH" &&
    segments.length === 5 &&
    segments[0] === "api" &&
    segments[1] === "engagement" &&
    segments[2] === "touchpoints" &&
    segments[4] === "status"
  ) {
    try {
      const payload = await readJson<UpdateTouchpointStatusRequest>(request);
      const touchpoint = await services.updateTouchpointStatus.execute(
        segments[3],
        payload.status,
        payload.lastWorkflowRunPublicId
      );

      if (touchpoint === null) {
        json(response, 404, { code: "touchpoint_not_found", message: "Touchpoint was not found." });
        return;
      }

      json(response, 200, touchpoint);
      return;
    } catch (error) {
      handleDomainError(response, error);
      return;
    }
  }

  json(response, 404, {
    code: "route_not_found",
    message: "Route was not found."
  });
}
