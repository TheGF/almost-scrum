// Package core provides basic functionality for Almost Scrum
package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Store structure
type Store struct {
	path string
}

// StoreItem is the result of List operation
type StoreItem struct {
	path    string
	modtime time.Time
	dir     bool
}

// List the content of a store
func List(s Store) []StoreItem {
	log.Printf("List content of store at %s", s.path)
	var list []StoreItem = make([]StoreItem, 0, 100)

	rootLen := len(s.path)
	filepath.Walk(s.path, func(parent string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		var dir = info.IsDir()
		var path = fmt.Sprintf("%s/%s", parent, info.Name())[1+rootLen:]
		if path == ".." {
			return nil
		}

		list = append(list, StoreItem{
			path:    path,
			modtime: info.ModTime(),
			dir:     dir,
		})
		log.Printf("List lenght: %d", len(list))
		return nil
	})
	log.Printf("Final list lenght: %d", len(list))
	return list
}

// Set a story in the Store
func Set(s Store, path string, story Story) {

}
