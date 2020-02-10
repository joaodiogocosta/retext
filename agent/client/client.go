package client

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ConnectionAdapter interface {
	Connect(*Session, chan []byte)
}

type Connection struct {
	SendCh chan []byte
	adapter ConnectionAdapter
}

type Session struct {
	id string
	token string
}

const (
	WebsocketAdapter = iota
	DryAdapter = iota
)

func getAdapter(useAdapter int) ConnectionAdapter {
	var adapter ConnectionAdapter

	switch adp := useAdapter; adp {
	case DryAdapter:
		adapter = &ConnDryAdapter{}
	case WebsocketAdapter:
		adapter = &ConnWsAdapter{}
	default:
		adapter = &ConnWsAdapter{}
	}

	return adapter
}

func Connect(session *Session, useAdapter int) *Connection {
	adapter := getAdapter(useAdapter)

	conn := &Connection{
		SendCh: make(chan []byte),
		adapter: adapter,
	}

	go conn.adapter.Connect(session, conn.SendCh)

	return conn
}

func NewSession() *Session {
	url := url.URL{Scheme: "http", Host: "localhost:4001", Path: "/sessions"}
	requestBody := bytes.NewBuffer([]byte{})
	resp, err := http.Post(url.String(), "text/plain", requestBody)

	if err != nil {
		log.Println("Error getting new session", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading body", err)
	}

	session_values := strings.Split(string(body), "|")
	return &Session{ id: session_values[0], token: session_values[1] }
}
