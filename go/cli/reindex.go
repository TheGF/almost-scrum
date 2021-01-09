package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"time"
)

func processReIndex(projectPath string, args []string) {
	project := getProject(projectPath)
	if len(args) == 1 && args[0] == "full" {
		_ = core.ClearIndex(project)
	}
	start := time.Now()
	err := core.ReIndex(project)
	abortIf(err, "")
	elapsed := time.Since(start)

	color.Green("Reindex completed in %s: %d stop words, %d indexes",
		elapsed, len(project.Index.StopWords), len(project.Index.Ids))
}
