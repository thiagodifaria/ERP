import pg from "pg";
import { CreateCampaign } from "../application/create-campaign.js";
import { CreateTemplate } from "../application/create-template.js";
import { CreateTouchpoint } from "../application/create-touchpoint.js";
import { CreateTouchpointDelivery } from "../application/create-touchpoint-delivery.js";
import { DispatchWorkflowTouchpoint } from "../application/dispatch-workflow-touchpoint.js";
import { GetCampaignByPublicId } from "../application/get-campaign-by-public-id.js";
import { GetDeliverySummary } from "../application/get-delivery-summary.js";
import { GetProviderEventByPublicId } from "../application/get-provider-event-by-public-id.js";
import { GetProviderEventSummary } from "../application/get-provider-event-summary.js";
import { GetTemplateByPublicId } from "../application/get-template-by-public-id.js";
import { GetTouchpointByPublicId } from "../application/get-touchpoint-by-public-id.js";
import { GetTouchpointSummary } from "../application/get-touchpoint-summary.js";
import { IngestProviderLead } from "../application/ingest-provider-lead.js";
import { ListCampaigns } from "../application/list-campaigns.js";
import { ListProviderCapabilities } from "../application/list-provider-capabilities.js";
import { ListProviderEvents } from "../application/list-provider-events.js";
import { ListTemplates } from "../application/list-templates.js";
import { ListTouchpointDeliveries } from "../application/list-touchpoint-deliveries.js";
import { ListTouchpoints } from "../application/list-touchpoints.js";
import { RegisterProviderEvent } from "../application/register-provider-event.js";
import { UpdateCampaignStatus } from "../application/update-campaign-status.js";
import { UpdateTemplateStatus } from "../application/update-template-status.js";
import { UpdateTouchpointDeliveryStatus } from "../application/update-touchpoint-delivery-status.js";
import { UpdateTouchpointStatus } from "../application/update-touchpoint-status.js";
import { InMemoryDeliveryRepository } from "../infrastructure/in-memory-delivery-repository.js";
import { InMemoryCampaignRepository } from "../infrastructure/in-memory-campaign-repository.js";
import { InMemoryCrmGateway, HttpCrmGateway } from "../infrastructure/crm-gateway.js";
import { InMemoryTemplateRepository } from "../infrastructure/in-memory-template-repository.js";
import { InMemoryTouchpointRepository } from "../infrastructure/in-memory-touchpoint-repository.js";
import { InMemoryProviderEventRepository } from "../infrastructure/in-memory-provider-event-repository.js";
import { PostgresDeliveryRepository } from "../infrastructure/postgres-delivery-repository.js";
import { PostgresCampaignRepository } from "../infrastructure/postgres-campaign-repository.js";
import { PostgresProviderEventRepository } from "../infrastructure/postgres-provider-event-repository.js";
import { PostgresTemplateRepository } from "../infrastructure/postgres-template-repository.js";
import { PostgresTouchpointRepository } from "../infrastructure/postgres-touchpoint-repository.js";
import { loadConfig } from "./env.js";

const { Pool } = pg;

export type ReadinessDependency = {
  name: string;
  status: string;
};

function providerStatus(capability: { configured: boolean; mode: string; status: string }): string {
  if (capability.configured) {
    return "ready";
  }

  if (capability.mode === "fallback" || capability.status === "fallback") {
    return "fallback";
  }

  if (capability.mode === "manual" || capability.status === "manual") {
    return "manual";
  }

  return "unconfigured";
}

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
const templateRepository =
  config.repositoryDriver === "postgres" && pool !== null
    ? new PostgresTemplateRepository(pool, config.bootstrapTenantSlug)
    : new InMemoryTemplateRepository(config.bootstrapTenantSlug);
const deliveryRepository =
  config.repositoryDriver === "postgres" && pool !== null
    ? new PostgresDeliveryRepository(
        pool,
        config.bootstrapTenantSlug,
        (tenantSlug?: string) => templateRepository.list({ tenantSlug: tenantSlug ?? config.bootstrapTenantSlug }),
        (tenantSlug?: string) => touchpointRepository.list({ tenantSlug: tenantSlug ?? config.bootstrapTenantSlug })
      )
    : new InMemoryDeliveryRepository(
        config.bootstrapTenantSlug,
        () => templateRepository.list({ tenantSlug: config.bootstrapTenantSlug }),
        () => touchpointRepository.list({ tenantSlug: config.bootstrapTenantSlug })
      );
const providerEventRepository =
  config.repositoryDriver === "postgres" && pool !== null
    ? new PostgresProviderEventRepository(pool, config.bootstrapTenantSlug)
    : new InMemoryProviderEventRepository(config.bootstrapTenantSlug);
