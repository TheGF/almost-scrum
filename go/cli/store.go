package cli

import (
	"almost-scrum/core"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func chooseStore(stores []string) string {

	cursorPos := sort.SearchStrings(stores, config.CurrentStore)
	prompt := promptui.Select{
		Label:     "Select the current store",
		Items:     stores,
		CursorPos: cursorPos,
	}

	_, selected, _ := prompt.Run()
	return selected
}

func processStore(projectPath string, args []string) {
	project := getProject(projectPath)
	stores, _ := core.ListStores(project)
	var selected string

	if len(args) != 0 {
		selected = args[0]
		if sort.SearchStrings(stores, selected) == len(stores) {
			color.Red("Store '%s' does not exist", selected)
			os.Exit(1)
		}
	} else {
		selected = chooseStore(stores)
	}

	config.CurrentStore = selected
	core.SaveConfig(config)
	color.Green("Current store '%s'", selected)
}
