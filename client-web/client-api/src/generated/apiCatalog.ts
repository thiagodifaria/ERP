/* eslint-disable */
// Gerado por scripts/generate-catalog.mjs. Nao edite manualmente.

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type ServiceContract = {
  slug: string;
  name: string;
  title: string;
  version: string;
  description: string;
  contractFile: string;
  endpointCount: number;
};

export type EndpointContract = {
  id: string;
  service: string;
  method: HttpMethod;
  path: string;
  tag: string;
  description: string;
  summary: string;
  source: string;
  hasBody: boolean;
  pathParams: string[];
};

export type EventSchemaContract = {
  fileName: string;
  name: string;
  source: string;
};

export const services: ServiceContract[] = [
  {
    "slug": "analytics",
    "name": "Analytics",
    "title": "ERP Analytics API",
    "version": "0.1.0",
    "description": "Executive reports, adapter catalog, SaaS control and contract governance.",
    "contractFile": "docs/contracts/http/analytics.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "billing",
    "name": "Billing",
    "title": "ERP Billing API",
    "version": "0.9.7",
    "description": "Subscription billing, payment recovery and gateway capabilities.",
    "contractFile": "docs/contracts/http/billing.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "catalog",
    "name": "Catalog",
    "title": "ERP Catalog API",
    "version": "0.2.0",
    "description": "Product and service catalog with categories, activation, versioned items, cursor pagination and bulk creation.",
    "contractFile": "docs/contracts/http/catalog.openapi.yaml",
    "endpointCount": 12
  },
  {
    "slug": "crm",
    "name": "CRM",
    "title": "ERP CRM API",
    "version": "0.2.0",
    "description": "CRM leads, customers, activity, attachments and commercial intelligence.",
    "contractFile": "docs/contracts/http/crm.openapi.yaml",
    "endpointCount": 5
  },
  {
    "slug": "documents",
    "name": "Documents",
    "title": "ERP Documents API",
    "version": "0.9.7",
    "description": "Attachment governance, upload orchestration and storage capabilities.",
    "contractFile": "docs/contracts/http/documents.openapi.yaml",
    "endpointCount": 10
  },
  {
    "slug": "edge",
    "name": "Edge",
    "title": "ERP Edge API",
    "version": "0.1.0",
    "description": "Aggregated operational cockpits for tenants, contracts and SaaS control.",
    "contractFile": "docs/contracts/http/edge.openapi.yaml",
    "endpointCount": 8
  },
  {
    "slug": "engagement",
    "name": "Engagement",
    "title": "ERP Engagement API",
    "version": "0.9.7",
    "description": "Omnichannel engagement, provider callbacks and campaign operations.",
    "contractFile": "docs/contracts/http/engagement.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "finance",
    "name": "Finance",
    "title": "ERP Finance API",
    "version": "0.4.0",
    "description": "Receivables, commission holds, cash control and cross-domain financial activity.",
    "contractFile": "docs/contracts/http/finance.openapi.yaml",
    "endpointCount": 5
  },
  {
    "slug": "fiscal",
    "name": "Fiscal",
    "title": "ERP Fiscal API",
    "version": "0.1.0",
    "description": "Fiscal profile, document operations, privacy rights and compliance governance.",
    "contractFile": "docs/contracts/http/fiscal.openapi.yaml",
    "endpointCount": 25
  },
  {
    "slug": "identity",
    "name": "Identity",
    "title": "ERP Identity API",
    "version": "0.5.0",
    "description": "Tenancy, access, sessions, invitations and tenant-scoped identity governance.",
    "contractFile": "docs/contracts/http/identity.openapi.yaml",
    "endpointCount": 6
  },
  {
    "slug": "notification",
    "name": "Notification",
    "title": "ERP Notification API",
    "version": "0.1.0",
    "description": "Internal notification center, preferences and reusable dispatch shape.",
    "contractFile": "docs/contracts/http/notification.openapi.yaml",
    "endpointCount": 7
  },
  {
    "slug": "platform-control",
    "name": "Platform Control",
    "title": "ERP Platform Control API",
    "version": "0.2.0",
    "description": "Tenant capabilities, entitlements, quotas, metering, lifecycle jobs and SaaS governance.",
    "contractFile": "docs/contracts/http/platform-control.openapi.yaml",
    "endpointCount": 40
  },
  {
    "slug": "rentals",
    "name": "Rentals",
    "title": "ERP Rentals API",
    "version": "0.8.0",
    "description": "Rental contracts, recurring charges, adjustments, rescission and attachment linkage.",
    "contractFile": "docs/contracts/http/rentals.openapi.yaml",
    "endpointCount": 4
  },
  {
    "slug": "sales",
    "name": "Sales",
    "title": "ERP Sales API",
    "version": "0.7.0",
    "description": "Opportunities, proposals, sales, invoices and commercial lifecycle control.",
    "contractFile": "docs/contracts/http/sales.openapi.yaml",
    "endpointCount": 6
  },
  {
    "slug": "simulation",
    "name": "Simulation",
    "title": "ERP Simulation API",
    "version": "0.7.0",
    "description": "Scenario simulation, load benchmark and cost estimation runtime.",
    "contractFile": "docs/contracts/http/simulation.openapi.yaml",
    "endpointCount": 3
  },
  {
    "slug": "supplier",
    "name": "Supplier",
    "title": "ERP Supplier API",
    "version": "0.1.0",
    "description": "Supplier directory, categories and payables-oriented vendor governance.",
    "contractFile": "docs/contracts/http/supplier.openapi.yaml",
    "endpointCount": 8
  },
  {
    "slug": "support",
    "name": "Support",
    "title": "ERP Support API",
    "version": "0.1.0",
    "description": "Queue-based support cases with SLA, comments and lifecycle history.",
    "contractFile": "docs/contracts/http/support.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "webhook-hub",
    "name": "Webhook Hub",
    "title": "ERP Webhook Hub API",
    "version": "0.9.7",
    "description": "Inbound webhook intake, DLQ and operator recovery surface.",
    "contractFile": "docs/contracts/http/webhook-hub.openapi.yaml",
    "endpointCount": 13
  },
  {
    "slug": "workflow-control",
    "name": "Workflow Control",
    "title": "ERP Workflow Control API",
    "version": "0.6.0",
    "description": "Workflow definition catalog, status, publication and action snapshots.",
    "contractFile": "docs/contracts/http/workflow-control.openapi.yaml",
    "endpointCount": 7
  },
  {
    "slug": "workflow-runtime",
    "name": "Workflow Runtime",
    "title": "ERP Workflow Runtime API",
    "version": "0.6.0",
    "description": "Workflow execution runtime with action snapshots, retries, delays and compensations.",
    "contractFile": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "endpointCount": 6
  }
];

