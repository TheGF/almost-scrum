package cli

import (
	"almost-scrum/core"
)

func openEditor(project *core.Project, board string, name string) {
	var editor = config.Editor

	p := core.GetTaskPath(project, board, name)
	openExternalEditor(editor, p)
	_ = core.ReIndex(project)
}
