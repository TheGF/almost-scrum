package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Project is the basic information about a scrum project.
type Project struct {
	Path string
}

func findProjectInside(path string) (Project, error) {
	files, _ := filepath.Glob(filepath.Join(path, ProjectConfigFile))
	switch len(files) {
	case 0:
		log.Infof("No projects found in %s", path)
		return Project{}, ErrNoFound
	case 1:
		projectPath := filepath.Dir(files[0])
		log.Infof("Project found in %s", projectPath)
		return Project{Path: projectPath}, nil
	default:
		log.Infof("Multiple projects found in %s", path)
		return Project{}, ErrTooMany
	}
}

func findProjectOutside(path, root string) (Project, error) {
	path, _ = filepath.Abs(path)
	fileInfo, err := os.Stat(filepath.Join(path, ProjectConfigFile))
	if err == nil && !fileInfo.IsDir() {
		log.Infof("Project found in %s", path)
		return Project{Path: path}, nil
	}

	if parent := filepath.Dir(path); parent != root && parent != path {
		log.Infof("Check in %s", parent)
		return findProjectOutside(parent, root)
	}

	return Project{}, ErrNoFound

}

// FindProject searches for a project inside path and its parents up to root.
// Usually, root can be an empty string.
func FindProject(path, root string) (Project, error) {
	if project, err := findProjectInside(path); err == nil {
		return project, nil
	}

	return findProjectOutside(path, root)
}

// InitProject initializes a new project in the specified directory
func InitProject(path string) (Project, error) {
	configPath := filepath.Join(path, ProjectConfigFile)
	if _, err := os.Stat(configPath); err == nil {
		log.Infof("Cannot initialize project. Project %s already exists", path)
		return Project{}, ErrExists
	}

	if err := ioutil.WriteFile(configPath, []byte("version: 1.0"), 0644); err != nil {
		return Project{}, err
	}

	for _, folder := range ProjectFolders {
		if err := os.MkdirAll(filepath.Join(path, folder), 0755); err != nil {
			return Project{}, err
		}
	}

	return Project{path}, nil
}

// GetStore returns the specified name
func GetStore(project Project, store string) (Store, error) {
	path := filepath.Join(project.Path, "stores", store)
	if fileInfo, err := os.Stat(path); err != nil || !fileInfo.IsDir() {
		return Store{}, ErrNoFound
	}
	return Store{Path: path}, nil
}

// GetUsers returns the project users
func GetUsers(project Project) []string {
	path := filepath.Join(project.Path, ProjectUsersFolder)
	file, _ := os.Open(path)
	users, _ := file.Readdirnames(0)
	return users
}
