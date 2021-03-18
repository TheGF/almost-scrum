package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"gopkg.in/yaml.v2"
)

type Sharing struct {
	ProjectId     string `yaml:"projectId"`
	ProjectConfig *core.ProjectConfigPublic `yaml:"projectConfig"`
	FedConfig     *Config `yaml:"fedConfig"`
}

func Share(project *core.Project, config *Config) (key string, token string, err error) {
	sharing := Sharing{
		ProjectId:     project.Config.UUID,
		ProjectConfig: &project.Config.Public,
		FedConfig:     config,
	}


	data, err := yaml.Marshal(&sharing)
	if err != nil {
		return "", "", err
	}

	key = core.GenerateRandomString(16)
	token, err = core.EncryptString(key, string(data))
	if err != nil {
		return "", "", err
	}
	return
}

func ShareWith(project *core.Project, exchanges[] string, removeCredentials bool) (key string, token string, err error) {
	config, err := ReadConfig(project, removeCredentials)
	if err != nil {
		return "", "", err
	}

	var s3Configs []transport.S3Config
	for _, c := range config.S3 {
		if _, found := core.FindStringInSlice(exchanges, c.Name); found {
			s3Configs = append(s3Configs, c)
		}
	}
	config.S3 = s3Configs

	var webDAVConfigs []transport.WebDAVConfig
	for _, c := range config.WebDAV {
		if _, found := core.FindStringInSlice(exchanges, c.Name); found {
			webDAVConfigs = append(webDAVConfigs, c)
		}
	}
	config.WebDAV = webDAVConfigs

	var ftpConfigs []transport.FTPConfig
	for _, c := range config.Ftp {
		if _, found := core.FindStringInSlice(exchanges, c.Name); found {
			ftpConfigs = append(ftpConfigs, c)
		}
	}
	config.Ftp = ftpConfigs

	return Share(project, config)
}