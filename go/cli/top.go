package cli

import (
	"almost-scrum/core"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func processTop(projectPath string, global bool, args []string) {
	n := 7
	if len(args) > 0 {
		n, _ = strconv.Atoi(args[0])
	}

	project := getProject(projectPath)
	board := getBoard(project, global)
	infos, err := core.SearchTask(project, board, true, args...)
	abortIf(err)

	color.Green("\n  %-40v%-20v%s", "Task", "Board", "Date")
	for i, info := range infos {
		tm := info.ModTime.Format(time.RFC822)
		color.Yellow("  %-40v%-20v%s", info.Name, info.Board, tm)
		if i > n {
			break
		}
	}
}
