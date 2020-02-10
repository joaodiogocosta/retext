defmodule Collector.Sessions do
  use GenServer

  def start_link(_) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  @impl true
  def init(state) do
    {:ok, state}
  end

  @impl true
	def handle_call(:create, _from, state) do
    id = UUID.uuid4()
    token = UUID.uuid4()
    state = Map.put_new(state, id, token)
    {:reply, {:ok, id, token}, state}
	end

  @impl true
	def handle_call({:get, id, token}, _from, state) do
    with {:ok, found_token} <- Map.fetch(state, id),
      true <- token == found_token do
      {:reply, :ok, state}
    else
      _ -> {:reply, nil, state}
    end
	end

  @impl true
	def handle_call({:remove, id}, _from, state) do
	  state = Map.delete(state, id)
    {:reply, :ok, state}
	end

  def create() do
    GenServer.call(__MODULE__, :create)
  end

  def get(id, token) do
    GenServer.call(__MODULE__, {:get, id, token})
  end

  def remove(id) do
    GenServer.call(__MODULE__, {:remove, id})
  end
end
