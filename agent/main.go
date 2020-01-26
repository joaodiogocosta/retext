package main

import (
	"encoding/json"
	"time"
	"github.com/joaodiogocosta/retext/client"
	"github.com/joaodiogocosta/retext/watcher"
)

const UPDATE = "UPDATE"

type EntryEvent struct {
	Name string `json:"name"`
	Entries []*watcher.Entry `json:"entries"`
}

func NewEntryEvent(name string, entries []*watcher.Entry) EntryEvent {
	return EntryEvent{
		Name: name,
		Entries: entries,
	}
}

func EntryEventToJson(event EntryEvent) []byte {
	message, err := json.Marshal(event)

	if err != nil {
		panic(err)
	}

	return message
}

func main() {
	conn := client.Connect()

	rootEntries := watcher.GetRootEntries()
	entryEvent := NewEntryEvent(UPDATE, rootEntries)

	conn.SendCh <- EntryEventToJson(entryEvent)
	
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
