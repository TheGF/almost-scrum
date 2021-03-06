package library

import (
	"almost-scrum/core"
	"github.com/code-to-go/fed/extfs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
)

// Item contains basic information about items in the library.
type Item struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Mime    string    `json:"mime"`
	Dir     bool      `json:"dir"`
	Empty   bool      `json:"empty"`
	Parent  string    `json:"parent"`
	core.FileAttr
}

func getItem(path string, fileInfo os.FileInfo) Item {
	path = filepath.Join(path, fileInfo.Name())
	mime, _ := mimetype.DetectFile(path)

	var size int64
	if fileInfo.IsDir() {
		files, _ := ioutil.ReadDir(path)
		cnt := 0
		for _, file := range files {
			if !strings.HasPrefix(".", file.Name()) {
				cnt++
			}
		}
		size = int64(cnt)
	} else {
		size = fileInfo.Size()
	}

	var attr core.FileAttr
	extfs.Get(path, attr)
	return Item{
		Name:     fileInfo.Name(),
		Size:     size,
		ModTime:  fileInfo.ModTime(),
		Mime:     mime.String(),
		FileAttr: attr,
		Dir:      fileInfo.IsDir(),
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
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		logrus.Warnf("List - Cannot files content of folder %s: %v", path, err)
	}

	versions := make(map[string]versionInfo)
	for _, file := range files {
		filterVersioned(project, path, file, versions)
	}

	items := make([]Item, 0, len(files))
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
	err = os.Rename(oldPath, path)
	if err != nil {
		return err
	}
	extfs.Move(oldPath, path)
	return nil
}

func CreateFolder(project *core.Project, path string, owner string) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return extfs.Set(path, &core.FileAttr{Owner: owner}, true)
}

func DeleteFile(project *core.Project, path string) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	err := os.Remove(path)
	if err != nil {
		return err
	}
	//attr, _ := fs.GetExtendedAttr(path)
	//attr.Modified = time.Now()
	//_ = fs.SetExtendedAttr(path, attr)
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

func SetFileInLibrary(project *core.Project, path string, reader io.ReadCloser,
	owner string, public bool) (string, error) {
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

	if err := project.Fed.SetTracked(path, public); err != nil {
		logrus.Warnf("Cannot set tracker for file %s: %v", path, err)
	}

	if err := extfs.Set(path, core.FileAttr{
		Owner:  owner,
		Public: public,
	}, true); err != nil {
		logrus.Warnf("Cannot set owner for file %s: %v", path, owner)
	}
	logrus.Debugf("set owner of file %s tp %s", path, owner)
	return path, err
}

func SetVisibility(project *core.Project, path string, public bool) error {
	fullPath := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	files, err := ioutil.ReadDir(fullPath)
	if err == nil {
		for _, file := range files {
			if err := SetVisibility(project, filepath.Join(path, file.Name()), public); err != nil {
				return err
			}
		}
	}

	err = project.Fed.SetTracked(fullPath, public)
	logrus.Infof("set visibility of file %s tp %t", path, public)
	return err
}

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
