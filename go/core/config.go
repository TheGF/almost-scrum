package core

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

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
	User         string
	Passwords    map[string]string
	Projects     map[string]string
	Secret       string
	OwnerLock    bool
	UseGitNative bool
	LdapConfig   *LdapConfig
}

var defaultConfig = Config{
	Editor:       "",
	User:         "",
	Passwords:    map[string]string{},
	Projects:     map[string]string{},
	Secret:       getSecret(),
	UseGitNative: HasGitNative(),
}

func getSecret() string {
	key := [48]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(key[:])
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	configFolder := filepath.Join(home, ".config")
	_ = os.MkdirAll(configFolder, os.FileMode(0755))
	return filepath.Join(configFolder, "almost-scrum.yaml")
}

//LoadConfig returns the global configuration
func LoadConfig() *Config {
	configPath := getConfigPath()

	var config Config
	err := ReadYaml(configPath, &config)
	if err != nil {
		logrus.Warnf("Cannot read global configuration %s: %v", configPath, err)

		SaveConfig(&defaultConfig)
		SetPassword("admin", "changeme")
		return LoadConfig()
	} else {
		logrus.Debugf("Successfully loaded config from %s: %v", configPath, config)
	}
	return &config
}

//SaveConfig saves the global configuration
func SaveConfig(config *Config) {
	configPath := getConfigPath()
	err := WriteYaml(configPath, config)
	if err != nil {
		logrus.Panicf("Cannot save global configuration in %s: %v", configPath, err)
		panic(err)
	}
	logrus.Debugf("Config saved to %s", configPath)
}
