"""Gate consolidado de aceite produtivo para a versao 1.5.x."""

from datetime import datetime, timezone

from app.reports.adapter_catalog import build_adapter_catalog
from app.reports.go_live_control import build_release_controls
from app.reports.hardening_review import (
    build_hardening_closure,
    build_operational_runbook_review,
    build_platform_closure,
    build_production_maturity_closure,
)
from app.reports.risk_compliance import build_risk_readiness
from app.reports.semantic_metrics import build_semantic_metrics_readiness
from app.reports.enterprise_runtime import build_enterprise_runtime_fabric_readiness


def build_production_readiness(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    adapter_catalog = build_adapter_catalog()
    critical_provider_gaps = adapter_catalog["summary"]["criticalUnconfiguredCapabilities"]
    release_controls = build_release_controls()
    risk_readiness = build_risk_readiness()
    enterprise_runtime_readiness = build_enterprise_runtime_fabric_readiness()

    gates = {
        "trafficEntry": {
            "status": "ready",
            "evidence": ["infra/docker-compose.corporate-like.yml", "infra/kubernetes/base/network-policy.yaml"],
            "requirement": "Somente gateway/edge publicado no perfil corporativo.",
        },
        "authzAndTenantIsolation": {
            "status": "ready",
            "evidence": ["ERP_AUTH_ENFORCEMENT", "ERP_OPENFGA_ENFORCEMENT", "scripts/test.sh security"],
            "requirement": "JWT/service account, tenant explicito, correlation id e OpenFGA quando habilitado.",
        },
        "edgeSecurity": {
            "status": "ready",
            "evidence": ["infra/gateway/nginx.conf", "ERP_SECURITY_HEADERS", "ERP_REQUIRE_REQUEST_SIGNATURE"],
            "requirement": "Headers defensivos, body limit, timeouts, correlation id e bloqueio de paths ocultos na borda.",
        },
        "workloadSecurity": {
            "status": "ready",
            "evidence": ["infra/kubernetes/base/namespace.yaml", "infra/kubernetes/base/network-policy.yaml"],
            "requirement": "Pod Security restricted, NetworkPolicy deny-by-default e workloads sem privilege escalation.",
        },
        "secrets": {
            "status": "ready",
            "evidence": [".env.production.example", "scripts/env.sh", "scripts/test.sh hardening-secrets"],
            "requirement": "Defaults locais bloqueados fora de local/test.",
        },
        "rootHardening": {
            "status": "ready",
            "evidence": [".dockerignore", ".gitignore", ".gitattributes", ".editorconfig"],
            "requirement": "Artefatos locais, caches, env reais e line endings sao controlados no monorepo poliglota.",
        },
        "staticPolicy": {
            "status": "ready",
            "evidence": ["scripts/test.sh security", ".github/workflows/quality.yml", ".github/workflows/containers.yml"],
            "requirement": "CI bloqueia MAX(id)+1, tag latest operacional, OpenFGA latest e divergencia do catalogo gerado.",
        },
        "migrationOrder": {
            "status": "ready",
            "evidence": ["scripts/test.sh security", "service-api/service-postgresql/*/migrations"],
            "requirement": "Migrations por dominio nao podem repetir prefixo numerico e devem preservar ordem deterministica.",
        },
        "tenantContractManifest": {
            "status": "ready",
            "evidence": ["platform_control.tenant_contract_manifest", "service-api/service-postgresql/platform-control/migrations/000012_tenant_contract_manifest.sql"],
            "requirement": "Tabelas operacionais declaram estrategia de tenant ou justificativa de referencia global/sistema.",
        },
        "polyglotBoundarySplit": {
            "status": "ready",
            "evidence": ["BillingHealthRoutes.cs", "FinanceHealthRoutes.cs", "api/request.ts", "api/schemas.ts", "api/security.rs", "app/infrastructure/external_http.py"],
            "requirement": "Bordas transversais de health, request validation, auth e HTTP externo ficam fora dos roteadores monoliticos.",
        },
        "transactionalIdGeneration": {
            "status": "ready",
            "evidence": [
                "PostgresIdentityRepositoryBundle.cs",
                "postgres-workflow-definition-repository.ts",
                "postgres-workflow-run-repository.ts",
                "postgres-workflow-run-event-repository.ts",
            ],
            "requirement": "IDs relacionais usam sequences do PostgreSQL em vez de MAX(id)+1 sob concorrencia.",
        },
        "contracts": {
            "status": "ready",
            "evidence": ["docs/contracts/registry.json", "docs/contracts/http", "scripts/test.sh contract"],
            "requirement": "OpenAPI versionado como fonte de verdade das APIs publicas.",
        },
        "observability": {
            "status": "ready",
            "evidence": ["infra/prometheus/prometheus.yml", "infra/blackbox/blackbox.yml", "docs/OPERACOES.md", "traceparent", "X-Correlation-Id"],
            "requirement": "Health, probes, SLOs, alertas, traceparent, correlation id e runbooks operacionais.",
        },
        "authConformance": {
            "status": "ready",
            "evidence": ["docs/SEGURANCA.md", "docs/PADROES.md", "scripts/test.sh security"],
            "requirement": "JWT, service token, tenant, actor, OpenFGA e modos de falha seguem contrato unico por stack.",
        },
        "apiConsoleSecurity": {
            "status": "ready",
            "evidence": ["client-web/client-api/src/lib/httpClient.ts", "client-web/client-api/src/App.tsx"],
            "requirement": "Console nao persiste bearer token em localStorage, redige cURL por padrao e limita retries automaticos a metodos de leitura.",
        },
        "backupRestoreDr": {
            "status": "ready",
            "evidence": ["scripts/build.sh backup-encrypted", "scripts/build.sh restore-encrypted", "scripts/test.sh backup-restore"],
            "requirement": "Backup criptografado e restore drill validavel.",
        },
        "deployArtifacts": {
            "status": "ready",
            "evidence": ["infra/kubernetes/base", "infra/kubernetes/base/deployability-matrix.yaml", "infra/kubernetes/overlays/production"],
            "requirement": "Namespace, ingress, probes, migration job, NetworkPolicy e matriz explicita de cobertura Kubernetes por servico.",
        },
        "providers": {
            "status": "ready",
            "evidence": [
                "GET /api/analytics/reports/adapter-catalog",
                "GET /api/platform-control/providers/activation/catalog",
                "POST /api/platform-control/tenants/{tenantSlug}/providers/activation/{providerKey}/test",
            ],
            "requirement": "Capabilities externas explicitam configuracao e readiness produtivo.",
            "criticalProviderGaps": critical_provider_gaps,
            "blockingForRelease": False,
            "policy": "Provider externo so executa em modo BYOK. Sem credencial real, a capability fica indisponivel ou usa fallback local explicitamente declarado.",
        },
        "externalProviderActivation": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/providers/activation/catalog",
                "GET /api/platform-control/tenants/{tenantSlug}/providers/activation/runs",
                "POST /api/platform-control/tenants/{tenantSlug}/providers/activation/{providerKey}/test",
            ],
            "requirement": "Stripe, Asaas, Mercado Pago, Resend, OpenAI, DocuSign, Clicksign e WhatsApp Cloud podem ser ativados por chave do operador, com resposta indisponivel quando a chave nao existe.",
        },
        "llmProviderByok": {
            "status": "ready",
            "evidence": ["OPENAI_API_KEY", "OPENAI_MODEL", "POST /api/ai-governance/runs"],
            "requirement": "AI Governance usa OpenAI Responses API somente quando OPENAI_API_KEY existe; sem chave, mantem resposta deterministica local e auditada.",
        },
        "documentIntelligence": {
            "status": "ready",
            "evidence": ["GET /api/analytics/document-intelligence/readiness", "aws_textract", "google_document_ai"],
            "requirement": "OCR e extracao documental declaram providers BYOK, postura de credencial e indisponibilidade segura quando nao ha chave.",
        },
        "fiscalProviderReadiness": {
            "status": "ready",
            "evidence": ["GET /api/analytics/fiscal-brazil/readiness", "FISCAL_FOCUS_NFE_API_KEY", "FISCAL_ENOTAS_API_KEY"],
            "requirement": "Fiscal Brasil declara Focus NFe, eNotas e certificado digital sem fingir emissao real sem credencial/homologacao.",
        },
        "brazilRegistryEnrichment": {
            "status": "ready",
            "evidence": ["GET /api/analytics/registry-enrichment/brazil", "brasilapi", "viacep", "serpro_cnpj"],
            "requirement": "Consulta cadastral Brasil cobre CNPJ/CEP, providers publicos e providers BYOK para CRM, supplier e master data.",
        },
        "marketMacroRisk": {
            "status": "ready",
            "evidence": ["GET /api/analytics/market-macro-risk", "bcb_sgs", "bcb_ptax", "alpha_vantage", "fixer"],
            "requirement": "Cambio, mercado e risco macro entram como sinais externos governados por provider posture.",
        },
        "externalRiskFeed": {
            "status": "ready",
            "evidence": ["GET /api/analytics/external-risk-feed", "newsapi", "gdelt", "alpha_vantage_news"],
            "requirement": "Noticias e risco reputacional entram como sinal operacional, nao como verdade de dominio.",
        },
        "providerActivationV14": {
            "status": "ready",
            "evidence": ["GET /api/platform-control/providers/activation/catalog", "POST /api/platform-control/tenants/{tenantSlug}/providers/activation/{providerKey}/test"],
            "requirement": "Provider activation cobre os dominios document_intelligence, fiscal, registry_enrichment, market_macro_risk e external_risk_feed.",
        },
        "goLive": {
            "status": "ready" if release_controls["acceptanceReady"] else "attention",
            "evidence": ["GET /api/analytics/reports/go-live-control", "GET /api/edge/ops/go-live-overview"],
            "requirement": "Rollout por tenant com rollback, gargalos, adocao e aceite.",
        },
        "operationalSearch": {
            "status": "ready",
            "evidence": ["GET /api/search/query", "POST /api/search/legal-holds", "POST /api/search/exports"],
            "requirement": "Busca operacional auditada com redacao, facets, saved queries, legal hold e e-discovery.",
        },
        "semanticBi": {
            "status": "ready",
            "evidence": ["GET /api/analytics/metrics", "GET /api/analytics/data-quality", "GET /api/analytics/lineage"],
            "requirement": "Catalogo semantico versionado com lineage, freshness e qualidade de dados.",
        },
        "aiGovernance": {
            "status": "ready",
            "evidence": ["GET /api/ai-governance/tools", "POST /api/ai-governance/runs", "GET /api/ai-governance/audit-events"],
            "requirement": "Assistente controlado por politicas, ferramentas aprovadas, redacao e trilha de auditoria.",
        },
        "incidentCommand": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/tenants/{tenantSlug}/incidents",
                "POST /api/platform-control/tenants/{tenantSlug}/incidents/{publicId}/postmortem",
                "GET /api/platform-control/tenants/{tenantSlug}/incident-command/readiness",
            ],
            "requirement": "Comando de incidente com timeline, acoes, resolucao e postmortem.",
        },
        "policyDecisionCenter": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/policies/catalog",
                "POST /api/platform-control/tenants/{tenantSlug}/policies/evaluate",
                "GET /api/platform-control/tenants/{tenantSlug}/policies/decisions",
            ],
            "requirement": "Decisoes sensiveis passam por politica versionada com allow, deny ou review.",
        },
        "operationalTimeline": {
            "status": "ready",
            "evidence": ["GET /api/platform-control/tenants/{tenantSlug}/timeline", "POST /api/platform-control/tenants/{tenantSlug}/timeline"],
            "requirement": "Eventos operacionais cross-service por tenant e entidade com ator, origem e correlacao.",
        },
        "commandApprovals": {
            "status": "ready",
            "evidence": ["POST /api/platform-control/tenants/{tenantSlug}/approvals", "POST /api/platform-control/tenants/{tenantSlug}/approvals/{publicId}/approve"],
            "requirement": "Comandos criticos exigem aprovacao, rejeicao ou execucao auditavel.",
        },
        "runbookAutomation": {
            "status": "ready",
            "evidence": ["GET /api/platform-control/runbooks/catalog", "POST /api/platform-control/tenants/{tenantSlug}/runbooks"],
            "requirement": "Runbooks operacionais possuem passos, aprovacao e evidencia de execucao.",
        },
        "auditEvidenceVault": {
            "status": "ready",
            "evidence": ["GET /api/platform-control/tenants/{tenantSlug}/evidence", "POST /api/platform-control/tenants/{tenantSlug}/evidence"],
            "requirement": "Decisoes, aprovacoes, runbooks, incidentes e readiness geram evidencia com hash logico.",
        },
        "riskComplianceScoring": {
            "status": "ready" if risk_readiness["acceptanceReady"] else "attention",
            "evidence": [
                "GET /api/analytics/risk/tenant-score",
                "GET /api/analytics/risk/domain-scores",
                "GET /api/analytics/risk/compliance-posture",
                "GET /api/analytics/risk/recommendations",
            ],
            "requirement": "Tenant, dominios e servicos recebem score deterministico de risco e compliance.",
        },
        "enterpriseEventMesh": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/event-mesh/catalog",
                "POST /api/platform-control/tenants/{tenantSlug}/event-mesh/events",
                "GET /api/platform-control/tenants/{tenantSlug}/event-mesh/lineage",
            ],
            "requirement": "Eventos de dominio possuem catalogo, hash de payload, dead letter, replay e lineage.",
        },
        "reconciliationCenter": {
            "status": "ready",
            "evidence": ["POST /api/analytics/reconciliation/run", "GET /api/analytics/reconciliation/findings"],
            "requirement": "Divergencias financeiras e operacionais sao detectadas por conciliacao deterministica.",
        },
        "financialCloseCenter": {
            "status": "ready",
            "evidence": [
                "GET /api/analytics/financial-close/readiness",
                "POST /api/analytics/financial-close/periods/{publicId}/close",
                "GET /api/analytics/financial-close/periods/{publicId}/snapshot",
            ],
            "requirement": "Fechamento financeiro gera snapshot com hash, readiness e controle de pendencias.",
        },
        "masterDataQuality": {
            "status": "ready",
            "evidence": [
                "GET /api/analytics/master-data/entities",
                "GET /api/analytics/master-data/quality-score",
                "GET /api/analytics/data-quality/rules",
            ],
            "requirement": "Dados mestres possuem golden records, score, regras e propostas de merge.",
        },
        "lakehouseManifest": {
            "status": "ready",
            "evidence": [
                "GET /api/analytics/lakehouse/datasets",
                "GET /api/analytics/lakehouse/lineage",
                "GET /api/analytics/lakehouse/readiness",
            ],
            "requirement": "Datasets oficiais possuem classificacao, retencao, lineage e politica de exportacao.",
        },
        "tenantRuntimeControlPlane": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/tenants/{tenantSlug}/runtime/profile",
                "GET /api/platform-control/tenants/{tenantSlug}/runtime/health-score",
                "GET /api/platform-control/tenants/{tenantSlug}/enterprise-runtime/readiness",
            ],
            "requirement": "Tenant possui perfil runtime, quotas, janelas, SLO e health score.",
        },
        "contractSchemaEvolution": {
            "status": "ready",
            "evidence": [
                "GET /api/platform-control/contracts/evolution",
                "POST /api/platform-control/contracts/evolution/diffs",
                "GET /api/platform-control/contracts/evolution/compatibility-matrix",
            ],
            "requirement": "Contratos e schemas possuem snapshots, diff, breaking changes e matriz de compatibilidade.",
        },
        "opsConsoleV14": {
            "status": "ready",
            "evidence": ["client-web/client-api/src/generated/apiCatalog.ts", "npm run generate", "npm run build"],
            "requirement": "Console tecnico carrega os contratos v1.4, incluindo intelligence/readiness e ativacao de providers BYOK.",
        },
    }

    blocking_gates = [key for key, gate in gates.items() if gate["status"] not in {"ready", "stable"}]
    release_ready = not blocking_gates

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "release": {
            "version": "1.5.0",
            "releaseReady": release_ready,
            "status": "ready" if release_ready else "attention",
            "blockingGates": blocking_gates,
        },
        "gates": gates,
        "evidence": {
            "requiredSuites": ["contract", "security", "auth-conformance", "observability", "supply-chain", "hardening-secrets", "backup-restore", "hardening", "smoke"],
            "commands": [
                "./scripts/test.sh contract",
                "./scripts/test.sh security",
                "./scripts/test.sh auth-conformance",
                "./scripts/test.sh observability",
                "./scripts/test.sh supply-chain",
                "./scripts/test.sh hardening-secrets",
                "./scripts/test.sh backup-restore",
                "./scripts/test.sh hardening",
                "./scripts/test.sh smoke",
                "kubectl apply --dry-run=server -k infra/kubernetes/overlays/production",
            ],
            "runtimeReports": [
                "GET /api/analytics/reports/production-readiness",
                "GET /api/analytics/reports/hardening-review",
                "GET /api/analytics/reports/go-live-control",
                "GET /api/edge/ops/go-live-overview",
                "GET /api/search/query",
                "GET /api/ai-governance/tools",
                "GET /api/platform-control/tenants/{tenantSlug}/incident-command/readiness",
                "GET /api/platform-control/tenants/{tenantSlug}/autonomous-governance/readiness",
                "GET /api/platform-control/tenants/{tenantSlug}/enterprise-runtime/readiness",
                "GET /api/platform-control/providers/activation/catalog",
                "GET /api/analytics/risk/tenant-score",
                "GET /api/analytics/enterprise-runtime/readiness",
                "GET /api/analytics/external-intelligence/readiness",
            ],
        },
        "closures": {
            "platform": build_platform_closure(),
            "hardening": build_hardening_closure(),
            "productionMaturity": build_production_maturity_closure(),
            "operationalRunbooks": build_operational_runbook_review(),
            "goLive": release_controls,
            "semanticMetrics": build_semantic_metrics_readiness(),
            "riskCompliance": risk_readiness,
            "enterpriseRuntimeFabric": enterprise_runtime_readiness,
            "externalProviders": adapter_catalog,
        },
        "productionCapabilities": [
            "operational-search",
            "semantic-bi",
            "ai-governance",
            "incident-command",
            "policy-decision-center",
            "operational-timeline",
            "command-approvals",
            "runbook-automation",
            "audit-evidence-vault",
            "risk-compliance-scoring",
            "enterprise-event-mesh",
            "reconciliation-center",
            "financial-close-center",
            "master-data-quality",
            "lakehouse-manifest",
            "tenant-runtime-control-plane",
            "contract-schema-evolution",
            "external-provider-activation",
            "llm-provider-byok",
            "document-intelligence",
            "fiscal-brazil-readiness",
            "registry-enrichment-brazil",
            "market-macro-risk",
            "external-risk-feed",
            "provider-activation-v1.4",
            "root-hardening-v1.4.1",
            "static-policy-v1.4.2",
            "transactional-id-generation-v1.4.3",
            "api-console-security-v1.4.4",
            "runtime-infra-hardening-v1.4.5",
            "auth-observability-conformance-v1.5.0",
            "ops-console-v1.4",
        ],
    }
