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

var emptyTask = core.Task{
	Description: "Replace with the task description",
	Properties:    map[string]string {
		"points": "0",
		"owner": "",
	},
	Parts:       []core.Part{},
	Attachments: []string{},
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

	name := core.NewTaskName(project, title)
	user := core.GetSystemUser()
	task := emptyTask
	task.Properties["owner"] = "@"+user

	err := core.SetTask(project, board, name, &task)
	abortIf(err)
	openEditor(project, board, name)
}
