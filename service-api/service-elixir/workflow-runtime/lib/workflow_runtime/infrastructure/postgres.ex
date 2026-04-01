defmodule WorkflowRuntime.Infrastructure.Postgres do
  @moduledoc false

  def child_spec(_opts) do
    Postgrex.child_spec(connection_options())
  end

  def enabled? do
    Application.fetch_env!(:workflow_runtime, :repository_driver) == "postgres"
  end

  def ready? do
    Process.whereis(__MODULE__) != nil
  end

  def query!(statement, params \\ []) do
    Postgrex.query!(__MODULE__, statement, params)
  end

  def transaction(fun) do
    Postgrex.transaction(__MODULE__, fun)
  end

  defp connection_options do
    [
      name: __MODULE__,
      hostname: Application.fetch_env!(:workflow_runtime, :postgres_host),
      port: Application.fetch_env!(:workflow_runtime, :postgres_port),
      database: Application.fetch_env!(:workflow_runtime, :postgres_database),
      username: Application.fetch_env!(:workflow_runtime, :postgres_username),
      password: Application.fetch_env!(:workflow_runtime, :postgres_password),
      ssl: Application.fetch_env!(:workflow_runtime, :postgres_ssl_mode) == "require"
    ]
  end
end
