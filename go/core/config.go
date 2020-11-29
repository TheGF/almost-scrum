package core

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//Config is the global configuration stored in the user's home directory.
type Config struct {
	Editor       string
	Passwords    map[string]string
	Projects     map[string]string
	Secret       string
	CurrentStore string
}

var defaultConfig = Config{
	Editor:       "",
	Passwords:    map[string]string{},
	Projects:     map[string]string{},
	Secret:       getSecret(),
	CurrentStore: "backlog",
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

// SetPassword add a user with password to the global configuration.
func SetPassword(user, password string) error {
	config := LoadConfig()

	if password == "" {
		delete(config.Passwords, user)
		SaveConfig(config)
		return nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Errorf("SetUser - Cannot save user %s and password %s: %v", user, err)
		return err
	}
	config.Passwords[user] = hex.EncodeToString(bytes)
	SaveConfig(config)
	log.Debugf("SetPassword - set password for user %s", user)
	return nil
}

//CheckUser checks if a user has expected password
func CheckUser(user, password string) bool {
	config := LoadConfig()
	hash, _ := hex.DecodeString(config.Passwords[user])

	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
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
	log.Debugf("Config saved to %s", configPath, config)
}
