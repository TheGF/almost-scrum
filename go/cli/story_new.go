package cli

import (
	"almost-scrum/core"
	"errors"
	"os/user"
	"unicode"

	"github.com/manifoldco/promptui"
)

func validateTitle(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return errors.New("Title can only contain letters (no digits, no specials)")
		}
	}
	return nil
}

func getTitle(args []string) string {
	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		prompt := promptui.Prompt{
			Label:    "Title",
			Validate: validateTitle,
		}

		title, _ = prompt.Run()
	}
	return title
}

var emptyStory = core.Story{
	Description: "",
	Points:      0,
	Owner:       "",
	Tasks:       []core.Task{},
	TimeEntries: []core.TimeEntry{},
	Attachments: []string{},
}

func processNew(projectPath string, args []string) {
	title := getTitle(args)
	project := getProject(projectPath)
	store := getCurrentStore(project)
	name := core.GetStoryName(project, title)

	user, _ := user.Current()
	story := emptyStory
	story.Owner = user.Name

	core.SetStory(store, name, &story)
	openEditor(store, name)
}
