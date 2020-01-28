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

type Entry struct {
	Path string `json:"path"`
	IsDir bool `json:"isDir"`
}

var ignored = ".git"

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

func Watch(path string, cb func(string, Entry)) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:

				if !ok {
					fmt.Println("not ok")
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					entry := Entry{Path: event.Name, IsDir: false}
					cb(ActionUpdate, entry)
				}

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
