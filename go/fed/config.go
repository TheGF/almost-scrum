package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"almost-scrum/fs"
	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const configFile = "fed.yaml"

type Config struct {
	UUID          string                   `json:"uuid" yaml:"uuid"`
	Secret        string                   `json:"secret" yaml:"secret"`
	ReconnectTime time.Duration            `json:"reconnectTime" yaml:"reconnectTime"`
	PollTime      time.Duration            `json:"pollTime" yaml:"pollTime"`
	Span          int                      `json:"span" yaml:"span"`
	LastExport    time.Time                `json:"lastExport" yaml:"lastExport"`
	LastSync      time.Time                `json:"lastSync" yaml:"lastSync"`
	S3            []transport.S3Config     `json:"s3" yaml:"s3"`
	WebDAV        []transport.WebDAVConfig `json:"webDAV" yaml:"webdav"`
	Ftp           []transport.FTPConfig    `json:"ftp" yaml:"ftp"`
	USB           []transport.USBConfig    `json:"usb" yaml:"usb"`
}

func ReadConfig(project *core.Project, removeSecret bool) (*Config, error) {
	var config Config
	path := filepath.Join(project.Path, core.ProjectFedFolder, configFile)
	if err := fs.ReadYaml(path, &config); err == nil {
		if config.Span < 1 {
			logrus.Warnf("invalid span value %d in %s; default to 10", config.Span, path)
			config.Span = 10
		}

		if config.PollTime < 10*time.Second {
			logrus.Warnf("invalid pool time %s in %s; default to 10m", config.PollTime, path)
			config.PollTime = 10 * time.Minute
		}
	} else if os.IsNotExist(err) {
		config = Config{
			UUID:          uuid2.New().String(),
			Secret:        core.GenerateRandomString(32),
			ReconnectTime: 10 * time.Minute,
			PollTime:      time.Minute,
			Span:          10,
			LastExport:    time.Time{},
			WebDAV: []transport.WebDAVConfig{},
			USB: []transport.USBConfig{},
			Ftp: []transport.FTPConfig{},
			S3: []transport.S3Config{},
		}

		_ = os.MkdirAll(filepath.Dir(path), 0755)
		if err = WriteConfig(project, &config); err != nil {
			return nil, err
		}
	} else {
		logrus.Errorf("cannot read fed config %s: %v", path, err)
		return nil, err
	}

	if !removeSecret {
		return &config, nil
	}

	config.Secret = ""
	transport.RemoveS3Secret(config.S3...)
	transport.RemoveWebDAVSecret(config.WebDAV...)
	transport.RemoveFTPSecret(config.Ftp...)
	return &config, nil
}

func WriteConfig(project *core.Project, config *Config) error {
	path := filepath.Join(project.Path, core.ProjectFedFolder, configFile)
	return fs.WriteYaml(path, config)
}
