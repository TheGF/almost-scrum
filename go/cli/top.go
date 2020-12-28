package cli

import (
	"almost-scrum/core"
	"strconv"

	"github.com/fatih/color"
)

func processTop(projectPath string, args []string) {
	n := 7
	if len(args) > 0 {
		n, _ = strconv.Atoi(args[0])
	}

	project := getProject(projectPath)
	board := project.Config.CurrentBoard
	infos, err := core.SearchTask(project, board, true, args...)
	abortIf(err)

	for i, info := range infos {
		color.Yellow(info.Name)
		if i > n {
			break
		}
	}
}
