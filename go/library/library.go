package library

import (
	"almost-scrum/core"
	"almost-scrum/fs"
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
	Empty   bool      `json:"empty"`
	Public  bool      `json:"public_"`
	Owner   string    `json:"owner"`
	Parent  string    `json:"parent"`
}

func getItem(path string, fileInfo os.FileInfo) Item {
	path = filepath.Join(path, fileInfo.Name())
	mime, _ := mimetype.DetectFile(path)

	var size int64
	if fileInfo.IsDir() {
		files, _ := ioutil.ReadDir(path)
		size = int64(len(files))
	} else {
		size = fileInfo.Size()
	}

	extendedAttr, err := fs.GetExtendedAttr(path)
	if err != nil {
		logrus.Warnf("Cannot get extended attr for %s: %v", path, err)
		return Item{}
	}
	return Item{
		Name:    fileInfo.Name(),
		Size:    size,
		ModTime: fileInfo.ModTime(),
		Mime:    mime.String(),
		Public:  extendedAttr.Public,
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
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		logrus.Warnf("List - Cannot files content of folder %s: %v", path, err)
	}

	versions := make(map[string]versionInfo)
	for _, file := range files {
		if file.Name() != fs.AttrsFileName {
			filterVersioned(project, path, file, versions)
		}
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
	xAttr, err := fs.GetExtendedAttr(oldPath)
	if err != nil {
		return err
	}

	err = os.Rename(oldPath, path)
	if err != nil {
		return err
	}
	_ = fs.SetExtendedAttr(oldPath, xAttr)
	_ = fs.SetExtendedAttr(path, nil)
	return nil
}

func CreateFolder(project *core.Project, path string, owner string) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return fs.SetExtendedAttr(path, &fs.ExtendedAttr{
		Owner: owner,
	})
}

func DeleteFile(project *core.Project, path string) error {
	path = filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	//if recursive {
	//	_ = fs.SetExtendedAttr(path, nil)
	//	return os.RemoveAll(path)
	//} else {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	attr, _ := fs.GetExtendedAttr(path)
	attr.Deleted = time.Now()
	_ = fs.SetExtendedAttr(path, attr)
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

	if err := fs.SetExtendedAttr(path, &fs.ExtendedAttr{
		Owner:   owner,
		Origin:  nil,
		Public:  public,
		Deleted: time.Time{},
	}); err != nil {
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
	attr, _ := fs.GetExtendedAttr(fullPath)
	if attr.Public != public {
		attr.Public = public
		if err = fs.SetExtendedAttr(fullPath, attr); err == nil {
			tm := time.Now()
			_ = os.Chtimes(fullPath, tm, tm)
		}
		logrus.Infof("set visibility of file %s tp %t", path, public)
	}
	return nil
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
