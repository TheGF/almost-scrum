package cli

import (
	"almost-scrum/core"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"strings"
	"time"
)

func getTaskRow(info core.TaskInfo, selected bool) string {
	var tick = " "
	if selected {
		tick = "✔"
	}
	tm := info.ModTime.Format(time.RFC822)
	return fmt.Sprintf("%s %-40v%-20v%s",
		tick, info.Name, info.Board, tm)
}

func selectTasks(project *core.Project, global bool) []core.TaskInfo {
	var board string
	if global {
		board = ""
	} else {
		board = project.Config.Public.CurrentBoard
	}

	infos, _ := core.ListTasks(project, board, "")
	rows := make([]string, 0, len(infos))
	rows = append(rows, "Done")
	for _, info := range infos {
		rows = append(rows, getTaskRow(info, false))
	}

	for {
		prompt := promptui.Select{
			Label:        "Choose the tasks you worked on for this commit. Select Done to move on",
			Items:        rows,
			HideSelected: true,
		}
		index, selected, err := prompt.Run()
		if err != nil {
			return []core.TaskInfo{}
		}

		if selected == "Done" {
			break
		}

		rows[index] = getTaskRow(infos[index-1], !strings.HasPrefix(selected, "✔"))
	}

	selected := make([]core.TaskInfo,0)
	for index, name := range rows {
		if strings.HasPrefix(name, "✔") {
			selected = append(selected, infos[index-1])
		}
	}

	return selected
}

func addComment(project *core.Project, info core.TaskInfo, commitInfo core.CommitInfo) {
	task, err := core.GetTask(project, info.Board, info.Name)
	if err != nil {
		return
	}

	color.Green("Comment changes in: %s/%s", info.Board, info.Name)

	if len(task.Parts) > 0 {
		color.Green("Progress Summary")
		for idx, part := range task.Parts {
			color.Green("    %d. %s", idx, part.Description)
		}
	}

	comments := []string{"Generic Progress", "Bugfix"}
	prompt := promptui.SelectWithAdd{
		Label:    "Choose a comment or add a custom",
		Items:    comments,
		AddLabel: "Add custom comment",
	}

	_, comment, err := prompt.Run()
	if err != nil {
		return
	}

	commitInfo.Body[info.Name] = comment
}

func printStatus(status core.GitStatus) {
	untrackedFiles := make([]string, 0)
	color.Green("Changes staged for commit:")
	for file, change := range status.Files {
		if change != core.GitUntracked {
			color.Red("       %s %s", change, file)
		} else {
			untrackedFiles = append(untrackedFiles, file)
		}
	}
	color.Green("\nUntracked files: %s", strings.Join(untrackedFiles, " "))
	color.Green("\nScrum files: %s", strings.Join(status.AshFiles, " "))
}

func processCommit(projectPath string, global bool) {
	project := getProject(projectPath)

	git := core.GetGitClient(project)
	status, err := git.GetStatus(project)
	abortIf(err, "Ops. Something went wrong with your Git Repo. Check integrity with Git."+
		"Error is: %v")

	if len(status.Files) == 0 && len(status.AshFiles) == 0 {
		color.Green("Nothing to commit. Bye")
		return
	}

	printStatus(status)
	commitInfo := core.CommitInfo{
		User: core.GetSystemUser(),
		Body: map[string]string{},
	}

	prompt := promptui.Prompt{
		Label: "Enter a Commit Header",
	}
	header, err := prompt.Run()
	abortIf(err, "")
	commitInfo.Header = header

	selectedTasks := selectTasks(project, global)
	for _, selectedTask := range selectedTasks {
		addComment(project, selectedTask, commitInfo)
	}

	for _, file := range status.AshFiles {
		commitInfo.Files = append(commitInfo.Files, file)
	}
	for file, change := range status.Files {
		if change != core.GitUntracked {
			commitInfo.Files = append(commitInfo.Files, file)
		}
	}

	hash, err := git.Commit(project, commitInfo)
	abortIf(err, "")

	color.Green("Commit completed. Hash: %v", hash)
	color.Green("You can now push the update to the master")
}
