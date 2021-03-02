package library

import (
	"almost-scrum/core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
)

// LibraryItem contains basic information about items in the library.
type Item struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Mime    string    `json:"mime"`
	Dir     bool      `json:"dir"`
	Owner   string    `json:"owner"`
	Parent  string    `json:"parent"`
}

func getItem(path string, fileInfo os.FileInfo) Item {
	name := fileInfo.Name()
	mime, _ := mimetype.DetectFile(filepath.Join(path, name))

	extendedAttr, err := getExtendedAttr(path, name)
	if err != nil {
		logrus.Warnf("Cannot get extended attr for %s: %v", path, err)
		return Item{}
	}
	return Item{
		Name:    fileInfo.Name(),
		Size:    fileInfo.Size(),
		ModTime: fileInfo.ModTime(),
		Mime:    mime.String(),
		Owner:   extendedAttr.Owner,
		Dir:     fileInfo.IsDir(),
	}
}

func IsDir(project *core.Project, path string) (bool, error) {
	p := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if fileInfo, err := os.Stat(p); err != nil {
		return false, err
	} else {
		return fileInfo.IsDir(), nil
	}
}

// List returns the content of the specified path in the library.
func List(project *core.Project, path string) ([]Item, error) {

	absPath, err := AbsPath(project, path)
	if err != nil {
		logrus.Warnf("List - Cannot get info about file %s: %v", path, err)
		return nil, err
	}
	fileInfos, err := ioutil.ReadDir(absPath)
	if err != nil {
		logrus.Warnf("List - Cannot fileInfos content of folder %s: %v", path, err)
	}

	versions := make(map[string]versionInfo)
	for _, fileInfo := range fileInfos {
		filterVersioned(project, path, fileInfo, versions)
	}

	items := make([]Item, 0, len(fileInfos))
	for _, versionInfo := range versions {
		items = append(items, getItem(absPath, versionInfo.fileInfo))
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ModTime.After(items[j].ModTime)
	})

	logrus.Debugf("List - Path %s is absPath folder: %v", path, items)
	return items, nil
}

func MoveFile(project *core.Project, oldPath string, path string) error {
	var err error
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if oldPath, err = AbsPath(project, oldPath); err != nil {
		return err
	}
	dir, name := filepath.Split(oldPath)
	xAttr, err := getExtendedAttr(dir, name); if err != nil {
		return err
	}

	err = os.Rename(oldPath, path); if err != nil {
		return err
	}
	dir_, name_ := filepath.Split(path)
	setExtendedAttr(dir_, name_, xAttr)
	setExtendedAttr(dir, name, nil)
	return nil
}

func CreateFolder(project *core.Project, path string, owner string) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	dir, name := filepath.Split(path)
	return setOwner(dir, name, owner)
}

func DeleteFile(project *core.Project, path string, recursive bool) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if recursive {
		return os.RemoveAll(path)
	} else {
		return os.Remove(path)
	}
	dir, name := filepath.Split(path)
	setExtendedAttr(dir, name, nil)
	return nil
}

// AbsPath returns the absolute path for a resource stored in the library.
func AbsPath(project *core.Project, path string) (string, error) {
	p := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	return filepath.Abs(p)
}

func ArchivePath(project *core.Project, path string) (string, error) {
	p := filepath.Join(project.Path, core.ProjectArchiveFolder, path)
	return filepath.Abs(p)
}


func SetFileInLibrary(project *core.Project, path string, reader io.ReadCloser, owner string) (string, error) {
	var err error
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)

	writer, err := os.Create(path)
	if err != nil {
		logrus.Warnf("Cannot open file %s for writing: %v", path, err)
		return "", err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		logrus.Warnf("cannot write file %s in library: %v", path, err)
		return "", err
	}
	logrus.Debugf("successfully set file %s in library", path)

	dir, name := filepath.Split(path)
	if err := setOwner(dir, name, owner); err != nil {
		logrus.Warnf("Cannot set owner for file %s: %v", path, owner)
	}
	logrus.Debugf("set owner of file %s tp %s", path, owner)
	return path, err
}

// AbsPath returns the absolute path for a resource stored in the library.
func GetItems(project *core.Project, paths []string) ([]Item, error) {
	items := make([]Item, 0, len(paths))

	for _, path := range paths {
		fullPath := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		mime, _ := mimetype.DetectFile(path)
		items = append(items, Item{
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
