defmodule Collector.Endpoint do
  @moduledoc """
    A Plug responsible for logging request info, parsing request
    body's as JSON, matching routes, and dispatching responses.
  """

  use Plug.Router
  alias Collector.Sessions

  plug(Plug.Logger)
  plug(:match)
  plug(:dispatch)

  post "/sessions" do
    {:ok, id, session} = Sessions.create()
    send_resp(conn, 200, "#{id}|#{session}")
  end

  match _ do
    send_resp(conn, 404, "Not Found")
  end
end
