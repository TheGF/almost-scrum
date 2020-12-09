package cli

import (
	"almost-scrum/core"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func processEdit(projectPath string, args []string) {
	project := getProject(projectPath)
	store := getCurrentStore(project)

	stores, err := core.WalkStore(store)
	if err != nil {
		color.Red("Something went wrong while reading the detault store")
		os.Exit(1)
	}

	var names []string = make([]string, 0, len(stores))
	filter := ""
	if len(args) != 0 {
		filter = args[0]
	}
	for _, s := range stores {
		if strings.Contains(s.Name, filter) {
			names = append(names, s.Name)
		}
	}
	if len(names) == 0 {
		color.Yellow("No matching story for '%s'", filter)
		os.Exit(1)
	}
	if len(names) == 1 {
		path := core.GetStoryAbsPath(store, names[0])
		openEditor(path)
		return
	}

	prompt := promptui.Select{
		Label: "Select the story (CTRL+C to exit)",
		Items: names,
	}
	_, selected, _ := prompt.Run()
	if selected != "" {
		path := core.GetStoryAbsPath(store, selected)
		openEditor(path)
	}
}
