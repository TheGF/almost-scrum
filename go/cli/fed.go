package cli

import (
	"almost-scrum/core"
	"github.com/code-to-go/fed/transport"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"time"
)


func syncCommand(projectPath string, args []string) {
	project := getProject(projectPath)
	now := time.Now()

	f := project.Fed

	connected := false
	for _, c := range transport.ListExchanges(f.GetTransport()) {
		connected = connected || c == ""
	}
	if !connected {
		color.Red("The project has no active transport. Use 'scrum fed claim' to add exchanges")
		os.Exit(0)
	}

	f.Sync()
	state := f.GetState(now)

	color.Green("Updates")
	for _, update := range state.Updates {
		if update.Deleted {
			color.Green("%s > Deleted by %s", update.Path, update.User)
		} else {
			color.Green("%s > Updated by %s", update.Path, update.User)
		}
	}

	for _, parked := range state.Parked {
		if parked.Deleted {
			color.Red("%s > Deletion by %s parked", parked.Path, parked.User)
		} else {
			color.Red("%s > Update by %s parked", parked.Path, parked.User)
		}
	}

	color.Green("Sent")
	for _, sent := range state.Sent {
		color.Green("%s > Update by %s sent to the federation", sent.Path, sent.User)
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
