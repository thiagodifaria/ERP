import type { EndpointContract } from "../generated/apiCatalog";

export type RuntimeEnvironment = {
  id: string;
  name: string;
  mode: "proxy" | "direct";
  tenantSlug: string;
  bearerToken: string;
  correlationId: string;
  idempotencyKey: string;
  baseUrls: Record<string, string>;
};

export type RequestResult = {
  status: number;
  statusText: string;
  durationMs: number;
  headers: Record<string, string>;
  body: unknown;
  url: string;
  errorKind?: "auth" | "http" | "network" | "timeout" | "validation";
  error?: string;
};

export const servicesList = [
  "edge",
  "identity",
  "crm",
  "sales",
  "finance",
  "billing",
  "documents",
  "rentals",
  "workflow-control",
  "workflow-runtime",
  "engagement",
  "analytics",
  "simulation",
  "catalog",
  "platform-control",
  "support",
  "supplier",
  "accounting",
  "inventory",
  "procurement",
  "banking",
  "search",
  "ai-governance",
  "notification",
  "fiscal",
  "webhook-hub"
];

export const localBaseUrls: Record<string, string> = {
  edge: "http://localhost:8080",
  identity: "http://localhost:8081",
  "webhook-hub": "http://localhost:8082",
  crm: "http://localhost:8083",
  "workflow-control": "http://localhost:8084",
  "workflow-runtime": "http://localhost:8085",
  analytics: "http://localhost:8086",
  sales: "http://localhost:8087",
  engagement: "http://localhost:8088",
  finance: "http://localhost:8092",
  documents: "http://localhost:8093",
  simulation: "http://localhost:8094",
  billing: "http://localhost:8095",
  rentals: "http://localhost:8096",
  catalog: "http://localhost:8097",
  "platform-control": "http://localhost:8098",
  support: "http://localhost:8099",
  supplier: "http://localhost:8100",
  notification: "http://localhost:8101",
  fiscal: "http://localhost:8102",
  accounting: "http://localhost:8103",
  inventory: "http://localhost:8104",
  procurement: "http://localhost:8105",
  banking: "http://localhost:8106",
  search: "http://localhost:8107",
  "ai-governance": "http://localhost:8108"
};

export const defaultEnvironment: RuntimeEnvironment = {
  id: "local-docker",
  name: "Local Docker",
  mode: "proxy",
  tenantSlug: "bootstrap-ops",
  bearerToken: "",
  correlationId: "console-local",
  idempotencyKey: "console-request",
  baseUrls: localBaseUrls
};

export function defaultBodyFor(endpoint: EndpointContract, tenantSlug: string): string {
  if (!endpoint.hasBody) return "";

  const path = endpoint.path.toLowerCase();
  const common = { tenantSlug };

  if (path.includes("/api/crm/leads")) {
    return stringify({
      name: "Lead Console ERP",
      email: "lead.console@bootstrap-ops.local",
      source: "api-console",
      ownerType: "user",
      ownerPublicId: "0195e7a0-7a9c-7c1f-8a44-4a6e70000301"
    });
  }

  if (path.includes("/api/billing/plans")) {
    return stringify({
      code: "console-growth",
      name: "Console Growth",
      amountCents: 9900,
      intervalUnit: "monthly",
      intervalCount: 1,
      gracePeriodDays: 5,
      maxRetries: 2
    });
  }

  if (path.includes("/api/platform-control")) {
    if (path.includes("/providers/activation/")) {
      return stringify({
        actor: "console@bootstrap-ops.local",
        action: "connection_test",
        payload: { source: "api-console" }
      });
    }

    return stringify({
      requestedBy: "console@bootstrap-ops.local",
      payload: { source: "api-console" }
    });
  }

  if (path.includes("/api/documents")) {
    return stringify({
      ...common,
      ownerType: "crm.lead",
      ownerPublicId: "0195e7a0-7a9c-7c1f-8a44-4a6e70000341",
      fileName: "console-evidence.pdf",
      contentType: "application/pdf",
      storageDriver: "manual",
      source: "api-console",
      uploadedBy: "api-console"
    });
  }

  return stringify({
    ...common,
    source: "api-console",
    notes: "Payload inicial editável gerado pelo ERP - Control Console."
  });
}

export function stringify(value: unknown): string {
  return JSON.stringify(value, null, 2);
}
