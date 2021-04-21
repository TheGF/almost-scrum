package cli

import (
	"almost-scrum/fed"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func syncCommand(projectPath string, args []string) {
	project := getProject(projectPath)

	//since := time.ModTime{}
	//if len(args) > 1 {
	//	if days, err := strconv.Atoi(args[1]); err != nil {
	//		since = time.Now().Add(-time.Duration(days) * time.Hour * 24)
	//	}
	//}

	status := fed.GetStatus(project)
	if len(status.Exchanges) == 0 {
		color.Red("The project is not federated. Use 'ash fed claim' to add exchanges")
		os.Exit(0)
	}

	color.Green("Looking for files to export")
	stats, err := fed.Pull(project, time.Time{})
	abortIf(err, "")

	for _, stat := range stats {
		if stat.Error != nil {
			color.Red("%s !Error: %v", stat.Error)
			continue
		}

		color.Green("%s [%s] -> %s", strings.Join(stat.Locs, ","))
		if len(stat.Issues) > 0 {
			color.Red("%s [%s] !Issues: %v", stat.Issues)
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

	if err := fed.Join(project, key, token); err != nil {
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
	case "claim": joinCommand(projectPath, args)
	default: color.Red("unknown command %s", subCmd)
	}

}
