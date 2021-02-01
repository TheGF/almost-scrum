package core

import (
	"io"
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
	Parent  string    `json:"parent"`
}

// ListLibrary returns the content of the specified path in the library.
func ListLibrary(project *Project, path string) ([]LibraryItem, string, error) {

	path, _ = GetPathInLibrary(project, path)
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

func MoveFileInLibrary(project *Project, oldPath string, path string) error {
	oldPath = filepath.Join(project.Path, ProjectLibraryFolder, oldPath)
	path = filepath.Join(project.Path, ProjectLibraryFolder, path)

	return os.Rename(oldPath, path)
}

func CreateFolderInLibrary(project *Project, path string) error {
	path = filepath.Join(project.Path, ProjectLibraryFolder, path)
	return os.MkdirAll(path, 0755)
}

func DeleteFileFromLibrary(project *Project, path string, recursive bool) error {
	path = filepath.Join(project.Path, ProjectLibraryFolder, path)
	if recursive {
		return os.RemoveAll(path)
	} else {
		return os.Remove(path)
	}
}

// GetPathInLibrary returns the absolute path for a resource stored in the library.
func GetPathInLibrary(project *Project, path string) (string, error) {
	path = filepath.Join(project.Path, ProjectLibraryFolder, path)
	return filepath.Abs(path)
}

func SetFileInLibrary(project *Project, path string, reader io.ReadCloser) error {
	path = filepath.Join(project.Path, ProjectLibraryFolder, path)
	writer, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, reader)
	return err
}

// GetPathInLibrary returns the absolute path for a resource stored in the library.
func GetLibraryItems(project *Project, paths []string) ([]LibraryItem, error) {
	items := make([]LibraryItem, 0, len(paths))

	for _, path := range paths {
		fullPath := filepath.Join(project.Path, ProjectLibraryFolder, path)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		mime, _ := mimetype.DetectFile(path)
		items = append(items, LibraryItem{
			Name:    fileInfo.Name(),
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
			Mime:    mime.String(),
			Dir:     fileInfo.IsDir(),
			Parent:  filepath.Dir(path),
		})
	}

	return items, nil
}
