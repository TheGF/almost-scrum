package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const configFile = "fed.yaml"

type Config struct {
	Secret        string                   `yaml:"secret"`
	ReconnectTime time.Duration            `yaml:"reconnectTime"`
	PollTime      time.Duration            `yaml:"pollTime"`
	Span          int                      `yaml:"span"`
	LastExport    time.Time                `yaml:"lastExport"`
	Ftp           []transport.FTPConfig    `yaml:"ftp"`
	WebDAV        []transport.WebDAVConfig `yaml:"webdav"`
}

func ReadConfig(project *core.Project) (*Config, error) {
	var config Config
	path := filepath.Join(project.Path, core.ProjectFedFolder, configFile)
	if err := core.ReadYaml(path, &config); err == nil {
		if config.Span < 1 {
			logrus.Warnf("invalid span value %d in %s; default to 10", config.Span, path)
			config.Span = 10
		}

		if config.PollTime < 10*time.Second {
			logrus.Warnf("invalid pool time %s in %s; default to 10m", config.PollTime, path)
			config.PollTime = 10 * time.Minute
		}
		return &config, nil
	} else if os.IsNotExist(err) {
		config = Config{
			Secret:        core.GenerateRandomString(64),
			ReconnectTime: 10 * time.Minute,
			PollTime:      time.Minute,
			Span:          10,
			LastExport:    time.Time{},
		}

		_ = os.MkdirAll(filepath.Dir(path), 0755)
		err = WriteConfig(project, &config)
		return &config, err
	} else {
		logrus.Errorf("cannot read fed config %s: %v", path, err)
		return nil, err
	}
}

func WriteConfig(project *core.Project, config *Config) error {
	path := filepath.Join(project.Path, core.ProjectFedFolder, configFile)
	return core.WriteYaml(path, config)
}
