import Config

config :workflow_runtime,
  http_host: System.get_env("WORKFLOW_RUNTIME_HTTP_HOST", "0.0.0.0"),
  http_port: String.to_integer(System.get_env("WORKFLOW_RUNTIME_HTTP_PORT", "8085")),
  bootstrap_tenant_slug: System.get_env("WORKFLOW_RUNTIME_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
  repository_driver: System.get_env("WORKFLOW_RUNTIME_REPOSITORY_DRIVER", "memory"),
  postgres_host: System.get_env("WORKFLOW_RUNTIME_POSTGRES_HOST", "service-postgresql"),
  postgres_port: String.to_integer(System.get_env("WORKFLOW_RUNTIME_POSTGRES_PORT", "5432")),
  postgres_database: System.get_env("WORKFLOW_RUNTIME_POSTGRES_DB", "erp"),
  postgres_username: System.get_env("WORKFLOW_RUNTIME_POSTGRES_USER", "erp"),
  postgres_password: System.get_env("WORKFLOW_RUNTIME_POSTGRES_PASSWORD", "erp"),
  postgres_ssl_mode: System.get_env("WORKFLOW_RUNTIME_POSTGRES_SSL_MODE", "disable")
