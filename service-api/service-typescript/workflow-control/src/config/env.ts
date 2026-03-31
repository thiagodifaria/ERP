export type WorkflowControlConfig = {
  serviceName: "workflow-control";
  repositoryDriver: "memory" | "postgres";
  bootstrapTenantSlug: string;
  postgresHost: string;
  postgresPort: string;
  postgresDatabase: string;
  postgresUser: string;
  postgresPassword: string;
  postgresSslMode: string;
};

function envOrDefault(key: string, fallback: string): string {
  return process.env[key] ?? fallback;
}

export function loadConfig(): WorkflowControlConfig {
  return {
    serviceName: "workflow-control",
    repositoryDriver: envOrDefault("WORKFLOW_CONTROL_REPOSITORY_DRIVER", "memory") as "memory" | "postgres",
    bootstrapTenantSlug: envOrDefault("WORKFLOW_CONTROL_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
    postgresHost: envOrDefault("WORKFLOW_CONTROL_POSTGRES_HOST", envOrDefault("ERP_POSTGRES_HOST", "localhost")),
    postgresPort: envOrDefault("WORKFLOW_CONTROL_POSTGRES_PORT", "5432"),
    postgresDatabase: envOrDefault("WORKFLOW_CONTROL_POSTGRES_DB", envOrDefault("ERP_POSTGRES_DB", "erp")),
    postgresUser: envOrDefault("WORKFLOW_CONTROL_POSTGRES_USER", envOrDefault("ERP_POSTGRES_USER", "erp")),
    postgresPassword: envOrDefault("WORKFLOW_CONTROL_POSTGRES_PASSWORD", envOrDefault("ERP_POSTGRES_PASSWORD", "erp")),
    postgresSslMode: envOrDefault("WORKFLOW_CONTROL_POSTGRES_SSL_MODE", "disable")
  };
}
