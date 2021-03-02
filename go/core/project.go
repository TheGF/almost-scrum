package core

import (
	"almost-scrum/assets"
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type ProjectConfigPublic struct {
	CurrentBoard    string              `json:"currentStore" yaml:"currentStore"`
	BoardTypes      map[string][]string `json:"boardTypes" yaml:"boardTypes"`
	IncludeLibInGit bool                `json:"includeLibInGit" yaml:"includeLibInGit"`
	UseGitNative    bool                `json:"useGitNative" yaml:"useGitNative"`
}

type ProjectConfig struct {
	CipherKey string              `yaml:"cipherKey"`
	UUID      string              `yaml:"uuid"`
	Public    ProjectConfigPublic `yaml:"public"`
}

type PropertyKind string

const (
	KindString PropertyKind = "String"
	KindEnum   PropertyKind = "Enum"
	KindBool   PropertyKind = "Bool"
	KindUser   PropertyKind = "User"
	KindTag    PropertyKind = "Tag"
)

// Project is the basic information about a scrum project.
type Project struct {
	Path   string
	Config ProjectConfig
	Models []Model
	//	EncryptionSeed string
	Index      *Index
	TasksCount int
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
	return nil, os.ErrNotExist
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

	project := &Project{
		Path:       path,
		Config:     projectConfig,
		TasksCount: 0,
		Models:     models,
	}

	infos, err := ListTasks(project, "", "")
	if err != nil {
		return nil, err
	}
	project.TasksCount = len(infos)

	logrus.Debugf("FindProject - Project found in %s", path)
	return project, nil
}

func GetGitClient(project *Project) GitClient {
	if project.Config.Public.UseGitNative {
		return GitNative{}
	} else {
		return GoGit{}
	}
}

func GetGitClientFromGlobalConfig() GitClient {
	config := ReadConfig()
	if config.UseGitNative {
		return GitNative{}
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

func initConfig(path string) error {
	config := ReadConfig()

	projectConfig, err := ReadProjectConfig(path)
	if os.IsExist(err) {
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
	if err := WriteProjectConfig(path, &projectConfig); err != nil {
		logrus.Errorf("InitProject - Cannot create config file in %s", path)
		return err
	}
	return nil
}

func unzipTemplates(path string, templates []string) error {
	for _, template := range templates {
		template = fmt.Sprintf("%s%s.zip", ProjectTemplatesPath, template)
		templateData, err := assets.Asset(template)
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

	if err := unzipTemplates(path, templates); err != nil {
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
	globalConfig := ReadConfig()
	delete(globalConfig.Projects, filepath.Base(project.Path))
	WriteConfig(globalConfig)

	return nil
}

func NameProject(project *Project, name string) {
	config := ReadConfig()

	path, found := config.Projects[name]
	if !found || path != project.Path {
		config.Projects[name] = project.Path
		WriteConfig(config)
	}
}

func EncryptStringForProject(project *Project, value string) (string, error) {
	c, err := aes.NewCipher([]byte(project.Config.CipherKey))
	if err != nil {
		return "", err
	}
	out := ""
	buf := make([]byte, c.BlockSize())
	for l := 0; l < len(value); l += len(buf) {
		end := l + len(buf)
		if end > len(value) {
			end = len(value)
			buf[end-l] = 0
		}
		copy(buf, value[l:end])
		c.Encrypt(buf, buf)
		out += hex.EncodeToString(buf)
	}

	return out, nil
}

func DecryptStringForProject(project *Project, value string) (string, error) {
	c, err := aes.NewCipher([]byte(project.Config.CipherKey))
	if err != nil {
		return "", err
	}

	ciphertext, _ := hex.DecodeString(value)
	blockSize := c.BlockSize()
	for l := 0; l < len(ciphertext); l += blockSize {
		c.Decrypt(ciphertext[l:l+blockSize], ciphertext[l:l+blockSize])
	}

	idx := bytes.IndexByte(ciphertext, 0)
	return string(ciphertext[0:idx]), nil
}
