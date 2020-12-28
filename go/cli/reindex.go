package cli

import "almost-scrum/core"

func processReIndex(projectPath string, args []string) {
	project := getProject(projectPath)
	if len(args) == 1 && args[0] == "full" {
		core.ClearIndex(project)
	}
	err := core.ReIndex(project)
	abortIf(err)
}
