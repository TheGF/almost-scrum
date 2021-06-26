package core

import (
	"almost-scrum/fs"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ListBoards returns the boards in the project
func ListBoards(project *Project) ([]string, error) {
	p := filepath.Join(project.Path, "boards")
	infos, err := ioutil.ReadDir(p)
	if err != nil {
		logrus.Warnf("Cannot list store folder: %v", err)
		return nil, err
	}

	stores := make([]string, 0, len(infos))
	for _, info := range infos {
		stores = append(stores, info.Name())
	}
	return stores, nil
}

// CreateBoard creates a new store inside a project
func CreateBoard(project *Project, name string) error {
	p := filepath.Join(project.Path, "boards", name)
	return os.MkdirAll(p, 0777)
}

// DeleteBoard deletes an empty board
func DeleteBoard(project *Project, name string) error {
	p := filepath.Join(project.Path, "boards", name)
	return os.Remove(p)
}

// RenameBoard renames a board
func RenameBoard(project *Project, oldName string, newName string) error {
	p := filepath.Join(project.Path, "boards", oldName)
	np := filepath.Join(project.Path, "boards", newName)
	return os.Rename(p, np)
}

type BoardProperties struct {
	TaskTypes []string `json:"taskTypes"`
}

func GetBoardProperties(project *Project, name string) (BoardProperties, error) {
	var boardProperties BoardProperties = BoardProperties{TaskTypes: make([]string, 0)}


	p := filepath.Join(project.Path, "boards", name, ".board.yaml")
	fs.ReadYaml(p, &boardProperties)
	return boardProperties, nil
}

func SetBoardProperties(project *Project, name string, boardProperties BoardProperties) error {
	p := filepath.Join(project.Path, "boards", name, ".board.yaml")
	return fs.WriteYaml(p, &boardProperties)
}
