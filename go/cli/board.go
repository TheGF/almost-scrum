package cli

import (
	"almost-scrum/core"
	"os"
	"sort"

	"github.com/fatih/color"
)

func listBoard(project *core.Project, args []string) {
	boards, _ := core.ListBoards(project)
	var selected string

	if len(args) != 0 {
		selected = args[0]
		if sort.SearchStrings(boards, selected) == len(boards) {
			color.Red("Board '%s' does not exist", selected)
			os.Exit(1)
		}
	} else {
		selected = chooseBoard(project)
	}

	config := getProjectConfig(project)
	config.CurrentBoard = selected
	err := core.WriteProjectConfig(project.Path, &config)
	abortIf(err)

	color.Green("Current board '%s'", selected)

}

func processBoard(projectPath string, args []string) {
	project := getProject(projectPath)

	if len(args) == 2 && args[0] == "new" {
		err := core.CreateBoard(project, args[1])
		abortIf(err)
		color.Green("Board %s created", args[1])
	} else {
		listBoard(project, args)
	}
}
