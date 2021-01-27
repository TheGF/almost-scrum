package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
)

func chooseTemplate() string {
	templates := []string {"Empty"}
	templates = append(templates, core.ListProjectTemplates()...)

	prompt := promptui.Select{
		Label: "Choose the project template",
		Items: templates,
	}

	_, selected, _ := prompt.Run()
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

	template := chooseTemplate()
	var project *core.Project
	if template == "Empty" {
		project, err = core.InitProject(projectPath)
	} else {
		project, err = core.InitProjectFromTemplate(projectPath, template)
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
