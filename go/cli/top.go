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

	store := getCurrentStore(getProject(projectPath))
	items, err := core.WalkStore(store)
	if err != nil {
		color.Red("Something went wrong while readind store '%s': %v",
			config.CurrentStore, err)
	}

	for i, item := range items {
		color.Yellow(item.Name)
		if i > n {
			break
		}
	}
}
