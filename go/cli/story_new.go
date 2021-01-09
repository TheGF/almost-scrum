package cli

import (
	"almost-scrum/core"
	"errors"
	"unicode"

	"github.com/manifoldco/promptui"
)

func validateTitle(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) && r != ' '{
			return errors.New("title can only contain letters (no digits, no specials)")
		}
	}
	return nil
}

func chooseTitle(title string) string {
	prompt := promptui.Prompt{
		Label:    "Title",
		Validate: validateTitle,
		Default: title,
	}

	title, _ = prompt.Run()
	return title
}


func processNew(projectPath string, args []string) {
	project := getProject(projectPath)
	board := chooseBoard(project)

	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = chooseTitle("")
	}
	if title == "" {
		return
	}

	_, name, err := core.CreateTask(project, board, title, core.GetSystemUser())
	abortIf(err, "")
	openEditor(project, board, name)
}
