package cli

import (
	"almost-scrum/core"
	"almost-scrum/fed"
	"github.com/manifoldco/promptui"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func syncCommand(projectPath string, args []string) {
	project := getProject(projectPath)

	since := time.Time{}
	if len(args) > 1 {
		if days, err := strconv.Atoi(args[1]); err != nil {
			since = time.Now().Add(-time.Duration(days) * time.Hour * 24)
		}
	}

	status := fed.GetStatus(project)
	if len(status.Exchanges) == 0 {
		color.Red("The project is not federated. Use 'ash fed claim' to add exchanges")
		os.Exit(0)
	}

	color.Green("Looking for files to export")
	files, err := fed.Export(project, core.GetSystemUser(), time.Time{})
	abortIf(err, "")
	color.Green("The following files will be exported: %s", strings.Join(files, ","))

	if failedExchanges, err := fed.Sync(project, since); err != nil {
		color.Red("cannot sync project with federation: %v", err)
	} else {
		color.Green("synchronization with federation completed, %d exchanges ok, %d failed",
			len(status.Exchanges), failedExchanges)
	}

	diff, err := fed.GetDiffs(project)
	abortIf(err, "")

	print(diff)
	//for _, d := range diff {
	//	d.
	//}
	//
	//color.Green("The following files will be imported:" strings.Join( ","))

	//fed.ImportDiff(project, diff)

}

func claimCommand(projectPath string, args []string) {
	if len(args) == 1 {
		color.Red("destination folder required!")
		os.Exit(1)
	}
	folder := args[1]
	prompt := promptui.Prompt{
		Label: "Enter the token",
	}
	token, _ := prompt.Run()

	prompt = promptui.Prompt{
		Label: "Enter the decryption key",
	}
	key, _ := prompt.Run()

	project, err := fed.ClaimInvite(fed.Invite{Key: key, Token: token}, folder)
	if err != nil {
		color.Red("Cannot create project from invite: %v", err)
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
	case "sync": syncCommand(projectPath, args)
	case "claim": claimCommand(projectPath, args)
	default: color.Red("unknown command %s", subCmd)
	}

}
