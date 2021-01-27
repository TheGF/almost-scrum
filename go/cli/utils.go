package cli

import (
	"almost-scrum/core"
	"github.com/manifoldco/promptui"
	"os"
	"sort"

	"github.com/fatih/color"
)

func abortIf(err error, msg string) {
	if err != nil {
		if msg == "" {
			msg = "Wow. Something went wrong: %v"
		}
		color.Red(msg, err)
		os.Exit(1)
	}
}

func getProject(projectPath string) *core.Project {
	project, err := core.FindProject(projectPath)
	if err != nil {
		color.Red("No project found. Make sure a project exists in current directory" +
			" or specify a project location with the parameter -p")
		os.Exit(1)
	}

	current := core.GetSystemUser()
	for _, user := range core.GetUserList(project) {
		if user == current {
			return project
		}
	}
	color.Red("I found the project but you are not a user. Bye")
	os.Exit(1)

	return project
}

func getProjectConfig(project *core.Project) core.ProjectConfig {
	config, err := core.ReadProjectConfig(project.Path)
	if err != nil {
		color.Red("No project found. Make sure a project exists in current directory" +
			" or specify a project location with the parameter -p")
		os.Exit(1)
	}
	return config
}

func getBoard(project *core.Project, global bool) string {
	if global {
		return ""
	} else {
		return project.Config.CurrentBoard
	}
}

func chooseBoard(project *core.Project) string {
	boards, err := core.ListBoards(project)
	abortIf(err, "")

	cursorPos := sort.SearchStrings(boards, project.Config.CurrentBoard)
	prompt := promptui.Select{
		Label:     "Choose a board",
		Items:     boards,
		CursorPos: cursorPos,
	}

	_, selected, _ := prompt.Run()
	return selected
}

func confirmAction(message string, a ...interface{}) bool {
	color.Yellow(message, a...)
	prompt := promptui.Prompt{
		Label: "Type yes to confirm",
	}
	confirm, _ := prompt.Run()
	return confirm == "yes"
}
