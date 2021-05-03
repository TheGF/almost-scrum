package cli

import (
	"almost-scrum/core"
	"github.com/code-to-go/fed"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
)

func printEvent(event fed.Event) {
	switch event.EventType {
	case fed.ExportLog:
		color.Green("-> %s", event.Loc)
	case fed.MergeEvent:
		color.Green("<- %s", event.Loc)
	}
}

func syncCommand(projectPath string, args []string) {
	project := getProject(projectPath)

	f := project.Fed
	f.Watch(printEvent)
	f.Pull()
	f.Push(core.GetSystemUser())

	state := f.GetState()
	if len(state.NetStats) == 0 {
		color.Red("The project has no active transport. Use 'scrum fed claim' to add exchanges")
		os.Exit(0)
	}

	color.Green("Available updates")
	for loc, update := range state.Updates {

		if update.Delete {
			color.Green("%s > Deletion", loc)
			continue
		}

		switch update.State {
		case fed.New:
			color.Green("%s > New file", loc)
		case fed.Newer:
			color.Green("%s > Updated file", loc)
		case fed.Conflict:
			color.Green("%s > Conflict", loc)
		}
	}

}

func joinCommand(projectPath string, args []string) {
	project := getProject(projectPath)

	prompt := promptui.Prompt{
		Label: "Enter the token",
	}
	token, _ := prompt.Run()

	prompt = promptui.Prompt{
		Label: "Enter the decryption key",
	}
	key, _ := prompt.Run()

	if err := core.JoinFed(project, key, token); err != nil {
		color.Red("Cannot join federation from invite: %v", err)
	} else {
		color.Green("Invite accepted. Project %s created in %s",
			project.Config.Public.Name, project.Path)
	}
}

func processFed(projectPath string, args []string) {
	if len(args) == 0 {
		color.Red("subcommand required!")
		os.Exit(1)
	}

	subCmd := args[0]
	switch subCmd {
	case "sync":
		syncCommand(projectPath, args)
	case "claim":
		joinCommand(projectPath, args)
	default:
		color.Red("unknown command %s", subCmd)
	}

}
