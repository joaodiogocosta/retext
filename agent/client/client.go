package client

import (
	"fmt"
)

type Connection struct {
	SendCh chan []byte
}

func (conn *Connection) Listen() {
	for message := range conn.SendCh {
		fmt.Println(string(message))
	}
}

func Connect() *Connection {
	conn := &Connection{
		SendCh: make(chan []byte),
	}

	go conn.Listen()

	return conn
}
