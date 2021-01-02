package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"os"
	"path/filepath"
)

func processInit(projectPath string, _ []string) {

	projectPath, err := filepath.Abs(projectPath)
	abortIf(err)

	_, err = os.Stat(filepath.Join(projectPath, core.GitFolder))
	if err == nil {
		if confirmAction("Found a Git repository in %s. Do you want to connect Ash to Git?", projectPath) {
			projectPath = filepath.Join(projectPath, core.ProjectFolder)
		}
	}

	if !confirmAction("Do you want to create a project in %s", projectPath) {
		return
	}
	project, err := core.InitProject(projectPath)
	abortIf(err)

	config, err := core.ReadProjectConfig(projectPath)
	abortIf(err)

	err = core.WriteProjectConfig(projectPath, &config)
	abortIf(err)

	err = core.SetUserInfo(project, core.GetSystemUser(), &core.UserInfo{})
	abortIf(err)

	color.Green("Project initialized successfully in %s", projectPath)
}
