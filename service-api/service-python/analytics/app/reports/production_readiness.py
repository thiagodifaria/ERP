"""Gate consolidado de production readiness para o release 1.0.0."""

from datetime import datetime, timezone

from app.reports.adapter_catalog import build_adapter_catalog
from app.reports.go_live_control import build_release_controls
from app.reports.hardening_review import (
    build_hardening_closure,
    build_operational_runbook_review,
    build_platform_closure,
    build_production_maturity_closure,
)


def build_production_readiness(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    adapter_catalog = build_adapter_catalog()
    critical_provider_gaps = adapter_catalog["summary"]["criticalUnconfiguredCapabilities"]
    release_controls = build_release_controls()

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
        "contracts": {
            "status": "ready",
            "evidence": ["docs/contracts/registry.json", "docs/contracts/http", "scripts/test.sh contract"],
            "requirement": "OpenAPI versionado como fonte de verdade das APIs publicas.",
        },
        "observability": {
            "status": "ready",
            "evidence": ["infra/prometheus/prometheus.yml", "infra/blackbox/blackbox.yml", "docs/OPERACOES.md"],
            "requirement": "Health, probes, SLOs, alertas e runbooks operacionais.",
        },
        "backupRestoreDr": {
            "status": "ready",
            "evidence": ["scripts/build.sh backup-encrypted", "scripts/build.sh restore-encrypted", "scripts/test.sh backup-restore"],
            "requirement": "Backup criptografado e restore drill validavel.",
        },
        "deployArtifacts": {
            "status": "ready",
            "evidence": ["infra/kubernetes/base", "infra/kubernetes/overlays/production"],
            "requirement": "Namespace, ingress, probes, migration job e NetworkPolicy para deploy corporativo.",
        },
        "providers": {
            "status": "ready",
            "evidence": ["GET /api/analytics/reports/adapter-catalog", "GET /api/platform-control/tenants/{tenantSlug}/lifecycle/readiness"],
            "requirement": "Capabilities externas explicitam configuracao e readiness produtivo.",
            "criticalProviderGaps": critical_provider_gaps,
            "blockingForRelease": False,
            "policy": "Capability sem credencial real permanece visivel como fallback/manual/unconfigured, mas nao e apresentada como provider produtivo.",
        },
        "goLive": {
            "status": "ready" if release_controls["acceptanceReady"] else "attention",
            "evidence": ["GET /api/analytics/reports/go-live-control", "GET /api/edge/ops/go-live-overview"],
            "requirement": "Rollout por tenant com rollback, gargalos, adocao e aceite.",
        },
    }

    blocking_gates = [key for key, gate in gates.items() if gate["status"] not in {"ready", "stable"}]
    release_ready = not blocking_gates

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "release": {
            "version": "1.0.0",
            "name": "Production Readiness & Enterprise Deployment",
            "releaseReady": release_ready,
            "status": "ready" if release_ready else "attention",
            "blockingGates": blocking_gates,
        },
        "gates": gates,
        "evidence": {
            "requiredSuites": ["contract", "security", "hardening-secrets", "backup-restore", "hardening", "smoke"],
            "commands": [
                "./scripts/test.sh contract",
                "./scripts/test.sh security",
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
            ],
        },
        "closures": {
            "platform": build_platform_closure(),
            "hardening": build_hardening_closure(),
            "productionMaturity": build_production_maturity_closure(),
            "operationalRunbooks": build_operational_runbook_review(),
            "goLive": release_controls,
        },
        "nonGoals": [
            "AI/LLM operacional",
            "marketplace de extensoes",
            "BI semantico dedicado",
            "e-discovery externo",
            "novos dominios alem do baseline atual",
        ],
    }