export const endpoints: EndpointContract[] = [
  {
    "id": "analytics:GET:/health/live",
    "service": "analytics",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço analytics.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/health/ready",
    "service": "analytics",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço analytics.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/health/details",
    "service": "analytics",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço analytics.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/adapter-catalog",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/adapter-catalog",
    "tag": "Analytics",
    "description": "Read external adapter capability catalog",
    "summary": "Read external adapter capability catalog",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/integration-readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/integration-readiness",
    "tag": "Analytics",
    "description": "Read external integration readiness",
    "summary": "Read external integration readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/saas-control",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/saas-control",
    "tag": "Analytics",
    "description": "Read SaaS control posture by tenant",
    "summary": "Read SaaS control posture by tenant",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/contract-governance",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/contract-governance",
    "tag": "Analytics",
    "description": "Read contract governance posture",
    "summary": "Read contract governance posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/hardening-review",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/hardening-review",
    "tag": "Analytics",
    "description": "Read hardening review",
    "summary": "Read hardening review",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/core-operations",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/core-operations",
    "tag": "Analytics",
    "description": "Read core product operations",
    "summary": "Read core product operations",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/relationship-intelligence",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/relationship-intelligence",
    "tag": "Analytics",
    "description": "Read relationship intelligence",
    "summary": "Read relationship intelligence",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/compliance-control",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/compliance-control",
    "tag": "Analytics",
    "description": "Read fiscal and privacy compliance control",
    "summary": "Read fiscal and privacy compliance control",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/go-live-control",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/go-live-control",
    "tag": "Analytics",
    "description": "Read go-live rollout control",
    "summary": "Read go-live rollout control",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/live",
    "service": "billing",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço billing.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/ready",
    "service": "billing",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço billing.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/details",
    "service": "billing",
    "method": "GET",
    "path": "/health/details",
    "tag": "Billing",
    "description": "Return readiness details and gateway posture",
    "summary": "Return readiness details and gateway posture",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/gateways",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/gateways",
    "tag": "Billing",
    "description": "List gateway capabilities and Pix posture",
    "summary": "List gateway capabilities and Pix posture",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/gateways/{provider}",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/gateways/{provider}",
    "tag": "Billing",
    "description": "Read one gateway capability",
    "summary": "Read one gateway capability",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "provider"
    ]
  },
  {
    "id": "billing:GET:/api/billing/plans",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/plans",
    "tag": "Billing",
    "description": "List billing plans including flat, hybrid and usage-based pricing",
    "summary": "List billing plans including flat, hybrid and usage-based pricing",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:POST:/api/billing/plans",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/plans",
    "tag": "Billing",
    "description": "Create billing plan",
    "summary": "Create billing plan",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/subscriptions",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions",
    "tag": "Billing",
    "description": "List subscriptions",
    "summary": "List subscriptions",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:POST:/api/billing/subscriptions",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/subscriptions",
    "tag": "Billing",
    "description": "Create subscription",
    "summary": "Create subscription",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/subscriptions/{publicId}/usage-pricing",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions/{publicId}/usage-pricing",
    "tag": "Billing",
    "description": "Project usage-based charge for one subscription",
    "summary": "Project usage-based charge for one subscription",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/invoices/{publicId}/attempts",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/invoices/{publicId}/attempts",
    "tag": "Billing",
    "description": "Create payment attempt with idempotency support",
    "summary": "Create payment attempt with idempotency support",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "catalog:GET:/health/live",
    "service": "catalog",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço catalog.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/health/ready",
    "service": "catalog",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço catalog.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/health/details",
    "service": "catalog",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço catalog.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/capabilities",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/capabilities",
    "tag": "Catalog",
    "description": "Read catalog capability posture",
    "summary": "Read catalog capability posture",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/consumers",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/consumers",
    "tag": "Catalog",
    "description": "Read catalog consumer contracts across core domains",
    "summary": "Read catalog consumer contracts across core domains",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/categories",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/categories",
    "tag": "Catalog",
    "description": "List categories by tenant",
    "summary": "List categories by tenant",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:POST:/api/catalog/categories",
    "service": "catalog",
    "method": "POST",
    "path": "/api/catalog/categories",
    "tag": "Catalog",
    "description": "Create one category",
    "summary": "Create one category",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/categories/page",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/categories/page",
    "tag": "Catalog",
    "description": "Cursor-based category listing",
    "summary": "Cursor-based category listing",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/items",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/items",
    "tag": "Catalog",
    "description": "List catalog items",
    "summary": "List catalog items",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:POST:/api/catalog/items",
    "service": "catalog",
    "method": "POST",
    "path": "/api/catalog/items",
    "tag": "Catalog",
    "description": "Create one catalog item",
    "summary": "Create one catalog item",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/items/page",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/items/page",
    "tag": "Catalog",
    "description": "Cursor-based item listing",
    "summary": "Cursor-based item listing",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "catalog:POST:/api/catalog/items/bulk",
    "service": "catalog",
    "method": "POST",
    "path": "/api/catalog/items/bulk",
    "tag": "Catalog",
    "description": "Bulk create catalog items with partial success",
    "summary": "Bulk create catalog items with partial success",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "catalog:GET:/api/catalog/items/{publicId}",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/items/{publicId}",
    "tag": "Catalog",
    "description": "Read one catalog item",
    "summary": "Read one catalog item",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "catalog:PATCH:/api/catalog/items/{publicId}",
    "service": "catalog",
    "method": "PATCH",
    "path": "/api/catalog/items/{publicId}",
    "tag": "Catalog",
    "description": "Update active state, price and attributes",
    "summary": "Update active state, price and attributes",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "catalog:GET:/api/catalog/items/{publicId}/versions",
    "service": "catalog",
    "method": "GET",
    "path": "/api/catalog/items/{publicId}/versions",
    "tag": "Catalog",
    "description": "Read catalog item version history",
    "summary": "Read catalog item version history",
    "source": "docs/contracts/http/catalog.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/health/live",
    "service": "crm",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço crm.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:GET:/health/ready",
    "service": "crm",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço crm.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:GET:/health/details",
    "service": "crm",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço crm.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/enrichment/cnpj/capabilities",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/enrichment/cnpj/capabilities",
    "tag": "Crm",
    "description": "Read CNPJ enrichment provider capabilities",
    "summary": "Read CNPJ enrichment provider capabilities",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:POST:/api/crm/enrichment/cnpj/lookup",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/enrichment/cnpj/lookup",
    "tag": "Crm",
    "description": "Lookup and enrich one CNPJ through provider contract",
    "summary": "Lookup and enrich one CNPJ through provider contract",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/pipeline/config",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/pipeline/config",
    "tag": "Crm",
    "description": "Read tenant pipeline configuration",
    "summary": "Read tenant pipeline configuration",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:PUT:/api/crm/pipeline/config",
    "service": "crm",
    "method": "PUT",
    "path": "/api/crm/pipeline/config",
    "tag": "Crm",
    "description": "Upsert tenant pipeline configuration",
    "summary": "Upsert tenant pipeline configuration",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/leads/intelligence/summary",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/intelligence/summary",
    "tag": "Crm",
    "description": "Read lead scoring and pipeline intelligence summary",
    "summary": "Read lead scoring and pipeline intelligence summary",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/health/live",
    "service": "documents",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço documents.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/health/ready",
    "service": "documents",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço documents.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/health/details",
    "service": "documents",
    "method": "GET",
    "path": "/health/details",
    "tag": "Documents",
    "description": "Return runtime readiness and storage posture",
    "summary": "Return runtime readiness and storage posture",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/signing/capabilities",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/signing/capabilities",
    "tag": "Documents",
    "description": "List digital signature capabilities",
    "summary": "List digital signature capabilities",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/signing/capabilities/{provider}",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/signing/capabilities/{provider}",
    "tag": "Documents",
    "description": "Read one signing capability",
    "summary": "Read one signing capability",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "provider"
    ]
  },
  {
    "id": "documents:POST:/api/documents/signing/requests",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/signing/requests",
    "tag": "Documents",
    "description": "Queue one digital signature request",
    "summary": "Queue one digital signature request",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/storage/capabilities",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/storage/capabilities",
    "tag": "Documents",
    "description": "List storage capability registry",
    "summary": "List storage capability registry",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/storage/capabilities/{provider}",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/storage/capabilities/{provider}",
    "tag": "Documents",
    "description": "Read one storage capability",
    "summary": "Read one storage capability",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "provider"
    ]
  },
  {
    "id": "documents:GET:/api/documents/attachments",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/attachments",
    "tag": "Documents",
    "description": "List attachments",
    "summary": "List attachments",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:POST:/api/documents/attachments",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/attachments",
    "tag": "Documents",
    "description": "Create attachment metadata",
    "summary": "Create attachment metadata",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/attachments/{publicId}/versions",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/attachments/{publicId}/versions",
    "tag": "Documents",
    "description": "List attachment versions",
    "summary": "List attachment versions",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:POST:/api/documents/attachments/{publicId}/versions",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/attachments/{publicId}/versions",
    "tag": "Documents",
    "description": "Append attachment version",
    "summary": "Append attachment version",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "edge:GET:/health/live",
    "service": "edge",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço edge.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/health/ready",
    "service": "edge",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço edge.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/health/details",
    "service": "edge",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço edge.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/core-operations",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/core-operations",
    "tag": "Edge",
    "description": "Read executive core product cockpit",
    "summary": "Read executive core product cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/relationship-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/relationship-overview",
    "tag": "Edge",
    "description": "Read executive relationship cockpit",
    "summary": "Read executive relationship cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/compliance-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/compliance-overview",
    "tag": "Edge",
    "description": "Read executive compliance cockpit",
    "summary": "Read executive compliance cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/go-live-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/go-live-overview",
    "tag": "Edge",
    "description": "Read executive go-live cockpit",
    "summary": "Read executive go-live cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/integrations-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/integrations-overview",
    "tag": "Edge",
    "description": "Read executive integrations cockpit",
    "summary": "Read executive integrations cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/saas-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/saas-overview",
    "tag": "Edge",
    "description": "Read executive SaaS cockpit",
    "summary": "Read executive SaaS cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/contracts-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/contracts-overview",
    "tag": "Edge",
    "description": "Read executive contracts cockpit",
    "summary": "Read executive contracts cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/hardening-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/hardening-overview",
    "tag": "Edge",
    "description": "Read executive hardening cockpit",
    "summary": "Read executive hardening cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/health/live",
    "service": "engagement",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço engagement.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/health/ready",
    "service": "engagement",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço engagement.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/health/details",
    "service": "engagement",
    "method": "GET",
    "path": "/health/details",
    "tag": "Engagement",
    "description": "Return readiness details for engagement runtime",
    "summary": "Return readiness details for engagement runtime",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/api/engagement/providers",
    "service": "engagement",
    "method": "GET",
    "path": "/api/engagement/providers",
    "tag": "Engagement",
    "description": "List provider capabilities and fallback posture",
    "summary": "List provider capabilities and fallback posture",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/api/engagement/providers/{provider}",
    "service": "engagement",
    "method": "GET",
    "path": "/api/engagement/providers/{provider}",
    "tag": "Engagement",
    "description": "Read one provider capability",
    "summary": "Read one provider capability",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "provider"
    ]
  },
  {
    "id": "engagement:POST:/api/engagement/providers/meta-ads/leads",
    "service": "engagement",
    "method": "POST",
    "path": "/api/engagement/providers/meta-ads/leads",
    "tag": "Engagement",
    "description": "Ingest inbound lead from Meta Ads",
    "summary": "Ingest inbound lead from Meta Ads",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "engagement:POST:/api/engagement/providers/resend/events",
    "service": "engagement",
    "method": "POST",
    "path": "/api/engagement/providers/resend/events",
    "tag": "Engagement",
    "description": "Register Resend callback event",
    "summary": "Register Resend callback event",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "engagement:POST:/api/engagement/providers/whatsapp-cloud/events",
    "service": "engagement",
    "method": "POST",
    "path": "/api/engagement/providers/whatsapp-cloud/events",
    "tag": "Engagement",
    "description": "Register WhatsApp callback event",
    "summary": "Register WhatsApp callback event",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "engagement:POST:/api/engagement/providers/telegram-bot/events",
    "service": "engagement",
    "method": "POST",
    "path": "/api/engagement/providers/telegram-bot/events",
    "tag": "Engagement",
    "description": "Register Telegram callback event",
    "summary": "Register Telegram callback event",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/api/engagement/provider-events",
    "service": "engagement",
    "method": "GET",
    "path": "/api/engagement/provider-events",
    "tag": "Engagement",
    "description": "List provider events",
    "summary": "List provider events",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "engagement:GET:/api/engagement/provider-events/{publicId}",
    "service": "engagement",
    "method": "GET",
    "path": "/api/engagement/provider-events/{publicId}",
    "tag": "Engagement",
    "description": "Read one provider event",
    "summary": "Read one provider event",
    "source": "docs/contracts/http/engagement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:GET:/health/live",
    "service": "finance",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço finance.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/health/ready",
    "service": "finance",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço finance.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/health/details",
    "service": "finance",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço finance.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/receivable-projections",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/receivable-projections",
    "tag": "Finance",
    "description": "List receivable projections",
    "summary": "List receivable projections",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/receivable-projections/sync",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/receivable-projections/sync",
    "tag": "Finance",
    "description": "Sync projections from sales and rentals",
    "summary": "Sync projections from sales and rentals",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/commission-holds",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/commission-holds",
    "tag": "Finance",
    "description": "List commission holds",
    "summary": "List commission holds",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/commission-holds/{publicId}/release",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/commission-holds/{publicId}/release",
    "tag": "Finance",
    "description": "Release one commission hold",
    "summary": "Release one commission hold",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:GET:/api/finance/activity",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/activity",
    "tag": "Finance",
    "description": "List finance operational activity",
    "summary": "List finance operational activity",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/health/live",
    "service": "fiscal",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço fiscal.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/health/ready",
    "service": "fiscal",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço fiscal.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/health/details",
    "service": "fiscal",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço fiscal.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/capabilities",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/capabilities",
    "tag": "Fiscal",
    "description": "Read fiscal capability registry",
    "summary": "Read fiscal capability registry",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/companies/{companyPublicId}/profile",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/companies/{companyPublicId}/profile",
    "tag": "Fiscal",
    "description": "Read fiscal company profile",
    "summary": "Read fiscal company profile",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "companyPublicId"
    ]
  },
  {
    "id": "fiscal:PUT:/api/fiscal/companies/{companyPublicId}/profile",
    "service": "fiscal",
    "method": "PUT",
    "path": "/api/fiscal/companies/{companyPublicId}/profile",
    "tag": "Fiscal",
    "description": "Upsert fiscal company profile",
    "summary": "Upsert fiscal company profile",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "companyPublicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/companies/{companyPublicId}/retention-policies",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/companies/{companyPublicId}/retention-policies",
    "tag": "Fiscal",
    "description": "List retention policies by company",
    "summary": "List retention policies by company",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "companyPublicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/companies/{companyPublicId}/retention-execution",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/companies/{companyPublicId}/retention-execution",
    "tag": "Fiscal",
    "description": "Read retention execution plan for one company",
    "summary": "Read retention execution plan for one company",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "companyPublicId"
    ]
  },
  {
    "id": "fiscal:POST:/api/fiscal/companies/{companyPublicId}/retention-execution/execute",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/companies/{companyPublicId}/retention-execution/execute",
    "tag": "Fiscal",
    "description": "Execute retention and anonymization plan",
    "summary": "Execute retention and anonymization plan",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "companyPublicId"
    ]
  },
  {
    "id": "fiscal:PUT:/api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}",
    "service": "fiscal",
    "method": "PUT",
    "path": "/api/fiscal/companies/{companyPublicId}/retention-policies/{dataDomain}",
    "tag": "Fiscal",
    "description": "Upsert retention policy for one data domain",
    "summary": "Upsert retention policy for one data domain",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "companyPublicId",
      "dataDomain"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/documents",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/documents",
    "tag": "Fiscal",
    "description": "List fiscal documents",
    "summary": "List fiscal documents",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/documents",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/documents",
    "tag": "Fiscal",
    "description": "Issue one fiscal document",
    "summary": "Issue one fiscal document",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/documents/{publicId}",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/documents/{publicId}",
    "tag": "Fiscal",
    "description": "Read one fiscal document",
    "summary": "Read one fiscal document",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:POST:/api/fiscal/documents/{publicId}/cancel",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/documents/{publicId}/cancel",
    "tag": "Fiscal",
    "description": "Cancel one fiscal document",
    "summary": "Cancel one fiscal document",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:POST:/api/fiscal/documents/{publicId}/correction-letter",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/documents/{publicId}/correction-letter",
    "tag": "Fiscal",
    "description": "Register correction letter for one fiscal document",
    "summary": "Register correction letter for one fiscal document",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:POST:/api/fiscal/documents/{publicId}/invalidate",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/documents/{publicId}/invalidate",
    "tag": "Fiscal",
    "description": "Register invalidation for one fiscal document",
    "summary": "Register invalidation for one fiscal document",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/documents/{publicId}/events",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/documents/{publicId}/events",
    "tag": "Fiscal",
    "description": "List fiscal document audit events",
    "summary": "List fiscal document audit events",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/privacy-requests",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/privacy-requests",
    "tag": "Fiscal",
    "description": "List privacy requests",
    "summary": "List privacy requests",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/privacy-requests",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/privacy-requests",
    "tag": "Fiscal",
    "description": "Create privacy request",
    "summary": "Create privacy request",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/privacy-requests/{publicId}",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/privacy-requests/{publicId}",
    "tag": "Fiscal",
    "description": "Read one privacy request",
    "summary": "Read one privacy request",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/privacy-requests/{publicId}/export-package",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/privacy-requests/{publicId}/export-package",
    "tag": "Fiscal",
    "description": "Build export package for one privacy request",
    "summary": "Build export package for one privacy request",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:POST:/api/fiscal/privacy-requests/{publicId}/execute",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/privacy-requests/{publicId}/execute",
    "tag": "Fiscal",
    "description": "Execute one privacy request with audit trail",
    "summary": "Execute one privacy request with audit trail",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:PATCH:/api/fiscal/privacy-requests/{publicId}/status",
    "service": "fiscal",
    "method": "PATCH",
    "path": "/api/fiscal/privacy-requests/{publicId}/status",
    "tag": "Fiscal",
    "description": "Transition privacy request lifecycle status",
    "summary": "Transition privacy request lifecycle status",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/consents",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/consents",
    "tag": "Fiscal",
    "description": "List consent ledger",
    "summary": "List consent ledger",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/consents",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/consents",
    "tag": "Fiscal",
    "description": "Create consent record",
    "summary": "Create consent record",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:PATCH:/api/fiscal/consents/{publicId}",
    "service": "fiscal",
    "method": "PATCH",
    "path": "/api/fiscal/consents/{publicId}",
    "tag": "Fiscal",
    "description": "Transition consent status",
    "summary": "Transition consent status",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "fiscal:GET:/api/fiscal/audit-events",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/audit-events",
    "tag": "Fiscal",
    "description": "List fiscal audit events",
    "summary": "List fiscal audit events",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/compliance/summary",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/compliance/summary",
    "tag": "Fiscal",
    "description": "Read fiscal compliance summary",
    "summary": "Read fiscal compliance summary",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/health/live",
    "service": "identity",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço identity.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/health/ready",
    "service": "identity",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço identity.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/health/details",
    "service": "identity",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço identity.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/api/identity/tenants",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants",
    "tag": "Identity",
    "description": "List tenants",
    "summary": "List tenants",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/tenants",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants",
    "tag": "Identity",
    "description": "Create tenant",
    "summary": "Create tenant",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/snapshot",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/snapshot",
    "tag": "Identity",
    "description": "Read one tenant snapshot",
    "summary": "Read one tenant snapshot",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:POST:/api/identity/sessions/login",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/sessions/login",
    "tag": "Identity",
    "description": "Authenticate identity session",
    "summary": "Authenticate identity session",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/sessions/refresh",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/sessions/refresh",
    "tag": "Identity",
    "description": "Refresh identity session",
    "summary": "Refresh identity session",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/invitations",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/invitations",
    "tag": "Identity",
    "description": "Create invitation",
    "summary": "Create invitation",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "notification:GET:/health/live",
    "service": "notification",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço notification.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "notification:GET:/health/ready",
    "service": "notification",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço notification.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "notification:GET:/health/details",
    "service": "notification",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço notification.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "notification:GET:/api/notification/capabilities",
    "service": "notification",
    "method": "GET",
    "path": "/api/notification/capabilities",
    "tag": "Notification",
    "description": "Read notification capability catalog",
    "summary": "Read notification capability catalog",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "notification:GET:/api/notification/preferences/{userPublicId}",
    "service": "notification",
    "method": "GET",
    "path": "/api/notification/preferences/{userPublicId}",
    "tag": "Notification",
    "description": "Read one user notification preference",
    "summary": "Read one user notification preference",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "userPublicId"
    ]
  },
  {
    "id": "notification:PUT:/api/notification/preferences/{userPublicId}",
    "service": "notification",
    "method": "PUT",
    "path": "/api/notification/preferences/{userPublicId}",
    "tag": "Notification",
    "description": "Upsert one user notification preference",
    "summary": "Upsert one user notification preference",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "userPublicId"
    ]
  },
  {
    "id": "notification:GET:/api/notification/center",
    "service": "notification",
    "method": "GET",
    "path": "/api/notification/center",
    "tag": "Notification",
    "description": "List notification center items with cursor filters",
    "summary": "List notification center items with cursor filters",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "notification:POST:/api/notification/center",
    "service": "notification",
    "method": "POST",
    "path": "/api/notification/center",
    "tag": "Notification",
    "description": "Create one notification center item",
    "summary": "Create one notification center item",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "notification:PATCH:/api/notification/center/{publicId}/status",
    "service": "notification",
    "method": "PATCH",
    "path": "/api/notification/center/{publicId}/status",
    "tag": "Notification",
    "description": "Transition notification status",
    "summary": "Transition notification status",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "notification:GET:/api/notification/summary",
    "service": "notification",
    "method": "GET",
    "path": "/api/notification/summary",
    "tag": "Notification",
    "description": "Read notification summary",
    "summary": "Read notification summary",
    "source": "docs/contracts/http/notification.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/health/live",
    "service": "platform-control",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço platform-control.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/health/ready",
    "service": "platform-control",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço platform-control.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/health/details",
    "service": "platform-control",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço platform-control.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/capabilities/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/capabilities/catalog",
    "tag": "Platform Control",
    "description": "List platform capability catalog",
    "summary": "List platform capability catalog",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/providers/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/providers/catalog",
    "tag": "Platform Control",
    "description": "List provider capability catalog and environment posture",
    "summary": "List provider capability catalog and environment posture",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/entitlements",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/entitlements",
    "tag": "Platform Control",
    "description": "List tenant entitlements with cursor pagination",
    "summary": "List tenant entitlements with cursor pagination",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/feature-flags",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/feature-flags",
    "tag": "Platform Control",
    "description": "List tenant feature flags with capability metadata",
    "summary": "List tenant feature flags with capability metadata",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/entitlements/{capabilityKey}",
    "tag": "Platform Control",
    "description": "Upsert one entitlement",
    "summary": "Upsert one entitlement",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "capabilityKey"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/feature-flags/{capabilityKey}",
    "tag": "Platform Control",
    "description": "Upsert one feature flag using entitlement governance",
    "summary": "Upsert one feature flag using entitlement governance",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "capabilityKey"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/entitlements/bulk",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/entitlements/bulk",
    "tag": "Platform Control",
    "description": "Bulk upsert entitlements with partial success",
    "summary": "Bulk upsert entitlements with partial success",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/provider-defaults",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/provider-defaults",
    "tag": "Platform Control",
    "description": "List provider defaults selected for one tenant",
    "summary": "List provider defaults selected for one tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/provider-defaults/{capabilityKey}",
    "tag": "Platform Control",
    "description": "Upsert provider default for one tenant capability",
    "summary": "Upsert provider default for one tenant capability",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "capabilityKey"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/quotas",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/quotas",
    "tag": "Platform Control",
    "description": "List quotas by tenant",
    "summary": "List quotas by tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/quotas/{metricKey}",
    "tag": "Platform Control",
    "description": "Upsert one quota",
    "summary": "Upsert one quota",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "metricKey"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/quotas/bulk",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/quotas/bulk",
    "tag": "Platform Control",
    "description": "Bulk upsert quotas with partial success",
    "summary": "Bulk upsert quotas with partial success",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/blocks",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/blocks",
    "tag": "Platform Control",
    "description": "List tenant blocks",
    "summary": "List tenant blocks",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/blocks/{blockKey}",
    "tag": "Platform Control",
    "description": "Upsert tenant block",
    "summary": "Upsert tenant block",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "blockKey"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/metering",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/metering",
    "tag": "Platform Control",
    "description": "Read metering snapshots and summary with cursor pagination",
    "summary": "Read metering snapshots and summary with cursor pagination",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/metering/snapshots",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/metering/snapshots",
    "tag": "Platform Control",
    "description": "Create one usage snapshot",
    "summary": "Create one usage snapshot",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/usage-summary",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/usage-summary",
    "tag": "Platform Control",
    "description": "Read quota and metering utilization summary",
    "summary": "Read quota and metering utilization summary",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/lifecycle/readiness",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/readiness",
    "tag": "Platform Control",
    "description": "Read tenant lifecycle readiness and provider posture",
    "summary": "Read tenant lifecycle readiness and provider posture",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs",
    "tag": "Platform Control",
    "description": "List onboarding and offboarding jobs with cursor pagination",
    "summary": "List onboarding and offboarding jobs with cursor pagination",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}",
    "tag": "Platform Control",
    "description": "Read one lifecycle job with audit events",
    "summary": "Read one lifecycle job with audit events",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding/preview",
    "tag": "Platform Control",
    "description": "Preview onboarding plan, provider defaults and lifecycle readiness",
    "summary": "Preview onboarding plan, provider defaults and lifecycle readiness",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/onboarding",
    "tag": "Platform Control",
    "description": "Queue onboarding job with Idempotency-Key and 202 Accepted",
    "summary": "Queue onboarding job with Idempotency-Key and 202 Accepted",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding/preview",
    "tag": "Platform Control",
    "description": "Preview offboarding plan, retention posture and lifecycle readiness",
    "summary": "Preview offboarding plan, retention posture and lifecycle readiness",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/offboarding",
    "tag": "Platform Control",
    "description": "Queue offboarding job with Idempotency-Key and 202 Accepted",
    "summary": "Queue offboarding job with Idempotency-Key and 202 Accepted",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/start",
    "tag": "Platform Control",
    "description": "Transition lifecycle job to running",
    "summary": "Transition lifecycle job to running",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/complete",
    "tag": "Platform Control",
    "description": "Transition lifecycle job to completed",
    "summary": "Transition lifecycle job to completed",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/fail",
    "tag": "Platform Control",
    "description": "Transition lifecycle job to failed",
    "summary": "Transition lifecycle job to failed",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/lifecycle/jobs/{publicId}/cancel",
    "tag": "Platform Control",
    "description": "Transition lifecycle job to cancelled",
    "summary": "Transition lifecycle job to cancelled",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/readiness",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/readiness",
    "tag": "Platform Control",
    "description": "Read go-live rollout readiness by tenant",
    "summary": "Read go-live rollout readiness by tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/adoption",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/adoption",
    "tag": "Platform Control",
    "description": "Read tenant go-live adoption baseline and gap",
    "summary": "Read tenant go-live adoption baseline and gap",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/bottlenecks",
    "tag": "Platform Control",
    "description": "List go-live bottlenecks and operational blockers",
    "summary": "List go-live bottlenecks and operational blockers",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/playbook",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/playbook",
    "tag": "Platform Control",
    "description": "Read rollout and rollback playbook for one tenant",
    "summary": "Read rollout and rollback playbook for one tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/adjustments",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/adjustments",
    "tag": "Platform Control",
    "description": "List recommended go-live adjustments",
    "summary": "List recommended go-live adjustments",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/adjustments/apply",
    "tag": "Platform Control",
    "description": "Apply one go-live operational adjustment",
    "summary": "Apply one go-live operational adjustment",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts",
    "tag": "Platform Control",
    "description": "List go-live rollouts",
    "summary": "List go-live rollouts",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts",
    "tag": "Platform Control",
    "description": "Create one go-live rollout",
    "summary": "Create one go-live rollout",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}",
    "tag": "Platform Control",
    "description": "Read one go-live rollout with events",
    "summary": "Read one go-live rollout with events",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/start",
    "tag": "Platform Control",
    "description": "Transition go-live rollout to running",
    "summary": "Transition go-live rollout to running",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/complete",
    "tag": "Platform Control",
    "description": "Transition go-live rollout to completed",
    "summary": "Transition go-live rollout to completed",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/go-live/rollouts/{publicId}/rollback",
    "tag": "Platform Control",
    "description": "Roll back one go-live rollout",
    "summary": "Roll back one go-live rollout",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "rentals:GET:/health/live",
    "service": "rentals",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço rentals.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "rentals:GET:/health/ready",
    "service": "rentals",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço rentals.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "rentals:GET:/health/details",
    "service": "rentals",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço rentals.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "rentals:GET:/api/rentals/contracts",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts",
    "tag": "Rentals",
    "description": "List rental contracts",
    "summary": "List rental contracts",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "rentals:POST:/api/rentals/contracts",
    "service": "rentals",
    "method": "POST",
    "path": "/api/rentals/contracts",
    "tag": "Rentals",
    "description": "Create rental contract",
    "summary": "Create rental contract",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "rentals:GET:/api/rentals/contracts/{publicId}/charges",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/{publicId}/charges",
    "tag": "Rentals",
    "description": "List contract charges",
    "summary": "List contract charges",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:PATCH:/api/rentals/contracts/{publicId}/charges/{chargePublicId}/status",
    "service": "rentals",
    "method": "PATCH",
    "path": "/api/rentals/contracts/{publicId}/charges/{chargePublicId}/status",
    "tag": "Rentals",
    "description": "Update charge status",
    "summary": "Update charge status",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId",
      "chargePublicId"
    ]
  },
  {
    "id": "sales:GET:/health/live",
    "service": "sales",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço sales.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/health/ready",
    "service": "sales",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço sales.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/health/details",
    "service": "sales",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço sales.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/opportunities",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/opportunities",
    "tag": "Sales",
    "description": "List opportunities",
    "summary": "List opportunities",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:POST:/api/sales/opportunities",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/opportunities",
    "tag": "Sales",
    "description": "Create opportunity",
    "summary": "Create opportunity",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/proposals",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/proposals",
    "tag": "Sales",
    "description": "List proposals",
    "summary": "List proposals",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:POST:/api/sales/proposals",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/proposals",
    "tag": "Sales",
    "description": "Create proposal",
    "summary": "Create proposal",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/sales",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales",
    "tag": "Sales",
    "description": "List sales",
    "summary": "List sales",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/invoices",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/invoices",
    "tag": "Sales",
    "description": "List commercial invoices",
    "summary": "List commercial invoices",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/health/live",
    "service": "simulation",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço simulation.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/health/ready",
    "service": "simulation",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço simulation.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/health/details",
    "service": "simulation",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço simulation.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/api/simulation/scenarios",
    "service": "simulation",
    "method": "GET",
    "path": "/api/simulation/scenarios",
    "tag": "Simulation",
    "description": "List scenarios",
    "summary": "List scenarios",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:POST:/api/simulation/scenarios",
    "service": "simulation",
    "method": "POST",
    "path": "/api/simulation/scenarios",
    "tag": "Simulation",
    "description": "Create scenario run",
    "summary": "Create scenario run",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "simulation:POST:/api/simulation/benchmarks/load",
    "service": "simulation",
    "method": "POST",
    "path": "/api/simulation/benchmarks/load",
    "tag": "Simulation",
    "description": "Execute one load benchmark run",
    "summary": "Execute one load benchmark run",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/health/live",
    "service": "supplier",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço supplier.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/health/ready",
    "service": "supplier",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço supplier.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/health/details",
    "service": "supplier",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço supplier.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/api/supplier/capabilities",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/capabilities",
    "tag": "Supplier",
    "description": "Read supplier capability catalog",
    "summary": "Read supplier capability catalog",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/api/supplier/categories",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/categories",
    "tag": "Supplier",
    "description": "List supplier categories",
    "summary": "List supplier categories",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:PUT:/api/supplier/categories/{categoryKey}",
    "service": "supplier",
    "method": "PUT",
    "path": "/api/supplier/categories/{categoryKey}",
    "tag": "Supplier",
    "description": "Upsert one supplier category",
    "summary": "Upsert one supplier category",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "categoryKey"
    ]
  },
  {
    "id": "supplier:GET:/api/supplier/suppliers",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/suppliers",
    "tag": "Supplier",
    "description": "List suppliers by tenant and status",
    "summary": "List suppliers by tenant and status",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:POST:/api/supplier/suppliers",
    "service": "supplier",
    "method": "POST",
    "path": "/api/supplier/suppliers",
    "tag": "Supplier",
    "description": "Create one supplier",
    "summary": "Create one supplier",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/api/supplier/suppliers/summary",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/suppliers/summary",
    "tag": "Supplier",
    "description": "Read supplier summary",
    "summary": "Read supplier summary",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:GET:/api/supplier/suppliers/{publicId}",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/suppliers/{publicId}",
    "tag": "Supplier",
    "description": "Read one supplier",
    "summary": "Read one supplier",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "supplier:PATCH:/api/supplier/suppliers/{publicId}",
    "service": "supplier",
    "method": "PATCH",
    "path": "/api/supplier/suppliers/{publicId}",
    "tag": "Supplier",
    "description": "Update one supplier",
    "summary": "Update one supplier",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "support:GET:/health/live",
    "service": "support",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço support.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:GET:/health/ready",
    "service": "support",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço support.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:GET:/health/details",
    "service": "support",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço support.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:GET:/api/support/capabilities",
    "service": "support",
    "method": "GET",
    "path": "/api/support/capabilities",
    "tag": "Support",
    "description": "Read support capability catalog",
    "summary": "Read support capability catalog",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:GET:/api/support/queues",
    "service": "support",
    "method": "GET",
    "path": "/api/support/queues",
    "tag": "Support",
    "description": "List support queues by tenant",
    "summary": "List support queues by tenant",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:PUT:/api/support/queues/{queueKey}",
    "service": "support",
    "method": "PUT",
    "path": "/api/support/queues/{queueKey}",
    "tag": "Support",
    "description": "Upsert one support queue",
    "summary": "Upsert one support queue",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "queueKey"
    ]
  },
  {
    "id": "support:GET:/api/support/cases",
    "service": "support",
    "method": "GET",
    "path": "/api/support/cases",
    "tag": "Support",
    "description": "List support cases with cursor filters",
    "summary": "List support cases with cursor filters",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:POST:/api/support/cases",
    "service": "support",
    "method": "POST",
    "path": "/api/support/cases",
    "tag": "Support",
    "description": "Create one support case",
    "summary": "Create one support case",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "support:GET:/api/support/cases/summary",
    "service": "support",
    "method": "GET",
    "path": "/api/support/cases/summary",
    "tag": "Support",
    "description": "Read support case summary",
    "summary": "Read support case summary",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:GET:/api/support/cases/{publicId}",
    "service": "support",
    "method": "GET",
    "path": "/api/support/cases/{publicId}",
    "tag": "Support",
    "description": "Read one support case",
    "summary": "Read one support case",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "support:PATCH:/api/support/cases/{publicId}/status",
    "service": "support",
    "method": "PATCH",
    "path": "/api/support/cases/{publicId}/status",
    "tag": "Support",
    "description": "Transition support case status",
    "summary": "Transition support case status",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "support:POST:/api/support/cases/{publicId}/comments",
    "service": "support",
    "method": "POST",
    "path": "/api/support/cases/{publicId}/comments",
    "tag": "Support",
    "description": "Append comment to support case",
    "summary": "Append comment to support case",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:GET:/health/live",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço webhook-hub.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/health/ready",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço webhook-hub.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/health/details",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/health/details",
    "tag": "Webhook Hub",
    "description": "Return readiness details for webhook runtime",
    "summary": "Return readiness details for webhook runtime",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/capabilities",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/capabilities",
    "tag": "Webhook Hub",
    "description": "Read outbound webhook capability posture",
    "summary": "Read outbound webhook capability posture",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/outbound-endpoints",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/outbound-endpoints",
    "tag": "Webhook Hub",
    "description": "List tenant outbound endpoints",
    "summary": "List tenant outbound endpoints",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/outbound-endpoints",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/outbound-endpoints",
    "tag": "Webhook Hub",
    "description": "Register one tenant outbound endpoint",
    "summary": "Register one tenant outbound endpoint",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/outbound-endpoints/{publicId}",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/outbound-endpoints/{publicId}",
    "tag": "Webhook Hub",
    "description": "Read one tenant outbound endpoint",
    "summary": "Read one tenant outbound endpoint",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/outbound-endpoints/{publicId}/deliveries",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/outbound-endpoints/{publicId}/deliveries",
    "tag": "Webhook Hub",
    "description": "List outbound delivery log for one endpoint",
    "summary": "List outbound delivery log for one endpoint",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/outbound-endpoints/{publicId}/deliveries",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/outbound-endpoints/{publicId}/deliveries",
    "tag": "Webhook Hub",
    "description": "Register one outbound delivery attempt",
    "summary": "Register one outbound delivery attempt",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/outbound-endpoints/{publicId}/deliveries/{deliveryPublicId}/dead-letter",
    "tag": "Webhook Hub",
    "description": "Move one outbound delivery to dead letter",
    "summary": "Move one outbound delivery to dead letter",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId",
      "deliveryPublicId"
    ]
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/events",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/events",
    "tag": "Webhook Hub",
    "description": "List inbound webhook events",
    "summary": "List inbound webhook events",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events",
    "tag": "Webhook Hub",
    "description": "Register inbound webhook event",
    "summary": "Register inbound webhook event",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/events/summary",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/events/summary",
    "tag": "Webhook Hub",
    "description": "Aggregate inbound webhook state",
    "summary": "Aggregate inbound webhook state",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/dead-letter",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/dead-letter",
    "tag": "Webhook Hub",
    "description": "Move event to dead letter queue",
    "summary": "Move event to dead letter queue",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/requeue",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/requeue",
    "tag": "Webhook Hub",
    "description": "Requeue dead-letter event",
    "summary": "Requeue dead-letter event",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:GET:/health/live",
    "service": "workflow-control",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço workflow-control.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/health/ready",
    "service": "workflow-control",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço workflow-control.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/health/details",
    "service": "workflow-control",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço workflow-control.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/definitions",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions",
    "tag": "Workflow Control",
    "description": "List workflow definitions",
    "summary": "List workflow definitions",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/definitions",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/definitions",
    "tag": "Workflow Control",
    "description": "Create workflow definition",
    "summary": "Create workflow definition",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/definitions/{key}",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions/{key}",
    "tag": "Workflow Control",
    "description": "Read one workflow definition",
    "summary": "Read one workflow definition",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:PATCH:/api/workflow-control/definitions/{key}",
    "service": "workflow-control",
    "method": "PATCH",
    "path": "/api/workflow-control/definitions/{key}",
    "tag": "Workflow Control",
    "description": "Update one workflow definition",
    "summary": "Update one workflow definition",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:PATCH:/api/workflow-control/definitions/{key}/status",
    "service": "workflow-control",
    "method": "PATCH",
    "path": "/api/workflow-control/definitions/{key}/status",
    "tag": "Workflow Control",
    "description": "Update workflow definition status",
    "summary": "Update workflow definition status",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/capabilities/triggers",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/capabilities/triggers",
    "tag": "Workflow Control",
    "description": "List workflow trigger catalog",
    "summary": "List workflow trigger catalog",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/capabilities/actions",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/capabilities/actions",
    "tag": "Workflow Control",
    "description": "List workflow action catalog",
    "summary": "List workflow action catalog",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/health/live",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço workflow-runtime.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/health/ready",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço workflow-runtime.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/health/details",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço workflow-runtime.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions",
    "tag": "Workflow Runtime",
    "description": "List workflow executions",
    "summary": "List workflow executions",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions",
    "tag": "Workflow Runtime",
    "description": "Create workflow execution",
    "summary": "Create workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/{publicId}",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/{publicId}",
    "tag": "Workflow Runtime",
    "description": "Read one workflow execution",
    "summary": "Read one workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/{publicId}/actions",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/{publicId}/actions",
    "tag": "Workflow Runtime",
    "description": "List execution action snapshots",
    "summary": "List execution action snapshots",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/advance",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/advance",
    "tag": "Workflow Runtime",
    "description": "Advance one workflow execution",
    "summary": "Advance one workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:GET:/api/workflow-runtime/capabilities",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/capabilities",
    "tag": "Workflow Runtime",
    "description": "List runtime capabilities",
    "summary": "List runtime capabilities",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  }
];

