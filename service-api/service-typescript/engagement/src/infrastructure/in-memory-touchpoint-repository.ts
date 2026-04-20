import { randomUUID } from "node:crypto";
import { campaignChannels } from "../domain/campaign.js";
import {
  CreateTouchpointInput,
  Touchpoint,
  TouchpointFilters,
  TouchpointStatus,
  TouchpointSummary,
  touchpointStatuses
} from "../domain/touchpoint.js";
import { TouchpointRepository } from "../domain/touchpoint-repository.js";

function bootstrapTouchpoints(tenantSlug: string): Touchpoint[] {
  return [
    {
      id: 1,
      publicId: "00000000-0000-0000-0000-00000000e101",
      tenantSlug,
      campaignPublicId: "00000000-0000-0000-0000-00000000c101",
      campaignKey: "lead-follow-up-campaign",
      leadPublicId: "00000000-0000-0000-0000-000000008850",
      channel: "whatsapp",
      contactValue: "+5531999999999",
      source: "crm",
      status: "responded",
      workflowDefinitionKey: "lead-follow-up",
      lastWorkflowRunPublicId: "00000000-0000-0000-0000-000000000301",
      createdBy: "bootstrap-owner",
      notes: "Lead respondeu ao contato inicial.",
      createdAt: "2026-04-20T00:10:00.000Z",
      updatedAt: "2026-04-20T00:11:00.000Z"
    }
  ];
}

export class InMemoryTouchpointRepository implements TouchpointRepository {
  private readonly touchpoints: Touchpoint[];

  constructor(tenantSlug: string) {
    this.touchpoints = bootstrapTouchpoints(tenantSlug);
  }

  async list(filters: TouchpointFilters = {}): Promise<Touchpoint[]> {
    return this.touchpoints.filter((touchpoint) => {
      if (filters.tenantSlug && touchpoint.tenantSlug !== filters.tenantSlug) {
        return false;
      }

      if (filters.campaignPublicId && touchpoint.campaignPublicId !== filters.campaignPublicId) {
        return false;
      }

      if (filters.status && touchpoint.status !== filters.status) {
        return false;
      }

      if (filters.channel && touchpoint.channel !== filters.channel) {
        return false;
      }

      if (filters.leadPublicId && touchpoint.leadPublicId !== filters.leadPublicId) {
        return false;
      }

      return true;
    });
  }

  async getByPublicId(publicId: string): Promise<Touchpoint | null> {
    return this.touchpoints.find((touchpoint) => touchpoint.publicId === publicId) ?? null;
  }

  async create(
    input: CreateTouchpointInput & {
      campaignKey: string;
      channel: Touchpoint["channel"];
      workflowDefinitionKey: string | null;
    }
  ): Promise<Touchpoint> {
    const now = new Date().toISOString();
    const touchpoint: Touchpoint = {
      id: this.touchpoints.length + 1,
      publicId: randomUUID(),
      tenantSlug: input.tenantSlug,
      campaignPublicId: input.campaignPublicId,
      campaignKey: input.campaignKey,
      leadPublicId: input.leadPublicId,
      channel: input.channel,
      contactValue: input.contactValue,
      source: input.source,
      status: "queued",
      workflowDefinitionKey: input.workflowDefinitionKey,
      lastWorkflowRunPublicId: null,
      createdBy: input.createdBy,
      notes: input.notes,
      createdAt: now,
      updatedAt: now
    };

    this.touchpoints.push(touchpoint);
    return touchpoint;
  }

  async updateStatus(
    publicId: string,
    status: TouchpointStatus,
    lastWorkflowRunPublicId?: string | null
  ): Promise<Touchpoint | null> {
    const touchpoint = this.touchpoints.find((item) => item.publicId === publicId);

    if (!touchpoint) {
      return null;
    }

    touchpoint.status = status;
    touchpoint.lastWorkflowRunPublicId = lastWorkflowRunPublicId ?? touchpoint.lastWorkflowRunPublicId;
    touchpoint.updatedAt = new Date().toISOString();
    return touchpoint;
  }

  async getSummary(filters: TouchpointFilters = {}): Promise<TouchpointSummary> {
    const touchpoints = await this.list(filters);
    const campaignIds = new Set(touchpoints.map((touchpoint) => touchpoint.campaignPublicId));
    const byStatus = Object.fromEntries(touchpointStatuses.map((status) => [status, 0])) as Record<TouchpointStatus, number>;
    const byChannel = Object.fromEntries(campaignChannels.map((channel) => [channel, 0])) as Record<Touchpoint["channel"], number>;

    for (const touchpoint of touchpoints) {
      byStatus[touchpoint.status] += 1;
      byChannel[touchpoint.channel] += 1;
    }

    const responded = byStatus.responded + byStatus.converted;
    const converted = byStatus.converted;
    const failed = byStatus.failed;
    const sentBase = byStatus.sent + byStatus.delivered + byStatus.responded + byStatus.converted;
    const responseBase = byStatus.responded + byStatus.converted + byStatus.failed;

    return {
      tenantSlug: filters.tenantSlug ?? "global",
      generatedAt: new Date().toISOString(),
      totals: {
        campaigns: campaignIds.size,
        touchpoints: touchpoints.length,
        workflowConfigured: touchpoints.filter((touchpoint) => touchpoint.workflowDefinitionKey !== null).length,
        workflowDispatched: touchpoints.filter((touchpoint) => touchpoint.lastWorkflowRunPublicId !== null).length
      },
      byStatus,
      byChannel,
      outcomes: {
        responded,
        converted,
        failed,
        responseRate: sentBase > 0 ? Number((responded / sentBase).toFixed(4)) : 0,
        conversionRate: responseBase > 0 ? Number((converted / responseBase).toFixed(4)) : 0
      }
    };
  }
}
