package cli

import (
	"almost-scrum/core"
	"strconv"

	"github.com/fatih/color"
)

func processTop(projectPath string, args []string) {
	project, err := core.FindProject(projectPath, "")
	if err != nil {
		panic(err)
	}

	n := 7
	if len(args) > 0 {
		n, _ = strconv.Atoi(args[0])
	}

	config := core.LoadConfig()
	store, err := core.GetStore(project, config.CurrentStore)

	for _, item := range core.ListStore(store)[0:n] {
		color.Yellow(item.Path)
	}
}
