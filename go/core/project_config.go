package core

import "strings"

type ProjectConfigPublic struct {
	CurrentBoard    string              `json:"currentStore" yaml:"currentStore"`
	BoardTypes      map[string][]string `json:"boardTypes" yaml:"boardTypes"`
	IncludeLibInGit bool                `json:"includeLibInGit" yaml:"includeLibInGit"`
	UseGitNative    bool                `json:"useGitNative" yaml:"useGitNative"`
}

type ProjectConfig struct {
	CipherKey string                 `yaml:"cipherKey"`
	UUID      string                 `yaml:"uuid"`
	Public    ProjectConfigPublic    `yaml:"public"`
	Settings  map[string]interface{} `yaml:parts`
}

func FilterConfigParts(config *ProjectConfig, prefix string) map[string]interface{}{
	var filtered map[string]interface{}
	for key, value := range config.Settings {
		if strings.HasPrefix(key, prefix) {
			filtered[key] = value
		}
	}
	return filtered
}