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
    "slug": "accounting",
    "name": "Business Operating System Accounting",
    "title": "Business Operating System Accounting API",
    "version": "0.1.0",
    "description": "Managerial accounting, chart of accounts, immutable journal entries and accounting close",
    "contractFile": "docs/contracts/http/accounting.openapi.yaml",
    "endpointCount": 25
  },
  {
    "slug": "ai-governance",
    "name": "Business Operating System AI Governance",
    "title": "Business Operating System AI Governance API",
    "version": "0.1.0",
    "description": "Governed assistant runs, tool registry, policy enforcement and prompt audit.",
    "contractFile": "docs/contracts/http/ai-governance.openapi.yaml",
    "endpointCount": 6
  },
  {
    "slug": "analytics",
    "name": "Business Operating System Analytics",
    "title": "Business Operating System Analytics API",
    "version": "0.1.0",
    "description": "Executive reports, adapter catalog, SaaS control and contract governance.",
    "contractFile": "docs/contracts/http/analytics.openapi.yaml",
    "endpointCount": 62
  },
  {
    "slug": "banking",
    "name": "Business Operating System Banking",
    "title": "Business Operating System Banking API",
    "version": "0.1.0",
    "description": "Brazilian banking hub for bank accounts, CNAB, boleto, Pix, statements and reconciliation",
    "contractFile": "docs/contracts/http/banking.openapi.yaml",
    "endpointCount": 33
  },
  {
    "slug": "billing",
    "name": "Business Operating System Billing",
    "title": "Business Operating System Billing API",
    "version": "1.0.0",
    "description": "Subscription billing, invoices, gateway capabilities, payment attempts, webhooks and recovery operations.",
    "contractFile": "docs/contracts/http/billing.openapi.yaml",
    "endpointCount": 31
  },
  {
    "slug": "catalog",
    "name": "Business Operating System Catalog",
    "title": "Business Operating System Catalog API",
    "version": "0.2.0",
    "description": "Product and service catalog with categories, activation, versioned items, cursor pagination and bulk creation.",
    "contractFile": "docs/contracts/http/catalog.openapi.yaml",
    "endpointCount": 12
  },
  {
    "slug": "crm",
    "name": "Business Operating System CRM",
    "title": "Business Operating System CRM API",
    "version": "1.0.0",
    "description": "Leads, customers, ownership, history, attachments, pipeline intelligence and CRM outbox.",
    "contractFile": "docs/contracts/http/crm.openapi.yaml",
    "endpointCount": 26
  },
  {
    "slug": "documents",
    "name": "Business Operating System Documents",
    "title": "Business Operating System Documents API",
    "version": "0.9.7",
    "description": "Attachment governance, upload orchestration and storage capabilities.",
    "contractFile": "docs/contracts/http/documents.openapi.yaml",
    "endpointCount": 19
  },
  {
    "slug": "edge",
    "name": "Business Operating System Edge",
    "title": "Business Operating System Edge API",
    "version": "0.1.0",
    "description": "Aggregated operational cockpits for tenants, contracts and SaaS control.",
    "contractFile": "docs/contracts/http/edge.openapi.yaml",
    "endpointCount": 19
  },
  {
    "slug": "engagement",
    "name": "Business Operating System Engagement",
    "title": "Business Operating System Engagement API",
    "version": "0.9.7",
    "description": "Omnichannel engagement, provider callbacks and campaign operations.",
    "contractFile": "docs/contracts/http/engagement.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "finance",
    "name": "Business Operating System Finance",
    "title": "Business Operating System Finance API",
    "version": "1.0.0",
    "description": "Receivables, projections, settlements, commission lifecycle, payables, costs, treasury, period closures and financial activity.",
    "contractFile": "docs/contracts/http/finance.openapi.yaml",
    "endpointCount": 26
  },
  {
    "slug": "fiscal",
    "name": "Business Operating System Fiscal",
    "title": "Business Operating System Fiscal API",
    "version": "0.1.0",
    "description": "Fiscal profile, document operations, privacy rights and compliance governance.",
    "contractFile": "docs/contracts/http/fiscal.openapi.yaml",
    "endpointCount": 37
  },
  {
    "slug": "identity",
    "name": "Business Operating System Identity",
    "title": "Business Operating System Identity API",
    "version": "0.9.0",
    "description": "Tenancy, access, sessions, invitations, MFA, recovery and tenant-scoped identity governance.",
    "contractFile": "docs/contracts/http/identity.openapi.yaml",
    "endpointCount": 46
  },
  {
    "slug": "inventory",
    "name": "Business Operating System Inventory",
    "title": "Business Operating System Inventory API",
    "version": "0.1.0",
    "description": "Inventory balances, warehouse locations, movements, reservations and cycle counts",
    "contractFile": "docs/contracts/http/inventory.openapi.yaml",
    "endpointCount": 23
  },
  {
    "slug": "notification",
    "name": "Business Operating System Notification",
    "title": "Business Operating System Notification API",
    "version": "0.1.0",
    "description": "Internal notification center, preferences and reusable dispatch shape.",
    "contractFile": "docs/contracts/http/notification.openapi.yaml",
    "endpointCount": 8
  },
  {
    "slug": "platform-control",
    "name": "Business Operating System Platform Control",
    "title": "Business Operating System Platform Control API",
    "version": "0.2.0",
    "description": "Tenant capabilities, entitlements, quotas, metering, lifecycle jobs and SaaS governance.",
    "contractFile": "docs/contracts/http/platform-control.openapi.yaml",
    "endpointCount": 92
  },
  {
    "slug": "procurement",
    "name": "Business Operating System Procurement",
    "title": "Business Operating System Procurement API",
    "version": "0.1.0",
    "description": "Purchase requisitions, quotations, purchase orders, receipts and three-way matching",
    "contractFile": "docs/contracts/http/procurement.openapi.yaml",
    "endpointCount": 25
  },
  {
    "slug": "rentals",
    "name": "Business Operating System Rentals",
    "title": "Business Operating System Rentals API",
    "version": "0.9.0",
    "description": "Rental contracts, recurring charges, adjustments, rescission and attachment linkage.",
    "contractFile": "docs/contracts/http/rentals.openapi.yaml",
    "endpointCount": 12
  },
  {
    "slug": "sales",
    "name": "Business Operating System Sales",
    "title": "Business Operating System Sales API",
    "version": "1.0.0",
    "description": "Opportunities, proposals, sales, invoices, installments, commissions, pending items, renegotiations and commercial outbox.",
    "contractFile": "docs/contracts/http/sales.openapi.yaml",
    "endpointCount": 37
  },
  {
    "slug": "search",
    "name": "Business Operating System Operational Search",
    "title": "Business Operating System Operational Search API",
    "version": "0.1.0",
    "description": "Unified operational search, audit, legal hold and e-discovery workflows.",
    "contractFile": "docs/contracts/http/search.openapi.yaml",
    "endpointCount": 9
  },
  {
    "slug": "simulation",
    "name": "Business Operating System Simulation",
    "title": "Business Operating System Simulation API",
    "version": "0.8.0",
    "description": "Scenario simulation, load benchmark and cost estimation runtime.",
    "contractFile": "docs/contracts/http/simulation.openapi.yaml",
    "endpointCount": 6
  },
  {
    "slug": "supplier",
    "name": "Business Operating System Supplier",
    "title": "Business Operating System Supplier API",
    "version": "0.1.0",
    "description": "Supplier directory, categories and payables-oriented vendor governance.",
    "contractFile": "docs/contracts/http/supplier.openapi.yaml",
    "endpointCount": 10
  },
  {
    "slug": "support",
    "name": "Business Operating System Support",
    "title": "Business Operating System Support API",
    "version": "0.1.0",
    "description": "Queue-based support cases with SLA, comments and lifecycle history.",
    "contractFile": "docs/contracts/http/support.openapi.yaml",
    "endpointCount": 11
  },
  {
    "slug": "webhook-hub",
    "name": "Business Operating System Webhook Hub",
    "title": "Business Operating System Webhook Hub API",
    "version": "0.9.7",
    "description": "Inbound webhook intake, DLQ and operator recovery surface.",
    "contractFile": "docs/contracts/http/webhook-hub.openapi.yaml",
    "endpointCount": 22
  },
  {
    "slug": "workflow-control",
    "name": "Business Operating System Workflow Control",
    "title": "Business Operating System Workflow Control API",
    "version": "0.9.0",
    "description": "Workflow definition catalog, publication, versioning, execution control and audit events.",
    "contractFile": "docs/contracts/http/workflow-control.openapi.yaml",
    "endpointCount": 25
  },
  {
    "slug": "workflow-runtime",
    "name": "Business Operating System Workflow Runtime",
    "title": "Business Operating System Workflow Runtime API",
    "version": "0.9.0",
    "description": "Workflow execution runtime with transitions, action snapshots, retries, delays and compensations.",
    "contractFile": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "endpointCount": 15
  }
];

