package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

const pingPeriod = 5 * 1000 * 1000 * 1000;

type Connection struct {
	SendCh chan []byte
}

func ping(ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.PingMessage, []byte{})
}

func writeMessage(ws *websocket.Conn, message []byte) error {
	return ws.WriteMessage(websocket.TextMessage, message)
}

func (conn *Connection) Open() {
	u := url.URL{Scheme: "ws", Host: "localhost:4001", Path: "/ws"}
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println("dial:", err)
	}

	defer ws.Close()

	ws.SetPongHandler(func(string) error {
		log.Println("pong")
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		defer close(done)
		for {
			_, receivedMessage, err := ws.ReadMessage()

			if err != nil {
				log.Println("received message:", err)
				return
			}

			log.Printf("recv: %s", receivedMessage)
		}
	}()

	for {
		select {
		case <-done:
			if err := ws.Close(); err != nil {
				log.Println("Error closing websocket: ", err)
			}
			return
		case message := <-conn.SendCh:
			if err := writeMessage(ws, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := ping(ws); err != nil {
				return
			}
		}
	}
}

func Connect() *Connection {
	conn := &Connection{
		SendCh: make(chan []byte),
	}

	go conn.Open()

	return conn
}
