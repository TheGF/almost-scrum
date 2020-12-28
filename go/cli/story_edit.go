package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
)

func chooseTask(project core.Project, board string, keys... string) string {
	infos, err := core.SearchTask(project, board, true, keys...)
	abortIf(err)

	names := make([]string, 0, len(infos))
	for _, info := range infos {
		names = append(names, info.Name)
	}
	if len(names) == 0 {
		color.Yellow("Ops. Look like there are no tasks you can edit")
		return ""
	}
	if len(names) == 1 {
		color.Yellow("Only one task found: %s. Taking you there", names[0])
		return names[0]
	}

	prompt := promptui.Select{
		Label: "Select the task (CTRL+C to exit)",
		Items: names,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func chooseUser(project core.Project) string {
	prompt := promptui.Select{
		Label: "Select a user (CTRL+C to exit)",
		Items: project.Config.Users,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func processEdit(projectPath string, args []string) {
	project := getProject(projectPath)
	user := getCurrentUser()
	board := project.Config.CurrentBoard
	
	filter := append(args, "@"+user)
	name := chooseTask(project, project.Config.CurrentBoard, filter...)
	if name == "" {
		return
	}

	if _, err := core.GetTask(project, board, name); err != nil {
		color.Red( "The story is corrupted. This may happen after a Git pull." +
			" Do you want to continue and open an editor?")
		prompt := promptui.Prompt{
			Label: "Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer != "yes" {
			return
		}
	}

	if name != "" {
		openEditor(project, board, name)
	}
}

func processOwner(projectPath string, args []string) {
	project := getProject(projectPath)
	board := project.Config.CurrentBoard

	name := chooseTask(project, board, args...)
	if name == "" {
		return
	}

	user := getCurrentUser()
	task, err := core.GetTask(project, board, name)
	if err != nil {
		color.Red( "The task is corrupted. This may happen after a Git pull." +
			" Do you want to fix manually in an editor?")
		prompt := promptui.Prompt{
			Label: "Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer == "yes" {
			openEditor(project, board, name)
		}
		return
	}

	owner := task.Features["owner"]
	if owner != "@"+user && owner != ""{
		prompt := promptui.Prompt{
			Label: "You are not the owner of the task and you should not change ownership." +
				"It may corrupt the task. Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer != "yes" {
			color.Green("Good choice. Ask %s to change the ownership", owner)
			os.Exit(1)
		} else {
			color.Red("Ok, but please inform %s", owner)
		}
	}

	owner = chooseUser(project)
	if owner == "" {
		return
	}
	task.Features["owner"] = "@"+owner
	err = core.SetTask(project, board, name, &task)
	abortIf(err)
	err = core.ReIndex(project)
	abortIf(err)

	color.Green("Task %s assigned to %s", name, owner)
}

func processTouch(projectPath string, args []string) {
	project := getProject(projectPath)
	board := project.Config.CurrentBoard

	name := chooseTask(project, board, args...)
	if name == "" {
		return
	}

	if _, err := core.GetTask(project, board, name); err != nil {
		color.Red( "The story is corrupted. This may happen after a Git pull." +
			" Do you want to continue and open an editor?")
		prompt := promptui.Prompt{
			Label: "Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer != "yes" {
			return
		}
	}

	if name != "" {
		err := core.TouchTask(project, board, name)
		abortIf(err)
	}
}
