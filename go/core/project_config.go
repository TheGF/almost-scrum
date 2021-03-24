package core

import "strings"

type ProjectConfigPublic struct {
	Name            string              `json:"name" yaml:"name"`
	CurrentBoard    string              `json:"currentStore" yaml:"currentStore"`
	BoardTypes      map[string][]string `json:"boardTypes" yaml:"boardTypes"`
	IncludeLibInGit bool                `json:"includeLibInGit" yaml:"includeLibInGit"`
	UseGitNative    bool                `json:"useGitNative" yaml:"useGitNative"`
}

type ProjectConfig struct {
	CipherKey string                 `json:"cipherKey" yaml:"cipherKey"`
	UUID      string                 `json:"uuid" yaml:"uuid"`
	Public    ProjectConfigPublic    `json:"public" yaml:"public"`
	Settings  map[string]interface{} `json:"settings" yaml:"settings""`
}

func FilterConfigParts(config *ProjectConfig, prefix string) map[string]interface{} {
	var filtered map[string]interface{}
	for key, value := range config.Settings {
		if strings.HasPrefix(key, prefix) {
			filtered[key] = value
		}
	}
	return filtered
}
