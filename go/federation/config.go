package federation

import (
	"almost-scrum/core"
	"os"
	"path/filepath"
	"time"
)

const FederationConfigFile = "fed.yaml"

type Config struct {
	LastImport time.Time
	LastExport time.Time
	Secret string
	Ftp []FTPConfig `yaml:"ftp"`
}

func ReadConfig(project *core.Project) (*Config, error) {
	var config Config
	path := filepath.Join(project.Path, core.ProjectFedFolder, FederationConfigFile)
	if err := core.ReadYaml(path, &config); err == nil {
		return &config, nil
	} else if os.IsNotExist(err) {
		return &Config{
			LastImport: time.Unix(0,0),
			Secret: core.GenerateRandomString(64),
		}, nil
	} else {
		return nil, err
	}
}

func WriteConfig(project *core.Project, config *Config) error {
	path := filepath.Join(project.Path, core.ProjectFedFolder, FederationConfigFile)
	return core.WriteYaml(path, config)
}
