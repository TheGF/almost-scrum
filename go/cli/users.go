package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"sort"
)

func processUsers(projectPath string, args []string) {
	project := getProject(projectPath)
	config, _ := core.ReadProjectConfig(project.Path)

	switch len(args) {
	case 0:
		for _, user :=  range config.Users {
			color.Green(user)
		}
		return
	case 1:
		return
	}

	cmd := args[0]
	user := args[1]

	switch cmd {
	case "add":
		if idx := sort.SearchStrings(config.Users, user); idx == len(config.Users) {
			config.Users = append(config.Users, user)
			if err := core.WriteProjectConfig(project.Path, &config); err != nil {
				color.Red("Something went write while saving the project config: %v", err)
			}
		}
	}
}
