import pg from "pg";
import { CreateCampaign } from "../application/create-campaign.js";
import { CreateTouchpoint } from "../application/create-touchpoint.js";
import { GetCampaignByPublicId } from "../application/get-campaign-by-public-id.js";
import { GetTouchpointByPublicId } from "../application/get-touchpoint-by-public-id.js";
import { GetTouchpointSummary } from "../application/get-touchpoint-summary.js";
import { ListCampaigns } from "../application/list-campaigns.js";
import { ListTouchpoints } from "../application/list-touchpoints.js";
import { UpdateCampaignStatus } from "../application/update-campaign-status.js";
import { UpdateTouchpointStatus } from "../application/update-touchpoint-status.js";
import { InMemoryCampaignRepository } from "../infrastructure/in-memory-campaign-repository.js";
import { InMemoryTouchpointRepository } from "../infrastructure/in-memory-touchpoint-repository.js";
import { PostgresCampaignRepository } from "../infrastructure/postgres-campaign-repository.js";
import { PostgresTouchpointRepository } from "../infrastructure/postgres-touchpoint-repository.js";
import { loadConfig } from "./env.js";

const { Pool } = pg;

export type ReadinessDependency = {
  name: string;
  status: string;
};

const config = loadConfig();
const pool =
  config.repositoryDriver === "postgres"
    ? new Pool({
        host: config.postgresHost,
        port: Number(config.postgresPort),
        database: config.postgresDatabase,
        user: config.postgresUser,
        password: config.postgresPassword,
        ssl: config.postgresSslMode === "disable" ? false : { rejectUnauthorized: false }
      })
    : null;
const campaignRepository =
  config.repositoryDriver === "postgres" && pool !== null
    ? new PostgresCampaignRepository(pool, config.bootstrapTenantSlug)
    : new InMemoryCampaignRepository(config.bootstrapTenantSlug);
const touchpointRepository =
  config.repositoryDriver === "postgres" && pool !== null
    ? new PostgresTouchpointRepository(pool, config.bootstrapTenantSlug)
    : new InMemoryTouchpointRepository(config.bootstrapTenantSlug);

export const repositories = {
  campaigns: campaignRepository,
  touchpoints: touchpointRepository
};

export const services = {
  listCampaigns: new ListCampaigns(campaignRepository),
  getCampaignByPublicId: new GetCampaignByPublicId(campaignRepository),
  createCampaign: new CreateCampaign(campaignRepository),
  updateCampaignStatus: new UpdateCampaignStatus(campaignRepository),
  listTouchpoints: new ListTouchpoints(touchpointRepository),
  getTouchpointByPublicId: new GetTouchpointByPublicId(touchpointRepository),
  createTouchpoint: new CreateTouchpoint(campaignRepository, touchpointRepository),
  updateTouchpointStatus: new UpdateTouchpointStatus(touchpointRepository),
  getTouchpointSummary: new GetTouchpointSummary(touchpointRepository)
};

export const runtime = {
  config,
  async readinessDependencies(): Promise<ReadinessDependency[]> {
    if (config.repositoryDriver !== "postgres") {
      return [
        { name: "router", status: "ready" },
        { name: "campaign-catalog", status: "ready" },
        { name: "touchpoints", status: "ready" }
      ];
    }

    try {
      await campaignRepository.list({ tenantSlug: config.bootstrapTenantSlug });
      await touchpointRepository.list({ tenantSlug: config.bootstrapTenantSlug });

      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "ready" },
        { name: "campaign-catalog", status: "ready" },
        { name: "touchpoints", status: "ready" }
      ];
    } catch {
      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "not_ready" },
        { name: "campaign-catalog", status: "not_ready" },
        { name: "touchpoints", status: "not_ready" }
      ];
    }
  }
};
