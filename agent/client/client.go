package client

type ConnectionAdapter interface {
	Connect(chan []byte)
}

type Connection struct {
	SendCh chan []byte
	adapter ConnectionAdapter
}

const (
	WebsocketAdapter = iota
	DryAdapter = iota
)

func Connect(useAdapter int) *Connection {
	var adapter ConnectionAdapter

	switch adp := useAdapter; adp {
	case DryAdapter:
		adapter = &ConnDryAdapter{}
	case WebsocketAdapter:
		adapter = &ConnWsAdapter{}
	default:
		adapter = &ConnWsAdapter{}
	}

	conn := &Connection{
		SendCh: make(chan []byte),
		adapter: adapter,
	}

	go conn.adapter.Connect(conn.SendCh)

	return conn
}
