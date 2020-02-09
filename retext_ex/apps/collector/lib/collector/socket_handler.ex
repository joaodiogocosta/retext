defmodule Collector.SocketHandler do
  require Logger

  @behaviour :cowboy_websocket
  @idle_timeout 60000

  def init(request, state) do
    {:cowboy_websocket, request, state, %{idle_timeout: @idle_timeout}}
  end

  def websocket_init(state) do
    {:ok, state}
  end

  def websocket_handle(:ping, state) do
    # It's not necessary to reply with a :pong frame
    # because cowboy does it automatically for us
    {:ok, state}
  end

  def websocket_handle({:text, message}, state) do
    IO.puts(message)
    {:ok, state}
  end

  def websocket_handle(any, state) do
    Logger.info("Received unknown message: #{IO.inspect(any)}")
    {:ok, state}
  end

  def websocket_info(info, state) do
    # IO.puts("state", state)
    {:reply, {:text, info}, state}
  end

  def websocket_terminate(_reason, _req, _state) do
    :ok
  end
end
