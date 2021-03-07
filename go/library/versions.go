package library

import (
	"almost-scrum/attributes"
	"almost-scrum/core"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"

	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var versionsRegex = regexp.MustCompile(`(.*?)((\d+\.)+\d+)?(\.\w*)?$`)
func parsePath(path string) (dir string, prefix string, version string, ext string, err error) {
	dir = filepath.Dir(path)
	name := filepath.Base(path)
	match := versionsRegex.FindStringSubmatch(name)
	if len(match) != 5 {
		err = os.ErrInvalid
		logrus.Errorf("Cannot parse %s: %v", path, err)
		return
	}

	prefix = match[1]
	version = match[2]
	ext = match[4]
	err = nil

	return
}


type versionInfo struct {
	fileInfo os.FileInfo
	version  string
}

func archiveFile(project *core.Project, path string, fileInfo os.FileInfo) error {
	name := fileInfo.Name()
	parent := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	source := filepath.Join(parent, name)
	sourceXAttr, err := attributes.GetExtendedAttr(parent, name)
	if err != nil {
		return err
	}

	parent_ := filepath.Join(project.Path, core.ProjectArchiveFolder, path)
	dest := filepath.Join(parent_, name)
	err = os.Rename(source, dest)
	if err != nil {
		logrus.Errorf("Cannot archive file %s: %v", source, err)
		return err
	}

	_ = attributes.SetExtendedAttr(parent_, name, sourceXAttr)
	_ = attributes.SetExtendedAttr(parent, name, nil)
	logrus.Debugf("Archived file %s/%s to %s", parent, name, parent_)
	return err
}

func filterVersioned(project *core.Project, parent string, fileInfo os.FileInfo, latest map[string]versionInfo) {
	name := fileInfo.Name()
	_, prefix, version_, ext, err := parsePath(name)
	if err != nil {
		return
	}

	if version_ == "" {
		latest[name] = versionInfo{
			fileInfo: fileInfo,
			version:  version_,
		}
		return
	}

	id := fmt.Sprintf("%s*%s", prefix, ext)
	if last, found := latest[id]; found {
		v1, _ := version.NewVersion(version_)
		v2, _ := version.NewVersion(last.version)

		if v1.LessThan(v2) {
			_ = archiveFile(project, parent, fileInfo)
		} else {
			_ = archiveFile(project, parent, last.fileInfo)
			latest[id] = versionInfo{
				fileInfo: fileInfo,
				version:  version_,
			}
		}
	} else {
		latest[id] = versionInfo{
			fileInfo: fileInfo,
			version:  version_,
		}
	}
}

func GetPreviousVersions(project *core.Project, path string) ([]Item, error) {
	parent, prefix, _, ext, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	absPath := filepath.Join(project.Path, core.ProjectArchiveFolder, parent)
	fileInfos, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)
	for _, fileInfo := range fileInfos {
		_, prefix_, _, ext_, err := parsePath(fileInfo.Name())

		logrus.Debugf("Matching %s...%s with %s", prefix_, ext_, absPath)
		if err == nil && prefix == prefix_ && ext == ext_ {
			items = append(items, getItem(absPath, fileInfo))
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name > items[i].Name
	})

	return items, nil
}
