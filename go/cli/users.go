package cli

import (
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
)


func processUsers(projectPath string, args []string) {
	project := getProject(projectPath)

	if len(args) == 0 {
		users := core.GetUserList(project)

		for _, user := range users {
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
		abortIf(core.SetUserInfo(project, user, &core.UserInfo{}))
	case "del":
		users := core.GetUserList(project)
		if len(users) == 1 {
			color.Red("The project has only one user. Cannot delete further")
			os.Exit(1)
		}

		prompt := promptui.Select{
			Label: "Choose the user to remove (CTRL+C to exit)",
			Items: users,
		}
		_, user, err := prompt.Run()
		if err == nil {
			abortIf(core.DelUserInfo(project, user))
		}
	}
	color.Green("Users updated!")
}
