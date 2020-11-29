package cli

import (
	"almost-scrum/core"
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func listStores(project core.Project) {
	stores, err := core.ListStores(project)
	if err != nil {
		panic(err)
	}
	config := core.LoadConfig()

	fmt.Print("Stores:\n")
	for _, store := range stores {
		if config.CurrentStore == store {
			color.Green("\t%s <current>\n", store)
		} else {
			color.Yellow("\t%s\n", store)
		}
	}
	return
}

func setCurrentStore(project core.Project) {
	config := core.LoadConfig()
	stores, err := core.ListStores(project)
	if err != nil {
		panic(err)
	}
	fmt.Print("Stores:\n")

	cursorPos := sort.Search(len(stores), func(i int) bool {
		return stores[i] == config.CurrentStore
	})
	prompt := promptui.Select{
		Label:     "Select the current store",
		Items:     stores,
		CursorPos: cursorPos,
	}

	_, selected, _ := prompt.Run()

	config.CurrentStore = selected
	core.SaveConfig(config)

}

func processStores(projectPath string, args []string) {
	project, err := core.FindProject(projectPath, "")
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		listStores(project)
		return
	}
	switch args[0] {
	case "cd":
		setCurrentStore(project)
	}
}
