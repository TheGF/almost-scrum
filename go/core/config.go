package core

import (
	"crypto/rand"
	"encoding/hex"
	uuid2 "github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"

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

//Config is the global configuration stored in the user's home directory.
type Config struct {
	Editor       string
	UUID         string
	User         string
	Passwords    map[string]string
	Projects     map[string]string
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
			UUID:         uuid2.New().String(),
			Passwords:    map[string]string{},
			Projects:     map[string]string{},
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
	logrus.Debugf("Config saved to %s", configPath)
}
