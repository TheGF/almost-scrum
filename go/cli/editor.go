package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
)

func openEditor(project *core.Project, board string, name string) {
	var editor = config.Editor

	p := core.GetTaskPath(project, board, name)
	if err := core.RunProgram(editor, p); err != nil {
		color.Red("Something went wrong: %v", err)
		os.Exit(1)
	}

	prompt := promptui.Prompt{Label: "press enter to reindex and complete"}
	prompt.Run()
	_ = core.ReIndex(project)
}
