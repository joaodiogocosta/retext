defmodule Collector.SocketHandler do
  require Logger
  alias Collector.Sessions

  @behaviour :cowboy_websocket
  @idle_timeout 60000

  @impl true
  def init(request, state) do
    case authorize(request) do
      {:ok, id} -> authorized(id, request, %{})
      _ -> unauthorized(request, state)
    end
  end

  @impl true
  def websocket_init(state) do
    {:ok, state}
  end

  @impl true
  def websocket_handle(:ping, state) do
    # It's not necessary to reply with a :pong frame
    # because cowboy does it automatically for us
    {:ok, state}
  end

  @impl true
  def websocket_handle({:text, message}, state) do
    IO.puts(message)
    {:ok, state}
  end

  @impl true
  def websocket_handle(any, state) do
    Logger.info("Received unknown message: #{IO.inspect(any)}")
    {:ok, state}
  end

  @impl true
  def websocket_info(info, state) do
    {:reply, {:text, info}, state}
  end

  @impl true
  def terminate(_reason, _partialReq, %{ id: id }) do
    Sessions.remove(id)
    :ok
  end

  @impl true
  def terminate(_, _, _) do
    :ok
  end

  def authorize(%{ headers: %{ "x-id" => id, "x-token" => token } }) do
    case Sessions.get(id, token) do
      :ok -> {:ok, id}
      _ -> :error
    end
  end

  def authorize(_) do
    :error
  end

  def authorized(id, request, state, opts \\ %{}) do
    state = Map.put(state, :id, id)
    opts = Map.put_new(opts, :idle_timeout, @idle_timeout)
    {:cowboy_websocket, request, state, opts}
  end

  def unauthorized(request, state) do
    req = :cowboy_req.reply(401, request)
    {:ok, req, state}
  end
end
