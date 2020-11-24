package core

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

//Config is the global configuration stored in the user's home directory.
type Config struct {
	Editor    string
	Passwords map[string]string
	Projects  map[string]string
	Secret    string
}

var defaultConfig = Config{
	Editor:    "",
	Passwords: map[string]string{},
	Projects:  map[string]string{},
	Secret:    getSecret(),
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
	os.MkdirAll(configFolder, os.FileMode(0755))
	return filepath.Join(configFolder, "almost-scrum.yaml")
}

//LoadConfig returns the global configuration
func LoadConfig() *Config {
	configPath := getConfigPath()

	var config Config
	err := ReadYaml(configPath, &config)
	if err != nil {
		log.Warnf("Cannot read global configuration %s: %v", configPath, err)
		SaveConfig(&defaultConfig)
		return &defaultConfig
	}
	return &config
}

//SaveConfig saves the global configuration
func SaveConfig(config *Config) {
	configPath := getConfigPath()
	err := WriteYaml(configPath, config)
	if err != nil {
		log.Panicf("Cannot save global configuration in %s: %v", configPath, err)
		panic(err)
	}
}
