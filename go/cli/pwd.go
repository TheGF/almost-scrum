package cli

import (
	"almost-scrum/core"
	"errors"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
)

func validate(input string) error {
	if 0 < len(input) && len(input) < 3 {
		return errors.New("password must have more than 2 characters")
	}
	return nil
}

func processPwd(args []string) {
	if len(args) == 0 {
		color.Red("User is required")
		usage()
		os.Exit(1)
	}

	user := args[0]
	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	password, err := prompt.Run()
	abortIf(err)

	err = core.SetPassword(user, password)
	abortIf(err)
}
