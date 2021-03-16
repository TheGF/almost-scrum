package fed

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/alexmullins/zip"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)


type MergeMatch string

const (
	New      MergeMatch = "new"
	Outdated            = "outdated"
	Update              = "update"
	Conflict            = "conflict"
)

type MergeStrategy string

const (
	Extract  MergeStrategy = "extract"
	Ignore                 = "ignore"
)

type MergeItem struct {
	Match    MergeMatch    `json:"match"`
	Strategy MergeStrategy `json:"strategy"`
}

type Diff struct {
	Input        string               `json:"input"`
	Name         string               `json:"name"`
	Header       Header               `json:"header"`
	CreationTime time.Time            `json:"creationTime"`
	Items        map[string]MergeItem `json:"items"`
}

func parseComment(comment string) (owner string, origin []byte) {
	parts := strings.Split(comment, ",")
	switch len(parts) {
	case 0:
		return "", nil
	case 1:
		return parts[0], nil
	default:
		origin, _ := hex.DecodeString(parts[1])
		return parts[0], origin
	}
}


func readHeader(project *core.Project, source *zip.File) (Header, error) {
	var header Header
	var content []byte
	rc, err := source.Open()
	if err != nil {
		return header, err
	}
	defer rc.Close()

	content, err = ioutil.ReadAll(rc)
	if err != nil {
		return header, err
	}
	if content == nil {
		return header, ErrFedCorrupted
	}
	if err := json.Unmarshal(content, &header); err != nil {
		return header, ErrFedCorrupted
	}

	if header.ProjectID != project.Config.UUID {
		return header, ErrFedCorrupted
	}
	return header, nil
}

func matchFile(project *core.Project, source *zip.File) (match MergeMatch, owner string, err error) {
	dest := filepath.Join(project.Path, source.Name)
	owner, origin := parseComment(source.Comment)

	stat, err := os.Stat(dest)
	if err == nil {
		if !stat.ModTime().Before(source.ModTime()) {
			return Outdated, owner, nil
		}

		hash, err := fs.GetHash(dest)
		if err != nil {
			return "", "", err
		}

		if bytes.Compare(hash, origin) != 0 {
			return Conflict, owner, nil
		} else {
			return Update, owner, nil
		}
	} else if os.IsNotExist(err) {
		attr, err := fs.GetExtendedAttr(dest)

		if err == nil && attr.Deleted.After(source.ModTime()) {
			return Outdated, "", err
		} else {
			return New, owner, nil
		}
	} else {
		return "", "", err
	}
}

func CreateDiff(project *core.Project, file string) (*Diff, error) {
	r, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	log := Diff{
		Input:        filepath.Base(file),
		Header:       Header{},
		CreationTime: time.Time{},
		Items:        map[string]MergeItem{},
	}
	for _, f := range r.File {
		if f.Name == HeaderFile {
			if header, err := readHeader(project, f); err != nil {
				return nil, err
			} else {
				log.Header = header
			}
			continue
		}
		match, _, err := matchFile(project, f)
		if err != nil {
			return nil, err
		}
		var strategy MergeStrategy
		switch match {
		case Outdated, Conflict:
			strategy = Ignore
		case New, Update:
			strategy = Extract
		}

		log.Items[f.Name] = MergeItem{
			Match:    match,
			Strategy: strategy,
		}
	}

	return &log, nil
}

func GetDiff(project *core.Project) ([]*Diff, error) {
	var diffs []*Diff

	state, err := GetSignal(project)
	if err != nil {
		return nil, err
	}
	state.inUse.Add(1)
	defer state.inUse.Done()

	err = filepath.Walk(state.local, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), "exp-") {
			log, err := CreateDiff(project, path)
			if err == nil {
				diffs = append(diffs, log)
			} else {
				logrus.Warnf("cannot create log for %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}


	return diffs, nil
}

