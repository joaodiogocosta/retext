package watcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Entry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	IsDir bool `json:"isDir"`
}

var ignored = ".git"

func getPath(rootPath string, filename string) string {
	return filepath.Join(rootPath, filename)
}

func GetRootEntries() []*Entry {
	rootPath := "./"
	files, err := ioutil.ReadDir(rootPath)

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

		entry := &Entry{Name: filename, Path: getPath(rootPath, filename), IsDir: f.IsDir()}
		rootEntries = append(rootEntries, entry)
	}

	return rootEntries
}