const crmGateway =
  config.repositoryDriver === "postgres"
    ? new HttpCrmGateway(config.crmBaseUrl)
    : new InMemoryCrmGateway();

export const repositories = {
  campaigns: campaignRepository,
  touchpoints: touchpointRepository,
  templates: templateRepository,
  deliveries: deliveryRepository,
  providerEvents: providerEventRepository
};

export const services = {
  listCampaigns: new ListCampaigns(campaignRepository),
  getCampaignByPublicId: new GetCampaignByPublicId(campaignRepository),
  createCampaign: new CreateCampaign(campaignRepository),
  updateCampaignStatus: new UpdateCampaignStatus(campaignRepository),
  listTemplates: new ListTemplates(templateRepository),
  getTemplateByPublicId: new GetTemplateByPublicId(templateRepository),
  createTemplate: new CreateTemplate(templateRepository),
  updateTemplateStatus: new UpdateTemplateStatus(templateRepository),
  listTouchpoints: new ListTouchpoints(touchpointRepository),
  getTouchpointByPublicId: new GetTouchpointByPublicId(touchpointRepository),
  createTouchpoint: new CreateTouchpoint(campaignRepository, touchpointRepository),
  updateTouchpointStatus: new UpdateTouchpointStatus(touchpointRepository),
  getTouchpointSummary: new GetTouchpointSummary(touchpointRepository),
  listTouchpointDeliveries: new ListTouchpointDeliveries(deliveryRepository),
  createTouchpointDelivery: new CreateTouchpointDelivery(deliveryRepository),
  updateTouchpointDeliveryStatus: new UpdateTouchpointDeliveryStatus(deliveryRepository),
  getDeliverySummary: new GetDeliverySummary(deliveryRepository),
  listProviderCapabilities: new ListProviderCapabilities({
    resendApiKey: config.resendApiKey,
    whatsappAccessToken: config.whatsappAccessToken,
    telegramBotToken: config.telegramBotToken,
    metaAdsAccessToken: config.metaAdsAccessToken
  }),
  listProviderEvents: new ListProviderEvents(providerEventRepository),
  getProviderEventByPublicId: new GetProviderEventByPublicId(providerEventRepository),
  getProviderEventSummary: new GetProviderEventSummary(providerEventRepository),
  ingestProviderLead: new IngestProviderLead(campaignRepository, touchpointRepository, providerEventRepository, crmGateway),
  dispatchWorkflowTouchpoint: new DispatchWorkflowTouchpoint(campaignRepository, templateRepository, touchpointRepository, deliveryRepository, providerEventRepository),
  registerProviderEvent: new RegisterProviderEvent(touchpointRepository, deliveryRepository, providerEventRepository)
};

export const runtime = {
  config,
  async readinessDependencies(): Promise<ReadinessDependency[]> {
    const providerCapabilities = await services.listProviderCapabilities.execute();
    const providerDependencies = providerCapabilities.map((capability) => ({
      name: `provider:${capability.provider}`,
      status: providerStatus(capability)
    }));

    if (config.repositoryDriver !== "postgres") {
      return [
        { name: "router", status: "ready" },
        { name: "campaign-catalog", status: "ready" },
        { name: "templates", status: "ready" },
        { name: "touchpoints", status: "ready" },
        { name: "deliveries", status: "ready" },
        { name: "provider-events", status: "ready" },
        { name: "crm-gateway", status: "ready" },
        ...providerDependencies
      ];
    }

    try {
      await campaignRepository.list({ tenantSlug: config.bootstrapTenantSlug });
      await templateRepository.list({ tenantSlug: config.bootstrapTenantSlug });
      await touchpointRepository.list({ tenantSlug: config.bootstrapTenantSlug });
      await deliveryRepository.getSummary({ tenantSlug: config.bootstrapTenantSlug });
      await providerEventRepository.getSummary({ tenantSlug: config.bootstrapTenantSlug });

      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "ready" },
        { name: "campaign-catalog", status: "ready" },
        { name: "templates", status: "ready" },
        { name: "touchpoints", status: "ready" },
        { name: "deliveries", status: "ready" },
        { name: "provider-events", status: "ready" },
        { name: "crm-gateway", status: "ready" },
        ...providerDependencies
      ];
    } catch {
      return [
        { name: "router", status: "ready" },
        { name: "postgresql", status: "not_ready" },
        { name: "campaign-catalog", status: "not_ready" },
        { name: "templates", status: "not_ready" },
        { name: "touchpoints", status: "not_ready" },
        { name: "deliveries", status: "not_ready" },
        { name: "provider-events", status: "not_ready" },
        { name: "crm-gateway", status: "not_ready" },
        ...providerDependencies
      ];
    }
  }
};
