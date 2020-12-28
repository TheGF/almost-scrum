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

var emptyTask = core.Task{
	Description: "Replace with the task description",
	Features:    map[string]string {
		"points": "0",
		"owner": "",
	},
	Tasks:       []core.Step{},
	TimeEntries: []core.TimeEntry{},
	Attachments: []string{},
}

func processNew(projectPath string, args []string) {
	project := getProject(projectPath)
	board := chooseBoard(project)

	title := getTitle(args)
	if title == "" {
		return
	}

	name := core.NewTaskName(project, title)
	user := getCurrentUser()
	task := emptyTask
	task.Features["owner"] = "@"+user

	err := core.SetTask(project, board, name, &task)
	abortIf(err)
	openEditor(project, board, name)
}
