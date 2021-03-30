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


func extractItem(project *core.Project, source *zip.File) (string, error) {
	owner, _ := parseComment(source.Comment)

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

	origin, _ := fs.GetHash(dest)
	_ = fs.SetExtendedAttr(dest, &fs.ExtendedAttr{
		Owner:  owner,
		Origin: origin,
	})

	return source.Name, err
}


func ImportDiff(project *core.Project, diff *Diff) ([]string, error) {
	file := filepath.Join(project.Path, core.ProjectFedFilesFolder, diff.Input)
	r, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var names []string
	for _, f := range r.File {
		if item, found := diff.Items[f.Name]; found && item.Strategy == Extract {
			if name, err := extractItem(project, f); err != nil {
				return nil, err
			} else {
				names = append(names, name)
			}
		}
	}
	return names, nil
}

func Import(project *core.Project, diffs []*Diff) ([]string, error) {
	var names []string
	for _, diff := range diffs {
		if names_, err := ImportDiff(project, diff); err != nil {
			return nil, err
		} else {
			names = append(names, names_...)
		}
	}
	return names, nil
}
