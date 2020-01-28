package main

import (
	"encoding/json"
	"github.com/joaodiogocosta/retext/client"
	"github.com/joaodiogocosta/retext/watcher"
)

type EntryEvent struct {
	Action string `json:"name"`
	Entries []*watcher.Entry `json:"entries"`
}

func NewEntryEvent(action string, entries []*watcher.Entry) EntryEvent {
	return EntryEvent{
		Action: action,
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
	rootPath := ".."
	conn := client.Connect()

	rootEntries := watcher.GetRootEntries(rootPath)
	entryEvent := NewEntryEvent(watcher.ActionUpdate, rootEntries)
	conn.SendCh <- EntryEventToJson(entryEvent)

	watcher.Watch(rootPath, func(action string, entry watcher.Entry) {
		entryEvent := NewEntryEvent(action, []*watcher.Entry{&entry})
		conn.SendCh <- EntryEventToJson(entryEvent)
	})
}
