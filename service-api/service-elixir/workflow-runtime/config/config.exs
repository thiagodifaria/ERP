import Config

config :workflow_runtime,
  http_host: System.get_env("WORKFLOW_RUNTIME_HTTP_HOST", "0.0.0.0"),
  http_port: String.to_integer(System.get_env("WORKFLOW_RUNTIME_HTTP_PORT", "8085")),
  bootstrap_tenant_slug: System.get_env("WORKFLOW_RUNTIME_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops")
