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

type UpdateState string

const (
	New      UpdateState = "new"
	Newer                = "newer"
	Older                = "older"
	Conflict             = "conflict"
)

type MergeStrategy string

const (
	Extract MergeStrategy = "extract"
	Ignore                = "ignore"
)

type MergeItem struct {
	Match    UpdateState   `json:"match"`
	Strategy MergeStrategy `json:"strategy"`
}

type Diff struct {
	Input        string               `json:"input"`
	Name         string               `json:"name"`
	Header       Header               `json:"header"`
	CreationTime time.Time            `json:"creationTime"`
	Items        map[string]MergeItem `json:"items"`
}

type Update struct {
	Loc     string      `json:"loc"`
	State   UpdateState `json:"state"`
	Source  string      `json:"source"`
	Owner   string      `json:"owner"`
	ModTime time.Time   `json:"modTime"`
}

type Status struct {
	Exchanges  map[string]bool        `json:"exchanges"`
	Updates    []Update               `json:"updates"`
	Throughput map[string]*Throughput `json:"throughput"`
}

func parseComment(comment string) (owner string, origin []byte, hash []byte) {
	parts := strings.Split(comment, ",")
	switch len(parts) {
	case 0:
		return "", nil, nil
	case 1:
		return parts[0], nil, nil
	case 2:
		origin, _ := hex.DecodeString(parts[1])
		return parts[0], origin, nil
	default:
		origin, _ := hex.DecodeString(parts[1])
		hash, _ := hex.DecodeString(parts[2])
		return parts[0], origin, hash
	}
}

func readHeader(source *zip.File) (Header, error) {
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

	return header, nil
}

//func matchFile(project *core.Project, source *zip.File) (match UpdateState, owner string, err error) {
//	dest := filepath.Join(project.Path, source.Name)
//	owner, origin, hash := parseComment(source.Comment)
//
//	stat, err := os.Stat(dest)
//	if err == nil {
//		if !stat.ModTime().Before(source.ModTime()) {
//			return Outdated, owner, nil
//		}
//
//		hash_, err := fs.GetHash(dest)
//		if err != nil {
//			return "", "", err
//		}
//
//		if bytes.Compare(hash_, hash) == 0 {
//			return Outdated, owner, nil
//		} else if bytes.Compare(hash_, origin) != 0 {
//			return Conflict, owner, nil
//		} else {
//			return Update, owner, nil
//		}
//	} else if os.IsNotExist(err) {
//		attr, err := fs.GetExtendedAttr(dest)
//
//		if err == nil && attr.Modified.After(source.ModTime()) {
//			return Outdated, "", err
//		} else {
//			return New, owner, nil
//		}
//	} else {
//		return "", "", err
//	}
//}

//func CreateDiff(project *core.Project, base string, file string) (*Diff, error) {
//	r, err := zip.OpenReader(file)
//	if err != nil {
//		return nil, err
//	}
//	defer r.Close()
//
//	input, _ := filepath.Rel(base, file)
//	diff := Diff{
//		Input:        input,
//		Header:       Header{},
//		CreationTime: time.ModTime{},
//		Items:        map[string]MergeItem{},
//	}
//	for _, f := range r.File {
//		if f.Name == HeaderFile {
//			if header, err := readHeader(f); err != nil {
//				return nil, err
//			} else {
//				diff.Header = header
//			}
//			continue
//		}
//		match, _, err := getFileInfo(project, f)
//		if err != nil {
//			return nil, err
//		}
//		var strategy MergeStrategy
//		switch match {
//		case Outdated, Conflict:
//			strategy = Ignore
//		case New, Update:
//			strategy = Extract
//		}
//
//		diff.Items[f.Name] = MergeItem{
//			Match:    match,
//			Strategy: strategy,
//		}
//	}
//
//	return &diff, nil
//}
//

func findUpdateByLoc(updates []Update, loc string) *Update {
	for idx, update := range updates {
		if update.Loc == loc {
			return &updates[idx]
		}
	}
	return nil
}

func getUpdateForZipItem(project *core.Project, source string, zipItem *zip.File, updates []Update) ([]Update, error) {
	loc := zipItem.Name
	modTime := zipItem.ModTime()
	owner, importHash, exportHash := parseComment(zipItem.Comment)

	update := findUpdateByLoc(updates, loc)
	if update != nil && update.ModTime.After(modTime) {
		return updates, nil
	}

	var state UpdateState
	dest := filepath.Join(project.Path, zipItem.Name)
	stat, err := os.Stat(dest)
	if err == nil {
		if !stat.ModTime().Before(zipItem.ModTime()) {
			state = Older
		} else {
			hash, err := fs.GetHash(dest)
			if err != nil {
				return nil, err
			}

			switch {
			case uint64(stat.Size()) == zipItem.UncompressedSize64 && bytes.Compare(hash, exportHash) == 0:
				return updates, nil
			case bytes.Compare(hash, importHash) != 0:
				state = Conflict
			default:
				state = Newer
			}
		}
	} else if os.IsNotExist(err) {
		attr, _ := fs.GetExtendedAttr(dest)

		switch {
		case attr == nil:
			state = New
		case attr.Modified.Before(modTime):
			state = Newer
		default:
			state = Older
		}
	} else {
		return nil, err
	}

	if update == nil {
		updates = append(updates, Update{Loc: loc})
		update = &updates[len(updates)-1]
	}
	update.Source = source
	update.ModTime = modTime
	update.Owner = owner
	update.State = state

	return updates, nil
}

func getUpdatesInZip(project *core.Project, base string, file string, updates []Update) ([]Update, error) {
	r, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	source, _ := filepath.Rel(base, file)
	for _, zipItem := range r.File {
		if zipItem.Name == HeaderFile {
			if _, err := readHeader(zipItem); err != nil {
				return nil, err
			}
			continue
		}
		updates, err = getUpdateForZipItem(project, source, zipItem, updates)
		if err != nil {
			return nil, err
		}
	}
	return updates, nil
}

func getUpdates(project *core.Project) ([]Update, error) {
	var updates []Update

	connection, err := Connect(project)
	if err != nil {
		return nil, err
	}

	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	err = filepath.Walk(connection.local, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, exportItem := range syncItems {
			if strings.HasPrefix(info.Name(), exportItem.prefix) {
				updates, err = getUpdatesInZip(project, connection.local, path, updates)
				if err != nil {
					logrus.Errorf("cannot analyze %s: %v", path, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updates, nil
}

func GetStatus(project *core.Project) Status {
	connection, err := Connect(project)
	if err != nil {
		return Status{}
	}

	updates, _ := getUpdates(project)

	s := Status{
		Exchanges:  map[string]bool{},
		Updates:    updates,
		Throughput: connection.throughput,
	}
	for exchange, connected := range connection.exchanges {
		s.Exchanges[exchange.Name()] = connected
	}

	return s
}
