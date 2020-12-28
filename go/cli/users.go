package cli

import (
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)


func processUsers(projectPath string, args []string) {
	project := getProject(projectPath)

	if len(args) == 0 {
		for _, user := range project.Config.Users {
			fmt.Print(color.GreenString("%s ", user))
		}
		fmt.Println()
		return
	}

	cmd := args[0]
	switch cmd {
	case "add":
		user := ""
		if len(args) == 2 {
			user = args[1]
		} else {
			prompt := promptui.Prompt{Label: "Enter the new user"}
			user, _ = prompt.Run()
		}
		project.Config.Users = append(project.Config.Users, user)
	case "del":
		prompt := promptui.Select{
			Label: "Choose the user to remove (CTRL+C to exit)",
			Items: project.Config.Users,
		}
		selected, _, _ := prompt.Run()
		project.Config.Users = append(project.Config.Users[0:selected], project.Config.Users[selected+1:]...)
	}

	err := core.WriteProjectConfig(project.Path, &project.Config)
	abortIf(err)
	color.Green("Users updated!")
}
