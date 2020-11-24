// Package core provides basic functionality for Almost Scrum
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// Store structure
type Store struct {
	Path string
}

// StoreItem is the result of List operation
type StoreItem struct {
	Path    string    `json:"path"`
	ModTime time.Time `json:"modTime"`
	Dir     bool      `json:"dir"`
}

// ListStore the content of a store
func ListStore(s Store) []StoreItem {
	log.Debugf("ListStore - List content of store at %s", s.Path)
	var list []StoreItem = make([]StoreItem, 0, 100)

	rootLen := len(s.Path)
	filepath.Walk(s.Path, func(parent string, info os.FileInfo, err error) error {
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
			Path:    path,
			ModTime: info.ModTime(),
			Dir:     dir,
		})
		return nil
	})
	log.Debugf("ListStore - Found items: %v", list)
	return list
}

// GetStory a story in the Store
func GetStory(s Store, path string) (story Story, err error) {
	path = filepath.Join(s.Path, path)
	err = ReadYaml(path, &story)
	if err != nil {
		log.Errorf("Cannot read story %s: %v", path, err)
	} else {
		log.Debugf("Story %s read %+v", path, story)
	}
	return
}

//SetStory a story in the Store
func SetStory(s Store, path string, story *Story) (err error) {
	path = filepath.Join(s.Path, path)
	err = WriteYaml(path, story)
	if err != nil {
		log.Errorf("Cannot write story %s: %v", path, err)
	} else {
		log.Debugf("Story saved to %s", path)
	}
	return
}
