package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"net/url"
	"time"
)

type ConnWsAdapter struct {}

const pingPeriod = 5 * 1000 * 1000 * 1000;

func ping(ws *websocket.Conn) error {
	return ws.WriteMessage(websocket.PingMessage, []byte{})
}

func writeMessage(ws *websocket.Conn, message []byte) error {
	return ws.WriteMessage(websocket.TextMessage, message)
}

func (adapter *ConnWsAdapter) Connect(session *Session, sendCh chan []byte) {
	header := &http.Header{}
	header.Add("X-Id", session.id)
	header.Add("X-Token", session.token)

	u := url.URL{Scheme: "ws", Host: "localhost:4001", Path: "/ws"}
	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), *header)

	if err != nil {
		fmt.Println("dial:", err)

		if resp.StatusCode == 401 {
			log.Println("Unauthorized")
			return
		}
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
		case message := <-sendCh:
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
