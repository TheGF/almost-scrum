// Package core provides basic functionality for Almost Scrum
package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Store structure
type Store struct {
	Path string
}

// StoreItem is the result of List operation
type StoreItem struct {
	Name    string    `json:"name"`
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
func ListStore(s Store, path string) ([]StoreItem, error) {

	path = filepath.Join(s.Path, path)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Warnf("Error reading path store")
		return nil, err
	}

	var list []StoreItem = make([]StoreItem, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		list = append(list, StoreItem{
			Name:    fileInfo.Name(),
			ModTime: fileInfo.ModTime(),
			Dir:     fileInfo.IsDir(),
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ModTime.After(list[j].ModTime)
	})

	log.Debugf("ListStore - Found items: %v", list)
	return list, nil
}

func walkIn(root string, path string) ([]StoreItem, error) {
	items := make([]StoreItem, 0)
	logrus.Debugf("Walk into %s/%s", root, path)
	fileInfos, err := ioutil.ReadDir(filepath.Join(root, path))
	if err != nil {
		log.Warnf("Error reading path store")
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		name := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			subFolderItems, err := walkIn(root, name)
			if err != nil {
				return nil, err
			}
			items = append(items, subFolderItems...)
		} else {
			logrus.Debugf("Adding item %s to result", name)
			items = append(items, StoreItem{
				Name:    name,
				ModTime: fileInfo.ModTime(),
				Dir:     false,
			})
		}
	}
	return items, nil
}

// WalkStore returns the c the content of a store
func WalkStore(s Store) ([]StoreItem, error) {
	items, err := walkIn(s.Path, "")
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ModTime.After(items[j].ModTime)
	})

	return items, nil
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

// GetStoryAbsPath returns the absolute path of a story
func GetStoryAbsPath(s Store, path string) string {
	p, _ := filepath.Abs(filepath.Join(s.Path, path))
	return p
}

type _MetaItem struct {
	Owner string `json:"owner"`
}

type _Meta map[string]_MetaItem

// SetOwner sets the owner of the specified story
func setOwnership(s Store, path string, owner string) error {
	path = filepath.Join(s.Path, path)
	name := filepath.Base(path)
	folder := filepath.Dir(path)
	file := filepath.Join(folder, ".ash-meta.json")

	var meta _Meta
	err := ReadJSON(file, &meta)
	if err != nil {
		meta = map[string]_MetaItem{}
	}
	meta[name] = _MetaItem{
		Owner: owner,
	}
	return json.Unmarshal([]byte(owner), &meta)
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
