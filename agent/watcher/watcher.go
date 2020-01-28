package watcher

// IDEIA: Only tracks modified files by default.
// Tracks all files if options is passed

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path/filepath"
)

const ActionUpdate = "UPDATE"
const ActionRemove = "REMOVE"

type Entry struct {
	Path string `json:"path"`
	IsDir bool `json:"isDir"`
}

type EntryEvent struct {
	Action string `json:"name"`
	Entries []*Entry `json:"entries"`
}

var ignored = ".git"

func NewEntryEvent(action string, entries []*Entry) EntryEvent {
	return EntryEvent{
		Action: action,
		Entries: entries,
	}
}

func getPath(rootPath string, filename string) string {
	return filepath.Join(rootPath, filename)
}

func GetRootEntries(path string) []*Entry {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	var rootEntries []*Entry
	for _, f := range files {
		filename := f.Name()

		if filename == ignored {
			fmt.Printf("ignored: %s\n", filename)
			continue
		}

		entry := &Entry{Path: getPath(path, filename), IsDir: f.IsDir()}
		rootEntries = append(rootEntries, entry)
	}

	return rootEntries
}

func isUpdate(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write ||
	event.Op&fsnotify.Create == fsnotify.Create
}

func isRemove(event fsnotify.Event) bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove ||
	event.Op&fsnotify.Rename == fsnotify.Rename
}

func handleEvent(eventsCh chan EntryEvent, event fsnotify.Event) {
	if isUpdate(event) {
		entry := Entry{Path: event.Name, IsDir: false}
		entryEvent := NewEntryEvent(ActionUpdate, []*Entry{&entry})
		eventsCh <- entryEvent
	} else if isRemove(event) {
		entry := Entry{Path: event.Name}
		entryEvent := NewEntryEvent(ActionRemove, []*Entry{&entry})
		eventsCh <- entryEvent
	}
}

func Watch(path string, eventsCh chan EntryEvent) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case fsevent, ok := <-watcher.Events:

				if !ok {
					fmt.Println("not ok")
					return
				}

				handleEvent(eventsCh, fsevent)

			case err, ok := <-watcher.Errors:
				
				if !ok {
					fmt.Println("not ok error")
					return
				}

				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)

	if err != nil {
		log.Fatal(err)
	}
	
	select {}
}
