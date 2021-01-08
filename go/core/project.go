package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type ProjectConfig struct {
	CurrentBoard  string `yaml:"current_store"`
	PropertyModel []PropertyDef
}

type PropertyDef struct {
	Name    string   `json:"name" yaml:"name"`
	Kind    string   `json:"kind" yaml:"kind"`
	Values  []string `json:"values" yaml:"values"`
	Prefix  string   `json:"prefix" yaml:"prefix"`
	Default string   `json:"default" yaml:"default"`
}

// Project is the basic information about a scrum project.
type Project struct {
	Path   string
	Config ProjectConfig
	Index  *Index
}

// LoadTheProjectConfig
func ReadProjectConfig(path string) (ProjectConfig, error) {
	var projectConfig ProjectConfig
	err := ReadYaml(filepath.Join(path, ProjectConfigFile), &projectConfig)
	return projectConfig, err
}

func WriteProjectConfig(path string, config *ProjectConfig) error {
	return WriteYaml(filepath.Join(path, ProjectConfigFile), config)
}

// FindProject searches for a project inside path and its parents up to root.
// Usually, root can be an empty string.
func FindProject(path string) (*Project, error) {
	p, _ := FindFileUpwards(path, GitFolder)
	if p != "" {
		_, err := os.Stat(filepath.Join(p, ProjectFolder))
		if err == nil {
			return OpenProject(filepath.Join(p, ProjectFolder))
		}
	}

	p, _ = FindFileUpwards(path, ProjectConfigFile)
	if p != "" {
		return OpenProject(p)
	}
	return nil, ErrNoFound
}

// OpenProject checks if the given path contains a project and creates an instance of Project.
func OpenProject(path string) (*Project, error) {
	projectConfig, err := ReadProjectConfig(path)
	if err != nil {
		return nil, ErrNoFound
	}

	logrus.Debugf("FindProject - Project found in %s", path)
	return &Project{
		Path:   path,
		Config: projectConfig,
	}, nil
}

// InitProject initializes a new project in the specified directory
func InitProject(path string) (*Project, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if _, err := ReadProjectConfig(path); err == nil {
		logrus.Warnf("InitProject - Cannot initialize project. Project %s already exists", path)
		return &Project{Path: path}, ErrExists
	}

	// Create required folders
	for _, folder := range ProjectFolders {
		folder = filepath.Join(path, folder)
		if err := os.MkdirAll(folder, 0755); err != nil {
			logrus.Errorf("InitProject - Cannot create folder %s", folder)
			return nil, err
		}
	}

	// Create the project configuration
	projectConfig := ProjectConfig{
		CurrentBoard: "backlog",
		PropertyModel: []PropertyDef{
			{"Owner", "User", nil, "", ""},
			{"Status", "Tag", []string{"#Draft", "#Started", "#Done"},
				"", "#Draft"},
			{"Points", "Enum", []string{"1", "2", "3", "5", "7", "9", "12", "15", "21"},
				"", "3"},
		},
	}
	if err := WriteProjectConfig(path, &projectConfig); err != nil {
		logrus.Errorf("InitProject - Cannot create config file in %s", path)
		return nil, err
	}

	// Board a reference to the project in the global configuration
	globalConfig := LoadConfig()
	globalConfig.Projects[filepath.Base(path)] = path
	SaveConfig(globalConfig)

	return &Project{
		Path:   path,
		Config: projectConfig,
	}, nil
}

//NewTaskName browses all stories in all boards and returns the next possible id.
func NewTaskName(project *Project, title string) string {
	path := filepath.Join(project.Path, "boards")
	id := 1

	_ = filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}
		name := fileInfo.Name()
		firstDot := strings.IndexByte(name, '.')
		if firstDot > 0 {
			fileID, _ := strconv.Atoi(name[:firstDot])
			if id <= fileID {
				id = fileID + 1
			}
		}
		return nil
	})
	return fmt.Sprintf("%d.%s", id, title)
}

// ShredProject fully deletes all files in a project. Use with caution!
func ShredProject(project *Project) error {
	files := append([]string{}, ProjectFolders...)
	files = append(files, ProjectConfigFile)

	projectPath, err := filepath.Abs(project.Path)
	if err != nil {
		logrus.Errorf("ShredProject - Cannot resolve path %s", project.Path)
	}
	for _, file := range files {
		path := filepath.Join(projectPath, file)
		logrus.Debugf("ShredProject - Going to remove %s", path)
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	// Remove a reference to the project from the global configuration
	globalConfig := LoadConfig()
	delete(globalConfig.Projects, filepath.Base(project.Path))
	SaveConfig(globalConfig)

	return nil
}

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
