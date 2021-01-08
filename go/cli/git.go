package cli

import (
	"almost-scrum/core"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"strings"
)

func selectTasks(project *core.Project, global bool) []core.TaskInfo {
	var board string
	if global {
		board = ""
	} else {
		board = project.Config.CurrentBoard
	}

	infos, _ := core.ListTasks(project, board, "")
	names := make([]string, 0, len(infos))
	names = append(names, "Done")
	for _, info := range infos {
		names = append(names, info.Name)
	}

	for {
		prompt := promptui.Select{
			Label:    "Choose the tasks for the commit. Select Done to move on",
			Items:    names,
		}
		index, name, err := prompt.Run()
		if err != nil {
			return []core.TaskInfo{}
		}

		if name == "Done" {
			break
		}

		if strings.HasPrefix(name, "✔") {
			names[index] = name[1:]
		} else {
			names[index] = "✔" + name
		}
	}

	selected := []core.TaskInfo{}
	for index, name := range names {
		if strings.HasPrefix(name, "✔") {
			selected = append(selected, infos[index-1])
		}
	} 

	return selected
}

func addProgress(project *core.Project, info core.TaskInfo, commitInfo core.CommitInfo) {
	task, err := core.GetTask(project, info.Board, info.Name)
	if err != nil {
		return
	}

	actions := make([]string, 0, len(task.Parts))
	for _, part := range task.Parts {
		actions = append(actions, part.Description)
	}

	prompt := promptui.SelectWithAdd{
		Label:    "Choose an action or add a new one",
		Items:    actions,
		AddLabel: "Add custom action",
	}

	_, action, err := prompt.Run()
	if err != nil {
		return
	}

	commitInfo.Body = append(commitInfo.Body, core.GitMessagePart{
		Task:     info.Name,
		Progress:  []string{action},
	})
}

func processCommit(projectPath string, global bool, args []string) {
	project := getProject(projectPath)

	commitInfo := core.CommitInfo{
		User: core.GetSystemUser(),
		Body: []core.GitMessagePart{},
	}

	prompt := promptui.Prompt{
		Label: "Enter a Commit Header",
	}
	header, err := prompt.Run()
	abortIf(err)
	commitInfo.Header = header

	selectedTasks := selectTasks(project, global)
	for _, selectedTask := range selectedTasks {
		addProgress(project, selectedTask, commitInfo)
	}

	err = core.GitCommit(project, commitInfo)
	abortIf(err)

	color.Green("You can now use Git push to complete the update to the master")
}