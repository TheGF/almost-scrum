package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

// LibraryItem contains basic information about items in the library.
type LibraryItem struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Mime    string    `json:"mime"`
	Dir     bool      `json:"dir"`
}

// ListLibrary returns the content of the specified path in the library.
func ListLibrary(p Project, path string) ([]LibraryItem, string, error) {

	path, _ = GetPathInLibrary(p, path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Warnf("ListLibrary - Cannot get info about file %s: %v", path, err)
		return nil, "", err
	}
	if fileInfo.IsDir() {
		list, err := ioutil.ReadDir(path)
		if err != nil {
			log.Warnf("ListLibrary - Cannot list content of folder %s: %v", path, err)
		}

		items := make([]LibraryItem, 0, len(list))
		for _, fileInfo := range list {
			name := fileInfo.Name()
			mime, _ := mimetype.DetectFile(filepath.Join(path, name))

			items = append(items, LibraryItem{
				Name:    fileInfo.Name(),
				Size:    fileInfo.Size(),
				ModTime: fileInfo.ModTime(),
				Mime:    mime.String(),
				Dir:     fileInfo.IsDir(),
			})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].ModTime.After(items[j].ModTime)
		})

		log.Debugf("ListLibrary - Path %s is a folder: %v", path, items)
		return items, path, nil
	}

	log.Debugf("ListLibrary - Path %s is a file", path)
	return nil, path, nil
}

// GetPathInLibrary returns the absolute path for a resource stored in the library.
func GetPathInLibrary(p Project, path string) (string, error) {
	path = filepath.Join(p.Path, ProjectLibraryFolder, path)
	return filepath.Abs(path)
}
