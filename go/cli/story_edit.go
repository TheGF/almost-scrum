package cli

import (
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"time"
)

func chooseTask(project core.Project, board string, keys... string) core.TaskInfo {
	infos, err := core.SearchTask(project, board, true, keys...)
	abortIf(err)

	choices := make([]string, 0, len(infos))
	for _, info := range infos {
		tm := info.ModTime.Format(time.RFC822)
		choice := fmt.Sprintf("  %-40v%-20v%s", info.Name, info.Board, tm)
		choices = append(choices, choice)
	}
	if len(choices) == 0 {
		color.Yellow("Ops. Look like there are no tasks you can edit")
		return core.TaskInfo{}
	}
	prompt := promptui.Select{
		Label: "Select the task (CTRL+C to exit)",
		Items: choices,
	}
	selected, _, err := prompt.Run()
	if err != nil {
		return core.TaskInfo{}
	}
	return infos[selected]
}

func chooseUser(project core.Project) string {
	prompt := promptui.Select{
		Label: "Select a user (CTRL+C to exit)",
		Items: project.Config.Users,
	}
	_, selected, _ := prompt.Run()
	return selected
}

func processEdit(projectPath string, global bool, args []string) {
	project := getProject(projectPath)
	user := getCurrentUser()
	board := getBoard(project, global)

	filter := append(args, "@"+user)
	info := chooseTask(project, board, filter...)
	if info.Name == "" {
		return
	}

	if _, err := core.GetTask(project, info.Board, info.Name); err != nil {
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

	if info.Name != "" {
		openEditor(project, info.Board, info.Name)
	}
}

func processOwner(projectPath string, global bool, args []string) {
	project := getProject(projectPath)
	board := getBoard(project, global)

	info := chooseTask(project, board, args...)
	if info.Name == "" {
		return
	}

	user := getCurrentUser()
	task, err := core.GetTask(project, info.Board, info.Name)
	if err != nil {
		color.Red( "The task is corrupted. This may happen after a Git pull." +
			" Do you want to fix manually in an editor?")
		prompt := promptui.Prompt{
			Label: "Type 'yes' to confirm",
		}
		answer, _ := prompt.Run()
		if answer == "yes" {
			openEditor(project, info.Board, info.Name)
		}
		return
	}

	owner := task.Properties["owner"]
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
	task.Properties["owner"] = "@"+owner
	abortIf(core.SetTask(project, info.Board, info.Name, &task))
	abortIf(core.ReIndex(project))
	color.Green("Task %s assigned to %s", info.Name, owner)
}

func processTouch(projectPath string, global bool, args []string) {
	project := getProject(projectPath)
	board := getBoard(project, global)

	info := chooseTask(project, board, args...)
	if info.Name == "" {
		return
	}

	if _, err := core.GetTask(project, info.Board, info.Name); err != nil {
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

	if info.Name != "" {
		err := core.TouchTask(project, info.Board, info.Name)
		abortIf(err)
	}
}
