import { randomUUID } from "node:crypto";
import { Campaign, CampaignFilters, CampaignStatus, CreateCampaignInput } from "../domain/campaign.js";
import { CampaignRepository } from "../domain/campaign-repository.js";

function bootstrapCampaigns(tenantSlug: string): Campaign[] {
  return [
    {
      id: 1,
      publicId: "00000000-0000-0000-0000-00000000c101",
      tenantSlug,
      key: "lead-follow-up-campaign",
      name: "Lead Follow-Up Campaign",
      description: "Sequencia operacional para os primeiros contatos comerciais do lead capturado.",
      channel: "whatsapp",
      status: "active",
      touchpointGoal: "book-meeting",
      workflowDefinitionKey: "lead-follow-up",
      budgetCents: 95000,
      createdAt: "2026-04-20T00:00:00.000Z",
      updatedAt: "2026-04-20T00:00:00.000Z"
    },
    {
      id: 2,
      publicId: "00000000-0000-0000-0000-00000000c102",
      tenantSlug,
      key: "proposal-nurture-email",
      name: "Proposal Nurture Email",
      description: "Mantem follow-up de propostas abertas por email com cadencia leve.",
      channel: "email",
      status: "paused",
      touchpointGoal: "proposal-reminder",
      workflowDefinitionKey: "proposal-reminder",
      budgetCents: 35000,
      createdAt: "2026-04-20T00:05:00.000Z",
      updatedAt: "2026-04-20T00:05:00.000Z"
    }
  ];
}

export class InMemoryCampaignRepository implements CampaignRepository {
  private readonly campaigns: Campaign[];

  constructor(tenantSlug: string) {
    this.campaigns = bootstrapCampaigns(tenantSlug);
  }

  async list(filters: CampaignFilters = {}): Promise<Campaign[]> {
    return this.campaigns.filter((campaign) => {
      if (filters.tenantSlug && campaign.tenantSlug !== filters.tenantSlug) {
        return false;
      }

      if (filters.status && campaign.status !== filters.status) {
        return false;
      }

      if (filters.channel && campaign.channel !== filters.channel) {
        return false;
      }

      if (filters.q) {
        const query = filters.q.toLowerCase();
        const haystack = `${campaign.key} ${campaign.name} ${campaign.description}`.toLowerCase();
        return haystack.includes(query);
      }

      return true;
    });
  }

  async getByPublicId(publicId: string): Promise<Campaign | null> {
    return this.campaigns.find((campaign) => campaign.publicId === publicId) ?? null;
  }

  async create(input: CreateCampaignInput): Promise<Campaign> {
    if (this.campaigns.some((campaign) => campaign.tenantSlug === input.tenantSlug && campaign.key === input.key)) {
      throw new Error("campaign_key_conflict");
    }

    const now = new Date().toISOString();
    const campaign: Campaign = {
      id: this.campaigns.length + 1,
      publicId: randomUUID(),
      tenantSlug: input.tenantSlug,
      key: input.key,
      name: input.name,
      description: input.description,
      channel: input.channel,
      status: input.status ?? "draft",
      touchpointGoal: input.touchpointGoal,
      workflowDefinitionKey: input.workflowDefinitionKey ?? null,
      budgetCents: input.budgetCents,
      createdAt: now,
      updatedAt: now
    };

    this.campaigns.push(campaign);
    return campaign;
  }

  async updateStatus(publicId: string, status: CampaignStatus): Promise<Campaign | null> {
    const campaign = this.campaigns.find((item) => item.publicId === publicId);

    if (!campaign) {
      return null;
    }

    campaign.status = status;
    campaign.updatedAt = new Date().toISOString();
    return campaign;
  }
}
