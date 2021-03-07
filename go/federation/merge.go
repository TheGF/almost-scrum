package federation

import (
	"almost-scrum/attributes"
	"almost-scrum/core"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/alexmullins/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func verifyFile(project *core.Project, filename string) error {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	var header Header
	var content []byte
	for _, f := range r.File {
		if f.Name != FedHeaderFile {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		content, err = ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		break
	}

	if content == nil {
		return ErrFedCorrupted
	}
	if err := json.Unmarshal(content, &header); err != nil {
		return ErrFedCorrupted
	}

	if header.ProjectID != project.Config.UUID {
		return ErrFedCorrupted
	}
	return nil
}

func mergeItem(project *core.Project, source *zip.File) error {
	base := project.Path

	source.SetPassword(project.Config.CipherKey)
	name := source.Name

	dest := filepath.Join(base, name)
	stat, err := os.Stat(dest)
	if err == nil {
		if !stat.ModTime().Before(source.ModTime()) {
			return nil
		}

		hash, _ := hex.DecodeString(source.Comment)
		attr, err := attributes.GetExtendedAttr(filepath.Dir(dest), filepath.Base(dest))
		if err == nil && bytes.Compare(attr.Hash, hash) == 0 {
			return nil
		}
	}

	w, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	r, err := source.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return err
}

func MergeFile(project *core.Project, filename string) error {
	if err := verifyFile(project, filename); err != nil {
		os.Remove(filename)
		return err
	}

	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == FedHeaderFile {
			continue
		}
		mergeItem(project, f)
	}

	return nil
}

func MergeFiles(project *core.Project) error {
	path := filepath.Join(project.Path, core.ProjectFedFolder, "in")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		_ = MergeFile(project, filepath.Join(path, file.Name()))
	}
	return nil
}
