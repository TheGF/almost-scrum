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
	templates := []string {"Empty"}
	templates = append(templates, core.ListProjectTemplates()...)

	for true {
		prompt := promptui.Select{
			Label: "Choose the project template",
			Items: templates,
		}

		_, s, _ := prompt.Run()
		selected = append(selected, s)
	}
	return selected
}

func processInit(projectPath string, _ []string) {

	projectPath, err := filepath.Abs(projectPath)
	abortIf(err, "")

	_, err = os.Stat(filepath.Join(projectPath, core.GitFolder))
	if err == nil {
		if confirmAction("Found a Git repository in %s. Do you want to connect Ash to Git?", projectPath) {
			projectPath = filepath.Join(projectPath, core.ProjectFolder)
		}
	}

	if !confirmAction("Do you want to create a project in %s", projectPath) {
		return
	}

	templates := chooseTemplates()
	var project *core.Project
	if len(templates) == 0 {
		project, err = core.InitProject(projectPath)
	} else {
		project, err = core.InitProjectFromTemplate(projectPath, templates)
	}
	abortIf(err, "")

	//config, err := core.ReadProjectConfig(projectPath)
	//abortIf(err, "")
	//
	//err = core.WriteProjectConfig(projectPath, &config)
	//abortIf(err, "")

	err = core.SetUserInfo(project, core.GetSystemUser(), &core.UserInfo{})
	abortIf(err, "")

	color.Green("Project initialized successfully in %s", projectPath)
}
