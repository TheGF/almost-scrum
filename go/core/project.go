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

// Project is the basic information about a scrum project.
type Project struct {
	Path         string
	CurrentStore string
	Users        []string
}

type ProjectConfig struct {
	CurrentStore string   `yaml:"currentStore"`
	Users        []string `yaml:"users"`
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
func FindProject(path string) (Project, error) {
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
	return Project{}, ErrNoFound
}

// OpenProject checks if the given path contains a project and creates an instance of Project.
func OpenProject(path string) (Project, error) {
	projectConfig, err := ReadProjectConfig(path)
	if err != nil {
		return Project{}, ErrNoFound
	}

	logrus.Debugf("FindProject - Project found in %s", path)
	return Project{
		Path:         path,
		CurrentStore: projectConfig.CurrentStore,
		Users:        projectConfig.Users,
	}, nil
}

// InitProject initializes a new project in the specified directory
func InitProject(path string) (Project, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return Project{}, err
	}

	if _, err := ReadProjectConfig(path); err == nil {
		logrus.Warnf("InitProject - Cannot initialize project. Project %s already exists", path)
		return Project{Path: path}, ErrExists
	}

	// Create required folders
	for _, folder := range ProjectFolders {
		folder = filepath.Join(path, folder)
		if err := os.MkdirAll(folder, 0755); err != nil {
			logrus.Errorf("InitProject - Cannot create folder %s", folder)
			return Project{}, err
		}
	}

	// Create the project configuration
	projectConfig := ProjectConfig{
		CurrentStore: "backlog",
		Users:        []string{GetCurrentUser()},
	}
	if err := WriteProjectConfig(path, &projectConfig); err != nil {
		logrus.Errorf("InitProject - Cannot create config file in %s", path)
		return Project{}, err
	}

	// Store a reference to the project in the global configuration
	globalConfig := LoadConfig()
	globalConfig.Projects[filepath.Base(path)] = path
	SaveConfig(globalConfig)

	return Project{
		Path:         path,
		CurrentStore: projectConfig.CurrentStore,
		Users:        projectConfig.Users,
	}, nil
}

//GetStoryName browses all stories in all stores and returns the next possible id.
func GetStoryName(project Project, title string) string {
	path := filepath.Join(project.Path, "stores")
	id := 1

	filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
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
	return fmt.Sprintf("%d.%s.story", id, title)
}

// ShredProject fully deletes all files in a project. Use with caution!
func ShredProject(project Project) error {
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

// ListStores returns the stores in the project
func ListStores(project Project) ([]string, error) {
	storesPath := filepath.Join(project.Path, "stores")
	infos, err := ioutil.ReadDir(storesPath)
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
