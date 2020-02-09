defmodule Collector.Endpoint do
  @moduledoc """
    A Plug responsible for logging request info, parsing request
    body's as JSON, matching routes, and dispatching responses.
  """

  use Plug.Router

  # plug(Plug.Logger)
  plug(:match)
  plug(:dispatch)
  # we will only parse the request AFTER there is a route match.
  # plug(Plug.Parsers, parsers: [:json], json_decoder: Poison)
  # responsible for dispatching responses

  match _ do
    IO.puts("nothing to see")
    send_resp(conn, 404, "oops... Nothing here :(")
  end
end
