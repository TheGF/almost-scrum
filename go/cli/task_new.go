package cli

import (
	"almost-scrum/core"
	"errors"
	"github.com/fatih/color"
	"os"
	"unicode"

	"github.com/manifoldco/promptui"
)

func validateTitle(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) && r != ' ' {
			return errors.New("title can only contain letters (no digits, no specials)")
		}
	}
	return nil
}

func chooseTitle(title string) string {
	prompt := promptui.Prompt{
		Label:    "Title",
		Validate: validateTitle,
		Default:  title,
	}

	title, _ = prompt.Run()
	return title
}

func chooseType(project *core.Project) string {
	names := make([]string, 0)
	for _, model := range project.Models {
		names = append(names, model.Name)
	}
	switch len(names) {
	case 0:
		color.Red("No models available. The project is corrupted")
		os.Exit(1)
	case 1:
		return names[0]
	}

	prompt := promptui.Select{
		Label: "Choose a model",
		Items: names,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func processNew(projectPath string, args []string) {
	project := getProject(projectPath)
	board := chooseBoard(project)
	type_ := chooseType(project)

	var title string
	if len(args) > 0 {
		title = args[0]
	} else {
		title = chooseTitle("")
	}
	if title == "" {
		return
	}

	_, name, err := core.CreateTask(project, board, title, type_, core.GetSystemUser())
	abortIf(err, "")
	openEditor(project, board, name)
}