export const endpoints: EndpointContract[] = [
  {
    "id": "accounting:GET:/health/live",
    "service": "accounting",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço accounting.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/health/ready",
    "service": "accounting",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço accounting.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/health/details",
    "service": "accounting",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço accounting.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/capabilities",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/capabilities",
    "tag": "Accounting",
    "description": "Read accounting capability catalog",
    "summary": "Read accounting capability catalog",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/statements/management-summary",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/statements/management-summary",
    "tag": "Accounting",
    "description": "Read accounting operational summary",
    "summary": "Read accounting operational summary",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/accounts",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/accounts",
    "tag": "Accounting",
    "description": "List accounts",
    "summary": "List accounts",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:POST:/api/accounting/accounts",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/accounts",
    "tag": "Accounting",
    "description": "Create account",
    "summary": "Create account",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/accounts/{publicId}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/accounts/{publicId}",
    "tag": "Accounting",
    "description": "Read one account",
    "summary": "Read one account",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:PATCH:/api/accounting/accounts/{publicId}/status",
    "service": "accounting",
    "method": "PATCH",
    "path": "/api/accounting/accounts/{publicId}/status",
    "tag": "Accounting",
    "description": "Transition account status",
    "summary": "Transition account status",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:GET:/api/accounting/journal-entries",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/journal-entries",
    "tag": "Accounting",
    "description": "List journal-entries",
    "summary": "List journal-entries",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:POST:/api/accounting/journal-entries",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/journal-entries",
    "tag": "Accounting",
    "description": "Create journal entry",
    "summary": "Create journal entry",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/journal-entries/{publicId}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/journal-entries/{publicId}",
    "tag": "Accounting",
    "description": "Read one journal entry",
    "summary": "Read one journal entry",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:PATCH:/api/accounting/journal-entries/{publicId}/status",
    "service": "accounting",
    "method": "PATCH",
    "path": "/api/accounting/journal-entries/{publicId}/status",
    "tag": "Accounting",
    "description": "Transition journal entry status",
    "summary": "Transition journal entry status",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:GET:/api/accounting/period-closes",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/period-closes",
    "tag": "Accounting",
    "description": "List period-closes",
    "summary": "List period-closes",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:POST:/api/accounting/period-closes",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/period-closes",
    "tag": "Accounting",
    "description": "Create period close",
    "summary": "Create period close",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/period-closes/{publicId}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/period-closes/{publicId}",
    "tag": "Accounting",
    "description": "Read one period close",
    "summary": "Read one period close",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:PATCH:/api/accounting/period-closes/{publicId}/status",
    "service": "accounting",
    "method": "PATCH",
    "path": "/api/accounting/period-closes/{publicId}/status",
    "tag": "Accounting",
    "description": "Transition period close status",
    "summary": "Transition period close status",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:GET:/api/accounting/ledger",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/ledger",
    "tag": "Accounting",
    "description": "GET Accounting Ledger",
    "summary": "GET Accounting Ledger",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/statements/{statement_kind}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/statements/{statement_kind}",
    "tag": "Accounting",
    "description": "GET Statements",
    "summary": "GET Statements",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "statement_kind"
    ]
  },
  {
    "id": "accounting:POST:/api/accounting/posting-rules/apply",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/posting-rules/apply",
    "tag": "Accounting",
    "description": "POST Posting Rules Apply",
    "summary": "POST Posting Rules Apply",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/cost-centers",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/cost-centers",
    "tag": "Accounting",
    "description": "GET Accounting Cost Centers",
    "summary": "GET Accounting Cost Centers",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:POST:/api/accounting/cost-centers",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/cost-centers",
    "tag": "Accounting",
    "description": "Create cost center",
    "summary": "Create cost center",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/cost-centers/{publicId}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/cost-centers/{publicId}",
    "tag": "Accounting",
    "description": "GET Cost Centers",
    "summary": "GET Cost Centers",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:PATCH:/api/accounting/cost-centers/{publicId}/status",
    "service": "accounting",
    "method": "PATCH",
    "path": "/api/accounting/cost-centers/{publicId}/status",
    "tag": "Accounting",
    "description": "PATCH Status",
    "summary": "PATCH Status",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:GET:/api/accounting/posting-rules",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/posting-rules",
    "tag": "Accounting",
    "description": "GET Accounting Posting Rules",
    "summary": "GET Accounting Posting Rules",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "accounting:POST:/api/accounting/posting-rules",
    "service": "accounting",
    "method": "POST",
    "path": "/api/accounting/posting-rules",
    "tag": "Accounting",
    "description": "Create posting rule",
    "summary": "Create posting rule",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "accounting:GET:/api/accounting/posting-rules/{publicId}",
    "service": "accounting",
    "method": "GET",
    "path": "/api/accounting/posting-rules/{publicId}",
    "tag": "Accounting",
    "description": "GET Posting Rules",
    "summary": "GET Posting Rules",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "accounting:PATCH:/api/accounting/posting-rules/{publicId}/status",
    "service": "accounting",
    "method": "PATCH",
    "path": "/api/accounting/posting-rules/{publicId}/status",
    "tag": "Accounting",
    "description": "PATCH Status",
    "summary": "PATCH Status",
    "source": "docs/contracts/http/accounting.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "ai-governance:GET:/health/live",
    "service": "ai-governance",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço ai-governance.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:GET:/health/ready",
    "service": "ai-governance",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço ai-governance.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:GET:/health/details",
    "service": "ai-governance",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço ai-governance.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:GET:/api/ai-governance/tools",
    "service": "ai-governance",
    "method": "GET",
    "path": "/api/ai-governance/tools",
    "tag": "Ai Governance",
    "description": "List approved assistant tools",
    "summary": "List approved assistant tools",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:GET:/api/ai-governance/policies",
    "service": "ai-governance",
    "method": "GET",
    "path": "/api/ai-governance/policies",
    "tag": "Ai Governance",
    "description": "List AI governance policies",
    "summary": "List AI governance policies",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:POST:/api/ai-governance/runs",
    "service": "ai-governance",
    "method": "POST",
    "path": "/api/ai-governance/runs",
    "tag": "Ai Governance",
    "description": "Create a governed assistant run",
    "summary": "Create a governed assistant run",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "ai-governance:GET:/api/ai-governance/runs/{publicId}",
    "service": "ai-governance",
    "method": "GET",
    "path": "/api/ai-governance/runs/{publicId}",
    "tag": "Ai Governance",
    "description": "Read one assistant run and tool decisions",
    "summary": "Read one assistant run and tool decisions",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "ai-governance:GET:/api/ai-governance/audit-events",
    "service": "ai-governance",
    "method": "GET",
    "path": "/api/ai-governance/audit-events",
    "tag": "Ai Governance",
    "description": "List AI governance audit events",
    "summary": "List AI governance audit events",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "ai-governance:POST:/api/ai-governance/redaction/preview",
    "service": "ai-governance",
    "method": "POST",
    "path": "/api/ai-governance/redaction/preview",
    "tag": "Ai Governance",
    "description": "Preview prompt redaction before assistant execution",
    "summary": "Preview prompt redaction before assistant execution",
    "source": "docs/contracts/http/ai-governance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
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
    "id": "analytics:GET:/api/analytics/reports/service-pulse",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/service-pulse",
    "tag": "Analytics",
    "description": "Read service pulse report",
    "summary": "Read service pulse report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/tenant-360",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/tenant-360",
    "tag": "Analytics",
    "description": "Read tenant 360 report",
    "summary": "Read tenant 360 report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/pipeline-summary",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/pipeline-summary",
    "tag": "Analytics",
    "description": "Read pipeline summary report",
    "summary": "Read pipeline summary report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/workflow-definition-health",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/workflow-definition-health",
    "tag": "Analytics",
    "description": "Read workflow definition health report",
    "summary": "Read workflow definition health report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/automation-board",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/automation-board",
    "tag": "Analytics",
    "description": "Read automation board report",
    "summary": "Read automation board report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/sales-journey",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/sales-journey",
    "tag": "Analytics",
    "description": "Read sales journey report",
    "summary": "Read sales journey report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/revenue-operations",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/revenue-operations",
    "tag": "Analytics",
    "description": "Read revenue operations report",
    "summary": "Read revenue operations report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/finance-control",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/finance-control",
    "tag": "Analytics",
    "description": "Read finance control report",
    "summary": "Read finance control report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/collections-control",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/collections-control",
    "tag": "Analytics",
    "description": "Read collections control report",
    "summary": "Read collections control report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/document-governance",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/document-governance",
    "tag": "Analytics",
    "description": "Read document governance report",
    "summary": "Read document governance report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/rental-operations",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/rental-operations",
    "tag": "Analytics",
    "description": "Read rental operations report",
    "summary": "Read rental operations report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/engagement-operations",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/engagement-operations",
    "tag": "Analytics",
    "description": "Read engagement operations report",
    "summary": "Read engagement operations report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/cost-estimator",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/cost-estimator",
    "tag": "Analytics",
    "description": "Read cost estimator report",
    "summary": "Read cost estimator report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/delivery-reliability",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/delivery-reliability",
    "tag": "Analytics",
    "description": "Read delivery reliability report",
    "summary": "Read delivery reliability report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/load-benchmark",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/load-benchmark",
    "tag": "Analytics",
    "description": "Read load benchmark report",
    "summary": "Read load benchmark report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reports/platform-reliability",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/platform-reliability",
    "tag": "Analytics",
    "description": "Read platform reliability report",
    "summary": "Read platform reliability report",
    "source": "docs/contracts/http/analytics.openapi.yaml",
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
    "id": "analytics:GET:/api/analytics/reports/production-readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reports/production-readiness",
    "tag": "Analytics",
    "description": "Read production readiness for version 1.4.6",
    "summary": "Read production readiness for version 1.4.6",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/metrics",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/metrics",
    "tag": "Analytics",
    "description": "List semantic metric definitions",
    "summary": "List semantic metric definitions",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/metrics/{code}",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/metrics/{code}",
    "tag": "Analytics",
    "description": "Read one semantic metric definition",
    "summary": "Read one semantic metric definition",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "code"
    ]
  },
  {
    "id": "analytics:GET:/api/analytics/metrics/{code}/snapshots",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/metrics/{code}/snapshots",
    "tag": "Analytics",
    "description": "List semantic metric snapshots",
    "summary": "List semantic metric snapshots",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "code"
    ]
  },
  {
    "id": "analytics:GET:/api/analytics/datasets/freshness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/datasets/freshness",
    "tag": "Analytics",
    "description": "List dataset freshness posture",
    "summary": "List dataset freshness posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/data-quality",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/data-quality",
    "tag": "Analytics",
    "description": "List semantic data quality checks",
    "summary": "List semantic data quality checks",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/lineage",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lineage",
    "tag": "Analytics",
    "description": "Read semantic metric lineage graph",
    "summary": "Read semantic metric lineage graph",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/risk/tenant-score",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/risk/tenant-score",
    "tag": "Analytics",
    "description": "Read tenant risk and compliance score",
    "summary": "Read tenant risk and compliance score",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/risk/domain-scores",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/risk/domain-scores",
    "tag": "Analytics",
    "description": "Read domain risk scores",
    "summary": "Read domain risk scores",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/risk/service-scores",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/risk/service-scores",
    "tag": "Analytics",
    "description": "Read service risk scores",
    "summary": "Read service risk scores",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/risk/compliance-posture",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/risk/compliance-posture",
    "tag": "Analytics",
    "description": "Read compliance posture",
    "summary": "Read compliance posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/risk/recommendations",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/risk/recommendations",
    "tag": "Analytics",
    "description": "Read risk recommendations",
    "summary": "Read risk recommendations",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/reconciliation/findings",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/reconciliation/findings",
    "tag": "Analytics",
    "description": "List reconciliation findings",
    "summary": "List reconciliation findings",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:POST:/api/analytics/reconciliation/run",
    "service": "analytics",
    "method": "POST",
    "path": "/api/analytics/reconciliation/run",
    "tag": "Analytics",
    "description": "Run operational reconciliation",
    "summary": "Run operational reconciliation",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/financial-close/periods",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/financial-close/periods",
    "tag": "Analytics",
    "description": "List financial close periods",
    "summary": "List financial close periods",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:POST:/api/analytics/financial-close/periods",
    "service": "analytics",
    "method": "POST",
    "path": "/api/analytics/financial-close/periods",
    "tag": "Analytics",
    "description": "Create one financial close period",
    "summary": "Create one financial close period",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "analytics:POST:/api/analytics/financial-close/periods/{publicId}/close",
    "service": "analytics",
    "method": "POST",
    "path": "/api/analytics/financial-close/periods/{publicId}/close",
    "tag": "Analytics",
    "description": "Close one financial period",
    "summary": "Close one financial period",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "analytics:GET:/api/analytics/financial-close/periods/{publicId}/snapshot",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/financial-close/periods/{publicId}/snapshot",
    "tag": "Analytics",
    "description": "Read financial close snapshot",
    "summary": "Read financial close snapshot",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "analytics:GET:/api/analytics/financial-close/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/financial-close/readiness",
    "tag": "Analytics",
    "description": "Read financial close readiness",
    "summary": "Read financial close readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/master-data/entities",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/master-data/entities",
    "tag": "Analytics",
    "description": "List master data entities",
    "summary": "List master data entities",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/master-data/quality-score",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/master-data/quality-score",
    "tag": "Analytics",
    "description": "Read master data quality score",
    "summary": "Read master data quality score",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/master-data/duplicates",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/master-data/duplicates",
    "tag": "Analytics",
    "description": "List master data duplicate candidates",
    "summary": "List master data duplicate candidates",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:POST:/api/analytics/master-data/merge-proposals",
    "service": "analytics",
    "method": "POST",
    "path": "/api/analytics/master-data/merge-proposals",
    "tag": "Analytics",
    "description": "Create one master data merge proposal",
    "summary": "Create one master data merge proposal",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/data-quality/rules",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/data-quality/rules",
    "tag": "Analytics",
    "description": "List data quality rules",
    "summary": "List data quality rules",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/data-quality/findings",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/data-quality/findings",
    "tag": "Analytics",
    "description": "List data quality findings",
    "summary": "List data quality findings",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/lakehouse/datasets",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lakehouse/datasets",
    "tag": "Analytics",
    "description": "List lakehouse datasets",
    "summary": "List lakehouse datasets",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/lakehouse/datasets/{datasetKey}",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lakehouse/datasets/{datasetKey}",
    "tag": "Analytics",
    "description": "Read one lakehouse dataset",
    "summary": "Read one lakehouse dataset",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "datasetKey"
    ]
  },
  {
    "id": "analytics:GET:/api/analytics/lakehouse/lineage",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lakehouse/lineage",
    "tag": "Analytics",
    "description": "Read lakehouse lineage",
    "summary": "Read lakehouse lineage",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/lakehouse/export-policies",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lakehouse/export-policies",
    "tag": "Analytics",
    "description": "List lakehouse export policies",
    "summary": "List lakehouse export policies",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/lakehouse/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/lakehouse/readiness",
    "tag": "Analytics",
    "description": "Read lakehouse readiness",
    "summary": "Read lakehouse readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/enterprise-runtime/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/enterprise-runtime/readiness",
    "tag": "Analytics",
    "description": "Read enterprise runtime readiness",
    "summary": "Read enterprise runtime readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/external-intelligence/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/external-intelligence/readiness",
    "tag": "Analytics",
    "description": "Read external intelligence and verification readiness",
    "summary": "Read external intelligence and verification readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/document-intelligence/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/document-intelligence/readiness",
    "tag": "Analytics",
    "description": "Read OCR and document intelligence provider readiness",
    "summary": "Read OCR and document intelligence provider readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/fiscal-brazil/readiness",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/fiscal-brazil/readiness",
    "tag": "Analytics",
    "description": "Read Brazilian fiscal provider readiness",
    "summary": "Read Brazilian fiscal provider readiness",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/registry-enrichment/brazil",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/registry-enrichment/brazil",
    "tag": "Analytics",
    "description": "Read Brazilian registry enrichment posture",
    "summary": "Read Brazilian registry enrichment posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/market-macro-risk",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/market-macro-risk",
    "tag": "Analytics",
    "description": "Read FX, market and macro risk provider posture",
    "summary": "Read FX, market and macro risk provider posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "analytics:GET:/api/analytics/external-risk-feed",
    "service": "analytics",
    "method": "GET",
    "path": "/api/analytics/external-risk-feed",
    "tag": "Analytics",
    "description": "Read external news and risk feed posture",
    "summary": "Read external news and risk feed posture",
    "source": "docs/contracts/http/analytics.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/health/live",
    "service": "banking",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço banking.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/health/ready",
    "service": "banking",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço banking.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/health/details",
    "service": "banking",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço banking.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/capabilities",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/capabilities",
    "tag": "Banking",
    "description": "Read banking capability catalog",
    "summary": "Read banking capability catalog",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/reconciliation/summary",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/reconciliation/summary",
    "tag": "Banking",
    "description": "Read banking operational summary",
    "summary": "Read banking operational summary",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/bank-accounts",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/bank-accounts",
    "tag": "Banking",
    "description": "List bank-accounts",
    "summary": "List bank-accounts",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/bank-accounts",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/bank-accounts",
    "tag": "Banking",
    "description": "Create bank account",
    "summary": "Create bank account",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/bank-accounts/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/bank-accounts/{publicId}",
    "tag": "Banking",
    "description": "Read one bank account",
    "summary": "Read one bank account",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:PATCH:/api/banking/bank-accounts/{publicId}/status",
    "service": "banking",
    "method": "PATCH",
    "path": "/api/banking/bank-accounts/{publicId}/status",
    "tag": "Banking",
    "description": "Transition bank account status",
    "summary": "Transition bank account status",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:GET:/api/banking/cnab-files",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/cnab-files",
    "tag": "Banking",
    "description": "List cnab-files",
    "summary": "List cnab-files",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/cnab-files",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/cnab-files",
    "tag": "Banking",
    "description": "Create cnab file",
    "summary": "Create cnab file",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/cnab-files/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/cnab-files/{publicId}",
    "tag": "Banking",
    "description": "Read one cnab file",
    "summary": "Read one cnab file",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:PATCH:/api/banking/cnab-files/{publicId}/status",
    "service": "banking",
    "method": "PATCH",
    "path": "/api/banking/cnab-files/{publicId}/status",
    "tag": "Banking",
    "description": "Transition cnab file status",
    "summary": "Transition cnab file status",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:GET:/api/banking/boletos",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/boletos",
    "tag": "Banking",
    "description": "List boletos",
    "summary": "List boletos",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/boletos",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/boletos",
    "tag": "Banking",
    "description": "Create boleto",
    "summary": "Create boleto",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/boletos/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/boletos/{publicId}",
    "tag": "Banking",
    "description": "Read one boleto",
    "summary": "Read one boleto",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:PATCH:/api/banking/boletos/{publicId}/status",
    "service": "banking",
    "method": "PATCH",
    "path": "/api/banking/boletos/{publicId}/status",
    "tag": "Banking",
    "description": "Transition boleto status",
    "summary": "Transition boleto status",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:GET:/api/banking/pix-charges",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/pix-charges",
    "tag": "Banking",
    "description": "List pix-charges",
    "summary": "List pix-charges",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/pix-charges",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/pix-charges",
    "tag": "Banking",
    "description": "Create pix charge",
    "summary": "Create pix charge",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/pix-charges/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/pix-charges/{publicId}",
    "tag": "Banking",
    "description": "Read one pix charge",
    "summary": "Read one pix charge",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:PATCH:/api/banking/pix-charges/{publicId}/status",
    "service": "banking",
    "method": "PATCH",
    "path": "/api/banking/pix-charges/{publicId}/status",
    "tag": "Banking",
    "description": "Transition pix charge status",
    "summary": "Transition pix charge status",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:GET:/api/banking/reconciliations",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/reconciliations",
    "tag": "Banking",
    "description": "List reconciliations",
    "summary": "List reconciliations",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/reconciliations",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/reconciliations",
    "tag": "Banking",
    "description": "Create reconciliation",
    "summary": "Create reconciliation",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/reconciliations/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/reconciliations/{publicId}",
    "tag": "Banking",
    "description": "Read one reconciliation",
    "summary": "Read one reconciliation",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:PATCH:/api/banking/reconciliations/{publicId}/status",
    "service": "banking",
    "method": "PATCH",
    "path": "/api/banking/reconciliations/{publicId}/status",
    "tag": "Banking",
    "description": "Transition reconciliation status",
    "summary": "Transition reconciliation status",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:POST:/api/banking/cnab-files/parse-return",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/cnab-files/parse-return",
    "tag": "Banking",
    "description": "POST Cnab Files Parse Return",
    "summary": "POST Cnab Files Parse Return",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/reconciliations/run",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/reconciliations/run",
    "tag": "Banking",
    "description": "POST Reconciliations Run",
    "summary": "POST Reconciliations Run",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/bank-statements",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/bank-statements",
    "tag": "Banking",
    "description": "GET Banking Bank Statements",
    "summary": "GET Banking Bank Statements",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/bank-statements",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/bank-statements",
    "tag": "Banking",
    "description": "Create bank statement",
    "summary": "Create bank statement",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/bank-statements/{publicId}",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/bank-statements/{publicId}",
    "tag": "Banking",
    "description": "GET Bank Statements",
    "summary": "GET Bank Statements",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "banking:GET:/api/banking/pix-refunds",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/pix-refunds",
    "tag": "Banking",
    "description": "GET Banking Pix Refunds",
    "summary": "GET Banking Pix Refunds",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/pix-refunds",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/pix-refunds",
    "tag": "Banking",
    "description": "Create Pix refund",
    "summary": "Create Pix refund",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/pix-webhooks",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/pix-webhooks",
    "tag": "Banking",
    "description": "GET Banking Pix Webhooks",
    "summary": "GET Banking Pix Webhooks",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/pix-webhooks",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/pix-webhooks",
    "tag": "Banking",
    "description": "Ingest Pix webhook",
    "summary": "Ingest Pix webhook",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "banking:GET:/api/banking/open-finance-connections",
    "service": "banking",
    "method": "GET",
    "path": "/api/banking/open-finance-connections",
    "tag": "Banking",
    "description": "GET Banking Open Finance Connections",
    "summary": "GET Banking Open Finance Connections",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "banking:POST:/api/banking/open-finance-connections",
    "service": "banking",
    "method": "POST",
    "path": "/api/banking/open-finance-connections",
    "tag": "Banking",
    "description": "Create Open Finance connection",
    "summary": "Create Open Finance connection",
    "source": "docs/contracts/http/banking.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/events",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/events",
    "tag": "Billing",
    "description": "Read billing events",
    "summary": "Read billing events",
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
    "description": "Read billing gateways",
    "summary": "Read billing gateways",
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
    "description": "Read billing gateways provider",
    "summary": "Read billing gateways provider",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "provider"
    ]
  },
  {
    "id": "billing:GET:/api/billing/invoices",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/invoices",
    "tag": "Billing",
    "description": "Read billing invoices",
    "summary": "Read billing invoices",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/invoices/{publicId}",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/invoices/{publicId}",
    "tag": "Billing",
    "description": "Read billing invoices publicId",
    "summary": "Read billing invoices publicId",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/invoices/{publicId}/attempts",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/invoices/{publicId}/attempts",
    "tag": "Billing",
    "description": "Read billing invoices publicId attempts",
    "summary": "Read billing invoices publicId attempts",
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
    "description": "Create or execute billing invoices publicId attempts",
    "summary": "Create or execute billing invoices publicId attempts",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/invoices/{publicId}/recovery/open",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/invoices/{publicId}/recovery/open",
    "tag": "Billing",
    "description": "Create or execute billing invoices publicId recovery open",
    "summary": "Create or execute billing invoices publicId recovery open",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/plans",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/plans",
    "tag": "Billing",
    "description": "Read billing plans",
    "summary": "Read billing plans",
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
    "description": "Create or execute billing plans",
    "summary": "Create or execute billing plans",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/recovery/cases",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/recovery/cases",
    "tag": "Billing",
    "description": "Read billing recovery cases",
    "summary": "Read billing recovery cases",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/recovery/cases/{publicId}",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/recovery/cases/{publicId}",
    "tag": "Billing",
    "description": "Read billing recovery cases publicId",
    "summary": "Read billing recovery cases publicId",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/recovery/cases/{publicId}/actions",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/recovery/cases/{publicId}/actions",
    "tag": "Billing",
    "description": "Read billing recovery cases publicId actions",
    "summary": "Read billing recovery cases publicId actions",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/recovery/cases/{publicId}/close",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/recovery/cases/{publicId}/close",
    "tag": "Billing",
    "description": "Create or execute billing recovery cases publicId close",
    "summary": "Create or execute billing recovery cases publicId close",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/recovery/cases/{publicId}/promise",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/recovery/cases/{publicId}/promise",
    "tag": "Billing",
    "description": "Create or execute billing recovery cases publicId promise",
    "summary": "Create or execute billing recovery cases publicId promise",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/recovery/cases/{publicId}/touchpoints",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/recovery/cases/{publicId}/touchpoints",
    "tag": "Billing",
    "description": "Create or execute billing recovery cases publicId touchpoints",
    "summary": "Create or execute billing recovery cases publicId touchpoints",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/reports/operations",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/reports/operations",
    "tag": "Billing",
    "description": "Read billing reports operations",
    "summary": "Read billing reports operations",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/subscriptions",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions",
    "tag": "Billing",
    "description": "Read billing subscriptions",
    "summary": "Read billing subscriptions",
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
    "description": "Create or execute billing subscriptions",
    "summary": "Create or execute billing subscriptions",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/api/billing/subscriptions/{publicId}",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions/{publicId}",
    "tag": "Billing",
    "description": "Read billing subscriptions publicId",
    "summary": "Read billing subscriptions publicId",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/subscriptions/{publicId}/events",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions/{publicId}/events",
    "tag": "Billing",
    "description": "Read billing subscriptions publicId events",
    "summary": "Read billing subscriptions publicId events",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/subscriptions/{publicId}/invoices",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/subscriptions/{publicId}/invoices",
    "tag": "Billing",
    "description": "Create or execute billing subscriptions publicId invoices",
    "summary": "Create or execute billing subscriptions publicId invoices",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/subscriptions/{publicId}/reactivate",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/subscriptions/{publicId}/reactivate",
    "tag": "Billing",
    "description": "Create or execute billing subscriptions publicId reactivate",
    "summary": "Create or execute billing subscriptions publicId reactivate",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:POST:/api/billing/subscriptions/{publicId}/suspend",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/subscriptions/{publicId}/suspend",
    "tag": "Billing",
    "description": "Create or execute billing subscriptions publicId suspend",
    "summary": "Create or execute billing subscriptions publicId suspend",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/subscriptions/{publicId}/usage-pricing",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/subscriptions/{publicId}/usage-pricing",
    "tag": "Billing",
    "description": "Read billing subscriptions publicId usage pricing",
    "summary": "Read billing subscriptions publicId usage pricing",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "billing:GET:/api/billing/webhook-events/pending",
    "service": "billing",
    "method": "GET",
    "path": "/api/billing/webhook-events/pending",
    "tag": "Billing",
    "description": "Read billing webhook events pending",
    "summary": "Read billing webhook events pending",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:POST:/api/billing/webhook-events/process",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/webhook-events/process",
    "tag": "Billing",
    "description": "Create or execute billing webhook events process",
    "summary": "Create or execute billing webhook events process",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:POST:/api/billing/webhook-events/process-batch",
    "service": "billing",
    "method": "POST",
    "path": "/api/billing/webhook-events/process-batch",
    "tag": "Billing",
    "description": "Create or execute billing webhook events process batch",
    "summary": "Create or execute billing webhook events process batch",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/details",
    "service": "billing",
    "method": "GET",
    "path": "/health/details",
    "tag": "Billing",
    "description": "Read health details",
    "summary": "Read health details",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/live",
    "service": "billing",
    "method": "GET",
    "path": "/health/live",
    "tag": "Billing",
    "description": "Read health live",
    "summary": "Read health live",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "billing:GET:/health/ready",
    "service": "billing",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Billing",
    "description": "Read health ready",
    "summary": "Read health ready",
    "source": "docs/contracts/http/billing.openapi.yaml",
    "hasBody": false,
    "pathParams": []
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
    "id": "crm:GET:/api/crm/leads/summary",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/summary",
    "tag": "Crm",
    "description": "Read lead summary",
    "summary": "Read lead summary",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/leads",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads",
    "tag": "Crm",
    "description": "List leads",
    "summary": "List leads",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:POST:/api/crm/leads",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/leads",
    "tag": "Crm",
    "description": "Create lead",
    "summary": "Create lead",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/leads/export",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/export",
    "tag": "Crm",
    "description": "Export filtered leads",
    "summary": "Export filtered leads",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:POST:/api/crm/leads/bulk",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/leads/bulk",
    "tag": "Crm",
    "description": "Create leads in bulk",
    "summary": "Create leads in bulk",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/leads/{publicId}",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/{publicId}",
    "tag": "Crm",
    "description": "Read lead by public id",
    "summary": "Read lead by public id",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:PATCH:/api/crm/leads/{publicId}",
    "service": "crm",
    "method": "PATCH",
    "path": "/api/crm/leads/{publicId}",
    "tag": "Crm",
    "description": "Update lead profile",
    "summary": "Update lead profile",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:POST:/api/crm/leads/{publicId}/convert",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/leads/{publicId}/convert",
    "tag": "Crm",
    "description": "Convert lead to customer",
    "summary": "Convert lead to customer",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/leads/{publicId}/history",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/{publicId}/history",
    "tag": "Crm",
    "description": "List lead relationship history",
    "summary": "List lead relationship history",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/leads/{publicId}/notes",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/{publicId}/notes",
    "tag": "Crm",
    "description": "List lead notes",
    "summary": "List lead notes",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:POST:/api/crm/leads/{publicId}/notes",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/leads/{publicId}/notes",
    "tag": "Crm",
    "description": "Create lead note",
    "summary": "Create lead note",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/leads/{publicId}/attachments",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/leads/{publicId}/attachments",
    "tag": "Crm",
    "description": "List lead attachments",
    "summary": "List lead attachments",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:POST:/api/crm/leads/{publicId}/attachments",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/leads/{publicId}/attachments",
    "tag": "Crm",
    "description": "Create lead attachment metadata",
    "summary": "Create lead attachment metadata",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:PATCH:/api/crm/leads/{publicId}/owner",
    "service": "crm",
    "method": "PATCH",
    "path": "/api/crm/leads/{publicId}/owner",
    "tag": "Crm",
    "description": "Update lead owner",
    "summary": "Update lead owner",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:PATCH:/api/crm/leads/{publicId}/status",
    "service": "crm",
    "method": "PATCH",
    "path": "/api/crm/leads/{publicId}/status",
    "tag": "Crm",
    "description": "Update lead status",
    "summary": "Update lead status",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/customers",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/customers",
    "tag": "Crm",
    "description": "List customers",
    "summary": "List customers",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "crm:GET:/api/crm/customers/{publicId}",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/customers/{publicId}",
    "tag": "Crm",
    "description": "Read customer by public id",
    "summary": "Read customer by public id",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/customers/{publicId}/history",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/customers/{publicId}/history",
    "tag": "Crm",
    "description": "List customer relationship history",
    "summary": "List customer relationship history",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/customers/{publicId}/attachments",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/customers/{publicId}/attachments",
    "tag": "Crm",
    "description": "List customer attachments",
    "summary": "List customer attachments",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:POST:/api/crm/customers/{publicId}/attachments",
    "service": "crm",
    "method": "POST",
    "path": "/api/crm/customers/{publicId}/attachments",
    "tag": "Crm",
    "description": "Create customer attachment metadata",
    "summary": "Create customer attachment metadata",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "crm:GET:/api/crm/outbox/pending",
    "service": "crm",
    "method": "GET",
    "path": "/api/crm/outbox/pending",
    "tag": "Crm",
    "description": "List pending CRM outbox events",
    "summary": "List pending CRM outbox events",
    "source": "docs/contracts/http/crm.openapi.yaml",
    "hasBody": false,
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
    "id": "documents:GET:/api/documents/attachments/{publicId}",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/attachments/{publicId}",
    "tag": "Documents",
    "description": "Read attachment metadata",
    "summary": "Read attachment metadata",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
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
    "id": "documents:GET:/api/documents/attachments/{publicId}/download",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/attachments/{publicId}/download",
    "tag": "Documents",
    "description": "Redirect to storage when access token is valid and retention allows download",
    "summary": "Redirect to storage when access token is valid and retention allows download",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:POST:/api/documents/attachments/{publicId}/archive",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/attachments/{publicId}/archive",
    "tag": "Documents",
    "description": "Archive attachment metadata under retention governance",
    "summary": "Archive attachment metadata under retention governance",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:POST:/api/documents/attachments/{publicId}/access-links",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/attachments/{publicId}/access-links",
    "tag": "Documents",
    "description": "Create signed short-lived access link",
    "summary": "Create signed short-lived access link",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:POST:/api/documents/attachments/{publicId}/access-links/revoke",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/attachments/{publicId}/access-links/revoke",
    "tag": "Documents",
    "description": "Revoke a signed access link before expiration",
    "summary": "Revoke a signed access link before expiration",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:GET:/api/documents/audit-events",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/audit-events",
    "tag": "Documents",
    "description": "List document access audit events",
    "summary": "List document access audit events",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "documents:POST:/api/documents/upload-sessions",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/upload-sessions",
    "tag": "Documents",
    "description": "Create governed upload session",
    "summary": "Create governed upload session",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "documents:GET:/api/documents/upload-sessions/{publicId}",
    "service": "documents",
    "method": "GET",
    "path": "/api/documents/upload-sessions/{publicId}",
    "tag": "Documents",
    "description": "Read upload session",
    "summary": "Read upload session",
    "source": "docs/contracts/http/documents.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "documents:POST:/api/documents/upload-sessions/{publicId}/complete",
    "service": "documents",
    "method": "POST",
    "path": "/api/documents/upload-sessions/{publicId}/complete",
    "tag": "Documents",
    "description": "Complete upload session after malware scan gate",
    "summary": "Complete upload session after malware scan gate",
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
    "id": "edge:GET:/api/edge/ops/health",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/health",
    "tag": "Edge",
    "description": "Read edge dependency health",
    "summary": "Read edge dependency health",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/tenant-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/tenant-overview",
    "tag": "Edge",
    "description": "Read tenant operational overview",
    "summary": "Read tenant operational overview",
    "source": "docs/contracts/http/edge.openapi.yaml",
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
    "id": "edge:GET:/api/edge/ops/automation-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/automation-overview",
    "tag": "Edge",
    "description": "Read automation cockpit",
    "summary": "Read automation cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/engagement-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/engagement-overview",
    "tag": "Edge",
    "description": "Read engagement cockpit",
    "summary": "Read engagement cockpit",
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
    "id": "edge:GET:/api/edge/ops/documents-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/documents-overview",
    "tag": "Edge",
    "description": "Read documents cockpit",
    "summary": "Read documents cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/collections-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/collections-overview",
    "tag": "Edge",
    "description": "Read collections cockpit",
    "summary": "Read collections cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/platform-reliability",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/platform-reliability",
    "tag": "Edge",
    "description": "Read platform reliability cockpit",
    "summary": "Read platform reliability cockpit",
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
    "id": "edge:GET:/api/edge/ops/sales-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/sales-overview",
    "tag": "Edge",
    "description": "Read sales cockpit",
    "summary": "Read sales cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/revenue-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/revenue-overview",
    "tag": "Edge",
    "description": "Read revenue cockpit",
    "summary": "Read revenue cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/finance-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/finance-overview",
    "tag": "Edge",
    "description": "Read finance cockpit",
    "summary": "Read finance cockpit",
    "source": "docs/contracts/http/edge.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "edge:GET:/api/edge/ops/rentals-overview",
    "service": "edge",
    "method": "GET",
    "path": "/api/edge/ops/rentals-overview",
    "tag": "Edge",
    "description": "Read rentals cockpit",
    "summary": "Read rentals cockpit",
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
    "id": "finance:POST:/api/finance/projections/ingest",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/projections/ingest",
    "tag": "Finance",
    "description": "Ingest pending sales outbox events into finance projections",
    "summary": "Ingest pending sales outbox events into finance projections",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/projections",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/projections",
    "tag": "Finance",
    "description": "List receivable projections",
    "summary": "List receivable projections",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/projections/summary",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/projections/summary",
    "tag": "Finance",
    "description": "Read projection summary",
    "summary": "Read projection summary",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/operations/sync",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/operations/sync",
    "tag": "Finance",
    "description": "Sync operational receivables and commissions from sales and rentals",
    "summary": "Sync operational receivables and commissions from sales and rentals",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/receivables",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/receivables",
    "tag": "Finance",
    "description": "List operational receivables",
    "summary": "List operational receivables",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/receivables/{publicId}/settlements",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/receivables/{publicId}/settlements",
    "tag": "Finance",
    "description": "Settle receivable idempotently",
    "summary": "Settle receivable idempotently",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:GET:/api/finance/commissions",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/commissions",
    "tag": "Finance",
    "description": "List commission entries",
    "summary": "List commission entries",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/commissions/summary",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/commissions/summary",
    "tag": "Finance",
    "description": "Read commission summary",
    "summary": "Read commission summary",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/commissions/{publicId}/block",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/commissions/{publicId}/block",
    "tag": "Finance",
    "description": "Block commission",
    "summary": "Block commission",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:POST:/api/finance/commissions/{publicId}/release",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/commissions/{publicId}/release",
    "tag": "Finance",
    "description": "Release commission",
    "summary": "Release commission",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:GET:/api/finance/payables",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/payables",
    "tag": "Finance",
    "description": "List payables",
    "summary": "List payables",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/payables",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/payables",
    "tag": "Finance",
    "description": "Create payable",
    "summary": "Create payable",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:PATCH:/api/finance/payables/{publicId}/status",
    "service": "finance",
    "method": "PATCH",
    "path": "/api/finance/payables/{publicId}/status",
    "tag": "Finance",
    "description": "Update payable status",
    "summary": "Update payable status",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "finance:GET:/api/finance/costs",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/costs",
    "tag": "Finance",
    "description": "List cost entries",
    "summary": "List cost entries",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/costs",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/costs",
    "tag": "Finance",
    "description": "Create cost entry",
    "summary": "Create cost entry",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/cash-accounts",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/cash-accounts",
    "tag": "Finance",
    "description": "List cash accounts",
    "summary": "List cash accounts",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/cash-accounts",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/cash-accounts",
    "tag": "Finance",
    "description": "Create cash account",
    "summary": "Create cash account",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/treasury/sync",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/treasury/sync",
    "tag": "Finance",
    "description": "Sync treasury movements from finance operations",
    "summary": "Sync treasury movements from finance operations",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/cash-movements",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/cash-movements",
    "tag": "Finance",
    "description": "List cash movements",
    "summary": "List cash movements",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/cash-movements/summary",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/cash-movements/summary",
    "tag": "Finance",
    "description": "Read cash movement summary",
    "summary": "Read cash movement summary",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/reports/treasury",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/reports/treasury",
    "tag": "Finance",
    "description": "Read treasury operational report",
    "summary": "Read treasury operational report",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/period-closures",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/period-closures",
    "tag": "Finance",
    "description": "List period closures",
    "summary": "List period closures",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "finance:POST:/api/finance/period-closures",
    "service": "finance",
    "method": "POST",
    "path": "/api/finance/period-closures",
    "tag": "Finance",
    "description": "Close current financial period with immutable snapshot",
    "summary": "Close current financial period with immutable snapshot",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "finance:GET:/api/finance/period-closures/{periodKey}",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/period-closures/{periodKey}",
    "tag": "Finance",
    "description": "Read period closure snapshot",
    "summary": "Read period closure snapshot",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "periodKey"
    ]
  },
  {
    "id": "finance:GET:/api/finance/reports/operations",
    "service": "finance",
    "method": "GET",
    "path": "/api/finance/reports/operations",
    "tag": "Finance",
    "description": "Read finance operational report",
    "summary": "Read finance operational report",
    "source": "docs/contracts/http/finance.openapi.yaml",
    "hasBody": false,
    "pathParams": []
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
    "id": "fiscal:GET:/api/fiscal/issuance-queue",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/issuance-queue",
    "tag": "Fiscal",
    "description": "List NF-e/NFS-e issuance queue items",
    "summary": "List NF-e/NFS-e issuance queue items",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/issuance-queue",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/issuance-queue",
    "tag": "Fiscal",
    "description": "Create NF-e/NFS-e issuance queue item",
    "summary": "Create NF-e/NFS-e issuance queue item",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/certificates",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/certificates",
    "tag": "Fiscal",
    "description": "List fiscal certificate vault posture records",
    "summary": "List fiscal certificate vault posture records",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/certificates",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/certificates",
    "tag": "Fiscal",
    "description": "Register fiscal certificate vault posture record",
    "summary": "Register fiscal certificate vault posture record",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/sped-exports",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/sped-exports",
    "tag": "Fiscal",
    "description": "List SPED/EFD export jobs",
    "summary": "List SPED/EFD export jobs",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/sped-exports",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/sped-exports",
    "tag": "Fiscal",
    "description": "Create SPED/EFD export job",
    "summary": "Create SPED/EFD export job",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/contingency-plans",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/contingency-plans",
    "tag": "Fiscal",
    "description": "List fiscal contingency plans",
    "summary": "List fiscal contingency plans",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/contingency-plans",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/contingency-plans",
    "tag": "Fiscal",
    "description": "Create fiscal contingency plan",
    "summary": "Create fiscal contingency plan",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/artifacts",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/artifacts",
    "tag": "Fiscal",
    "description": "GET Fiscal Artifacts",
    "summary": "GET Fiscal Artifacts",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/artifacts",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/artifacts",
    "tag": "Fiscal",
    "description": "Register XML/PDF fiscal artifact",
    "summary": "Register XML/PDF fiscal artifact",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "fiscal:GET:/api/fiscal/reconciliations",
    "service": "fiscal",
    "method": "GET",
    "path": "/api/fiscal/reconciliations",
    "tag": "Fiscal",
    "description": "GET Fiscal Reconciliations",
    "summary": "GET Fiscal Reconciliations",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "fiscal:POST:/api/fiscal/reconciliations",
    "service": "fiscal",
    "method": "POST",
    "path": "/api/fiscal/reconciliations",
    "tag": "Fiscal",
    "description": "Run fiscal finance reconciliation",
    "summary": "Run fiscal finance reconciliation",
    "source": "docs/contracts/http/fiscal.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/invites/{inviteToken}/accept",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/invites/{inviteToken}/accept",
    "tag": "Identity",
    "description": "Create or execute identity invites inviteToken accept",
    "summary": "Create or execute identity invites inviteToken accept",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "inviteToken"
    ]
  },
  {
    "id": "identity:POST:/api/identity/password-recovery",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/password-recovery",
    "tag": "Identity",
    "description": "Create or execute identity password recovery",
    "summary": "Create or execute identity password recovery",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/password-recovery/{resetToken}/complete",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/password-recovery/{resetToken}/complete",
    "tag": "Identity",
    "description": "Create or execute identity password recovery resetToken complete",
    "summary": "Create or execute identity password recovery resetToken complete",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "resetToken"
    ]
  },
  {
    "id": "identity:POST:/api/identity/sessions/login",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/sessions/login",
    "tag": "Identity",
    "description": "Create or execute identity sessions login",
    "summary": "Create or execute identity sessions login",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:POST:/api/identity/sessions/logout",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/sessions/logout",
    "tag": "Identity",
    "description": "Create or execute identity sessions logout",
    "summary": "Create or execute identity sessions logout",
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
    "description": "Create or execute identity sessions refresh",
    "summary": "Create or execute identity sessions refresh",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:GET:/api/identity/tenants",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants",
    "tag": "Identity",
    "description": "Read identity tenants",
    "summary": "Read identity tenants",
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
    "description": "Create or execute identity tenants",
    "summary": "Create or execute identity tenants",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}",
    "tag": "Identity",
    "description": "Read identity tenants slug",
    "summary": "Read identity tenants slug",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/access",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/access",
    "tag": "Identity",
    "description": "Read identity tenants slug access",
    "summary": "Read identity tenants slug access",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/companies",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/companies",
    "tag": "Identity",
    "description": "Read identity tenants slug companies",
    "summary": "Read identity tenants slug companies",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/companies",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/companies",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug companies",
    "summary": "Create or execute identity tenants slug companies",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/companies/{companyPublicId}",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/companies/{companyPublicId}",
    "tag": "Identity",
    "description": "Read identity tenants slug companies companyPublicId",
    "summary": "Read identity tenants slug companies companyPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "companyPublicId"
    ]
  },
  {
    "id": "identity:PATCH:/api/identity/tenants/{slug}/companies/{companyPublicId}",
    "service": "identity",
    "method": "PATCH",
    "path": "/api/identity/tenants/{slug}/companies/{companyPublicId}",
    "tag": "Identity",
    "description": "Update identity tenants slug companies companyPublicId",
    "summary": "Update identity tenants slug companies companyPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "companyPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/invites",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/invites",
    "tag": "Identity",
    "description": "Read identity tenants slug invites",
    "summary": "Read identity tenants slug invites",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/invites",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/invites",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug invites",
    "summary": "Create or execute identity tenants slug invites",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/invites/{invitePublicId}",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/invites/{invitePublicId}",
    "tag": "Identity",
    "description": "Read identity tenants slug invites invitePublicId",
    "summary": "Read identity tenants slug invites invitePublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "invitePublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/invites/{invitePublicId}/cancel",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/invites/{invitePublicId}/cancel",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug invites invitePublicId cancel",
    "summary": "Create or execute identity tenants slug invites invitePublicId cancel",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "invitePublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/invites/{invitePublicId}/resend",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/invites/{invitePublicId}/resend",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug invites invitePublicId resend",
    "summary": "Create or execute identity tenants slug invites invitePublicId resend",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "invitePublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/roles",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/roles",
    "tag": "Identity",
    "description": "Read identity tenants slug roles",
    "summary": "Read identity tenants slug roles",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/security/audit",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/security/audit",
    "tag": "Identity",
    "description": "Read identity tenants slug security audit",
    "summary": "Read identity tenants slug security audit",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:DELETE:/api/identity/tenants/{slug}/sessions/{sessionPublicId}",
    "service": "identity",
    "method": "DELETE",
    "path": "/api/identity/tenants/{slug}/sessions/{sessionPublicId}",
    "tag": "Identity",
    "description": "Delete identity tenants slug sessions sessionPublicId",
    "summary": "Delete identity tenants slug sessions sessionPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "sessionPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/snapshot",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/snapshot",
    "tag": "Identity",
    "description": "Read identity tenants slug snapshot",
    "summary": "Read identity tenants slug snapshot",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/teams",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/teams",
    "tag": "Identity",
    "description": "Read identity tenants slug teams",
    "summary": "Read identity tenants slug teams",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/teams",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/teams",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug teams",
    "summary": "Create or execute identity tenants slug teams",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:PATCH:/api/identity/tenants/{slug}/teams/{teamPublicId}",
    "service": "identity",
    "method": "PATCH",
    "path": "/api/identity/tenants/{slug}/teams/{teamPublicId}",
    "tag": "Identity",
    "description": "Update identity tenants slug teams teamPublicId",
    "summary": "Update identity tenants slug teams teamPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "teamPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/teams/{teamPublicId}/members",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/teams/{teamPublicId}/members",
    "tag": "Identity",
    "description": "Read identity tenants slug teams teamPublicId members",
    "summary": "Read identity tenants slug teams teamPublicId members",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "teamPublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/teams/{teamPublicId}/members",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/teams/{teamPublicId}/members",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug teams teamPublicId members",
    "summary": "Create or execute identity tenants slug teams teamPublicId members",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "teamPublicId"
    ]
  },
  {
    "id": "identity:DELETE:/api/identity/tenants/{slug}/teams/{teamPublicId}/members/{userPublicId}",
    "service": "identity",
    "method": "DELETE",
    "path": "/api/identity/tenants/{slug}/teams/{teamPublicId}/members/{userPublicId}",
    "tag": "Identity",
    "description": "Delete identity tenants slug teams teamPublicId members userPublicId",
    "summary": "Delete identity tenants slug teams teamPublicId members userPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "teamPublicId",
      "userPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/users",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/users",
    "tag": "Identity",
    "description": "Read identity tenants slug users",
    "summary": "Read identity tenants slug users",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/users",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/users",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug users",
    "summary": "Create or execute identity tenants slug users",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/users/{userPublicId}",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}",
    "tag": "Identity",
    "description": "Read identity tenants slug users userPublicId",
    "summary": "Read identity tenants slug users userPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:PATCH:/api/identity/tenants/{slug}/users/{userPublicId}",
    "service": "identity",
    "method": "PATCH",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}",
    "tag": "Identity",
    "description": "Update identity tenants slug users userPublicId",
    "summary": "Update identity tenants slug users userPublicId",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:PATCH:/api/identity/tenants/{slug}/users/{userPublicId}/access",
    "service": "identity",
    "method": "PATCH",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/access",
    "tag": "Identity",
    "description": "Update identity tenants slug users userPublicId access",
    "summary": "Update identity tenants slug users userPublicId access",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:DELETE:/api/identity/tenants/{slug}/users/{userPublicId}/mfa",
    "service": "identity",
    "method": "DELETE",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/mfa",
    "tag": "Identity",
    "description": "Delete identity tenants slug users userPublicId mfa",
    "summary": "Delete identity tenants slug users userPublicId mfa",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/users/{userPublicId}/mfa",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/mfa",
    "tag": "Identity",
    "description": "Read identity tenants slug users userPublicId mfa",
    "summary": "Read identity tenants slug users userPublicId mfa",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/users/{userPublicId}/mfa/enroll",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/mfa/enroll",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug users userPublicId mfa enroll",
    "summary": "Create or execute identity tenants slug users userPublicId mfa enroll",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/users/{userPublicId}/mfa/verify",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/mfa/verify",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug users userPublicId mfa verify",
    "summary": "Create or execute identity tenants slug users userPublicId mfa verify",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/users/{userPublicId}/roles",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/roles",
    "tag": "Identity",
    "description": "Read identity tenants slug users userPublicId roles",
    "summary": "Read identity tenants slug users userPublicId roles",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:POST:/api/identity/tenants/{slug}/users/{userPublicId}/roles",
    "service": "identity",
    "method": "POST",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/roles",
    "tag": "Identity",
    "description": "Create or execute identity tenants slug users userPublicId roles",
    "summary": "Create or execute identity tenants slug users userPublicId roles",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:DELETE:/api/identity/tenants/{slug}/users/{userPublicId}/roles/{roleCode}",
    "service": "identity",
    "method": "DELETE",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/roles/{roleCode}",
    "tag": "Identity",
    "description": "Delete identity tenants slug users userPublicId roles roleCode",
    "summary": "Delete identity tenants slug users userPublicId roles roleCode",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId",
      "roleCode"
    ]
  },
  {
    "id": "identity:DELETE:/api/identity/tenants/{slug}/users/{userPublicId}/sessions",
    "service": "identity",
    "method": "DELETE",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/sessions",
    "tag": "Identity",
    "description": "Delete identity tenants slug users userPublicId sessions",
    "summary": "Delete identity tenants slug users userPublicId sessions",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:GET:/api/identity/tenants/{slug}/users/{userPublicId}/sessions",
    "service": "identity",
    "method": "GET",
    "path": "/api/identity/tenants/{slug}/users/{userPublicId}/sessions",
    "tag": "Identity",
    "description": "Read identity tenants slug users userPublicId sessions",
    "summary": "Read identity tenants slug users userPublicId sessions",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "slug",
      "userPublicId"
    ]
  },
  {
    "id": "identity:GET:/health/details",
    "service": "identity",
    "method": "GET",
    "path": "/health/details",
    "tag": "Identity",
    "description": "Read health details",
    "summary": "Read health details",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/health/live",
    "service": "identity",
    "method": "GET",
    "path": "/health/live",
    "tag": "Identity",
    "description": "Read health live",
    "summary": "Read health live",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "identity:GET:/health/ready",
    "service": "identity",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Identity",
    "description": "Read health ready",
    "summary": "Read health ready",
    "source": "docs/contracts/http/identity.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/health/live",
    "service": "inventory",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço inventory.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/health/ready",
    "service": "inventory",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço inventory.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/health/details",
    "service": "inventory",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço inventory.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/capabilities",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/capabilities",
    "tag": "Inventory",
    "description": "Read inventory capability catalog",
    "summary": "Read inventory capability catalog",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/summary",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/summary",
    "tag": "Inventory",
    "description": "Read inventory operational summary",
    "summary": "Read inventory operational summary",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/locations",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/locations",
    "tag": "Inventory",
    "description": "List locations",
    "summary": "List locations",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:POST:/api/inventory/locations",
    "service": "inventory",
    "method": "POST",
    "path": "/api/inventory/locations",
    "tag": "Inventory",
    "description": "Create location",
    "summary": "Create location",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/locations/{publicId}",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/locations/{publicId}",
    "tag": "Inventory",
    "description": "Read one location",
    "summary": "Read one location",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:PATCH:/api/inventory/locations/{publicId}/status",
    "service": "inventory",
    "method": "PATCH",
    "path": "/api/inventory/locations/{publicId}/status",
    "tag": "Inventory",
    "description": "Transition location status",
    "summary": "Transition location status",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:GET:/api/inventory/movements",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/movements",
    "tag": "Inventory",
    "description": "List movements",
    "summary": "List movements",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:POST:/api/inventory/movements",
    "service": "inventory",
    "method": "POST",
    "path": "/api/inventory/movements",
    "tag": "Inventory",
    "description": "Create movement",
    "summary": "Create movement",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/movements/{publicId}",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/movements/{publicId}",
    "tag": "Inventory",
    "description": "Read one movement",
    "summary": "Read one movement",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:PATCH:/api/inventory/movements/{publicId}/status",
    "service": "inventory",
    "method": "PATCH",
    "path": "/api/inventory/movements/{publicId}/status",
    "tag": "Inventory",
    "description": "Transition movement status",
    "summary": "Transition movement status",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:GET:/api/inventory/reservations",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/reservations",
    "tag": "Inventory",
    "description": "List reservations",
    "summary": "List reservations",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:POST:/api/inventory/reservations",
    "service": "inventory",
    "method": "POST",
    "path": "/api/inventory/reservations",
    "tag": "Inventory",
    "description": "Create reservation",
    "summary": "Create reservation",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/reservations/{publicId}",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/reservations/{publicId}",
    "tag": "Inventory",
    "description": "Read one reservation",
    "summary": "Read one reservation",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:PATCH:/api/inventory/reservations/{publicId}/status",
    "service": "inventory",
    "method": "PATCH",
    "path": "/api/inventory/reservations/{publicId}/status",
    "tag": "Inventory",
    "description": "Transition reservation status",
    "summary": "Transition reservation status",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:GET:/api/inventory/balances",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/balances",
    "tag": "Inventory",
    "description": "GET Inventory Balances",
    "summary": "GET Inventory Balances",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/costing/summary",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/costing/summary",
    "tag": "Inventory",
    "description": "GET Costing Summary",
    "summary": "GET Costing Summary",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/cycle-counts/variances",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/cycle-counts/variances",
    "tag": "Inventory",
    "description": "GET Cycle Counts Variances",
    "summary": "GET Cycle Counts Variances",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/cycle-counts",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/cycle-counts",
    "tag": "Inventory",
    "description": "GET Inventory Cycle Counts",
    "summary": "GET Inventory Cycle Counts",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:POST:/api/inventory/cycle-counts",
    "service": "inventory",
    "method": "POST",
    "path": "/api/inventory/cycle-counts",
    "tag": "Inventory",
    "description": "Create cycle count",
    "summary": "Create cycle count",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "inventory:GET:/api/inventory/cycle-counts/{publicId}",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/cycle-counts/{publicId}",
    "tag": "Inventory",
    "description": "GET Cycle Counts",
    "summary": "GET Cycle Counts",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:PATCH:/api/inventory/cycle-counts/{publicId}/status",
    "service": "inventory",
    "method": "PATCH",
    "path": "/api/inventory/cycle-counts/{publicId}/status",
    "tag": "Inventory",
    "description": "PATCH Status",
    "summary": "PATCH Status",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "inventory:GET:/api/inventory/cost-layers",
    "service": "inventory",
    "method": "GET",
    "path": "/api/inventory/cost-layers",
    "tag": "Inventory",
    "description": "GET Inventory Cost Layers",
    "summary": "GET Inventory Cost Layers",
    "source": "docs/contracts/http/inventory.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "inventory:POST:/api/inventory/cost-layers",
    "service": "inventory",
    "method": "POST",
    "path": "/api/inventory/cost-layers",
    "tag": "Inventory",
    "description": "Create cost layer",
    "summary": "Create cost layer",
    "source": "docs/contracts/http/inventory.openapi.yaml",
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
    "id": "notification:POST:/api/notification/center/bulk",
    "service": "notification",
    "method": "POST",
    "path": "/api/notification/center/bulk",
    "tag": "Notification",
    "description": "Bulk create notification center items with partial success",
    "summary": "Bulk create notification center items with partial success",
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
    "id": "platform-control:GET:/api/platform-control/providers/activation/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/providers/activation/catalog",
    "tag": "Platform Control",
    "description": "List BYOK external provider activation catalog",
    "summary": "List BYOK external provider activation catalog",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/providers/activation/runs",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/providers/activation/runs",
    "tag": "Platform Control",
    "description": "List external provider activation runs for one tenant",
    "summary": "List external provider activation runs for one tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/providers/activation/{provider_key}/test",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/providers/activation/{provider_key}/test",
    "tag": "Platform Control",
    "description": "Execute a BYOK provider activation test or supported provider action",
    "summary": "Execute a BYOK provider activation test or supported provider action",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "provider_key"
    ]
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
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/incidents",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents",
    "tag": "Platform Control",
    "description": "List operational incidents by tenant",
    "summary": "List operational incidents by tenant",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/incidents",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents",
    "tag": "Platform Control",
    "description": "Open one operational incident",
    "summary": "Open one operational incident",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/incident-command/readiness",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/incident-command/readiness",
    "tag": "Platform Control",
    "description": "Read incident command readiness",
    "summary": "Read incident command readiness",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}",
    "tag": "Platform Control",
    "description": "Read one incident with timeline, actions and postmortem",
    "summary": "Read one incident with timeline, actions and postmortem",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/timeline",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/timeline",
    "tag": "Platform Control",
    "description": "Append one incident timeline event",
    "summary": "Append one incident timeline event",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/actions",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/actions",
    "tag": "Platform Control",
    "description": "Create one incident action item",
    "summary": "Create one incident action item",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/resolve",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/resolve",
    "tag": "Platform Control",
    "description": "Resolve one incident",
    "summary": "Resolve one incident",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/postmortem",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/postmortem",
    "tag": "Platform Control",
    "description": "Create one incident postmortem",
    "summary": "Create one incident postmortem",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/policies/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/policies/catalog",
    "tag": "Platform Control",
    "description": "List policy catalog",
    "summary": "List policy catalog",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/policies/evaluate",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/policies/evaluate",
    "tag": "Platform Control",
    "description": "Evaluate one policy decision",
    "summary": "Evaluate one policy decision",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/policies/decisions",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/policies/decisions",
    "tag": "Platform Control",
    "description": "List policy decisions",
    "summary": "List policy decisions",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/timeline",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/timeline",
    "tag": "Platform Control",
    "description": "List operational timeline events",
    "summary": "List operational timeline events",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/timeline",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/timeline",
    "tag": "Platform Control",
    "description": "Record one operational timeline event",
    "summary": "Record one operational timeline event",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/approvals",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/approvals",
    "tag": "Platform Control",
    "description": "List command approvals",
    "summary": "List command approvals",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/approvals",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/approvals",
    "tag": "Platform Control",
    "description": "Create one command approval request",
    "summary": "Create one command approval request",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/approve",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/approve",
    "tag": "Platform Control",
    "description": "Approve one command approval request",
    "summary": "Approve one command approval request",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/reject",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/reject",
    "tag": "Platform Control",
    "description": "Reject one command approval request",
    "summary": "Reject one command approval request",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/execute",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/execute",
    "tag": "Platform Control",
    "description": "Execute one approved command request",
    "summary": "Execute one approved command request",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/runbooks/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/runbooks/catalog",
    "tag": "Platform Control",
    "description": "List runbook catalog",
    "summary": "List runbook catalog",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/runbooks",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/runbooks",
    "tag": "Platform Control",
    "description": "List runbook runs",
    "summary": "List runbook runs",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/runbooks",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/runbooks",
    "tag": "Platform Control",
    "description": "Create one runbook run",
    "summary": "Create one runbook run",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/start",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/start",
    "tag": "Platform Control",
    "description": "Start one approved runbook run",
    "summary": "Start one approved runbook run",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/complete-step",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/complete-step",
    "tag": "Platform Control",
    "description": "Complete current runbook step",
    "summary": "Complete current runbook step",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/cancel",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/runbooks/{publicId}/cancel",
    "tag": "Platform Control",
    "description": "Cancel one runbook run",
    "summary": "Cancel one runbook run",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/evidence",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/evidence",
    "tag": "Platform Control",
    "description": "List audit evidence records",
    "summary": "List audit evidence records",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/evidence",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/evidence",
    "tag": "Platform Control",
    "description": "Register one audit evidence record",
    "summary": "Register one audit evidence record",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/evidence/{publicId}",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/evidence/{publicId}",
    "tag": "Platform Control",
    "description": "Read one audit evidence record",
    "summary": "Read one audit evidence record",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/autonomous-governance/readiness",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/autonomous-governance/readiness",
    "tag": "Platform Control",
    "description": "Read governance readiness",
    "summary": "Read governance readiness",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/event-mesh/catalog",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/event-mesh/catalog",
    "tag": "Platform Control",
    "description": "List enterprise event mesh streams",
    "summary": "List enterprise event mesh streams",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/event-mesh/events",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/event-mesh/events",
    "tag": "Platform Control",
    "description": "List tenant event mesh events",
    "summary": "List tenant event mesh events",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/event-mesh/events",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/event-mesh/events",
    "tag": "Platform Control",
    "description": "Publish one tenant event mesh event",
    "summary": "Publish one tenant event mesh event",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/event-mesh/dead-letters",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/event-mesh/dead-letters",
    "tag": "Platform Control",
    "description": "List tenant event mesh dead letters",
    "summary": "List tenant event mesh dead letters",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/event-mesh/dead-letters/{publicId}/replay",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/event-mesh/dead-letters/{publicId}/replay",
    "tag": "Platform Control",
    "description": "Replay one event mesh dead letter",
    "summary": "Replay one event mesh dead letter",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/event-mesh/lineage",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/event-mesh/lineage",
    "tag": "Platform Control",
    "description": "Read event mesh lineage",
    "summary": "Read event mesh lineage",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/runtime/profile",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/profile",
    "tag": "Platform Control",
    "description": "Read tenant runtime profile",
    "summary": "Read tenant runtime profile",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/runtime/profile",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/profile",
    "tag": "Platform Control",
    "description": "Update tenant runtime profile",
    "summary": "Update tenant runtime profile",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/runtime/health-score",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/health-score",
    "tag": "Platform Control",
    "description": "Read tenant runtime health score",
    "summary": "Read tenant runtime health score",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/runtime/quotas",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/quotas",
    "tag": "Platform Control",
    "description": "List tenant runtime quotas",
    "summary": "List tenant runtime quotas",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:PUT:/api/platform-control/tenants/{tenantSlug}/runtime/quotas/{metricKey}",
    "service": "platform-control",
    "method": "PUT",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/quotas/{metricKey}",
    "tag": "Platform Control",
    "description": "Upsert one tenant runtime quota",
    "summary": "Upsert one tenant runtime quota",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug",
      "metricKey"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/runtime/maintenance-windows",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/maintenance-windows",
    "tag": "Platform Control",
    "description": "List tenant maintenance windows",
    "summary": "List tenant maintenance windows",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:POST:/api/platform-control/tenants/{tenantSlug}/runtime/maintenance-windows",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/tenants/{tenantSlug}/runtime/maintenance-windows",
    "tag": "Platform Control",
    "description": "Schedule one tenant maintenance window",
    "summary": "Schedule one tenant maintenance window",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/contracts/evolution",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/contracts/evolution",
    "tag": "Platform Control",
    "description": "List contract evolution catalog",
    "summary": "List contract evolution catalog",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:POST:/api/platform-control/contracts/evolution/snapshots",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/contracts/evolution/snapshots",
    "tag": "Platform Control",
    "description": "Create one contract snapshot",
    "summary": "Create one contract snapshot",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/contracts/evolution/diffs",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/contracts/evolution/diffs",
    "tag": "Platform Control",
    "description": "List contract diffs",
    "summary": "List contract diffs",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:POST:/api/platform-control/contracts/evolution/diffs",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/contracts/evolution/diffs",
    "tag": "Platform Control",
    "description": "Create one contract diff",
    "summary": "Create one contract diff",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/contracts/evolution/breaking-changes",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/contracts/evolution/breaking-changes",
    "tag": "Platform Control",
    "description": "List contract breaking changes",
    "summary": "List contract breaking changes",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:POST:/api/platform-control/contracts/evolution/breaking-changes/{publicId}/approve",
    "service": "platform-control",
    "method": "POST",
    "path": "/api/platform-control/contracts/evolution/breaking-changes/{publicId}/approve",
    "tag": "Platform Control",
    "description": "Approve one contract breaking change",
    "summary": "Approve one contract breaking change",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "platform-control:GET:/api/platform-control/contracts/evolution/compatibility-matrix",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/contracts/evolution/compatibility-matrix",
    "tag": "Platform Control",
    "description": "Read contract compatibility matrix",
    "summary": "Read contract compatibility matrix",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "platform-control:GET:/api/platform-control/tenants/{tenantSlug}/enterprise-runtime/readiness",
    "service": "platform-control",
    "method": "GET",
    "path": "/api/platform-control/tenants/{tenantSlug}/enterprise-runtime/readiness",
    "tag": "Platform Control",
    "description": "Read enterprise runtime readiness",
    "summary": "Read enterprise runtime readiness",
    "source": "docs/contracts/http/platform-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "tenantSlug"
    ]
  },
  {
    "id": "procurement:GET:/health/live",
    "service": "procurement",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço procurement.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/health/ready",
    "service": "procurement",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço procurement.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/health/details",
    "service": "procurement",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço procurement.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/capabilities",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/capabilities",
    "tag": "Procurement",
    "description": "Read procurement capability catalog",
    "summary": "Read procurement capability catalog",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/matching/summary",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/matching/summary",
    "tag": "Procurement",
    "description": "Read procurement operational summary",
    "summary": "Read procurement operational summary",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/requisitions",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/requisitions",
    "tag": "Procurement",
    "description": "List requisitions",
    "summary": "List requisitions",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/requisitions",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/requisitions",
    "tag": "Procurement",
    "description": "Create requisition",
    "summary": "Create requisition",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/requisitions/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/requisitions/{publicId}",
    "tag": "Procurement",
    "description": "Read one requisition",
    "summary": "Read one requisition",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:PATCH:/api/procurement/requisitions/{publicId}/status",
    "service": "procurement",
    "method": "PATCH",
    "path": "/api/procurement/requisitions/{publicId}/status",
    "tag": "Procurement",
    "description": "Transition requisition status",
    "summary": "Transition requisition status",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:GET:/api/procurement/quotations",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/quotations",
    "tag": "Procurement",
    "description": "List quotations",
    "summary": "List quotations",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/quotations",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/quotations",
    "tag": "Procurement",
    "description": "Create quotation",
    "summary": "Create quotation",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/quotations/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/quotations/{publicId}",
    "tag": "Procurement",
    "description": "Read one quotation",
    "summary": "Read one quotation",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:PATCH:/api/procurement/quotations/{publicId}/status",
    "service": "procurement",
    "method": "PATCH",
    "path": "/api/procurement/quotations/{publicId}/status",
    "tag": "Procurement",
    "description": "Transition quotation status",
    "summary": "Transition quotation status",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:GET:/api/procurement/purchase-orders",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/purchase-orders",
    "tag": "Procurement",
    "description": "List purchase-orders",
    "summary": "List purchase-orders",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/purchase-orders",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/purchase-orders",
    "tag": "Procurement",
    "description": "Create purchase order",
    "summary": "Create purchase order",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/purchase-orders/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/purchase-orders/{publicId}",
    "tag": "Procurement",
    "description": "Read one purchase order",
    "summary": "Read one purchase order",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:PATCH:/api/procurement/purchase-orders/{publicId}/status",
    "service": "procurement",
    "method": "PATCH",
    "path": "/api/procurement/purchase-orders/{publicId}/status",
    "tag": "Procurement",
    "description": "Transition purchase order status",
    "summary": "Transition purchase order status",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:GET:/api/procurement/receipts",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/receipts",
    "tag": "Procurement",
    "description": "List receipts",
    "summary": "List receipts",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/receipts",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/receipts",
    "tag": "Procurement",
    "description": "Create receipt",
    "summary": "Create receipt",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/receipts/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/receipts/{publicId}",
    "tag": "Procurement",
    "description": "Read one receipt",
    "summary": "Read one receipt",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:PATCH:/api/procurement/receipts/{publicId}/status",
    "service": "procurement",
    "method": "PATCH",
    "path": "/api/procurement/receipts/{publicId}/status",
    "tag": "Procurement",
    "description": "Transition receipt status",
    "summary": "Transition receipt status",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:POST:/api/procurement/approvals/apply",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/approvals/apply",
    "tag": "Procurement",
    "description": "POST Approvals Apply",
    "summary": "POST Approvals Apply",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/matching/three-way",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/matching/three-way",
    "tag": "Procurement",
    "description": "POST Matching Three Way",
    "summary": "POST Matching Three Way",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/approvals",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/approvals",
    "tag": "Procurement",
    "description": "GET Procurement Approvals",
    "summary": "GET Procurement Approvals",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:POST:/api/procurement/approvals",
    "service": "procurement",
    "method": "POST",
    "path": "/api/procurement/approvals",
    "tag": "Procurement",
    "description": "Create approval",
    "summary": "Create approval",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/approvals/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/approvals/{publicId}",
    "tag": "Procurement",
    "description": "GET Approvals",
    "summary": "GET Approvals",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "procurement:GET:/api/procurement/three-way-matches",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/three-way-matches",
    "tag": "Procurement",
    "description": "GET Procurement Three Way Matches",
    "summary": "GET Procurement Three Way Matches",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "procurement:GET:/api/procurement/three-way-matches/{publicId}",
    "service": "procurement",
    "method": "GET",
    "path": "/api/procurement/three-way-matches/{publicId}",
    "tag": "Procurement",
    "description": "GET Three Way Matches",
    "summary": "GET Three Way Matches",
    "source": "docs/contracts/http/procurement.openapi.yaml",
    "hasBody": false,
    "pathParams": [
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
    "id": "rentals:GET:/api/rentals/contracts/summary",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/summary",
    "tag": "Rentals",
    "description": "Read rental portfolio summary",
    "summary": "Read rental portfolio summary",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "rentals:GET:/api/rentals/contracts/{publicId}",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/{publicId}",
    "tag": "Rentals",
    "description": "Read one rental contract",
    "summary": "Read one rental contract",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
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
    "id": "rentals:GET:/api/rentals/contracts/{publicId}/history",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/{publicId}/history",
    "tag": "Rentals",
    "description": "List contract lifecycle history",
    "summary": "List contract lifecycle history",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:GET:/api/rentals/contracts/{publicId}/adjustments",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/{publicId}/adjustments",
    "tag": "Rentals",
    "description": "List contract adjustments",
    "summary": "List contract adjustments",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:POST:/api/rentals/contracts/{publicId}/adjustments",
    "service": "rentals",
    "method": "POST",
    "path": "/api/rentals/contracts/{publicId}/adjustments",
    "tag": "Rentals",
    "description": "Create contract adjustment",
    "summary": "Create contract adjustment",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:POST:/api/rentals/contracts/{publicId}/terminate",
    "service": "rentals",
    "method": "POST",
    "path": "/api/rentals/contracts/{publicId}/terminate",
    "tag": "Rentals",
    "description": "Terminate rental contract",
    "summary": "Terminate rental contract",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:GET:/api/rentals/contracts/{publicId}/attachments",
    "service": "rentals",
    "method": "GET",
    "path": "/api/rentals/contracts/{publicId}/attachments",
    "tag": "Rentals",
    "description": "List contract attachments",
    "summary": "List contract attachments",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "rentals:POST:/api/rentals/contracts/{publicId}/attachments",
    "service": "rentals",
    "method": "POST",
    "path": "/api/rentals/contracts/{publicId}/attachments",
    "tag": "Rentals",
    "description": "Link attachment to contract",
    "summary": "Link attachment to contract",
    "source": "docs/contracts/http/rentals.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
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
    "id": "sales:GET:/api/sales/opportunities/summary",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/opportunities/summary",
    "tag": "Sales",
    "description": "Read opportunity summary",
    "summary": "Read opportunity summary",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/opportunities/{publicId}",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/opportunities/{publicId}",
    "tag": "Sales",
    "description": "Read opportunity by public id",
    "summary": "Read opportunity by public id",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:PATCH:/api/sales/opportunities/{publicId}",
    "service": "sales",
    "method": "PATCH",
    "path": "/api/sales/opportunities/{publicId}",
    "tag": "Sales",
    "description": "Update opportunity profile",
    "summary": "Update opportunity profile",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/opportunities/{publicId}/history",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/opportunities/{publicId}/history",
    "tag": "Sales",
    "description": "List opportunity history",
    "summary": "List opportunity history",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:PATCH:/api/sales/opportunities/{publicId}/stage",
    "service": "sales",
    "method": "PATCH",
    "path": "/api/sales/opportunities/{publicId}/stage",
    "tag": "Sales",
    "description": "Update opportunity stage",
    "summary": "Update opportunity stage",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/opportunities/{publicId}/proposals",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/opportunities/{publicId}/proposals",
    "tag": "Sales",
    "description": "List proposals by opportunity",
    "summary": "List proposals by opportunity",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/opportunities/{publicId}/proposals",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/opportunities/{publicId}/proposals",
    "tag": "Sales",
    "description": "Create proposal for opportunity",
    "summary": "Create proposal for opportunity",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/proposals/{publicId}",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/proposals/{publicId}",
    "tag": "Sales",
    "description": "Read proposal by public id",
    "summary": "Read proposal by public id",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/proposals/{publicId}/history",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/proposals/{publicId}/history",
    "tag": "Sales",
    "description": "List proposal history",
    "summary": "List proposal history",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:PATCH:/api/sales/proposals/{publicId}/status",
    "service": "sales",
    "method": "PATCH",
    "path": "/api/sales/proposals/{publicId}/status",
    "tag": "Sales",
    "description": "Update proposal status",
    "summary": "Update proposal status",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/proposals/{publicId}/convert",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/proposals/{publicId}/convert",
    "tag": "Sales",
    "description": "Convert proposal to sale",
    "summary": "Convert proposal to sale",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
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
    "id": "sales:GET:/api/sales/sales/summary",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/summary",
    "tag": "Sales",
    "description": "Read sales summary",
    "summary": "Read sales summary",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}",
    "tag": "Sales",
    "description": "Read sale by public id",
    "summary": "Read sale by public id",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}/history",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}/history",
    "tag": "Sales",
    "description": "List sale history",
    "summary": "List sale history",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:PATCH:/api/sales/sales/{publicId}/status",
    "service": "sales",
    "method": "PATCH",
    "path": "/api/sales/sales/{publicId}/status",
    "tag": "Sales",
    "description": "Update sale status",
    "summary": "Update sale status",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/cancel",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/cancel",
    "tag": "Sales",
    "description": "Cancel sale",
    "summary": "Cancel sale",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}/installments",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}/installments",
    "tag": "Sales",
    "description": "List sale installments",
    "summary": "List sale installments",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/installments",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/installments",
    "tag": "Sales",
    "description": "Create installment schedule",
    "summary": "Create installment schedule",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}/commissions",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}/commissions",
    "tag": "Sales",
    "description": "List sale commissions",
    "summary": "List sale commissions",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/commissions",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/commissions",
    "tag": "Sales",
    "description": "Create sale commission",
    "summary": "Create sale commission",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/commissions/{commissionPublicId}/block",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/commissions/{commissionPublicId}/block",
    "tag": "Sales",
    "description": "Block commission",
    "summary": "Block commission",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId",
      "commissionPublicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/commissions/{commissionPublicId}/release",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/commissions/{commissionPublicId}/release",
    "tag": "Sales",
    "description": "Release commission",
    "summary": "Release commission",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId",
      "commissionPublicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}/pending-items",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}/pending-items",
    "tag": "Sales",
    "description": "List sale pending items",
    "summary": "List sale pending items",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/pending-items",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/pending-items",
    "tag": "Sales",
    "description": "Create sale pending item",
    "summary": "Create sale pending item",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/pending-items/{pendingItemPublicId}/resolve",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/pending-items/{pendingItemPublicId}/resolve",
    "tag": "Sales",
    "description": "Resolve pending item",
    "summary": "Resolve pending item",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId",
      "pendingItemPublicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/sales/{publicId}/renegotiations",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/sales/{publicId}/renegotiations",
    "tag": "Sales",
    "description": "List sale renegotiations",
    "summary": "List sale renegotiations",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/renegotiations",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/renegotiations",
    "tag": "Sales",
    "description": "Apply sale renegotiation",
    "summary": "Apply sale renegotiation",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:POST:/api/sales/sales/{publicId}/invoice",
    "service": "sales",
    "method": "POST",
    "path": "/api/sales/sales/{publicId}/invoice",
    "tag": "Sales",
    "description": "Create invoice for sale",
    "summary": "Create invoice for sale",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
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
    "id": "sales:GET:/api/sales/invoices/summary",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/invoices/summary",
    "tag": "Sales",
    "description": "Read invoice summary",
    "summary": "Read invoice summary",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "sales:GET:/api/sales/invoices/{publicId}",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/invoices/{publicId}",
    "tag": "Sales",
    "description": "Read invoice by public id",
    "summary": "Read invoice by public id",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/invoices/{publicId}/history",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/invoices/{publicId}/history",
    "tag": "Sales",
    "description": "List invoice history",
    "summary": "List invoice history",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:PATCH:/api/sales/invoices/{publicId}/status",
    "service": "sales",
    "method": "PATCH",
    "path": "/api/sales/invoices/{publicId}/status",
    "tag": "Sales",
    "description": "Update invoice status",
    "summary": "Update invoice status",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "sales:GET:/api/sales/outbox/pending",
    "service": "sales",
    "method": "GET",
    "path": "/api/sales/outbox/pending",
    "tag": "Sales",
    "description": "List pending sales outbox events",
    "summary": "List pending sales outbox events",
    "source": "docs/contracts/http/sales.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/health/live",
    "service": "search",
    "method": "GET",
    "path": "/health/live",
    "tag": "Health",
    "description": "Health live do serviço search.",
    "summary": "Health live",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/health/ready",
    "service": "search",
    "method": "GET",
    "path": "/health/ready",
    "tag": "Health",
    "description": "Health ready do serviço search.",
    "summary": "Health ready",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/health/details",
    "service": "search",
    "method": "GET",
    "path": "/health/details",
    "tag": "Health",
    "description": "Health details do serviço search.",
    "summary": "Health details",
    "source": "runtime",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/api/search/capabilities",
    "service": "search",
    "method": "GET",
    "path": "/api/search/capabilities",
    "tag": "Search",
    "description": "Read searchable domains and e-discovery capabilities",
    "summary": "Read searchable domains and e-discovery capabilities",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/api/search/query",
    "service": "search",
    "method": "GET",
    "path": "/api/search/query",
    "tag": "Search",
    "description": "Query operational records with tenant filtering and redaction",
    "summary": "Query operational records with tenant filtering and redaction",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:GET:/api/search/facets",
    "service": "search",
    "method": "GET",
    "path": "/api/search/facets",
    "tag": "Search",
    "description": "Read search facets by entity type and source service",
    "summary": "Read search facets by entity type and source service",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:POST:/api/search/saved-queries",
    "service": "search",
    "method": "POST",
    "path": "/api/search/saved-queries",
    "tag": "Search",
    "description": "Create a governed saved query",
    "summary": "Create a governed saved query",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "search:GET:/api/search/audit-events",
    "service": "search",
    "method": "GET",
    "path": "/api/search/audit-events",
    "tag": "Search",
    "description": "List query audit events",
    "summary": "List query audit events",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "search:POST:/api/search/discovery-cases",
    "service": "search",
    "method": "POST",
    "path": "/api/search/discovery-cases",
    "tag": "Search",
    "description": "Create an e-discovery case",
    "summary": "Create an e-discovery case",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "search:POST:/api/search/discovery-cases/{publicId}/items",
    "service": "search",
    "method": "POST",
    "path": "/api/search/discovery-cases/{publicId}/items",
    "tag": "Search",
    "description": "Add a record to an e-discovery case",
    "summary": "Add a record to an e-discovery case",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "search:POST:/api/search/legal-holds",
    "service": "search",
    "method": "POST",
    "path": "/api/search/legal-holds",
    "tag": "Search",
    "description": "Create a legal hold over an entity selector",
    "summary": "Create a legal hold over an entity selector",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "search:POST:/api/search/exports",
    "service": "search",
    "method": "POST",
    "path": "/api/search/exports",
    "tag": "Search",
    "description": "Create a controlled export request",
    "summary": "Create a controlled export request",
    "source": "docs/contracts/http/search.openapi.yaml",
    "hasBody": true,
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
    "id": "simulation:GET:/api/simulation/scenarios/catalog",
    "service": "simulation",
    "method": "GET",
    "path": "/api/simulation/scenarios/catalog",
    "tag": "Simulation",
    "description": "List available scenario templates",
    "summary": "List available scenario templates",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:POST:/api/simulation/scenarios/operational-load",
    "service": "simulation",
    "method": "POST",
    "path": "/api/simulation/scenarios/operational-load",
    "tag": "Simulation",
    "description": "Create operational load scenario",
    "summary": "Create operational load scenario",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/api/simulation/scenarios/runs",
    "service": "simulation",
    "method": "GET",
    "path": "/api/simulation/scenarios/runs",
    "tag": "Simulation",
    "description": "List scenario runs",
    "summary": "List scenario runs",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "simulation:GET:/api/simulation/scenarios/runs/{publicId}",
    "service": "simulation",
    "method": "GET",
    "path": "/api/simulation/scenarios/runs/{publicId}",
    "tag": "Simulation",
    "description": "Read one scenario run",
    "summary": "Read one scenario run",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
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
    "id": "simulation:GET:/api/simulation/benchmarks/runs",
    "service": "simulation",
    "method": "GET",
    "path": "/api/simulation/benchmarks/runs",
    "tag": "Simulation",
    "description": "List load benchmark runs",
    "summary": "List load benchmark runs",
    "source": "docs/contracts/http/simulation.openapi.yaml",
    "hasBody": false,
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
    "id": "supplier:GET:/api/supplier/suppliers/export",
    "service": "supplier",
    "method": "GET",
    "path": "/api/supplier/suppliers/export",
    "tag": "Supplier",
    "description": "Export suppliers by tenant and status",
    "summary": "Export suppliers by tenant and status",
    "source": "docs/contracts/http/supplier.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "supplier:POST:/api/supplier/suppliers/bulk",
    "service": "supplier",
    "method": "POST",
    "path": "/api/supplier/suppliers/bulk",
    "tag": "Supplier",
    "description": "Bulk create suppliers with partial success",
    "summary": "Bulk create suppliers with partial success",
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
    "id": "support:GET:/api/support/cases/export",
    "service": "support",
    "method": "GET",
    "path": "/api/support/cases/export",
    "tag": "Support",
    "description": "Export support cases with filters",
    "summary": "Export support cases with filters",
    "source": "docs/contracts/http/support.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "support:POST:/api/support/cases/bulk",
    "service": "support",
    "method": "POST",
    "path": "/api/support/cases/bulk",
    "tag": "Support",
    "description": "Bulk create support cases with partial success",
    "summary": "Bulk create support cases with partial success",
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
    "id": "webhook-hub:GET:/api/webhook-hub/events/dead-letter",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/events/dead-letter",
    "tag": "Webhook Hub",
    "description": "List dead-lettered inbound webhook events",
    "summary": "List dead-lettered inbound webhook events",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/events/{publicId}",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/events/{publicId}",
    "tag": "Webhook Hub",
    "description": "Read one inbound webhook event",
    "summary": "Read one inbound webhook event",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:GET:/api/webhook-hub/events/{publicId}/transitions",
    "service": "webhook-hub",
    "method": "GET",
    "path": "/api/webhook-hub/events/{publicId}/transitions",
    "tag": "Webhook Hub",
    "description": "List state transitions for one event",
    "summary": "List state transitions for one event",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/validate",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/validate",
    "tag": "Webhook Hub",
    "description": "Validate inbound webhook event",
    "summary": "Validate inbound webhook event",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/queue",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/queue",
    "tag": "Webhook Hub",
    "description": "Queue inbound webhook event for processing",
    "summary": "Queue inbound webhook event for processing",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/process",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/process",
    "tag": "Webhook Hub",
    "description": "Mark inbound webhook event as processing",
    "summary": "Mark inbound webhook event as processing",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/forward",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/forward",
    "tag": "Webhook Hub",
    "description": "Mark inbound webhook event as forwarded",
    "summary": "Mark inbound webhook event as forwarded",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/fail",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/fail",
    "tag": "Webhook Hub",
    "description": "Record inbound webhook processing failure",
    "summary": "Record inbound webhook processing failure",
    "source": "docs/contracts/http/webhook-hub.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
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
    "id": "webhook-hub:POST:/api/webhook-hub/events/{publicId}/reject",
    "service": "webhook-hub",
    "method": "POST",
    "path": "/api/webhook-hub/events/{publicId}/reject",
    "tag": "Webhook Hub",
    "description": "Reject inbound webhook event",
    "summary": "Reject inbound webhook event",
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
    "id": "workflow-control:GET:/api/workflow-control/catalog/triggers",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/catalog/triggers",
    "tag": "Workflow Control",
    "description": "List workflow trigger catalog",
    "summary": "List workflow trigger catalog",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/catalog/actions",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/catalog/actions",
    "tag": "Workflow Control",
    "description": "List workflow action catalog",
    "summary": "List workflow action catalog",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/editor",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/editor",
    "tag": "Workflow Control",
    "description": "Read editor capabilities and catalogs",
    "summary": "Read editor capabilities and catalogs",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
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
    "id": "workflow-control:GET:/api/workflow-control/definitions/{key}/versions",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions/{key}/versions",
    "tag": "Workflow Control",
    "description": "List workflow definition versions",
    "summary": "List workflow definition versions",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/definitions/{key}/versions",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/definitions/{key}/versions",
    "tag": "Workflow Control",
    "description": "Publish workflow definition version",
    "summary": "Publish workflow definition version",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/definitions/{key}/versions/summary",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions/{key}/versions/summary",
    "tag": "Workflow Control",
    "description": "Read workflow definition version summary",
    "summary": "Read workflow definition version summary",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/definitions/{key}/versions/current",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions/{key}/versions/current",
    "tag": "Workflow Control",
    "description": "Read current workflow definition version",
    "summary": "Read current workflow definition version",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "key"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/definitions/{key}/versions/{versionNumber}",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/definitions/{key}/versions/{versionNumber}",
    "tag": "Workflow Control",
    "description": "Read one workflow definition version",
    "summary": "Read one workflow definition version",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "key",
      "versionNumber"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/definitions/{key}/versions/{versionNumber}/restore",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/definitions/{key}/versions/{versionNumber}/restore",
    "tag": "Workflow Control",
    "description": "Restore workflow definition version",
    "summary": "Restore workflow definition version",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "key",
      "versionNumber"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/runs",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/runs",
    "tag": "Workflow Control",
    "description": "List workflow runs",
    "summary": "List workflow runs",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs",
    "tag": "Workflow Control",
    "description": "Create workflow run",
    "summary": "Create workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/runs/summary",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/runs/summary",
    "tag": "Workflow Control",
    "description": "Read workflow run summary",
    "summary": "Read workflow run summary",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/runs/{publicId}",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/runs/{publicId}",
    "tag": "Workflow Control",
    "description": "Read one workflow run",
    "summary": "Read one workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs/{publicId}/start",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs/{publicId}/start",
    "tag": "Workflow Control",
    "description": "Start workflow run",
    "summary": "Start workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs/{publicId}/complete",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs/{publicId}/complete",
    "tag": "Workflow Control",
    "description": "Complete workflow run",
    "summary": "Complete workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs/{publicId}/fail",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs/{publicId}/fail",
    "tag": "Workflow Control",
    "description": "Fail workflow run",
    "summary": "Fail workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs/{publicId}/cancel",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs/{publicId}/cancel",
    "tag": "Workflow Control",
    "description": "Cancel workflow run",
    "summary": "Cancel workflow run",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/runs/{publicId}/events",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/runs/{publicId}/events",
    "tag": "Workflow Control",
    "description": "List workflow run events",
    "summary": "List workflow run events",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:POST:/api/workflow-control/runs/{publicId}/events",
    "service": "workflow-control",
    "method": "POST",
    "path": "/api/workflow-control/runs/{publicId}/events",
    "tag": "Workflow Control",
    "description": "Append workflow run note",
    "summary": "Append workflow run note",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-control:GET:/api/workflow-control/runs/{publicId}/events/summary",
    "service": "workflow-control",
    "method": "GET",
    "path": "/api/workflow-control/runs/{publicId}/events/summary",
    "tag": "Workflow Control",
    "description": "Read workflow run event summary",
    "summary": "Read workflow run event summary",
    "source": "docs/contracts/http/workflow-control.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
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
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/summary",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/summary",
    "tag": "Workflow Runtime",
    "description": "Read execution summary",
    "summary": "Read execution summary",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": []
  },
  {
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/summary/by-workflow",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/summary/by-workflow",
    "tag": "Workflow Runtime",
    "description": "Read execution summary by workflow",
    "summary": "Read execution summary by workflow",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
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
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/{publicId}/transitions",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/{publicId}/transitions",
    "tag": "Workflow Runtime",
    "description": "List execution transitions",
    "summary": "List execution transitions",
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
    "id": "workflow-runtime:GET:/api/workflow-runtime/executions/{publicId}/timeline",
    "service": "workflow-runtime",
    "method": "GET",
    "path": "/api/workflow-runtime/executions/{publicId}/timeline",
    "tag": "Workflow Runtime",
    "description": "List execution timeline",
    "summary": "List execution timeline",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": false,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/start",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/start",
    "tag": "Workflow Runtime",
    "description": "Start workflow execution",
    "summary": "Start workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/complete",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/complete",
    "tag": "Workflow Runtime",
    "description": "Complete workflow execution",
    "summary": "Complete workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/fail",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/fail",
    "tag": "Workflow Runtime",
    "description": "Fail workflow execution",
    "summary": "Fail workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/cancel",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/cancel",
    "tag": "Workflow Runtime",
    "description": "Cancel workflow execution",
    "summary": "Cancel workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
    "pathParams": [
      "publicId"
    ]
  },
  {
    "id": "workflow-runtime:POST:/api/workflow-runtime/executions/{publicId}/retry",
    "service": "workflow-runtime",
    "method": "POST",
    "path": "/api/workflow-runtime/executions/{publicId}/retry",
    "tag": "Workflow Runtime",
    "description": "Retry workflow execution",
    "summary": "Retry workflow execution",
    "source": "docs/contracts/http/workflow-runtime.openapi.yaml",
    "hasBody": true,
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
    "fileName": "crm.outbox-event.schema.json",
    "name": "crm.outbox-event",
    "source": "docs/contracts/events/crm.outbox-event.schema.json"
  },
  {
    "fileName": "documents.extraction-completed.schema.json",
    "name": "documents.extraction-completed",
    "source": "docs/contracts/events/documents.extraction-completed.schema.json"
  },
  {
    "fileName": "documents.extraction-requested.schema.json",
    "name": "documents.extraction-requested",
    "source": "docs/contracts/events/documents.extraction-requested.schema.json"
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
    "fileName": "event-envelope.schema.json",
    "name": "event-envelope",
    "source": "docs/contracts/events/event-envelope.schema.json"
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
    "fileName": "sales.outbox-event.schema.json",
    "name": "sales.outbox-event",
    "source": "docs/contracts/events/sales.outbox-event.schema.json"
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
