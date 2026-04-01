defmodule WorkflowRuntime.Server do
  @moduledoc false

  def child_spec(_opts) do
    host = Application.fetch_env!(:workflow_runtime, :http_host)
    port = Application.fetch_env!(:workflow_runtime, :http_port)

    Plug.Cowboy.child_spec(
      scheme: :http,
      plug: WorkflowRuntime.Api.Router,
      options: [ip: parse_host(host), port: port]
    )
  end

  defp parse_host("0.0.0.0"), do: {0, 0, 0, 0}
  defp parse_host("127.0.0.1"), do: {127, 0, 0, 1}
  defp parse_host(_host), do: {0, 0, 0, 0}
end
