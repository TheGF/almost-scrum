package cli

import (
	"almost-scrum/core"
	"os"

	"github.com/fatih/color"
)

func getProject(projectPath string) core.Project {
	project, err := core.FindProject(projectPath, "")
	if err != nil {
		color.Red("No project found. Make sure a project exists in current directory" +
			" or specify a project location with the parameter -p")
		os.Exit(1)
	}
	return project
}

func getCurrentStore(project core.Project) core.Store {
	store, err := core.GetStore(project, config.CurrentStore)
	if err != nil {
		color.Red("Something went wrong trying to load current store '%s'"+
			"Try to set the default store", config.CurrentStore)
		os.Exit(1)
	}
	return store
}

func listCurrentStore(project core.Project) []core.StoreItem {
	list, err := core.ListStore(getCurrentStore(project), "")
	if err != nil {
		color.Red("Wow. Something went wrong: %v", err)
		os.Exit(1)
	}
	return list
}
