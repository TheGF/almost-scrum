package library

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"fmt"
	"github.com/code-to-go/fed/extfs"
	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type versionInfo struct {
	fileInfo os.FileInfo
	version  string
}

func archiveFile(project *core.Project, path string) (string, error) {
	source := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	dest := filepath.Join(project.Path, core.ProjectArchiveFolder, path)
	if err := os.Rename(source, dest); err != nil {
		logrus.Errorf("Cannot archive file %s: %v", source, err)
		return "", err
	}
	extfs.Move(source, dest)

	logrus.Debugf("Archived file %s to %s", source, dest)
	return dest, nil
}

func filterVersioned(project *core.Project, parent string, fileInfo os.FileInfo, latest map[string]versionInfo) {
	name := fileInfo.Name()
	_, prefix, version_, ext, err := fs.ParsePath(name)
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
		v1, _ := version.NewVersion(version_[1:])
		v2, _ := version.NewVersion(last.version[1:])

		if v1.LessThan(v2) {
			_, _ = archiveFile(project, filepath.Join(parent, fileInfo.Name()))
		} else {
			_, _ = archiveFile(project, filepath.Join(parent, last.fileInfo.Name()))
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
	parent, prefix, _, ext, err := fs.ParsePath(path)
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
		_, prefix_, _, ext_, err := fs.ParsePath(fileInfo.Name())

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

var versionRegex = regexp.MustCompile(`~(\d+\.)+(\d+)`)

func getNextVersion(version string, subversion bool) (string, error) {
	if version == "" {
		return "~0.1", nil
	}
	if subversion {
		return fmt.Sprintf("%s.1", version), nil
	}

	match := versionRegex.FindStringSubmatch(version)
	if len(match) != 3 {
		return "", os.ErrInvalid
	}
	lower, _ := strconv.Atoi(match[2])
	return fmt.Sprintf("~%s%d", match[1], lower+1), nil
}

func IncreaseVersion(project *core.Project, path string, owner string, public bool) (string, error) {
	dir, prefix, ver, ext, err := fs.ParsePath(path)
	if err != nil {
		return "", err
	}

	ver_, _ := getNextVersion(ver, false)
	path_ := filepath.Join(dir, fmt.Sprintf("%s%s%s", prefix, ver_, ext))

	fullPath := filepath.Join(project.Path, core.ProjectLibraryFolder, path)
	r, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	fullPath_ := filepath.Join(project.Path, core.ProjectLibraryFolder, path_)
	w, err := os.Create(fullPath_)
	if err != nil {
		return "", err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return "", err
	}

	_, _ = archiveFile(project, path)
	_ = extfs.Set(fullPath_, core.FileAttr{
		Owner:  owner,
		Public: public,
	}, true)
	_ = project.Fed.SetTracked(fullPath_, public)

	archivePath := filepath.Join(project.Path, core.ProjectArchiveFolder, path)
	_ = os.MkdirAll(filepath.Dir(archivePath), 0755)
	_ = os.Rename(fullPath, archivePath)

	return path_, nil
}
