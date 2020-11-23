// Package core provides basic functionality for Almost Scrum
package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
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

// List the content of a store
func List(s Store) []StoreItem {
	log.Printf("List content of store at %s", s.Path)
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
	return list
}

// Get a story in the Store
func Get(s Store, path string) (story Story, err error) {
	path = fmt.Sprintf("%s/%s", s.Path, path)
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof("Invalid file %s: %v", path, err)
		return
	}

	err = yaml.Unmarshal(d, &story)
	if err != nil {
		log.Infof("Invalid file %s: %v", path, err)
		return
	}
	return
}

//Set a story in the Store
func Set(s Store, path string, story *Story) (err error) {
	d, err := yaml.Marshal(&story)
	if err != nil {
		log.Infof("Cannot marshal story %s: %v", path, err)
		return
	}
	path = filepath.Join(s.Path, path)
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Infof("Cannot save story %s: %v", path, err)
		return
	}
	log.Infof("Story saved to %s", path)
	return
}
