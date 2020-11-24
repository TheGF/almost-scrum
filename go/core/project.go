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
	path, _ = filepath.Abs(path)
	foundPaths := make([]string, 0, 1)

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && info.Name() == ProjectConfigFile {
			parent := filepath.Dir(path)
			log.Debugf("FindProject - Project found in %s", parent)
			foundPaths = append(foundPaths, parent)
		}
		return nil
	})

	switch len(foundPaths) {
	case 0:
		log.Infof("No projects found in %s", path)
		return Project{}, ErrNoFound
	case 1:
		log.Infof("Project found in %s", foundPaths[0])
		return Project{Path: foundPaths[0]}, nil
	default:
		log.Infof("Multiple projects found in %s", path)
		return Project{}, ErrTooMany
	}
}

func findProjectOutside(path, root string) (Project, error) {
	path, _ = filepath.Abs(path)
	fileInfo, err := os.Stat(filepath.Join(path, ProjectConfigFile))
	if err == nil && !fileInfo.IsDir() {
		log.Debugf("FindProject - Project found in %s", path)
		return Project{Path: path}, nil
	}

	if parent := filepath.Dir(path); parent != root && parent != path {
		log.Debugf("FindProject - Check in %s", parent)
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
	path, err := filepath.Abs(path)
	if err != nil {
		return Project{}, err
	}

	configPath := filepath.Join(path, ProjectConfigFile)
	// Check that no project is already initialized
	if _, err := os.Stat(configPath); err == nil {
		log.Errorf("InitProject - Cannot initialize project. Project %s already exists", path)
		return Project{}, ErrExists
	}

	// Create required folders
	for _, folder := range ProjectFolders {
		folder = filepath.Join(path, folder)
		if err := os.MkdirAll(folder, 0755); err != nil {
			log.Errorf("InitProject - Cannot create folder %s", folder)
			return Project{}, err
		}
	}

	// Create the project configuration
	if err := ioutil.WriteFile(configPath, []byte("version: 1.0"), 0644); err != nil {
		log.Errorf("InitProject - Cannot create file %s", configPath)
		return Project{}, err
	}

	// Store a reference to the project in the global configuration
	gconfig := LoadConfig()
	gconfig.Projects[filepath.Base(path)] = path
	SaveConfig(gconfig)

	return Project{path}, nil
}

// ShredProject fully deletes all files in a project. Use with caution!
func ShredProject(project Project) error {
	files := append([]string{}, ProjectFolders...)
	files = append(files, ProjectConfigFile)

	projectPath, err := filepath.Abs(project.Path)
	if err != nil {
		log.Errorf("ShredProject - Cannot resolve path %s", project.Path)
	}
	for _, file := range files {
		path := filepath.Join(projectPath, file)
		log.Debugf("ShredProject - Going to remove %s", path)
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	// Remove a reference to the project from the global configuration
	gconfig := LoadConfig()
	delete(gconfig.Projects, filepath.Base(project.Path))
	SaveConfig(gconfig)

	return nil
}

// GetStore returns the specified store
func GetStore(project Project, store string) (Store, error) {
	path := filepath.Join(project.Path, "stores", store)
	if fileInfo, err := os.Stat(path); err != nil || !fileInfo.IsDir() {
		return Store{}, ErrNoFound
	}
	return Store{Path: path}, nil
}
