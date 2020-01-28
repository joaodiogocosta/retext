package watcher

// IDEIA: Only tracks modified files by default.
// Tracks all files if options is passed

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
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

// TODO: ignore dotfiles by default
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

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
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
	var action string
	entry := Entry{Path: event.Name}

	if isUpdate(event) {
		entry.IsDir = isDir(event.Name)
		action = ActionUpdate
	} else if isRemove(event) {
		action = ActionRemove
	} else {
		return
	}

	entryEvent := NewEntryEvent(action, []*Entry{&entry})
	eventsCh <- entryEvent
}

func watchEvents(fswatcher *fsnotify.Watcher, eventsCh chan EntryEvent) {
	for {
		select {
		case fsevent, ok := <- fswatcher.Events:

			if !ok {
				fmt.Println("not ok")
				return
			}

			handleEvent(eventsCh, fsevent)

		case err, ok := <-fswatcher.Errors:

			if !ok {
				fmt.Println("not ok error")
				return
			}

			fmt.Println("error:", err)
		}
	}
}

func Watch(path string, eventsCh chan EntryEvent) {
	fswatcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer fswatcher.Close()

	go watchEvents(fswatcher, eventsCh)

	err = fswatcher.Add(path)

	if err != nil {
		log.Fatal(err)
	}
	
	select {}
}
