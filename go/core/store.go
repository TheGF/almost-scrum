// Package core provides basic functionality for Almost Scrum
package core

import (
	"os"
	"path/filepath"
	"sort"
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

// GetStore returns the specified store
func GetStore(project Project, store string) (Store, error) {
	path := filepath.Join(project.Path, "stores", store)
	if fileInfo, err := os.Stat(path); err != nil || !fileInfo.IsDir() {
		return Store{}, ErrNoFound
	}
	return Store{Path: path}, nil
}

// CreateStore creates a new store inside a project
func CreateStore(project Project, name string) (Store, error) {
	path := filepath.Join(project.Path, "stores", name)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return Store{}, err
	}
	return Store{path}, nil
}

// CreateSprint creates a new store with name sprint-n, where n is next available integer.
func CreateSprint(project Project) (Store, error) {
	return CreateStore(project, "")
}

// TouchContent set the modified time to current time. It applies to stories and folders
func TouchContent(s Store, path string) error {
	currentTime := time.Now().Local()
	path = filepath.Join(s.Path, path)
	err := os.Chtimes(path, currentTime, currentTime)
	if err != nil {
		log.Errorf("TouchContent - Cannot touch %s", path)
		return err
	}
	log.Debugf("TouchContent - Touch %s", path)
	return nil
}

// CreateFolder creates a folder in the specified store.
func CreateFolder(s Store, path string) error {
	path = filepath.Join(s.Path, path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		log.Errorf("CreateFolder - Cannot create %s", path)
		return err
	}
	log.Debugf("CreateFolder - Created %s in store", path)
	return nil
}

// ListStore the content of a store
func ListStore(s Store) []StoreItem {
	log.Debugf("ListStore - List content of store at %s", s.Path)
	var list []StoreItem = make([]StoreItem, 0, 100)

	rootLen := len(s.Path)

	filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if len(path) == rootLen {
			return nil
		}
		path = path[rootLen:]
		log.Debugf("Check %s, %s", path, info.Name())
		if path == ".." {
			return nil
		}

		list = append(list, StoreItem{
			Path:    path,
			ModTime: info.ModTime(),
			Dir:     info.IsDir(),
		})
		return nil
	})

	sort.Slice(list, func(i, j int) bool {
		return list[i].ModTime.After(list[j].ModTime)
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
