package main

import (
	"encoding/json"
	"github.com/joaodiogocosta/retext/cli"
	"github.com/joaodiogocosta/retext/client"
	"github.com/joaodiogocosta/retext/watcher"
)

func toJson(event interface{}) []byte {
	message, err := json.Marshal(event)

	if err != nil {
		panic(err)
	}

	return message
}

func preloadTree(conn *client.Connection, rootPaths []string) {
	var rootEntries []*watcher.Entry
	for _, path := range rootPaths {
		rootEntries = append(rootEntries, watcher.GetRootEntries(path)...)
	}
	rootEntriesEvent := watcher.NewEntryEvent(watcher.ActionUpdate, rootEntries)
	conn.SendCh <- toJson(rootEntriesEvent)
}

func main() {
	args := cli.Parse()
	session := client.NewSession()
	conn := client.Connect(session, args.ConnectionAdapter)

	// if err != nil {
	// 	return
	// }

	entryEvents := make(chan watcher.EntryEvent)
	go func() {
		for entryEvent := range entryEvents {
			conn.SendCh <- toJson(entryEvent)
		}
	}()
	watcher.Watch(args.RootPaths, entryEvents)
}