export const eventSchemas: EventSchemaContract[] = [
  {
    "fileName": "catalog.item.schema.json",
    "name": "catalog.item",
    "source": "docs/contracts/events/catalog.item.schema.json"
  },
  {
    "fileName": "crm.cnpj-enrichment.schema.json",
    "name": "crm.cnpj-enrichment",
    "source": "docs/contracts/events/crm.cnpj-enrichment.schema.json"
  },
  {
    "fileName": "documents.signing-request.schema.json",
    "name": "documents.signing-request",
    "source": "docs/contracts/events/documents.signing-request.schema.json"
  },
  {
    "fileName": "engagement.provider-event.schema.json",
    "name": "engagement.provider-event",
    "source": "docs/contracts/events/engagement.provider-event.schema.json"
  },
  {
    "fileName": "fiscal.consent.schema.json",
    "name": "fiscal.consent",
    "source": "docs/contracts/events/fiscal.consent.schema.json"
  },
  {
    "fileName": "fiscal.document-event.schema.json",
    "name": "fiscal.document-event",
    "source": "docs/contracts/events/fiscal.document-event.schema.json"
  },
  {
    "fileName": "platform-control.go-live-rollout.schema.json",
    "name": "platform-control.go-live-rollout",
    "source": "docs/contracts/events/platform-control.go-live-rollout.schema.json"
  },
  {
    "fileName": "platform-control.lifecycle-job.schema.json",
    "name": "platform-control.lifecycle-job",
    "source": "docs/contracts/events/platform-control.lifecycle-job.schema.json"
  },
  {
    "fileName": "platform-control.quota.schema.json",
    "name": "platform-control.quota",
    "source": "docs/contracts/events/platform-control.quota.schema.json"
  },
  {
    "fileName": "support.case.schema.json",
    "name": "support.case",
    "source": "docs/contracts/events/support.case.schema.json"
  },
  {
    "fileName": "webhook-hub.inbound-event.schema.json",
    "name": "webhook-hub.inbound-event",
    "source": "docs/contracts/events/webhook-hub.inbound-event.schema.json"
  },
  {
    "fileName": "webhook-hub.outbound-delivery.schema.json",
    "name": "webhook-hub.outbound-delivery",
    "source": "docs/contracts/events/webhook-hub.outbound-delivery.schema.json"
  }
];
