package fed

import (
	"almost-scrum/core"
	"almost-scrum/fs"
	"github.com/alexmullins/zip"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)


func importZipItem(project *core.Project, source *zip.File) (string, error) {
	owner, _, exportHash := parseComment(source.Comment)

	dest := filepath.Join(project.Path, source.Name)
	_ = os.MkdirAll(filepath.Dir(dest), 0755)

	source.SetPassword(project.Config.CipherKey)
	r, err := source.Open()
	if err != nil {
		logrus.Errorf("cannot open zip item %s: %v", source.Name, err)
		return "", err
	}
	defer r.Close()

	w, err := os.Create(dest)
	if err != nil {
		logrus.Errorf("cannot create %s: %v", dest, err)
		return "", err
	}

	_, err = io.Copy(w, r)
	w.Close()
	if err != nil {
		logrus.Errorf("cannot write %s: %v", dest, err)
		return "", err
	}

	_ = fs.SetExtendedAttr(dest, &fs.ExtendedAttr{
		Owner:      owner,
		ImportHash: exportHash,
	})

	logrus.Debugf("import file %s: owner %s, hash %v", source.Name, owner, exportHash)
	return source.Name, err
}


func importFromSource(project *core.Project, source string, updates []Update) ([]string, error) {
	file := filepath.Join(project.Path, core.ProjectFedFilesFolder, source)
	r, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var locs []string
	for _, zipItem := range r.File {
		loc := zipItem.Name
		update := findUpdateByLoc(updates, loc)
		if update != nil && update.Source == source {
			if name, err := importZipItem(project, zipItem); err != nil {
				return nil, err
			} else {
				locs = append(locs, name)
			}
		}
	}
	return locs, nil
}

func Import(project *core.Project, updates []Update) ([]string, error) {
	var locs []string
	sources := map[string]bool{}

	for _, update := range updates {
		sources[update.Source] = true
	}

	for source := range sources {
		if names_, err := importFromSource(project, source, updates); err != nil {
			return nil, err
		} else {
			locs = append(locs, names_...)
		}
	}
	return locs, nil
}
