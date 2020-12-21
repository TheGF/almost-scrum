package cli

import (
	"almost-scrum/core"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func chooseStory(store core.Store, args []string) string {
	stores, err := core.WalkStore(store)
	if err != nil {
		color.Red("Something went wrong while reading the default store")
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
		return names[0]
	}

	prompt := promptui.Select{
		Label: "Select the story (CTRL+C to exit)",
		Items: names,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func chooseUser(project core.Project) string {
	prompt := promptui.Select{
		Label: "Select a user (CTRL+C to exit)",
		Items: project.Users,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func processEdit(projectPath string, args []string) {
	project := getProject(projectPath)
	store := getCurrentStore(project)

	story := chooseStory(store, args)
	if story != "" {
		openEditor(store, story)
	}
}

func processGrant(projectPath string, args []string) {
	project := getProject(projectPath)
	store := getCurrentStore(project)

	path := chooseStory(store, args)
	if path == "" {
		return
	}

	user := core.GetCurrentUser()
	story, err := core.GetStory(store, path)
	if err != nil {
		color.Red("Something went wrong while reading story '%s': %v", path, err)
		os.Exit(1)
	}

	if story.Owner != user {
		prompt := promptui.Prompt{
			Label: "You are not the owner of the story and you should not change ownership." +
				"It may corrupt the story. Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer != "yes" {
			color.Green("Good choice. Ask %s to change the ownership", story.Owner)
			os.Exit(1)
		} else {
			color.Red("Ok, but please inform %s", story.Owner)
		}
	}

	owner := chooseUser(project)
	if owner == "" {
		return
	}
	story.Owner = owner
	core.SetStory(store, path, &story)
	color.Green("Story %s assigned to %s", path, owner)
}
