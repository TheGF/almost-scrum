package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"fmt"
	"github.com/joaojeronimo/go-crc16"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"regexp"
)

type Claim struct {
	UUID        string              `yaml:"uuid"`
	ProjectName string              `yaml:"name"`
	BoardTypes  map[string][]string `yaml:"boardTypes"`
	FedConfig   *Config             `yaml:"fedConfig"`
}

type Invite struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

func CreateInvite(project *core.Project, config *Config) (invite Invite, err error) {
	claim := Claim{
		UUID:        project.Config.UUID,
		ProjectName: project.Config.Public.Name,
		BoardTypes:  project.Config.Public.BoardTypes,
		FedConfig:   config,
	}

	data, err := yaml.Marshal(&claim)
	if err != nil {
		return Invite{}, err
	}

	key := core.GenerateRandomString(16)
	token, err := core.EncryptString(key, string(data))
	if err != nil {
		return Invite{}, err
	}

	//	logrus.Infof("created invite %s", string(data))

	return Invite{
		Key:   key,
		Token: token,
	}, nil
}

func CreateInviteForExchanges(project *core.Project, exchanges []string, removeCredentials bool) (invite Invite, err error) {
	config, err := ReadConfig(project, removeCredentials)
	if err != nil {
		return Invite{}, err
	}

	logrus.Infof("creating invite in project %s with exchanges %v", project.Config.UUID, exchanges)
	var s3Configs []transport.S3Config
	for _, c := range config.S3 {
		if _, found := core.FindStringInSlice(exchanges, c.Name); found {
			s3Configs = append(s3Configs, c)
		}
	}
	config.S3 = s3Configs

	var webDAVConfigs []transport.WebDAVConfig
	for _, c := range config.WebDAV {
		logrus.Infof("name %s", c.Name)
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

	return CreateInvite(project, config)
}

func GetClaimDest(base string, claim *Claim) (dest string, exist bool, err error) {
	if ref := core.FindProjInConfigByUUID(claim.UUID); ref != nil {
		return ref.Folder, true, nil
	}

	dest = filepath.Join(base, claim.ProjectName)
	stat, err := os.Stat(dest)

	if err == nil && stat.IsDir() {
		projectConfig, _ := core.ReadProjectConfig(dest)
		if projectConfig.UUID == claim.UUID {
			return dest, true,nil
		} else {
			crc := crc16.Crc16([]byte(claim.UUID))
			name := fmt.Sprintf("%s-%x", claim.ProjectName, crc)
			dest = filepath.Join(base, name)
			_, err = os.Stat(dest)
			return dest, err == nil, nil
		}
	} else if os.IsNotExist(err) {
		return dest, false, nil
	} else {
		return "", false, err
	}
}

var reWhiteSpace = regexp.MustCompile(`[\s\p{Zs}]+`)

func mergeFedConfig(project *core.Project, c *Config) {
	fedConfig, _ := ReadConfig(project, false)
	fedConfig.Secret = c.Secret
	fedConfig.Span = c.Span
	fedConfig.PollTime = c.PollTime
	fedConfig.ReconnectTime = c.ReconnectTime

	fedConfig.S3 = core.AppendOrUpdate(fedConfig.S3, c.S3, func(i, j int) bool {
		return fedConfig.S3[i].Name == c.S3[j].Name
	}).([]transport.S3Config)

	fedConfig.WebDAV = core.AppendOrUpdate(fedConfig.WebDAV, c.WebDAV, func(i, j int) bool {
		return fedConfig.WebDAV[i] == c.WebDAV[j]
	}).([]transport.WebDAVConfig)

	fedConfig.Ftp = core.AppendOrUpdate(fedConfig.Ftp, c.Ftp, func(i, j int) bool {
		return fedConfig.Ftp[i] == c.Ftp[j]
	}).([]transport.FTPConfig)

	WriteConfig(project, fedConfig)
}

func ClaimInvite(invite Invite, folder string) (*core.Project, error) {
	token := reWhiteSpace.ReplaceAllString(invite.Token, "")
	key := reWhiteSpace.ReplaceAllString(invite.Key, "")

	decrypted, err := core.DecryptString(key, token)
	if err != nil {
		return nil, err
	}
	var claim Claim

	err = yaml.Unmarshal([]byte(decrypted), &claim)
	if err != nil {
		return nil, err
	}

	dest, exist, err  := GetClaimDest(folder, &claim)
	if err != nil {
		return nil, err
	}

	var project *core.Project
	if exist {
		project, err = core.OpenProject(dest)
	} else {
		project, err = core.InitProject(dest, nil)
	}
	if err != nil {
		return nil, err
	}

	project.Config.UUID = claim.UUID
	project.Config.Public.BoardTypes = claim.BoardTypes


	mergeFedConfig(project, claim.FedConfig)

	if exist {
		logrus.Infof("succesfully updated project with %s(%s) from invite",
			project.Config.Public.Name, project.Config.UUID)
	} else {
		logrus.Infof("succesfully created project with %s(%s) from invite",
			project.Config.Public.Name, project.Config.UUID)
	}
	return project, nil
}
