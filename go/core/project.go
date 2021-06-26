package core

import (
	"almost-scrum/assets"
	"almost-scrum/fs"
	"fmt"
	"github.com/code-to-go/fed"
	uuid2 "github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type FileAttr struct {
	Owner  string `json:"owner"`
	Public bool   `json:"public_"`
}


// Project is the basic information about a scrum project.
type Project struct {
	Path   string
	Config ProjectConfig
	Models []Model
	//	EncryptionSeed string
	Index      *Index
	IndexMutex sync.Mutex
	TasksCount int
	Fed        fed.Connection
}

// LoadTheProjectConfig
func ReadProjectConfig(path string) (ProjectConfig, error) {
	var projectConfig ProjectConfig
	err := fs.ReadYaml(filepath.Join(path, ProjectConfigFile), &projectConfig)
	return projectConfig, err
}

func WriteProjectConfig(path string, config *ProjectConfig) error {
	return fs.WriteYaml(filepath.Join(path, ProjectConfigFile), config)
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
	return nil, os.ErrNotExist
}

func getFedConnection(path string, fedId string) (fed.Connection, error) {
	f, err := fed.Open(filepath.Join(path, "fed"), FileAttr{}, fed.WithFedID(fedId))
	if err != nil {
		return nil, err
	}

	f.Mount(filepath.Join(path, ProjectBoardsFolder), fed.TrackByDefault)
	f.Mount(filepath.Join(path, ProjectModelsFolder), fed.TrackByDefault)
	f.Mount(filepath.Join(path, ProjectChatFolder), fed.TrackByDefault)
	f.Mount(filepath.Join(path, ProjectLibraryFolder))
	return f, nil
}

// OpenProject checks if the given path contains a project and creates an instance of Project.
func OpenProject(path string) (*Project, error) {
	projectConfig, err := ReadProjectConfig(path)
	if err != nil {
		return nil, ErrNoFound
	}

	models, err := ReadModels(path)
	if err != nil {
		return nil, err
	}

	fedConnection, err := getFedConnection(path, projectConfig.UUID)
	if err != nil {
		return nil, err
	}

	project := &Project{
		Path:       path,
		Config:     projectConfig,
		TasksCount: 0,
		Models:     models,
		Fed:        fedConnection,
	}

	infos, err := ListTasks(project, "", "")
	if err != nil {
		return nil, err
	}
	project.TasksCount = len(infos)

	AddProjectRefToConfig(project)

	logrus.Debugf("FindProject - Project found in %s", path)
	return project, nil
}

func GetGitClient(project *Project) GitClient {
	if project.Config.Public.UseGitNative {
		return Native{}
	} else {
		return GoGit{}
	}
}

func GetGitClientFromGlobalConfig() GitClient {
	config := ReadConfig()
	if config.UseGitNative {
		return Native{}
	} else {
		return GoGit{}
	}
}

func ListProjectTemplates() []string {
	var templates []string

	for _, name := range assets.AssetNames() {
		if !strings.HasPrefix(name, ProjectTemplatesPath) || !strings.HasSuffix(name, ".zip") {
			continue
		}
		templates = append(templates, name[len(ProjectTemplatesPath):len(name)-len(".zip")])
	}
	return templates
}

func createRequiredFolders(path string) error {
	for _, folder := range ProjectFolders {
		folder = filepath.Join(path, folder)
		if err := os.MkdirAll(folder, 0755); err != nil {
			logrus.Errorf("InitProject - Cannot create folder %s", folder)
			return err
		}
	}
	return nil
}

func initConfig(folder string) error {
	config := ReadConfig()

	projectConfig, err := ReadProjectConfig(folder)
	if os.IsNotExist(err) {
		projectConfig = ProjectConfig{
			Public: ProjectConfigPublic{
				CurrentBoard:    "",
				BoardTypes:      make(map[string][]string),
				IncludeLibInGit: true,
				UseGitNative:    config.UseGitNative,
			},
		}
	} else if err != nil {
		return err
	}

	projectConfig.CipherKey = GenerateRandomString(64)
	projectConfig.UUID = uuid2.New().String()

	parent, name := filepath.Split(folder)
	if name == ProjectFolder {
		projectConfig.Public.Name = parent
	} else {
		projectConfig.Public.Name = name
	}

	if err := WriteProjectConfig(folder, &projectConfig); err != nil {
		logrus.Errorf("InitProject - Cannot create config file in %s", folder)
		return err
	}

	config.Projects = append(config.Projects, ProjectRef{
		UUID:   projectConfig.UUID,
		Name:   projectConfig.Public.Name,
		Folder: folder,
	})

	logrus.Infof("created project %s (id %s)", projectConfig.Public.Name, projectConfig.UUID)
	return nil
}

func UnzipProjectTemplates(path string, templates []string) error {
	for _, template := range templates {
		var err error
		var templateData []byte
		if strings.HasPrefix(template, "file:/") {
			file, err := os.Open(template[6:])
			if err != nil {
				return err
			}
			templateData, err = ioutil.ReadAll(file)
			file.Close()
		} else {
			template = fmt.Sprintf("%s%s.zip", ProjectTemplatesPath, template)
			templateData, err = assets.Asset(template)
		}
		if err != nil {
			logrus.Errorf("Cannot open template %s: %v", template, err)
			return err
		}

		if err := UnzipFile(templateData, path); err != nil {
			logrus.Errorf("Cannot unzip template %s: %v", template, err)
			return err
		}
	}
	return nil
}

func InitProject(path string, templates []string) (*Project, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := createRequiredFolders(path); err != nil {
		return nil, err
	}

	if _, err := os.Stat(filepath.Join(path, ProjectConfigFile)); !os.IsNotExist(err) {
		logrus.Warnf("InitProject - Cannot initialize project. Project %s already exists", path)
		return &Project{Path: path}, ErrExists
	}

	if err := UnzipProjectTemplates(path, templates); err != nil {
		return nil, err
	}

	if err := initConfig(path); err != nil {
		return nil, err
	}

	logrus.Infof("Successfully created project in path %s with templates %v", path, templates)
	return OpenProject(path)
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
	DeleteProjFromConfig(project.Config.UUID)

	return nil
}

func JoinFed(project *Project, key string, token string) error{
	c, err := project.Fed.Join(key, token)
	project.Fed = c
	return err
}
