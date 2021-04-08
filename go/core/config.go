package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	rand2 "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type LdapConfig struct {
	Host         string
	Base         string
	Port         int
	UseSSL       bool
	BindDN       string
	BindPassword string
	UserFilter   string
}

type ProjectRef struct {
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Folder string `json:"folder"`
}

//Config is the global configuration stored in the user's home directory.
type Config struct {
	Editor       string
	Host         string
	User         string
	Passwords    map[string]string
	Projects     []ProjectRef
	Secret       string
	OwnerLock    bool
	UseGitNative bool
	LdapConfig   *LdapConfig
}

func getSecret() string {
	key := [48]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key[:])
}

func HasGitNative() bool {
	out, err := UseCommand("git", "--version")
	return err == nil && strings.HasPrefix(out, "git")
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	configFolder := filepath.Join(home, ".config")
	_ = os.MkdirAll(configFolder, os.FileMode(0755))
	return filepath.Join(configFolder, "almost-scrum.yaml")
}


func generateHost() string {
	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))

	prefix := hostnames[r.Int() % len(hostnames)]
	num := r.Int31()
	return fmt.Sprintf("%s.%x", prefix, num)
}

//ReadConfig returns the global configuration
func ReadConfig() *Config {
	configPath := getConfigPath()

	var config Config
	err := ReadYaml(configPath, &config)
	if err != nil {
		logrus.Warnf("Cannot read global configuration %s: %v", configPath, err)

		WriteConfig(&Config{
			Editor:       "",
			User:         "",
			Host:         generateHost(),
			Passwords:    map[string]string{},
			Projects:     []ProjectRef{},
			Secret:       getSecret(),
			UseGitNative: HasGitNative(),
		})
		err := SetPassword(GetSystemUser(), "changeme")
		if err != nil {
			logrus.Fatalf("cannot set initial password for user %s in global config", GetSystemUser())
		}
		return ReadConfig()
	} else {
		logrus.Debugf("Successfully loaded config from %s: %v", configPath, config)
	}
	return &config
}

//WriteConfig saves the global configuration
func WriteConfig(config *Config) {
	configPath := getConfigPath()
	err := WriteYaml(configPath, config)
	if err != nil {
		logrus.Panicf("Cannot save global configuration in %s: %v", configPath, err)
		panic(err)
	}
	logrus.Debugf("config saved to %s", configPath)
}

func FindProjInConfigByName(config *Config, name string) *ProjectRef {
	for _, ref := range config.Projects {
		if ref.Name == name {
			return &ref
		}
	}
	return nil
}

func AddProjectRefToConfig(project *Project) bool {
	config := ReadConfig()
	for _, ref := range config.Projects {
		if ref.UUID == project.Config.UUID {
			return false
		}
	}
	config.Projects = append(config.Projects, ProjectRef{
		Name:   project.Config.Public.Name,
		UUID:   project.Config.UUID,
		Folder: project.Path,
	})
	WriteConfig(config)
	return true
}

func FindProjInConfigByUUID(uuid string) *ProjectRef {
	config := ReadConfig()
	for _, ref := range config.Projects {
		if ref.UUID == uuid {
			return &ref
		}
	}
	return nil
}

func DeleteProjFromConfig(uuid string) {
	config := ReadConfig()
	idx := -1
	for i, r := range config.Projects {
		if r.UUID == uuid {
			idx = i
			break
		}
	}
	if idx != -1 {
		config.Projects = append(config.Projects[0:idx], config.Projects[idx+1:]...)
		WriteConfig(config)
	}
}