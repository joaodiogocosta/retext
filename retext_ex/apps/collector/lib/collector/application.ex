defmodule Collector.Application do
  @moduledoc false

  use Application

  def start(_type, _args) do
    children = [
      websocket_server()
    ]

    opts = [strategy: :one_for_one, name: Collector.Supervisor]
    Supervisor.start_link(children, opts)
  end

  defp websocket_server do
    Plug.Cowboy.child_spec(
      scheme: :http,
      plug: Collector.Endpoint,
      options: [port: 4001, dispatch: dispatch()]
    )
  end

  defp dispatch do
    [{:_,[
      {"/ws", Collector.SocketHandler, []},
      {:_, Plug.Cowboy.Handler, {Collector.Endpoint, []}}
    ]}]
  end
end
