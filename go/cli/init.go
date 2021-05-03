package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

func chooseTemplates() []string {
	var selected []string
	templates := []string {"Done"}
	templates = append(templates, core.ListProjectTemplates()...)

	for true {
		prompt := promptui.Select{
			Label: "Choose the project template",
			Items: templates,
		}

		_, s, _ := prompt.Run()
		if s == "Done" {
			return selected
		} else {
			selected = append(selected, s)
		}
	}
	return selected
}

func processInit(projectPath string, _ []string) {

	projectPath, err := filepath.Abs(projectPath)
	abortIf(err, "")

	_, err = os.Stat(filepath.Join(projectPath, core.GitFolder))
	if err == nil {
		if confirmAction("Found a Git repository in %s. Do you want to connect Scrum to Git?", projectPath) {
			projectPath = filepath.Join(projectPath, core.ProjectFolder)
		}
	}

	if !confirmAction("Do you want to create a project in %s", projectPath) {
		return
	}

	templates := chooseTemplates()
	project, err := core.InitProject(projectPath, templates)

	abortIf(err, "")

	err = core.SetUserInfo(project, core.GetSystemUser(), &core.UserInfo{})
	abortIf(err, "")

	color.Green("Project initialized successfully in %s", projectPath)
}
